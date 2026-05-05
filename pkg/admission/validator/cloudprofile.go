// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
// SPDX-License-Identifier: Apache-2.0

package validator

import (
	"context"
	"fmt"

	extensionswebhook "github.com/gardener/gardener/extensions/pkg/webhook"
	gardencorev1beta1helper "github.com/gardener/gardener/pkg/api/core/v1beta1/helper"
	"github.com/gardener/gardener/pkg/apis/core"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/ironcore-dev/gardener-extension-provider-ironcore-metal/pkg/admission"
	ironcorevalidation "github.com/ironcore-dev/gardener-extension-provider-ironcore-metal/pkg/apis/metal/validation"
)

// NewCloudProfileValidator returns a new instance of a cloud profile validator.
func NewCloudProfileValidator(mgr manager.Manager) extensionswebhook.Validator {
	return &cloudProfileValidator{
		decoder: serializer.NewCodecFactory(mgr.GetScheme(), serializer.EnableStrict).UniversalDecoder(),
	}
}

type cloudProfileValidator struct {
	decoder runtime.Decoder
}

// Validate validates the given CloudProfile objects.
func (cp *cloudProfileValidator) Validate(_ context.Context, newObj, _ client.Object) error {
	cloudProfile, ok := newObj.(*core.CloudProfile)
	if !ok {
		return fmt.Errorf("wrong object type %T", newObj)
	}

	if cloudProfile.Spec.ProviderConfig == nil {
		return field.Required(specPath.Child("providerConfig"), "providerConfig must be set for Ironcore cloud profiles")
	}

	cpConfig, err := admission.DecodeCloudProfileConfig(cp.decoder, cloudProfile.Spec.ProviderConfig)
	if err != nil {
		return fmt.Errorf("could not decode providerConfig of CloudProfile %q: %w", cloudProfile.Name, err)
	}

	capabilityDefinitions, err := gardencorev1beta1helper.ConvertV1beta1CapabilityDefinitions(cloudProfile.Spec.MachineCapabilities)
	if err != nil {
		return field.InternalError(specPath.Child("machineCapabilities"), err)
	}

	return ironcorevalidation.ValidateCloudProfileConfig(cpConfig, cloudProfile.Spec.MachineImages, capabilityDefinitions, specPath).ToAggregate()
}
