// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package validator

import (
	"context"
	"fmt"

	extensionswebhook "github.com/gardener/gardener/extensions/pkg/webhook"
	"github.com/gardener/gardener/pkg/apis/core"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	metalvalidation "github.com/ironcore-dev/gardener-extension-provider-ironcore-metal/pkg/apis/metal/validation"
)

type secretBinding struct {
	apiReader client.Reader
}

// NewSecretBindingValidator returns a new instance of a secret binding validator.
func NewSecretBindingValidator(mgr manager.Manager) extensionswebhook.Validator {
	return &secretBinding{
		apiReader: mgr.GetAPIReader(),
	}
}

// Validate checks whether the given SecretBinding refers to a Secret with a valid metal service account.
func (sb *secretBinding) Validate(ctx context.Context, newObj, oldObj client.Object) error {
	secretBinding, ok := newObj.(*core.SecretBinding)
	if !ok {
		return fmt.Errorf("wrong object type %T", newObj)
	}

	if oldObj != nil {
		oldSecretBinding, ok := oldObj.(*core.SecretBinding)
		if !ok {
			return fmt.Errorf("wrong object type %T for old object", oldObj)
		}

		// If the provider type did not change, we exit early.
		if oldSecretBinding.Provider != nil && equality.Semantic.DeepEqual(secretBinding.Provider.Type, oldSecretBinding.Provider.Type) {
			return nil
		}
	}

	var (
		secret    = &corev1.Secret{}
		secretKey = client.ObjectKey{Namespace: secretBinding.SecretRef.Namespace, Name: secretBinding.SecretRef.Name}
	)
	// Explicitly use the client.Reader to prevent controller-runtime to start Informer for Secrets
	// under the hood. The latter increases the memory usage of the component.
	if err := sb.apiReader.Get(ctx, secretKey, secret); err != nil {
		return err
	}

	return metalvalidation.ValidateCloudProviderSecret(secret)
}
