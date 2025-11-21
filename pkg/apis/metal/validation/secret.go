// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package validation

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	apivalidation "k8s.io/apimachinery/pkg/api/validation"

	"github.com/ironcore-dev/gardener-extension-provider-ironcore-metal/pkg/metal"
)

// ValidateCloudProviderSecret checks whether the given secret contains a valid metal service account.
func ValidateCloudProviderSecret(secret *corev1.Secret) error {
	if _, ok := secret.Data[metal.TokenFieldName]; !ok {
		return fmt.Errorf("missing field: %s in cloud provider secret", metal.TokenFieldName)
	}
	namespace, ok := secret.Data[metal.NamespaceFieldName]
	if !ok {
		return fmt.Errorf("missing field: %s in cloud provider secret", metal.NamespaceFieldName)
	}
	if _, ok := secret.Data[metal.UsernameFieldName]; !ok {
		return fmt.Errorf("missing field: %s in cloud provider secret", metal.UsernameFieldName)
	}
	errs := apivalidation.ValidateNamespaceName(string(namespace), false)
	if len(errs) > 0 {
		return fmt.Errorf("invalid field: %s in cloud provider secret", metal.NamespaceFieldName)
	}

	return nil
}
