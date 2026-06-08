// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package cloudprovider

import (
	"context"
	"testing"

	extensionscontroller "github.com/gardener/gardener/extensions/pkg/controller"
	"github.com/gardener/gardener/extensions/pkg/webhook/cloudprovider"
	gcontext "github.com/gardener/gardener/extensions/pkg/webhook/context"
	gardencorev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	securityv1alpha1constants "github.com/gardener/gardener/pkg/apis/security/v1alpha1/constants"
	mockclient "github.com/gardener/gardener/third_party/mock/controller-runtime/client"
	mockmanager "github.com/gardener/gardener/third_party/mock/controller-runtime/manager"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/mock/gomock"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/clientcmd"

	api "github.com/ironcore-dev/gardener-extension-provider-ironcore-metal/pkg/apis/metal"
	"github.com/ironcore-dev/gardener-extension-provider-ironcore-metal/pkg/apis/metal/install"
)

const namespace = "test"

var testCloudProfileConfig = &api.CloudProfileConfig{
	TypeMeta:      metav1.TypeMeta{},
	MachineImages: []api.MachineImages{},
	RegionConfigs: []api.RegionConfig{
		{
			Name:                     "foo",
			Server:                   "https://localhost",
			CertificateAuthorityData: []byte("abcd1234"),
		},
	},
}

func TestController(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CloudProvider Webhook Suite")
}

var _ = Describe("Ensurer", func() {
	var (
		ctrl   *gomock.Controller
		ctx    = context.TODO()
		mgr    *mockmanager.MockManager
		c      *mockclient.MockClient
		scheme *runtime.Scheme

		eContextK8s = gcontext.NewInternalGardenContext(
			&extensionscontroller.Cluster{
				CloudProfile: &gardencorev1beta1.CloudProfile{
					Spec: gardencorev1beta1.CloudProfileSpec{
						ProviderConfig: &runtime.RawExtension{
							Object: testCloudProfileConfig,
						},
					},
				},
				Shoot: &gardencorev1beta1.Shoot{
					Spec: gardencorev1beta1.ShootSpec{
						Region: "foo",
						Kubernetes: gardencorev1beta1.Kubernetes{
							Version: "1.26.0",
						},
					},
				},
			},
		)
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		c = mockclient.NewMockClient(ctrl)
		scheme = &runtime.Scheme{}

		mgr = mockmanager.NewMockManager(ctrl)
		mgr.EXPECT().GetClient().Return(c)
		mgr.EXPECT().GetScheme().Return(scheme)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Describe("#EnsureCloudProviderSecret", func() {
		var (
			secret                 *corev1.Secret
			secretWithoutToken     *corev1.Secret
			secretWithoutNamespace *corev1.Secret
			secretWithoutUsername  *corev1.Secret
			ensurer                cloudprovider.Ensurer
		)

		BeforeEach(func() {
			secret = &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: namespace,
					Name:      "cloudprovider",
				},
				Data: map[string][]byte{
					"namespace": []byte("foo"),
					"token":     []byte("bar"),
					"username":  []byte("admin"),
				},
			}

			secretWithoutToken = &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: namespace,
					Name:      "cloudprovider",
				},
				Data: map[string][]byte{
					"token":    []byte("bar"),
					"username": []byte("admin"),
				},
			}

			secretWithoutNamespace = &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: namespace,
					Name:      "cloudprovider",
				},
				Data: map[string][]byte{
					"namespace": []byte("foo"),
					"username":  []byte("admin"),
				},
			}

			secretWithoutUsername = &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: namespace,
					Name:      "cloudprovider",
				},
				Data: map[string][]byte{
					"namespace": []byte("foo"),
					"token":     []byte("bar"),
				},
			}

			ensurer = NewEnsurer(logger, mgr, "")
		})

		It("should add a kubeconfig to the cloudprovider secret", func() {
			err := ensurer.EnsureCloudProviderSecret(ctx, eContextK8s, secret, nil)
			Expect(err).To(Not(HaveOccurred()))

			Expect(secret.Data).To(HaveKey("kubeconfig"))
			config, err := clientcmd.Load(secret.Data["kubeconfig"])
			Expect(err).NotTo(HaveOccurred())
			Expect(config.Clusters[config.CurrentContext].Server).To(Equal("https://localhost"))
			Expect(config.Clusters[config.CurrentContext].CertificateAuthorityData).To(Equal([]byte("abcd1234")))
			Expect(config.AuthInfos["admin"].Token).To(Equal("bar"))
		})

		It("should fail if the cloudprovider secret has no token", func() {
			err := ensurer.EnsureCloudProviderSecret(ctx, eContextK8s, secretWithoutToken, nil)
			Expect(err).To(HaveOccurred())
		})

		It("should fail if the cloudprovider secret has no namespace", func() {
			err := ensurer.EnsureCloudProviderSecret(ctx, eContextK8s, secretWithoutNamespace, nil)
			Expect(err).To(HaveOccurred())
		})

		It("should fail if the cloudprovider secret has no username", func() {
			err := ensurer.EnsureCloudProviderSecret(ctx, eContextK8s, secretWithoutUsername, nil)
			Expect(err).To(HaveOccurred())
		})
	})
})

