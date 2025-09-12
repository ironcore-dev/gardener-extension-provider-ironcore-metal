// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package validator

import (
	extensionswebhook "github.com/gardener/gardener/extensions/pkg/webhook"
	"github.com/gardener/gardener/pkg/apis/core"
	"github.com/gardener/gardener/pkg/apis/core/v1beta1/constants"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/ironcore-dev/gardener-extension-provider-ironcore-metal/pkg/metal"
)

const (
	// Name is a name for a validation webhook.
	Name = "validator"
	// SecretsValidatorName is the name of the secrets' validator.
	SecretsValidatorName = "secrets." + Name
)

var logger = log.Log.WithName("metal-validator-webhook")

// New creates a new validation webhook for `core.gardener.cloud` resources.
func New(mgr manager.Manager) (*extensionswebhook.Webhook, error) {
	logger.Info("Setting up webhook", "name", Name)

	return extensionswebhook.New(mgr, extensionswebhook.Args{
		Provider: metal.Type,
		Name:     Name,
		Path:     "/webhooks/validate",
		Validators: map[extensionswebhook.Validator][]extensionswebhook.Type{
			NewShootValidator(mgr):                  {{Obj: &core.Shoot{}}},
			NewCloudProfileValidator(mgr):           {{Obj: &core.CloudProfile{}}},
			NewNamespacedCloudProfileValidator(mgr): {{Obj: &core.NamespacedCloudProfile{}}},
			NewSecretBindingValidator(mgr):          {{Obj: &core.SecretBinding{}}},
			NewSeedValidator():                      {{Obj: &core.Seed{}}},
			NewBackupBucketValidator():              {{Obj: &core.BackupBucket{}}},
		},
		Target: extensionswebhook.TargetSeed,
		ObjectSelector: &metav1.LabelSelector{
			MatchLabels: map[string]string{constants.LabelExtensionProviderTypePrefix + metal.Type: "true"},
		},
	})
}

// NewSecretsWebhook creates a new validation webhook for Secrets.
func NewSecretsWebhook(mgr manager.Manager) (*extensionswebhook.Webhook, error) {
	logger.Info("Setting up webhook", "name", SecretsValidatorName)

	return extensionswebhook.New(mgr, extensionswebhook.Args{
		Provider: metal.Type,
		Name:     SecretsValidatorName,
		Path:     "/webhooks/validate/secrets",
		Validators: map[extensionswebhook.Validator][]extensionswebhook.Type{
			NewSecretValidator(): {{Obj: &corev1.Secret{}}},
		},
		Target: extensionswebhook.TargetSeed,
		ObjectSelector: &metav1.LabelSelector{
			MatchLabels: map[string]string{constants.LabelExtensionProviderTypePrefix + metal.Type: "true"},
		},
	})
}
