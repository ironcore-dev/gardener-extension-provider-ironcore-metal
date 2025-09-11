// SPDX-FileCopyrightText: 2025 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package validator

import (
	"context"
	"fmt"

	extensionswebhook "github.com/gardener/gardener/extensions/pkg/webhook"
	"github.com/gardener/gardener/pkg/apis/security"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	metalvalidation "github.com/ironcore-dev/gardener-extension-provider-ironcore-metal/pkg/apis/metal/validation"
)

type credentialsBinding struct {
	apiReader client.Reader
}

// NewCredentialsBindingValidator returns a new instance of a credentials binding validator.
func NewCredentialsBindingValidator(mgr manager.Manager) extensionswebhook.Validator {
	return &credentialsBinding{
		apiReader: mgr.GetAPIReader(),
	}
}

// Validate checks whether the given CredentialsBinding refers to a Secret with a valid metal service account.
func (cb *credentialsBinding) Validate(ctx context.Context, newObj, oldObj client.Object) error {
	credentialsBinding, ok := newObj.(*security.CredentialsBinding)
	if !ok {
		return fmt.Errorf("wrong object type %T", newObj)
	}

	if oldObj != nil {
		_, ok := oldObj.(*security.CredentialsBinding)
		if !ok {
			return fmt.Errorf("wrong object type %T for old object", oldObj)
		}

		// The relevant fields of the credentials binding are immutable so we can exit early on update
		return nil
	}

	var (
		credentialsKey = client.ObjectKey{Namespace: credentialsBinding.CredentialsRef.Namespace, Name: credentialsBinding.CredentialsRef.Name}
	)
	// Explicitly use the client.Reader to prevent controller-runtime to start Informer for Secrets
	// under the hood. The latter increases the memory usage of the component.
	switch {
	case credentialsBinding.CredentialsRef.APIVersion == corev1.SchemeGroupVersion.String() && credentialsBinding.CredentialsRef.Kind == "Secret":
		secret := &corev1.Secret{}
		if err := cb.apiReader.Get(ctx, credentialsKey, secret); err != nil {
			return err
		}

		return metalvalidation.ValidateCloudProviderSecret(secret)
	default:
		return fmt.Errorf("unsupported credentials reference: version %q, kind %q", credentialsBinding.CredentialsRef.APIVersion, credentialsBinding.CredentialsRef.Kind)
	}
}
