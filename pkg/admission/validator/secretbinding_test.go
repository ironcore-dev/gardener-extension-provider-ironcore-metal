// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package validator_test

import (
	"context"
	"fmt"

	extensionswebhook "github.com/gardener/gardener/extensions/pkg/webhook"
	"github.com/gardener/gardener/pkg/apis/core"
	mockclient "github.com/gardener/gardener/third_party/mock/controller-runtime/client"
	mockmanager "github.com/gardener/gardener/third_party/mock/controller-runtime/manager"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/ironcore-dev/gardener-extension-provider-ironcore-metal/pkg/admission/validator"
)

var _ = Describe("SecretBinding validator", func() {

	Describe("#Validate", func() {
		const (
			namespace = "garden-dev"
			name      = "my-provider-account"
		)

		var (
			secretBindingValidator extensionswebhook.Validator

			ctrl      *gomock.Controller
			apiReader *mockclient.MockReader

			secretBinding = &core.SecretBinding{
				SecretRef: corev1.SecretReference{
					Name:      name,
					Namespace: namespace,
				},
			}
			fakeErr = fmt.Errorf("fake err")

			mgr *mockmanager.MockManager
		)

		BeforeEach(func() {
			ctrl = gomock.NewController(GinkgoT())

			apiReader = mockclient.NewMockReader(ctrl)

			mgr = mockmanager.NewMockManager(ctrl)
			mgr.EXPECT().GetAPIReader().Return(apiReader)

			secretBindingValidator = validator.NewSecretBindingValidator(mgr)
		})

		AfterEach(func() {
			ctrl.Finish()
		})

		It("should return err when obj is not a SecretBinding", func() {
			err := secretBindingValidator.Validate(context.TODO(), &corev1.Secret{}, nil)
			Expect(err).To(MatchError("wrong object type *v1.Secret"))
		})

		It("should return err when oldObj is not a SecretBinding", func() {
			err := secretBindingValidator.Validate(context.TODO(), &core.SecretBinding{}, &corev1.Secret{})
			Expect(err).To(MatchError("wrong object type *v1.Secret for old object"))
		})

		It("should return err if it fails to get the corresponding Secret", func() {
			apiReader.EXPECT().Get(context.TODO(), client.ObjectKey{Namespace: namespace, Name: name}, gomock.AssignableToTypeOf(&corev1.Secret{})).Return(fakeErr)

			err := secretBindingValidator.Validate(context.TODO(), secretBinding, nil)
			Expect(err).To(MatchError(fakeErr))
		})

		It("should return err when the corresponding Secret is not valid", func() {
			apiReader.EXPECT().Get(context.TODO(), client.ObjectKey{Namespace: namespace, Name: name}, gomock.AssignableToTypeOf(&corev1.Secret{})).
				DoAndReturn(func(_ context.Context, _ client.ObjectKey, obj *corev1.Secret, _ ...client.GetOption) error {
					secret := &corev1.Secret{Data: map[string][]byte{
						"namespace": []byte("foo"),
					}}
					*obj = *secret
					return nil
				})

			err := secretBindingValidator.Validate(context.TODO(), secretBinding, nil)
			Expect(err).To(MatchError("missing field: token in cloud provider secret"))
		})

		It("should return nil when the corresponding Secret is valid", func() {
			apiReader.EXPECT().Get(context.TODO(), client.ObjectKey{Namespace: namespace, Name: name}, gomock.AssignableToTypeOf(&corev1.Secret{})).
				DoAndReturn(func(_ context.Context, _ client.ObjectKey, obj *corev1.Secret, _ ...client.GetOption) error {
					secret := &corev1.Secret{Data: map[string][]byte{
						"namespace": []byte("default"),
						"token":     []byte("abcd"),
						"username":  []byte("admin"),
					}}
					*obj = *secret
					return nil
				})

			err := secretBindingValidator.Validate(context.TODO(), secretBinding, nil)
			Expect(err).NotTo(HaveOccurred())
		})
	})

})
