// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package worker

import (
	"context"
	"fmt"

	"github.com/gardener/gardener/extensions/pkg/controller/worker"
	gardencorev1beta1helper "github.com/gardener/gardener/pkg/api/core/v1beta1/helper"
	gardencorev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	v1beta1constants "github.com/gardener/gardener/pkg/apis/core/v1beta1/constants"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"

	metalapi "github.com/ironcore-dev/gardener-extension-provider-ironcore-metal/pkg/apis/metal"
	"github.com/ironcore-dev/gardener-extension-provider-ironcore-metal/pkg/apis/metal/helper"
	apiv1alpha1 "github.com/ironcore-dev/gardener-extension-provider-ironcore-metal/pkg/apis/metal/v1alpha1"
)

// UpdateMachineImagesStatus updates the machine image status
// with the used machine images for the `Worker` resource.
func (w *workerDelegate) UpdateMachineImagesStatus(ctx context.Context) error {
	var machineImages []metalapi.MachineImage

	capabilityDefinitions := helper.NormalizeCapabilityDefinitions(w.cluster.CloudProfile.Spec.MachineCapabilities)

	for _, pool := range w.worker.Spec.Pools {
		arch := ptr.Deref[string](pool.Architecture, v1beta1constants.ArchitectureAMD64)

		machineTypeFromCloudProfile := gardencorev1beta1helper.FindMachineTypeByName(w.cluster.CloudProfile.Spec.MachineTypes, pool.MachineType)
		if machineTypeFromCloudProfile == nil {
			return fmt.Errorf("machine type %q not found in cloud profile %q", pool.MachineType, w.cluster.CloudProfile.Name)
		}

		machineTypeCapabilities := helper.NormalizeMachineTypeCapabilities(machineTypeFromCloudProfile.Capabilities, &arch, capabilityDefinitions)

		machineImage, err := w.selectMachineImageForWorkerPool(pool.MachineImage.Name, pool.MachineImage.Version, &arch, machineTypeCapabilities, capabilityDefinitions)
		if err != nil {
			return err
		}

		machineImages = ensureUniformMachineImages(machineImages, capabilityDefinitions)

		machineImages = appendMachineImage(machineImages, *machineImage, capabilityDefinitions)
	}

	// Decode the current worker provider status.
	workerStatus, err := w.decodeWorkerProviderStatus()
	if err != nil {
		return fmt.Errorf("unable to decode the worker provider status: %w", err)
	}

	workerStatus.MachineImages = machineImages

	return w.updateWorkerProviderStatus(ctx, workerStatus)
}

func (w *workerDelegate) selectMachineImageForWorkerPool(
	name, version string, workerArchitecture *string, machineCapabilities gardencorev1beta1.Capabilities,
	capabilityDefinitions []gardencorev1beta1.CapabilityDefinition,
) (*metalapi.MachineImage, error) {
	selectedMachineImage := &metalapi.MachineImage{
		Name:         name,
		Version:      version,
		Architecture: workerArchitecture,
	}
	if capabilitySet, err := helper.FindImageInCloudProfile(w.cloudProfileConfig, name, version, workerArchitecture, machineCapabilities, capabilityDefinitions); err == nil {
		selectedMachineImage.Capabilities = capabilitySet.Capabilities
		selectedMachineImage.Image = capabilitySet.Image
		return selectedMachineImage, nil
	}
	// Try to look up machine image in worker provider status as it was not found in componentconfig.
	if providerStatus := w.worker.Status.ProviderStatus; providerStatus != nil {
		workerStatus := &metalapi.WorkerStatus{}
		if _, _, err := w.decoder.Decode(providerStatus.Raw, nil, workerStatus); err != nil {
			return nil, fmt.Errorf("could not decode worker status of worker '%s': %w", client.ObjectKeyFromObject(w.worker), err)
		}

		return helper.FindImageInWorkerStatus(workerStatus.MachineImages, name, version, machineCapabilities, capabilityDefinitions)
	}

	return nil, worker.ErrorMachineImageNotFound(name, version, *workerArchitecture)
}

