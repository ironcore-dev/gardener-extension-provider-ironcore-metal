// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package validation

import (
	"fmt"

	gardenercore "github.com/gardener/gardener/pkg/apis/core"
	v1beta1constants "github.com/gardener/gardener/pkg/apis/core/v1beta1/constants"
	"github.com/gardener/gardener/pkg/utils"
	gutil "github.com/gardener/gardener/pkg/utils/gardener"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/utils/ptr"
	"k8s.io/utils/strings/slices"

	apismetal "github.com/ironcore-dev/gardener-extension-provider-ironcore-metal/pkg/apis/metal"
)

// ValidateCloudProfileConfig validates a CloudProfileConfig object.
func ValidateCloudProfileConfig(cpConfig *apismetal.CloudProfileConfig, machineImages []gardenercore.MachineImage, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	machineImagesPath := fldPath.Child("machineImages")

	// validate all provider images fields
	for i, machineImage := range cpConfig.MachineImages {
		idxPath := machineImagesPath.Index(i)
		allErrs = append(allErrs, ValidateProviderMachineImage(idxPath, machineImage)...)
	}
	allErrs = append(allErrs, validateProviderImagesMapping(cpConfig.MachineImages, machineImages, field.NewPath("spec").Child("machineImages"))...)

	return allErrs
}

// ValidateProviderMachineImage validates a CloudProfileConfig MachineImages entry.
func ValidateProviderMachineImage(validationPath *field.Path, machineImage apismetal.MachineImages) field.ErrorList {
	allErrs := field.ErrorList{}

	if len(machineImage.Name) == 0 {
		allErrs = append(allErrs, field.Required(validationPath.Child("name"), "must provide a name"))
	}

	if len(machineImage.Versions) == 0 {
		allErrs = append(allErrs, field.Required(validationPath.Child("versions"), fmt.Sprintf("must provide at least one version for machine image %q", machineImage.Name)))
	}
	for j, version := range machineImage.Versions {
		jdxPath := validationPath.Child("versions").Index(j)
		if len(version.Version) == 0 {
			allErrs = append(allErrs, field.Required(jdxPath.Child("version"), "must provide a version"))
		}
		if len(version.Image) == 0 {
			allErrs = append(allErrs, field.Required(jdxPath.Child("image"), "must provide an image"))
		}
		versionArch := ptr.Deref(version.Architecture, v1beta1constants.ArchitectureAMD64)
		if !slices.Contains(v1beta1constants.ValidArchitectures, versionArch) {
			allErrs = append(allErrs, field.NotSupported(jdxPath.Child("architecture"), versionArch, v1beta1constants.ValidArchitectures))
		}
	}

	return allErrs
}

// NewProviderImagesContext creates a new ImagesContext for provider images.
func NewProviderImagesContext(providerImages []apismetal.MachineImages) *gutil.ImagesContext[apismetal.MachineImages, apismetal.MachineImageVersion] {
	return gutil.NewImagesContext(
		utils.CreateMapFromSlice(providerImages, func(mi apismetal.MachineImages) string { return mi.Name }),
		func(mi apismetal.MachineImages) map[string]apismetal.MachineImageVersion {
			return utils.CreateMapFromSlice(mi.Versions, func(v apismetal.MachineImageVersion) string { return providerMachineImageKey(v) })
		},
	)
}

func providerMachineImageKey(v apismetal.MachineImageVersion) string {
	return VersionArchitectureKey(v.Version, ptr.Deref(v.Architecture, v1beta1constants.ArchitectureAMD64))
}

// VersionArchitectureKey returns a key for a version and architecture.
func VersionArchitectureKey(version, architecture string) string {
	return version + "-" + architecture
}

// verify that for each cp image a provider image exists
func validateProviderImagesMapping(cpConfigImages []apismetal.MachineImages, machineImages []gardenercore.MachineImage, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	providerImages := NewProviderImagesContext(cpConfigImages)

	// for each image in the CloudProfile, check if it exists in the CloudProfileConfig
	for idxImage, machineImage := range machineImages {
		if len(machineImage.Versions) == 0 {
			continue
		}
		machineImagePath := fldPath.Index(idxImage)
		if _, existsInParent := providerImages.GetImage(machineImage.Name); !existsInParent {
			allErrs = append(allErrs, field.Required(machineImagePath, fmt.Sprintf("must provide a provider image mapping for image %q", machineImage.Name)))
			continue
		}

		// validate that for each version and architecture of an image in the cloud profile a
		// corresponding provider specific image in the cloud profile config exists
		for versionIdx, version := range machineImage.Versions {
			imageVersionPath := machineImagePath.Child("versions").Index(versionIdx)
			for _, expectedArchitecture := range version.Architectures {
				if _, exists := providerImages.GetImageVersion(machineImage.Name, VersionArchitectureKey(version.Version, expectedArchitecture)); !exists {
					allErrs = append(allErrs, field.Required(imageVersionPath,
						fmt.Sprintf("must provide an image mapping for version %q and architecture: %s", version.Version, expectedArchitecture)))
				}
			}
		}
	}

	return allErrs
}
