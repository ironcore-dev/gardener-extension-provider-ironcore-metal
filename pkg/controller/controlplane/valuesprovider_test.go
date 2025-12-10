// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package controlplane

import (
	"os"
	"path/filepath"

	"github.com/gardener/gardener/extensions/pkg/controller"
	gardencorev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	v1beta1constants "github.com/gardener/gardener/pkg/apis/core/v1beta1/constants"
	extensionsv1alpha1 "github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"
	secretsmanager "github.com/gardener/gardener/pkg/utils/secrets/manager"
	fakesecretsmanager "github.com/gardener/gardener/pkg/utils/secrets/manager/fake"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	fakeclient "sigs.k8s.io/controller-runtime/pkg/client/fake"
	. "sigs.k8s.io/controller-runtime/pkg/envtest/komega"

	apismetal "github.com/ironcore-dev/gardener-extension-provider-ironcore-metal/pkg/apis/metal"
	"github.com/ironcore-dev/gardener-extension-provider-ironcore-metal/pkg/internal"
	"github.com/ironcore-dev/gardener-extension-provider-ironcore-metal/pkg/metal"
)

var _ = Describe("Valueprovider Reconcile", func() {
	ns, vp, cluster := SetupTest()

	var (
		fakeClient         client.Client
		fakeSecretsManager secretsmanager.Interface
	)

	BeforeEach(func(ctx SpecContext) {
		curDir, err := os.Getwd()
		Expect(err).NotTo(HaveOccurred())
		Expect(os.Chdir(filepath.Join("..", "..", ".."))).To(Succeed())
		DeferCleanup(os.Chdir, curDir)

		fakeClient = fakeclient.NewClientBuilder().Build()
		fakeSecretsManager = fakesecretsmanager.New(fakeClient, ns.Name)
		Expect(fakeClient.Create(ctx, &corev1.Secret{ObjectMeta: metav1.ObjectMeta{Name: "cloud-controller-manager-server", Namespace: ns.Name}})).To(Succeed())
	})

	Describe("#GetConfigChartValues", func() {
		It("should return correct config chart values", func(ctx SpecContext) {
			cp := &extensionsv1alpha1.ControlPlane{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "control-plane",
					Namespace: ns.Name,
				},
				Spec: extensionsv1alpha1.ControlPlaneSpec{
					Region: "foo",
					SecretRef: corev1.SecretReference{
						Name:      "my-infra-creds",
						Namespace: ns.Name,
					},
					DefaultSpec: extensionsv1alpha1.DefaultSpec{
						Type: metal.Type,
						ProviderConfig: &runtime.RawExtension{
							Raw: encode(&apismetal.ControlPlaneConfig{
								CloudControllerManager: &apismetal.CloudControllerManagerConfig{
									FeatureGates: map[string]bool{
										"CustomResourceValidation": true,
									},
								},
							}),
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, cp)).To(Succeed())

			By("ensuring that the provider ConfigMap has been created")
			config := &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: ns.Name,
					Name:      internal.CloudProviderConfigMapName,
				},
			}
			Eventually(Get(config)).Should(Succeed())
			Expect(config.Data).To(HaveKey("cloudprovider.conf"))
			cloudProviderConfig := map[string]any{}
			Expect(yaml.Unmarshal([]byte(config.Data["cloudprovider.conf"]), &cloudProviderConfig)).NotTo(HaveOccurred())
			Expect(cloudProviderConfig["clusterName"]).To(Equal(cluster.Name))
			networkingConfig, ok := cloudProviderConfig[metal.CloudControllerManagerNetworkingKeyName].(map[string]any)
			Expect(ok).To(BeTrue())
			Expect(networkingConfig[metal.CloudControllerManagerNodeAddressesConfigKeyName]).To(BeTrue())
		})
	})

	Describe("#GetConfigChartValues", func() {
		It("should return correct config chart values for disabled CCM address config ", func(ctx SpecContext) {
			cp := &extensionsv1alpha1.ControlPlane{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "control-plane",
					Namespace: ns.Name,
				},
				Spec: extensionsv1alpha1.ControlPlaneSpec{
					Region: "foo",
					SecretRef: corev1.SecretReference{
						Name:      "my-infra-creds",
						Namespace: ns.Name,
					},
					DefaultSpec: extensionsv1alpha1.DefaultSpec{
						Type: metal.Type,
						ProviderConfig: &runtime.RawExtension{
							Raw: encode(&apismetal.ControlPlaneConfig{
								CloudControllerManager: &apismetal.CloudControllerManagerConfig{
									FeatureGates: map[string]bool{
										"CustomResourceValidation": true,
									},
									Networking: &apismetal.CloudControllerNetworking{
										ConfigureNodeAddresses: false,
									},
								},
							}),
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, cp)).To(Succeed())

			By("ensuring that the provider ConfigMap has been created")
			config := &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: ns.Name,
					Name:      internal.CloudProviderConfigMapName,
				},
			}
			Eventually(Get(config)).Should(Succeed())
			Expect(config.Data).To(HaveKey("cloudprovider.conf"))
			cloudProviderConfig := map[string]any{}
			Expect(yaml.Unmarshal([]byte(config.Data["cloudprovider.conf"]), &cloudProviderConfig)).NotTo(HaveOccurred())
			Expect(cloudProviderConfig["clusterName"]).To(Equal(cluster.Name))
			networkingConfig, ok := cloudProviderConfig[metal.CloudControllerManagerNetworkingKeyName].(map[string]any)
			Expect(ok).To(BeTrue())
			Expect(networkingConfig[metal.CloudControllerManagerNodeAddressesConfigKeyName]).To(BeFalse())
		})
	})

	Describe("#GetConfigChartValues", func() {
		It("should return correct config chart values for ipamKind address config ", func(ctx SpecContext) {
			cp := &extensionsv1alpha1.ControlPlane{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "control-plane",
					Namespace: ns.Name,
				},
				Spec: extensionsv1alpha1.ControlPlaneSpec{
					Region: "foo",
					SecretRef: corev1.SecretReference{
						Name:      "my-infra-creds",
						Namespace: ns.Name,
					},
					DefaultSpec: extensionsv1alpha1.DefaultSpec{
						Type: metal.Type,
						ProviderConfig: &runtime.RawExtension{
							Raw: encode(&apismetal.ControlPlaneConfig{
								CloudControllerManager: &apismetal.CloudControllerManagerConfig{
									FeatureGates: map[string]bool{
										"CustomResourceValidation": true,
									},
									Networking: &apismetal.CloudControllerNetworking{
										ConfigureNodeAddresses: true,
										IPAMKind: &apismetal.IPAMKind{
											APIGroup: "ag",
											Kind:     "kind",
										},
									},
								},
							}),
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, cp)).To(Succeed())

			By("ensuring that the provider ConfigMap has been created")
			config := &corev1.ConfigMap{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: ns.Name,
					Name:      internal.CloudProviderConfigMapName,
				},
			}
			Eventually(Get(config)).Should(Succeed())
			Expect(config.Data).To(HaveKey("cloudprovider.conf"))
			cloudProviderConfig := map[string]any{}
			Expect(yaml.Unmarshal([]byte(config.Data["cloudprovider.conf"]), &cloudProviderConfig)).NotTo(HaveOccurred())
			Expect(cloudProviderConfig["clusterName"]).To(Equal(cluster.Name))
			networkingConfig, ok := cloudProviderConfig[metal.CloudControllerManagerNetworkingKeyName].(map[string]any)
			Expect(ok).To(BeTrue())
			Expect(networkingConfig[metal.CloudControllerManagerNodeAddressesConfigKeyName]).To(BeTrue())
			ipamKind, ok := networkingConfig[metal.CloudControllerManagerNodeIPAMKindKeyName].(map[string]any)
			Expect(ok).To(BeTrue())
			Expect(ipamKind).To(SatisfyAll(
				HaveKeyWithValue("apiGroup", "ag"),
				HaveKeyWithValue("kind", "kind"),
			))
		})
	})

	Describe("#GetControlPlaneShootCRDsChartValues", func() {
		It("should return correct config chart values", func(ctx SpecContext) {
			values, err := vp.GetControlPlaneShootCRDsChartValues(ctx, nil, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(values).To(Equal(map[string]any{}))
		})
	})

	Describe("#GetControlPlaneChartValues", func() {
		It("should return correct config chart values", func(ctx SpecContext) {
			cp := &extensionsv1alpha1.ControlPlane{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "control-plane",
					Namespace: ns.Name,
				},
				Spec: extensionsv1alpha1.ControlPlaneSpec{
					Region: "foo",
					SecretRef: corev1.SecretReference{
						Name:      "my-infra-creds",
						Namespace: ns.Name,
					},
					DefaultSpec: extensionsv1alpha1.DefaultSpec{
						Type: metal.Type,
						ProviderConfig: &runtime.RawExtension{
							Raw: encode(&apismetal.ControlPlaneConfig{
								CloudControllerManager: &apismetal.CloudControllerManagerConfig{
									FeatureGates: map[string]bool{
										"CustomResourceValidation": true,
									},
								},
								LoadBalancerConfig: &apismetal.LoadBalancerConfig{
									MetalLoadBalancerConfig: &apismetal.MetalLoadBalancerConfig{
										NodeCIDRMask:      80,
										AllocateNodeCIDRs: true,
										VNI:               12345,
										MetalBondServer:   "localhost:8080",
									},
								},
							}),
						},
					},
				},
			}
			providerCloudProfile := &apismetal.CloudProfileConfig{}
			providerCloudProfileJson, err := json.Marshal(providerCloudProfile)
			Expect(err).NotTo(HaveOccurred())
			networkProviderConfig := &unstructured.Unstructured{Object: map[string]any{
				"kind":       "FooNetworkConfig",
				"apiVersion": "v1alpha1",
				"overlay": map[string]any{
					"enabled": false,
				},
			}}
			networkProviderConfigData, err := runtime.Encode(unstructured.UnstructuredJSONScheme, networkProviderConfig)
			Expect(err).NotTo(HaveOccurred())
			cluster := &controller.Cluster{
				CloudProfile: &gardencorev1beta1.CloudProfile{
					Spec: gardencorev1beta1.CloudProfileSpec{
						ProviderConfig: &runtime.RawExtension{
							Raw: providerCloudProfileJson,
						},
					},
				},
				Shoot: &gardencorev1beta1.Shoot{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: ns.Name,
						Name:      "my-shoot",
					},
					Spec: gardencorev1beta1.ShootSpec{
						Networking: &gardencorev1beta1.Networking{
							ProviderConfig: &runtime.RawExtension{Raw: networkProviderConfigData},
							Pods:           ptr.To[string]("10.0.0.0/16"),
						},
						Kubernetes: gardencorev1beta1.Kubernetes{
							Version: "1.26.0",
							VerticalPodAutoscaler: &gardencorev1beta1.VerticalPodAutoscaler{
								Enabled: true,
							},
						},
					},
				},
				Seed: &gardencorev1beta1.Seed{
					ObjectMeta: metav1.ObjectMeta{
						Annotations: map[string]string{
							metal.LocalMetalAPIAnnotation: "true",
						},
					},
				},
			}

			checksums := map[string]string{
				metal.CloudProviderConfigName:            "8bafb35ff1ac60275d62e1cbd495aceb511fb354f74a20f7d06ecb48b3a68432",
				v1beta1constants.SecretNameCloudProvider: "abc",
			}
			values, err := vp.GetControlPlaneChartValues(ctx, cp, cluster, fakeSecretsManager, checksums, false)
			Expect(err).NotTo(HaveOccurred())
			Expect(values).To(Equal(map[string]any{
				"global": map[string]any{
					"genericTokenKubeconfigSecretName": "generic-token-kubeconfig",
				},
				"metal-load-balancer-controller-manager": map[string]any{
					"enabled":           true,
					"nodeCIDRMask":      int32(80),
					"allocateNodeCIDRs": true,
				},
				"cloud-controller-manager": map[string]any{
					"enabled":     true,
					"replicas":    1,
					"clusterName": ns.Name,
					"podAnnotations": map[string]any{
						"checksum/config-cloud-provider-config": "8bafb35ff1ac60275d62e1cbd495aceb511fb354f74a20f7d06ecb48b3a68432",
						"checksum/secret-cloudprovider":         "abc",
					},
					"podLabels": map[string]any{
						"maintenance.gardener.cloud/restart": "true",
						metal.AllowEgressToIstioIngressLabel: "allowed",
					},
					"tlsCipherSuites": []string{
						"TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256",
						"TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384",
						"TLS_AES_128_GCM_SHA256",
						"TLS_AES_256_GCM_SHA384",
						"TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256",
						"TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384",
						"TLS_CHACHA20_POLY1305_SHA256",
						"TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305",
						"TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305",
					},
					"secrets": map[string]any{
						"server": "cloud-controller-manager-server",
					},
					metal.CloudControllerManagerFeatureGatesKeyName: map[string]bool{
						"CustomResourceValidation": true,
					},
					"podNetwork":           "10.0.0.0/16",
					"configureCloudRoutes": true,
				},
			}))
		})
	})

	Describe("#GetControlPlaneShootChartValues", func() {
		It("should return correct shoot system chart values without metallb", func(ctx SpecContext) {
			cp := &extensionsv1alpha1.ControlPlane{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "control-plane",
					Namespace: ns.Name,
				},
				Spec: extensionsv1alpha1.ControlPlaneSpec{
					Region: "foo",
					SecretRef: corev1.SecretReference{
						Name:      "my-infra-creds",
						Namespace: ns.Name,
					},
					DefaultSpec: extensionsv1alpha1.DefaultSpec{
						Type: metal.Type,
						ProviderConfig: &runtime.RawExtension{
							Raw: encode(&apismetal.ControlPlaneConfig{
								CloudControllerManager: &apismetal.CloudControllerManagerConfig{
									FeatureGates: map[string]bool{
										"CustomResourceValidation": true,
									},
								},
							}),
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, cp)).To(Succeed())

			providerCloudProfile := &apismetal.CloudProfileConfig{}
			providerCloudProfileJson, err := json.Marshal(providerCloudProfile)
			Expect(err).NotTo(HaveOccurred())
			networkProviderConfig := &unstructured.Unstructured{Object: map[string]any{
				"kind":       "FooNetworkConfig",
				"apiVersion": "v1alpha1",
				"overlay": map[string]any{
					"enabled": false,
				},
			}}
			networkProviderConfigData, err := runtime.Encode(unstructured.UnstructuredJSONScheme, networkProviderConfig)
			Expect(err).NotTo(HaveOccurred())
			cluster := &controller.Cluster{
				CloudProfile: &gardencorev1beta1.CloudProfile{
					Spec: gardencorev1beta1.CloudProfileSpec{
						ProviderConfig: &runtime.RawExtension{
							Raw: providerCloudProfileJson,
						},
					},
				},
				Shoot: &gardencorev1beta1.Shoot{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: ns.Name,
						Name:      "my-shoot",
					},
					Spec: gardencorev1beta1.ShootSpec{
						Networking: &gardencorev1beta1.Networking{
							ProviderConfig: &runtime.RawExtension{Raw: networkProviderConfigData},
							Pods:           ptr.To[string]("10.0.0.0/16"),
						},
						Kubernetes: gardencorev1beta1.Kubernetes{
							Version: "1.26.0",
							VerticalPodAutoscaler: &gardencorev1beta1.VerticalPodAutoscaler{
								Enabled: true,
							},
						},
					},
				},
			}

			values, err := vp.GetControlPlaneShootChartValues(ctx, cp, cluster, fakeSecretsManager, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(values).To(Equal(map[string]any{
				"cloud-controller-manager": map[string]any{"enabled": true},
				"metallb": map[string]any{
					"enabled": false,
				},
				"calico-bgp": map[string]any{
					"enabled": false,
					"bgp": map[string]any{
						"enabled": false,
					},
				},
				"calico-ipam": map[string]any{
					"enabled": false,
				},
				"metal-load-balancer-controller-speaker": map[string]any{
					"enabled": false,
				},
			}))
		})
	})

	Describe("#GetControlPlaneShootChartValues", func() {
		It("should return correct shoot system chart values with metallb", func(ctx SpecContext) {
			cp := &extensionsv1alpha1.ControlPlane{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "control-plane",
					Namespace: ns.Name,
				},
				Spec: extensionsv1alpha1.ControlPlaneSpec{
					Region: "foo",
					SecretRef: corev1.SecretReference{
						Name:      "my-infra-creds",
						Namespace: ns.Name,
					},
					DefaultSpec: extensionsv1alpha1.DefaultSpec{
						Type: metal.Type,
						ProviderConfig: &runtime.RawExtension{
							Raw: encode(&apismetal.ControlPlaneConfig{
								CloudControllerManager: &apismetal.CloudControllerManagerConfig{
									FeatureGates: map[string]bool{
										"CustomResourceValidation": true,
									},
								},
								LoadBalancerConfig: &apismetal.LoadBalancerConfig{
									MetallbConfig: &apismetal.MetallbConfig{
										IPAddressPool: []string{"10.10.10.0/24", "10.20.20.10-10.20.20.30"},
									},
								},
							}),
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, cp)).To(Succeed())

			providerCloudProfile := &apismetal.CloudProfileConfig{}
			providerCloudProfileJson, err := json.Marshal(providerCloudProfile)
			Expect(err).NotTo(HaveOccurred())
			networkProviderConfig := &unstructured.Unstructured{Object: map[string]any{
				"kind":       "FooNetworkConfig",
				"apiVersion": "v1alpha1",
				"overlay": map[string]any{
					"enabled": false,
				},
			}}
			networkProviderConfigData, err := runtime.Encode(unstructured.UnstructuredJSONScheme, networkProviderConfig)
			Expect(err).NotTo(HaveOccurred())
			cluster := &controller.Cluster{
				CloudProfile: &gardencorev1beta1.CloudProfile{
					Spec: gardencorev1beta1.CloudProfileSpec{
						ProviderConfig: &runtime.RawExtension{
							Raw: providerCloudProfileJson,
						},
					},
				},
				Shoot: &gardencorev1beta1.Shoot{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: ns.Name,
						Name:      "my-shoot",
					},
					Spec: gardencorev1beta1.ShootSpec{
						Networking: &gardencorev1beta1.Networking{
							ProviderConfig: &runtime.RawExtension{Raw: networkProviderConfigData},
							Pods:           ptr.To[string]("10.0.0.0/16"),
						},
						Kubernetes: gardencorev1beta1.Kubernetes{
							Version: "1.26.0",
							VerticalPodAutoscaler: &gardencorev1beta1.VerticalPodAutoscaler{
								Enabled: true,
							},
						},
					},
				},
			}

			values, err := vp.GetControlPlaneShootChartValues(ctx, cp, cluster, fakeSecretsManager, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(values).To(Equal(map[string]any{
				"cloud-controller-manager": map[string]any{"enabled": true},
				"metallb": map[string]any{
					"enabled": true,
					"speaker": map[string]any{
						"enabled": false,
					},
					"l2Advertisement": map[string]any{
						"enabled": false,
					},
					"ipAddressPool": []string{"10.10.10.0/24", "10.20.20.10-10.20.20.30"},
				},
				"calico-bgp": map[string]any{
					"enabled": false,
					"bgp": map[string]any{
						"enabled": false,
					},
				},
				"calico-ipam": map[string]any{
					"enabled": false,
				},
				"metal-load-balancer-controller-speaker": map[string]any{
					"enabled": false,
				},
			}))
		})
	})

	Describe("#GetControlPlaneShootChartValues", func() {
		It("should return correct shoot system chart values with metal-load-balancer-controller", func(ctx SpecContext) {
			cp := &extensionsv1alpha1.ControlPlane{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "control-plane",
					Namespace: ns.Name,
				},
				Spec: extensionsv1alpha1.ControlPlaneSpec{
					Region: "foo",
					SecretRef: corev1.SecretReference{
						Name:      "my-infra-creds",
						Namespace: ns.Name,
					},
					DefaultSpec: extensionsv1alpha1.DefaultSpec{
						Type: metal.Type,
						ProviderConfig: &runtime.RawExtension{
							Raw: encode(&apismetal.ControlPlaneConfig{
								CloudControllerManager: &apismetal.CloudControllerManagerConfig{
									FeatureGates: map[string]bool{
										"CustomResourceValidation": true,
									},
								},
								LoadBalancerConfig: &apismetal.LoadBalancerConfig{
									MetalLoadBalancerConfig: &apismetal.MetalLoadBalancerConfig{
										VNI:             80,
										MetalBondServer: "localhost:8080",
									},
								},
							}),
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, cp)).To(Succeed())

			cluster := &controller.Cluster{
				Shoot: &gardencorev1beta1.Shoot{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: ns.Name,
						Name:      "my-shoot",
					},
					Spec: gardencorev1beta1.ShootSpec{
						Kubernetes: gardencorev1beta1.Kubernetes{
							Version: "1.26.0",
							VerticalPodAutoscaler: &gardencorev1beta1.VerticalPodAutoscaler{
								Enabled: true,
							},
						},
					},
				},
			}

			values, err := vp.GetControlPlaneShootChartValues(ctx, cp, cluster, fakeSecretsManager, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(values).To(Equal(map[string]any{
				"cloud-controller-manager": map[string]any{"enabled": true},
				"metallb": map[string]any{
					"enabled": false,
				},
				"calico-bgp": map[string]any{
					"enabled": false,
					"bgp": map[string]any{
						"enabled": false,
					},
				},
				"calico-ipam": map[string]any{
					"enabled": false,
				},
				"metal-load-balancer-controller-speaker": map[string]any{
					"enabled":         true,
					"vni":             int32(80),
					"metalBondServer": "localhost:8080",
				},
			}))
		})
	})

	Describe("#GetControlPlaneShootChartValues", func() {
		It("should return correct shoot system chart values with calico", func(ctx SpecContext) {
			cp := &extensionsv1alpha1.ControlPlane{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "control-plane",
					Namespace: ns.Name,
				},
				Spec: extensionsv1alpha1.ControlPlaneSpec{
					Region: "foo",
					SecretRef: corev1.SecretReference{
						Name:      "my-infra-creds",
						Namespace: ns.Name,
					},
					DefaultSpec: extensionsv1alpha1.DefaultSpec{
						Type: metal.Type,
						ProviderConfig: &runtime.RawExtension{
							Raw: encode(&apismetal.ControlPlaneConfig{
								CloudControllerManager: &apismetal.CloudControllerManagerConfig{
									FeatureGates: map[string]bool{
										"CustomResourceValidation": true,
									},
								},
								LoadBalancerConfig: &apismetal.LoadBalancerConfig{
									CalicoConfig: &apismetal.CalicoConfig{
										CalicoBgpConfig: &apismetal.CalicoBgpConfig{
											ASNumber:               12345,
											NodeToNodeMeshEnabled:  true,
											ServiceLoadBalancerIPs: []string{"10.10.10.0/24", "10.20.20.10-10.20.20.30"},
											ServiceClusterIPs:      []string{"10.10.10.0/24", "10.20.20.10-10.20.20.30"},
											ServiceExternalIPs:     []string{"10.10.10.0/24", "10.20.20.10-10.20.20.30"},
											BgpPeer: []apismetal.BgpPeer{
												{
													PeerIP:       "1.2.3.4",
													ASNumber:     12345,
													NodeSelector: "foo=bar",
												},
												{
													PeerIP:       "1.2.3.5",
													ASNumber:     12345,
													NodeSelector: "foo=bar",
												},
											},
										},
									},
								},
							}),
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, cp)).To(Succeed())

			providerCloudProfile := &apismetal.CloudProfileConfig{}
			providerCloudProfileJson, err := json.Marshal(providerCloudProfile)
			Expect(err).NotTo(HaveOccurred())
			networkProviderConfig := &unstructured.Unstructured{Object: map[string]any{
				"kind":       "FooNetworkConfig",
				"apiVersion": "v1alpha1",
				"overlay": map[string]any{
					"enabled": false,
				},
			}}
			networkProviderConfigData, err := runtime.Encode(unstructured.UnstructuredJSONScheme, networkProviderConfig)
			Expect(err).NotTo(HaveOccurred())
			cluster := &controller.Cluster{
				CloudProfile: &gardencorev1beta1.CloudProfile{
					Spec: gardencorev1beta1.CloudProfileSpec{
						ProviderConfig: &runtime.RawExtension{
							Raw: providerCloudProfileJson,
						},
					},
				},
				Shoot: &gardencorev1beta1.Shoot{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: ns.Name,
						Name:      "my-shoot",
					},
					Spec: gardencorev1beta1.ShootSpec{
						Networking: &gardencorev1beta1.Networking{
							Type:           ptr.To[string](metal.ShootCalicoNetworkType),
							ProviderConfig: &runtime.RawExtension{Raw: networkProviderConfigData},
							Pods:           ptr.To[string]("10.0.0.0/16"),
						},
						Kubernetes: gardencorev1beta1.Kubernetes{
							Version: "1.26.0",
							VerticalPodAutoscaler: &gardencorev1beta1.VerticalPodAutoscaler{
								Enabled: true,
							},
						},
					},
				},
			}

			values, err := vp.GetControlPlaneShootChartValues(ctx, cp, cluster, fakeSecretsManager, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(values).To(Equal(map[string]any{
				"cloud-controller-manager": map[string]any{"enabled": true},
				"metallb": map[string]any{
					"enabled": false,
				},
				"calico-bgp": map[string]any{
					"enabled": true,
					"bgp": map[string]any{
						"enabled":                true,
						"asNumber":               12345,
						"nodeToNodeMeshEnabled":  true,
						"serviceLoadBalancerIPs": []string{"10.10.10.0/24", "10.20.20.10-10.20.20.30"},
						"serviceExternalIPs":     []string{"10.10.10.0/24", "10.20.20.10-10.20.20.30"},
						"serviceClusterIPs":      []string{"10.10.10.0/24", "10.20.20.10-10.20.20.30"},
						"bgpPeer": []map[string]any{
							{
								"peerIP":       "1.2.3.4",
								"asNumber":     12345,
								"nodeSelector": "foo=bar",
							},
							{
								"peerIP":       "1.2.3.5",
								"asNumber":     12345,
								"nodeSelector": "foo=bar",
							},
						},
					},
				},
				"calico-ipam": map[string]any{
					"enabled": false,
				},
				"metal-load-balancer-controller-speaker": map[string]any{
					"enabled": false,
				},
			}))
		})
	})

	Describe("#GetControlPlaneShootChartValues", func() {
		It("should return correct shoot system chart values with calico bgp filters", func(ctx SpecContext) {
			cp := &extensionsv1alpha1.ControlPlane{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "control-plane",
					Namespace: ns.Name,
				},
				Spec: extensionsv1alpha1.ControlPlaneSpec{
					Region: "foo",
					SecretRef: corev1.SecretReference{
						Name:      "my-infra-creds",
						Namespace: ns.Name,
					},
					DefaultSpec: extensionsv1alpha1.DefaultSpec{
						Type: metal.Type,
						ProviderConfig: &runtime.RawExtension{
							Raw: encode(&apismetal.ControlPlaneConfig{
								CloudControllerManager: &apismetal.CloudControllerManagerConfig{
									FeatureGates: map[string]bool{
										"CustomResourceValidation": true,
									},
								},
								LoadBalancerConfig: &apismetal.LoadBalancerConfig{
									CalicoConfig: &apismetal.CalicoConfig{
										CalicoBgpConfig: &apismetal.CalicoBgpConfig{
											ASNumber:               12345,
											ServiceLoadBalancerIPs: []string{"10.10.10.0/24", "10.20.20.10-10.20.20.30"},
											ServiceClusterIPs:      []string{"10.10.10.0/24", "10.20.20.10-10.20.20.30"},
											ServiceExternalIPs:     []string{"10.10.10.0/24", "10.20.20.10-10.20.20.30"},
											BgpPeer: []apismetal.BgpPeer{
												{
													PeerIP:       "1.2.3.4",
													ASNumber:     12345,
													NodeSelector: "foo=bar",
													Filters: []string{
														"v4filter",
													},
												},
												{
													PeerIP:       "1.2.3.5",
													ASNumber:     12345,
													NodeSelector: "foo=bar",
													Filters: []string{
														"v6filter",
													},
												},
											},
											BGPFilter: []apismetal.BGPFilter{
												{
													Name: "v4filter",
													ImportV4: []apismetal.BGPFilterRule{
														{
															CIDR:          "10.10.10.0/24",
															Action:        "Deny",
															MatchOperator: "Equal",
														},
													},
													ExportV4: []apismetal.BGPFilterRule{
														{
															CIDR:          "10.10.20.0/24",
															Action:        "Deny",
															MatchOperator: "Equal",
														},
													},
												},
												{
													Name: "v6filter",
													ImportV6: []apismetal.BGPFilterRule{
														{
															CIDR:          "fd00:dead:beef:64:34::/80",
															Action:        "Accept",
															MatchOperator: "Equal",
														},
													},
													ExportV6: []apismetal.BGPFilterRule{
														{
															CIDR:          "fd00:dead:beef:64:35::/80",
															Action:        "Accept",
															MatchOperator: "Equal",
														},
													},
												},
											},
										},
									},
								},
							}),
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, cp)).To(Succeed())

			providerCloudProfile := &apismetal.CloudProfileConfig{}
			providerCloudProfileJson, err := json.Marshal(providerCloudProfile)
			Expect(err).NotTo(HaveOccurred())
			networkProviderConfig := &unstructured.Unstructured{Object: map[string]any{
				"kind":       "FooNetworkConfig",
				"apiVersion": "v1alpha1",
				"overlay": map[string]any{
					"enabled": false,
				},
			}}
			networkProviderConfigData, err := runtime.Encode(unstructured.UnstructuredJSONScheme, networkProviderConfig)
			Expect(err).NotTo(HaveOccurred())
			cluster := &controller.Cluster{
				CloudProfile: &gardencorev1beta1.CloudProfile{
					Spec: gardencorev1beta1.CloudProfileSpec{
						ProviderConfig: &runtime.RawExtension{
							Raw: providerCloudProfileJson,
						},
					},
				},
				Shoot: &gardencorev1beta1.Shoot{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: ns.Name,
						Name:      "my-shoot",
					},
					Spec: gardencorev1beta1.ShootSpec{
						Networking: &gardencorev1beta1.Networking{
							Type:           ptr.To[string](metal.ShootCalicoNetworkType),
							ProviderConfig: &runtime.RawExtension{Raw: networkProviderConfigData},
							Pods:           ptr.To[string]("10.0.0.0/16"),
						},
						Kubernetes: gardencorev1beta1.Kubernetes{
							Version: "1.26.0",
							VerticalPodAutoscaler: &gardencorev1beta1.VerticalPodAutoscaler{
								Enabled: true,
							},
						},
					},
				},
			}

			values, err := vp.GetControlPlaneShootChartValues(ctx, cp, cluster, fakeSecretsManager, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(values).To(Equal(map[string]any{
				"cloud-controller-manager": map[string]any{"enabled": true},
				"metallb": map[string]any{
					"enabled": false,
				},
				"metal-load-balancer-controller-speaker": map[string]any{
					"enabled": false,
				},
				"calico-ipam": map[string]any{
					"enabled": false,
				},
				"calico-bgp": map[string]any{
					"enabled": true,
					"bgp": map[string]any{
						"enabled":                true,
						"asNumber":               12345,
						"nodeToNodeMeshEnabled":  false,
						"serviceLoadBalancerIPs": []string{"10.10.10.0/24", "10.20.20.10-10.20.20.30"},
						"serviceExternalIPs":     []string{"10.10.10.0/24", "10.20.20.10-10.20.20.30"},
						"serviceClusterIPs":      []string{"10.10.10.0/24", "10.20.20.10-10.20.20.30"},
						"bgpPeer": []map[string]any{
							{
								"peerIP":       "1.2.3.4",
								"asNumber":     12345,
								"nodeSelector": "foo=bar",
								"filters": []string{
									"v4filter",
								},
							},
							{
								"peerIP":       "1.2.3.5",
								"asNumber":     12345,
								"nodeSelector": "foo=bar",
								"filters": []string{
									"v6filter",
								},
							},
						},
						"bgpFilter": []map[string]any{
							{
								"name": "v4filter",
								"importV4": []map[string]any{
									{
										"cidr":          "10.10.10.0/24",
										"action":        "Deny",
										"matchOperator": "Equal",
									},
								},
								"exportV4": []map[string]any{
									{
										"cidr":          "10.10.20.0/24",
										"action":        "Deny",
										"matchOperator": "Equal",
									},
								},
							},
							{
								"name": "v6filter",
								"importV6": []map[string]any{
									{
										"cidr":          "fd00:dead:beef:64:34::/80",
										"action":        "Accept",
										"matchOperator": "Equal",
									},
								},
								"exportV6": []map[string]any{
									{
										"cidr":          "fd00:dead:beef:64:35::/80",
										"action":        "Accept",
										"matchOperator": "Equal",
									},
								},
							},
						},
					},
				},
			}))
		})
	})

	Describe("#GetControlPlaneShootChartValues", func() {
		It("should return correct shoot system chart values with calico ipam/ippools", func(ctx SpecContext) {
			cp := &extensionsv1alpha1.ControlPlane{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "control-plane",
					Namespace: ns.Name,
				},
				Spec: extensionsv1alpha1.ControlPlaneSpec{
					Region: "foo",
					SecretRef: corev1.SecretReference{
						Name:      "my-infra-creds",
						Namespace: ns.Name,
					},
					DefaultSpec: extensionsv1alpha1.DefaultSpec{
						Type: metal.Type,
						ProviderConfig: &runtime.RawExtension{
							Raw: encode(&apismetal.ControlPlaneConfig{
								CloudControllerManager: &apismetal.CloudControllerManagerConfig{
									FeatureGates: map[string]bool{
										"CustomResourceValidation": true,
									},
								},
								LoadBalancerConfig: &apismetal.LoadBalancerConfig{
									CalicoConfig: &apismetal.CalicoConfig{
										CalicoIPPools: []apismetal.CalicoIPPool{
											{
												CIDR: "10.10.10.0/24",
											},
											{
												CIDR:                       "10.10.20.0/24",
												CalicoIPPoolAssignmentMode: apismetal.CalicoIPPoolAssignmentModeManual,
											},
										},
									},
								},
							}),
						},
					},
				},
			}
			Expect(k8sClient.Create(ctx, cp)).To(Succeed())

			providerCloudProfile := &apismetal.CloudProfileConfig{}
			providerCloudProfileJson, err := json.Marshal(providerCloudProfile)
			Expect(err).NotTo(HaveOccurred())
			networkProviderConfig := &unstructured.Unstructured{Object: map[string]any{
				"kind":       "FooNetworkConfig",
				"apiVersion": "v1alpha1",
				"overlay": map[string]any{
					"enabled": false,
				},
			}}
			networkProviderConfigData, err := runtime.Encode(unstructured.UnstructuredJSONScheme, networkProviderConfig)
			Expect(err).NotTo(HaveOccurred())
			cluster := &controller.Cluster{
				CloudProfile: &gardencorev1beta1.CloudProfile{
					Spec: gardencorev1beta1.CloudProfileSpec{
						ProviderConfig: &runtime.RawExtension{
							Raw: providerCloudProfileJson,
						},
					},
				},
				Shoot: &gardencorev1beta1.Shoot{
					ObjectMeta: metav1.ObjectMeta{
						Namespace: ns.Name,
						Name:      "my-shoot",
					},
					Spec: gardencorev1beta1.ShootSpec{
						Networking: &gardencorev1beta1.Networking{
							Type:           ptr.To[string](metal.ShootCalicoNetworkType),
							ProviderConfig: &runtime.RawExtension{Raw: networkProviderConfigData},
							Pods:           ptr.To[string]("10.0.0.0/16"),
						},
						Kubernetes: gardencorev1beta1.Kubernetes{
							Version: "1.26.0",
							VerticalPodAutoscaler: &gardencorev1beta1.VerticalPodAutoscaler{
								Enabled: true,
							},
						},
					},
				},
			}

			values, err := vp.GetControlPlaneShootChartValues(ctx, cp, cluster, fakeSecretsManager, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(values).To(Equal(map[string]any{
				"cloud-controller-manager": map[string]any{"enabled": true},
				"metallb": map[string]any{
					"enabled": false,
				},
				"calico-ipam": map[string]any{
					"enabled": true,
					"pools": []map[string]any{
						{
							"allowedUses":    []apismetal.CalicoIPPoolAllowedUse{apismetal.CalicoIPPoolAllowedUseLoadBalancer},
							"assignmentMode": apismetal.CalicoIPPoolAssignmentModeAutomatic,
							"cidr":           "10.10.10.0/24",
							"blockSize":      24,
							"disabled":       false,
						},
						{
							"allowedUses":    []apismetal.CalicoIPPoolAllowedUse{apismetal.CalicoIPPoolAllowedUseLoadBalancer},
							"assignmentMode": apismetal.CalicoIPPoolAssignmentModeManual,
							"cidr":           "10.10.20.0/24",
							"blockSize":      24,
							"disabled":       false,
						},
					},
				},
				"calico-bgp": map[string]any{
					"enabled": false,
					"bgp": map[string]any{
						"enabled": false,
					},
				},
				"metal-load-balancer-controller-speaker": map[string]any{
					"enabled": false,
				},
			}))
		})
	})
})

func encode(obj runtime.Object) []byte {
	data, _ := json.Marshal(obj)
	return data
}