// ensureUniformMachineImages ensures that all machine images are in the same format, either with or without Capabilities.
// Note: The original capabilityDefinition is required to determine which format to append to the worker status.
func ensureUniformMachineImages(images []metalapi.MachineImage, capabilityDefinitions []gardencorev1beta1.CapabilityDefinition) []metalapi.MachineImage {
	var uniformMachineImages []metalapi.MachineImage

	if len(capabilityDefinitions) == 0 {
		// transform images that were added with Capabilities to the legacy format without Capabilities
		for _, img := range images {
			if len(img.Capabilities) == 0 {
				// image is already legacy format
				uniformMachineImages = appendMachineImage(uniformMachineImages, img, capabilityDefinitions)
				continue
			}
			// transform to legacy format by using the Architecture capability if it exists
			var architecture *string
			if len(img.Capabilities[v1beta1constants.ArchitectureName]) > 0 {
				architecture = &img.Capabilities[v1beta1constants.ArchitectureName][0]
			}
			uniformMachineImages = appendMachineImage(uniformMachineImages, metalapi.MachineImage{
				Name:         img.Name,
				Version:      img.Version,
				Image:        img.Image,
				Architecture: architecture,
			}, capabilityDefinitions)
		}
		return uniformMachineImages
	}

	// transform images that were added without Capabilities to contain a MachineImageFlavor with defaulted Architecture
	for _, img := range images {
		if len(img.Capabilities) > 0 {
			// image is already in the new format with Capabilities
			uniformMachineImages = appendMachineImage(uniformMachineImages, img, capabilityDefinitions)
		} else {
			// transform image without Capabilities to capability format with defaulted Architecture
			architecture := ptr.Deref(img.Architecture, v1beta1constants.ArchitectureAMD64)
			uniformMachineImages = appendMachineImage(uniformMachineImages, metalapi.MachineImage{
				Name:         img.Name,
				Version:      img.Version,
				Image:        img.Image,
				Capabilities: gardencorev1beta1.Capabilities{v1beta1constants.ArchitectureName: []string{architecture}},
			}, capabilityDefinitions)
		}
	}
	return uniformMachineImages
}

// appendMachineImage appends a machine image to the list if it doesn't already exist with the same capabilities or architecture.
func appendMachineImage(machineImages []metalapi.MachineImage, machineImage metalapi.MachineImage, capabilityDefinitions []gardencorev1beta1.CapabilityDefinition) []metalapi.MachineImage {
	// support for cloudprofile machine images without capabilities
	if len(capabilityDefinitions) == 0 {
		for _, image := range machineImages {
			isArchEqual := ptr.Deref(image.Architecture, v1beta1constants.ArchitectureAMD64) == ptr.Deref(machineImage.Architecture, v1beta1constants.ArchitectureAMD64)
			if image.Name == machineImage.Name && image.Version == machineImage.Version && isArchEqual {
				// If the image already exists without capabilities, we can just return the existing list.
				return machineImages
			}
		}
		return append(machineImages, metalapi.MachineImage{
			Name:         machineImage.Name,
			Version:      machineImage.Version,
			Image:        machineImage.Image,
			Architecture: machineImage.Architecture,
		})
	}

	defaultedCapabilities := gardencorev1beta1.GetCapabilitiesWithAppliedDefaults(machineImage.Capabilities, capabilityDefinitions)

	for _, existingMachineImage := range machineImages {
		existingDefaultedCapabilities := gardencorev1beta1.GetCapabilitiesWithAppliedDefaults(existingMachineImage.Capabilities, capabilityDefinitions)
		if existingMachineImage.Name == machineImage.Name && existingMachineImage.Version == machineImage.Version && gardencorev1beta1helper.AreCapabilitiesEqual(defaultedCapabilities, existingDefaultedCapabilities) {
			// If the image already exists with the same capabilities return the existing list.
			return machineImages
		}
	}

	// If the image does not exist, we create a new machine image entry with the capabilities.
	return append(machineImages, metalapi.MachineImage{
		Name:         machineImage.Name,
		Version:      machineImage.Version,
		Image:        machineImage.Image,
		Capabilities: machineImage.Capabilities,
	})
}

func (w *workerDelegate) decodeWorkerProviderStatus() (*metalapi.WorkerStatus, error) {
	workerStatus := &metalapi.WorkerStatus{}

	if w.worker.Status.ProviderStatus == nil {
		return workerStatus, nil
	}

	if _, _, err := w.decoder.Decode(w.worker.Status.ProviderStatus.Raw, nil, workerStatus); err != nil {
		return nil, fmt.Errorf("could not decode WorkerStatus '%s': %w", client.ObjectKeyFromObject(w.worker), err)
	}

	return workerStatus, nil
}

func (w *workerDelegate) updateWorkerProviderStatus(ctx context.Context, workerStatus *metalapi.WorkerStatus) error {
	status := &apiv1alpha1.WorkerStatus{}

	if err := w.scheme.Convert(workerStatus, status, nil); err != nil {
		return err
	}

	status.TypeMeta = metav1.TypeMeta{
		APIVersion: apiv1alpha1.SchemeGroupVersion.String(),
		Kind:       "WorkerStatus",
	}

	patch := client.MergeFrom(w.worker.DeepCopy())
	w.worker.Status.ProviderStatus = &runtime.RawExtension{Object: status}
	return w.client.Status().Patch(ctx, w.worker, patch)
}
