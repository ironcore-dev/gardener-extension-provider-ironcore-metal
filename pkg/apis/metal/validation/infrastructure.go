// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package validation

import (
	apivalidation "k8s.io/apimachinery/pkg/api/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"

	apismetal "github.com/ironcore-dev/gardener-extension-provider-ironcore-metal/pkg/apis/metal"
)

// ValidateInfrastructureConfig validates a InfrastructureConfig object.
func ValidateInfrastructureConfig(infra *apismetal.InfrastructureConfig, nodesCIDR, podsCIDR, servicesCIDR *string, workloadIdentity bool, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if workloadIdentity {
		if infra == nil || infra.MetalApiNamespace == nil {
			allErrs = append(allErrs, field.Required(fldPath.Child("metalApiNamespace"), "must be set when using WorkloadIdentity auth flow"))
			return allErrs
		}

	}

	if infra != nil && infra.MetalApiNamespace != nil {
		if errs := apivalidation.ValidateNamespaceName(*infra.MetalApiNamespace, false); len(errs) > 0 {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("metalApiNamespace"), *infra.MetalApiNamespace, ""))
		}
	}

	return allErrs
}

// ValidateInfrastructureConfigUpdate validates a InfrastructureConfig object.
func ValidateInfrastructureConfigUpdate(oldConfig, newConfig *apismetal.InfrastructureConfig, fldPath *field.Path) field.ErrorList {
	return field.ErrorList{}
}
