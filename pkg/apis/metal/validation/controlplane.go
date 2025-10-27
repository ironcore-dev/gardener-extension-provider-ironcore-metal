// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package validation

import (
	featurevalidation "github.com/gardener/gardener/pkg/utils/validation/features"
	"k8s.io/apimachinery/pkg/util/validation/field"

	apismetal "github.com/ironcore-dev/gardener-extension-provider-ironcore-metal/pkg/apis/metal"
	"github.com/ironcore-dev/gardener-extension-provider-ironcore-metal/pkg/metal"
)

// ValidateControlPlaneConfig validates a ControlPlaneConfig object.
func ValidateControlPlaneConfig(controlPlaneConfig *apismetal.ControlPlaneConfig, version string, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if controlPlaneConfig.CloudControllerManager != nil {
		allErrs = append(allErrs, featurevalidation.ValidateFeatureGates(controlPlaneConfig.CloudControllerManager.FeatureGates, version, fldPath.Child("cloudControllerManager", metal.CloudControllerManagerFeatureGatesKeyName))...)
	}

	// TODO add validation for IPs

	return allErrs
}

// ValidateControlPlaneConfigUpdate validates a ControlPlaneConfig object.
func ValidateControlPlaneConfigUpdate(oldConfig, newConfig *apismetal.ControlPlaneConfig, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	return allErrs
}

func ValidateCalicoIPPoolAssignmentMode(mode apismetal.CalicoIPPoolAssignmentMode) bool {
	return mode == apismetal.CalicoIPPoolAssignmentModeAutomatic || mode == apismetal.CalicoIPPoolAssignmentModeManual
}

func ValidateCalicoIPPoolAllowedUses(mode apismetal.CalicoIPPoolAllowedUse) bool {
	return mode == apismetal.CalicoIPPoolAllowedUseLoadBalancer || mode == apismetal.CalicoIPPoolAllowedUseTunnel || mode == apismetal.CalicoIPPoolAllowedUseWorkload
}