var _ = Describe("Ensurer (workload identity)", func() {
	var (
		ctrl     *gomock.Controller
		ctx      = context.TODO()
		wiScheme *runtime.Scheme
		wiMgr    *mockmanager.MockManager
		wiC      *mockclient.MockClient
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		wiScheme = runtime.NewScheme()
		Expect(install.AddToScheme(wiScheme)).To(Succeed())

		wiC = mockclient.NewMockClient(ctrl)
		wiMgr = mockmanager.NewMockManager(ctrl)
		wiMgr.EXPECT().GetClient().Return(wiC)
		wiMgr.EXPECT().GetScheme().Return(wiScheme)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Describe("#EnsureCloudProviderSecret", func() {
		wiSecret := func() *corev1.Secret {
			return &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: namespace,
					Name:      "cloudprovider",
					Labels: map[string]string{
						securityv1alpha1constants.LabelPurpose: securityv1alpha1constants.LabelPurposeWorkloadIdentityTokenRequestor,
					},
				},
				Data: map[string][]byte{},
			}
		}

		wiCluster := gcontext.NewInternalGardenContext(
			&extensionscontroller.Cluster{
				CloudProfile: &gardencorev1beta1.CloudProfile{
					Spec: gardencorev1beta1.CloudProfileSpec{
						ProviderConfig: &runtime.RawExtension{
							Object: testCloudProfileConfig,
						},
					},
				},
				Shoot: &gardencorev1beta1.Shoot{
					Spec: gardencorev1beta1.ShootSpec{
						Region: "foo",
					},
				},
			},
		)

		It("should write a tokenFile kubeconfig when workload identity is active", func() {
			ensurer := NewEnsurer(logger, wiMgr, "my-namespace")
			secret := wiSecret()
			Expect(ensurer.EnsureCloudProviderSecret(ctx, wiCluster, secret, nil)).To(Succeed())

			Expect(secret.Data).To(HaveKey("kubeconfig"))
			config, err := clientcmd.Load(secret.Data["kubeconfig"])
			Expect(err).NotTo(HaveOccurred())
			Expect(config.Clusters[config.CurrentContext].Server).To(Equal("https://localhost"))
			Expect(config.Clusters[config.CurrentContext].CertificateAuthorityData).To(Equal([]byte("abcd1234")))
			Expect(config.AuthInfos["workload-identity"].Token).To(BeEmpty())
			Expect(config.AuthInfos["workload-identity"].TokenFile).To(Equal("/var/run/secrets/gardener.cloud/workload-identity/token"))
			Expect(config.Contexts[config.CurrentContext].Namespace).To(Equal("my-namespace"))
		})

		It("should fail when --metal-namespace flag is not set", func() {
			ensurer := NewEnsurer(logger, wiMgr, "")
			secret := wiSecret()
			Expect(ensurer.EnsureCloudProviderSecret(ctx, wiCluster, secret, nil)).To(MatchError(ContainSubstring("metal namespace is not configured")))
		})
	})
})
