// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package validation

import (
	featurevalidation "github.com/gardener/gardener/pkg/utils/validation/features"
	"k8s.io/apimachinery/pkg/util/validation/field"

	metalapi "github.com/ironcore-dev/gardener-extension-provider-ironcore-metal/pkg/apis/metal"
	"github.com/ironcore-dev/gardener-extension-provider-ironcore-metal/pkg/metal"
)

// ValidateControlPlaneConfig validates a ControlPlaneConfig object.
func ValidateControlPlaneConfig(controlPlaneConfig *metalapi.ControlPlaneConfig, version string, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if controlPlaneConfig.CloudControllerManager != nil {
		allErrs = append(allErrs, featurevalidation.ValidateFeatureGates(controlPlaneConfig.CloudControllerManager.FeatureGates, version, fldPath.Child("cloudControllerManager", metal.CloudControllerManagerFeatureGatesKeyName))...)
		if controlPlaneConfig.CloudControllerManager.PodPrefixSize < 0 {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("cloudControllerManager", "podPrefixSize"), controlPlaneConfig.CloudControllerManager.PodPrefixSize, "must be >= 0"))
		}
	}

	// TODO add validation for IPs

	return allErrs
}

// ValidateControlPlaneConfigUpdate validates a ControlPlaneConfig object.
func ValidateControlPlaneConfigUpdate(oldConfig, newConfig *metalapi.ControlPlaneConfig, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	return allErrs
}

func ValidateCalicoIPPoolAssignmentMode(mode metalapi.CalicoIPPoolAssignmentMode) bool {
	return mode == metalapi.CalicoIPPoolAssignmentModeAutomatic || mode == metalapi.CalicoIPPoolAssignmentModeManual
}

func ValidateCalicoIPPoolAllowedUses(mode metalapi.CalicoIPPoolAllowedUse) bool {
	return mode == metalapi.CalicoIPPoolAllowedUseLoadBalancer || mode == metalapi.CalicoIPPoolAllowedUseTunnel || mode == metalapi.CalicoIPPoolAllowedUseWorkload
}
