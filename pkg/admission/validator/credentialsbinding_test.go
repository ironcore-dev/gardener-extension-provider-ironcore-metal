// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package validator_test

import (
	"context"
	"errors"

	extensionswebhook "github.com/gardener/gardener/extensions/pkg/webhook"
	"github.com/gardener/gardener/pkg/apis/security"
	testutils "github.com/gardener/gardener/pkg/utils/test"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	fakeclient "sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/ironcore-dev/gardener-extension-provider-ironcore-metal/pkg/admission/validator"
)

var _ = Describe("CredentialsBinding validator", func() {
	Describe("#Validate", func() {
		const (
			namespace = "garden-dev"
			name      = "my-provider-account"
		)

		var (
			ctx                = context.TODO()
			credentialsBinding *security.CredentialsBinding

			scheme = func() *runtime.Scheme {
				s := runtime.NewScheme()
				Expect(corev1.AddToScheme(s)).To(Succeed())
				return s
			}()
		)

		BeforeEach(func() {
			credentialsBinding = &security.CredentialsBinding{
				CredentialsRef: corev1.ObjectReference{
					Name:       name,
					Namespace:  namespace,
					Kind:       "Secret",
					APIVersion: "v1",
				},
			}
		})

		newValidator := func(objects ...client.Object) extensionswebhook.Validator {
			b := fakeclient.NewClientBuilder().WithScheme(scheme)
			if len(objects) > 0 {
				b = b.WithObjects(objects...)
			}
			mgr := &testutils.FakeManager{APIReader: b.Build()}
			return validator.NewCredentialsBindingValidator(mgr)
		}

		It("should return err when obj is not a CredentialsBinding", func() {
			v := newValidator()
			err := v.Validate(ctx, &corev1.Secret{}, nil)
			Expect(err).To(MatchError("wrong object type *v1.Secret"))
		})

		It("should return err when oldObj is not a CredentialsBinding", func() {
			v := newValidator()
			err := v.Validate(ctx, &security.CredentialsBinding{}, &corev1.Secret{})
			Expect(err).To(MatchError("wrong object type *v1.Secret for old object"))
		})

		It("should return err if the CredentialsBinding references unknown credentials type", func() {
			v := newValidator()
			credentialsBinding.CredentialsRef.APIVersion = "unknown"
			err := v.Validate(ctx, credentialsBinding, nil)
			Expect(err).To(MatchError(errors.New(`unsupported credentials reference: version "unknown", kind "Secret"`)))
		})

		It("should return err if it fails to get the corresponding Secret", func() {
			// no secret pre-loaded → Get returns NotFound
			v := newValidator()
			err := v.Validate(ctx, credentialsBinding, nil)
			Expect(err).To(HaveOccurred())
		})

		It("should return err when the corresponding Secret is not valid", func() {
			secret := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace},
				Data: map[string][]byte{
					"foo": []byte("bar"),
				},
			}
			v := newValidator(secret)
			err := v.Validate(ctx, credentialsBinding, nil)
			Expect(err).To(HaveOccurred())
		})

		It("should succeed when the corresponding Secret is valid", func() {
			secret := &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{Name: name, Namespace: namespace},
				Data: map[string][]byte{
					"namespace": []byte("default"),
					"token":     []byte("abcd"),
					"username":  []byte("admin"),
				},
			}
			v := newValidator(secret)
			Expect(v.Validate(ctx, credentialsBinding, nil)).To(Succeed())
		})

		It("should succeed when the CredentialsBinding references a WorkloadIdentity", func() {
			credentialsBinding.CredentialsRef.APIVersion = "security.gardener.cloud/v1alpha1"
			credentialsBinding.CredentialsRef.Kind = "WorkloadIdentity"
			v := newValidator()
			Expect(v.Validate(ctx, credentialsBinding, nil)).To(Succeed())
		})

		It("should return nil when the CredentialsBinding did not change", func() {
			v := newValidator()
			old := credentialsBinding.DeepCopy()
			Expect(v.Validate(ctx, credentialsBinding, old)).To(Succeed())
		})
	})
})
