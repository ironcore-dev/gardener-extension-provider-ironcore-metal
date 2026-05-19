// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package validation

import (
	apivalidation "k8s.io/apimachinery/pkg/api/validation"
	"k8s.io/apimachinery/pkg/util/validation/field"

	metalapi "github.com/ironcore-dev/gardener-extension-provider-ironcore-metal/pkg/apis/metal"
)

// ValidateInfrastructureConfig validates a InfrastructureConfig object.
func ValidateInfrastructureConfig(infra *metalapi.InfrastructureConfig, nodesCIDR, podsCIDR, servicesCIDR *string, workloadIdentity bool, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if workloadIdentity {
		if infra == nil || infra.MetalNamespace == nil {
			allErrs = append(allErrs, field.Required(fldPath.Child("metalNamespace"), "must be set when using WorkloadIdentity auth flow"))
			return allErrs
		}

	}

	if infra != nil && infra.MetalNamespace != nil {
		if errs := apivalidation.ValidateNamespaceName(*infra.MetalNamespace, false); len(errs) > 0 {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("metalNamespace"), *infra.MetalNamespace, ""))
		}
	}

	return allErrs
}

// ValidateInfrastructureConfigUpdate validates a InfrastructureConfig object.
func ValidateInfrastructureConfigUpdate(oldConfig, newConfig *metalapi.InfrastructureConfig, fldPath *field.Path) field.ErrorList {
	return field.ErrorList{}
}
