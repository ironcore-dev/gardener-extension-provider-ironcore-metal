package metal

import (
	"context"

	v1beta1constants "github.com/gardener/gardener/pkg/apis/core/v1beta1/constants"
	securityv1alpha1constants "github.com/gardener/gardener/pkg/apis/security/v1alpha1/constants"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// IsWorkloadIdentityCloudProviderSecret reports whether the cloudprovider secret was created for WorkloadIdentity credentials.
func IsWorkloadIdentityCloudProviderSecret(secret *corev1.Secret) bool {
	if secret == nil || secret.Labels == nil {
		return false
	}

	return secret.Labels[securityv1alpha1constants.LabelPurpose] == securityv1alpha1constants.LabelPurposeWorkloadIdentityTokenRequestor
}

// IsWorkloadIdentityEnabledForShoot checks the cloudprovider secret on the seed for WorkloadIdentity labels.
func IsWorkloadIdentityEnabledForShoot(ctx context.Context, c client.Reader, shootNamespace string) (bool, error) {
	secret := &corev1.Secret{}
	if err := c.Get(ctx, client.ObjectKey{Namespace: shootNamespace, Name: v1beta1constants.SecretNameCloudProvider}, secret); err != nil {
		return false, client.IgnoreNotFound(err)
	}
	return IsWorkloadIdentityCloudProviderSecret(secret), nil
}
