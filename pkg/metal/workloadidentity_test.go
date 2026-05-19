// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package metal_test

import (
	"context"

	securityv1alpha1constants "github.com/gardener/gardener/pkg/apis/security/v1alpha1/constants"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	fakeclient "sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/ironcore-dev/gardener-extension-provider-ironcore-metal/pkg/metal"
)

var _ = Describe("Credentials", func() {
	Describe("#IsWorkloadIdentityCloudProviderSecret", func() {
		It("should return false for a nil secret", func() {
			Expect(metal.IsWorkloadIdentityCloudProviderSecret(nil)).To(BeFalse())
		})

		It("should return false for a secret with no labels", func() {
			secret := &corev1.Secret{}
			Expect(metal.IsWorkloadIdentityCloudProviderSecret(secret)).To(BeFalse())
		})

		It("should return false for a secret with an unrelated label", func() {
			secret := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"some-other-label": "value"},
				},
			}
			Expect(metal.IsWorkloadIdentityCloudProviderSecret(secret)).To(BeFalse())
		})

		It("should return false for a secret with the purpose label set to a different value", func() {
			secret := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						securityv1alpha1constants.LabelPurpose: "something-else",
					},
				},
			}
			Expect(metal.IsWorkloadIdentityCloudProviderSecret(secret)).To(BeFalse())
		})

		It("should return true for a secret with the workload-identity-token-requestor label", func() {
			secret := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						securityv1alpha1constants.LabelPurpose: securityv1alpha1constants.LabelPurposeWorkloadIdentityTokenRequestor,
					},
				},
			}
			Expect(metal.IsWorkloadIdentityCloudProviderSecret(secret)).To(BeTrue())
		})
	})

	Describe("#IsWorkloadIdentityEnabledForShoot", func() {
		const shootNamespace = "shoot--project--cluster"

		var (
			ctx    = context.TODO()
			scheme *runtime.Scheme
		)

		BeforeEach(func() {
			scheme = runtime.NewScheme()
			Expect(corev1.AddToScheme(scheme)).To(Succeed())
		})

		It("should return false when the cloudprovider secret does not exist", func() {
			c := fakeclient.NewClientBuilder().WithScheme(scheme).Build()
			enabled, err := metal.IsWorkloadIdentityEnabledForShoot(ctx, c, shootNamespace)
			Expect(err).NotTo(HaveOccurred())
			Expect(enabled).To(BeFalse())
		})

		It("should return false when the cloudprovider secret exists but has no WI label", func() {
			secret := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: shootNamespace,
					Name:      "cloudprovider",
				},
				Data: map[string][]byte{"token": []byte("t"), "namespace": []byte("ns"), "username": []byte("u")},
			}
			c := fakeclient.NewClientBuilder().WithScheme(scheme).WithObjects(secret).Build()
			enabled, err := metal.IsWorkloadIdentityEnabledForShoot(ctx, c, shootNamespace)
			Expect(err).NotTo(HaveOccurred())
			Expect(enabled).To(BeFalse())
		})

		It("should return true when the cloudprovider secret has the WI label", func() {
			secret := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: shootNamespace,
					Name:      "cloudprovider",
					Labels: map[string]string{
						securityv1alpha1constants.LabelPurpose: securityv1alpha1constants.LabelPurposeWorkloadIdentityTokenRequestor,
					},
				},
			}
			c := fakeclient.NewClientBuilder().WithScheme(scheme).WithObjects(secret).Build()
			enabled, err := metal.IsWorkloadIdentityEnabledForShoot(ctx, c, shootNamespace)
			Expect(err).NotTo(HaveOccurred())
			Expect(enabled).To(BeTrue())
		})
	})
})
