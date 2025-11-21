// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package worker

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	extensionscontroller "github.com/gardener/gardener/extensions/pkg/controller"
	gardencorev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	v1beta1constants "github.com/gardener/gardener/pkg/apis/core/v1beta1/constants"
	gardenerextensionv1alpha1 "github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"
	gardener "github.com/gardener/gardener/pkg/client/kubernetes"
	machinescheme "github.com/gardener/machine-controller-manager/pkg/client/clientset/versioned/scheme"
	"github.com/ironcore-dev/controller-utils/modutils"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/zap/zapcore"
	autoscalingv1 "k8s.io/api/autoscaling/v1"
	corev1 "k8s.io/api/core/v1"
	apiextensionsscheme "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/scheme"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiruntime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	"sigs.k8s.io/controller-runtime/pkg/envtest/komega"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/yaml"

	apiv1alpha1 "github.com/ironcore-dev/gardener-extension-provider-ironcore-metal/pkg/apis/metal/v1alpha1"
)

const (
	pollingInterval      = 50 * time.Millisecond
	eventuallyTimeout    = 10 * time.Second
	consistentlyDuration = 1 * time.Second
)

var (
	testEnv   *envtest.Environment
	cfg       *rest.Config
	k8sClient client.Client
)

// global Gardener resources used by delegates
var (
	technicalID = "shoot--test--cluster"

	shootVersionMajorMinor = "1.2"
	shootVersion           = shootVersionMajorMinor + ".3"
	userDataSecretName     = "userdata-secret-name"
	userDataSecretDataKey  = "userdata-secret-key"

	pool               gardenerextensionv1alpha1.WorkerPool
	cloudProfileConfig *apiv1alpha1.CloudProfileConfig

	clusterWithoutImages   *extensionscontroller.Cluster
	testCluster            *extensionscontroller.Cluster
	cloudProfileConfigJSON []byte

	workerConfig     *apiv1alpha1.WorkerConfig
	workerConfigJSON []byte

	w *gardenerextensionv1alpha1.Worker
)

func TestAPIs(t *testing.T) {
	SetDefaultConsistentlyPollingInterval(pollingInterval)
	SetDefaultEventuallyPollingInterval(pollingInterval)
	SetDefaultEventuallyTimeout(eventuallyTimeout)
	SetDefaultConsistentlyDuration(consistentlyDuration)

	RegisterFailHandler(Fail)
	RunSpecs(t, "Worker Controller Suite")
}

var _ = BeforeSuite(func() {
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true), zap.Level(zapcore.InfoLevel)))

	By("bootstrapping test environment")
	testEnv = &envtest.Environment{
		CRDDirectoryPaths: []string{
			modutils.Dir("github.com/gardener/machine-controller-manager", "kubernetes", "crds", "machine.sapcloud.io_machineclasses.yaml"),
			modutils.Dir("github.com/gardener/machine-controller-manager", "kubernetes", "crds", "machine.sapcloud.io_machinedeployments.yaml"),
			modutils.Dir("github.com/gardener/machine-controller-manager", "kubernetes", "crds", "machine.sapcloud.io_machines.yaml"),
			modutils.Dir("github.com/gardener/machine-controller-manager", "kubernetes", "crds", "machine.sapcloud.io_machinesets.yaml"),
			filepath.Join("..", "..", "..", "example", "20-crd-extensions.gardener.cloud_controlplanes.yaml"),
			filepath.Join("..", "..", "..", "example", "20-crd-extensions.gardener.cloud_workers.yaml"),
		},
		ErrorIfCRDPathMissing: true,

		// The BinaryAssetsDirectory is only required if you want to run the tests directly
		// without call the makefile target test. If not informed it will look for the
		// default path defined in controller-apiruntime which is /usr/local/kubebuilder/.
		// Note that you must have the required binaries setup under the bin directory to perform
		// the tests directly. When we run make test it will be setup and used automatically.
		BinaryAssetsDirectory: filepath.Join("..", "..", "..", "bin", "k8s",
			fmt.Sprintf("1.32.0-%s-%s", runtime.GOOS, runtime.GOARCH)),
	}

	var err error
	// cfg is defined in this file globally.
	cfg, err = testEnv.Start()
	Expect(err).NotTo(HaveOccurred())
	Expect(cfg).NotTo(BeNil())

	Expect(apiextensionsscheme.AddToScheme(scheme.Scheme)).To(Succeed())
	Expect(machinescheme.AddToScheme(scheme.Scheme)).To(Succeed())
	Expect(gardenerextensionv1alpha1.AddToScheme(scheme.Scheme)).To(Succeed())
	Expect(apiv1alpha1.AddToScheme(scheme.Scheme)).To(Succeed())

	// Init package-level k8sClient
	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
	Expect(err).NotTo(HaveOccurred())
	Expect(k8sClient).NotTo(BeNil())

	komega.SetClient(k8sClient)
})

func SetupTest() (*corev1.Namespace, *gardener.ChartApplier) {
	var chartApplier gardener.ChartApplier
	ns := &corev1.Namespace{}
	ign := &corev1.Secret{}

	BeforeEach(func(ctx SpecContext) {
		var err error
		*ns = corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: "testns-",
			},
		}
		Expect(k8sClient.Create(ctx, ns)).To(Succeed(), "failed to create test namespace")
		DeferCleanup(k8sClient.Delete, ns)

		chartApplier, err = gardener.NewChartApplierForConfig(cfg)
		Expect(err).NotTo(HaveOccurred())

		volumeName := "test-volume"
		volumeType := "fast"

		dataYml := map[string]any{
			"a": map[string]any{
				"b": "foo",
			},
		}
		yamlString, err := mapToString(dataYml)
		Expect(err).NotTo(HaveOccurred())

		dataYml2 := map[string]any{
			"a": map[string]any{
				"c": "bar",
			},
		}
		yamlString2, err := mapToString(dataYml2)
		Expect(err).NotTo(HaveOccurred())

		*ign = corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				GenerateName: v1beta1constants.ReferencedResourcesPrefix + "testign-",
				Namespace:    ns.Name,
			},
			Data: map[string][]byte{
				"ignition": []byte(yamlString2),
			},
		}
		Expect(k8sClient.Create(ctx, ign)).To(Succeed(), "failed to create test ignition secret")
		DeferCleanup(k8sClient.Delete, ign)

		workerConfig = &apiv1alpha1.WorkerConfig{
			ExtraServerLabels: map[string]string{
				"foo1": "bar1",
			},

			ExtraIgnition: &apiv1alpha1.IgnitionConfig{
				Raw:       yamlString,
				SecretRef: strings.TrimPrefix(ign.Name, v1beta1constants.ReferencedResourcesPrefix),
				Override:  true,
			},
			Metadata: map[string]string{
				"foo": "bar",
				"baz": "100",
			},
		}
		workerConfigJSON, _ = json.Marshal(workerConfig)

		userDataSecret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: ns.Name,
				Name:      userDataSecretName,
			},
			Data: map[string][]byte{
				userDataSecretDataKey: []byte("some-data"),
			},
		}
		Expect(k8sClient.Create(ctx, userDataSecret)).To(Succeed())
		DeferCleanup(k8sClient.Delete, userDataSecret)

		// define test resources
		pool = gardenerextensionv1alpha1.WorkerPool{
			MachineType:    "large",
			Maximum:        10,
			MaxSurge:       intstr.IntOrString{IntVal: 5},
			MaxUnavailable: intstr.IntOrString{IntVal: 2},
			Annotations:    map[string]string{"foo": "bar"},
			Labels:         map[string]string{"foo": "bar"},
			MachineImage: gardenerextensionv1alpha1.MachineImage{
				Name:    "my-os",
				Version: "1.0",
			},
			Minimum: 0,
			Name:    "pool",
			UserDataSecretRef: corev1.SecretKeySelector{
				LocalObjectReference: corev1.LocalObjectReference{Name: userDataSecretName},
				Key:                  userDataSecretDataKey,
			},
			Volume: &gardenerextensionv1alpha1.Volume{
				Name: &volumeName,
				Type: &volumeType,
				Size: "10Gi",
			},
			Zones:        []string{"zone1", "zone2"},
			Architecture: ptr.To[string]("amd64"),
			NodeTemplate: &gardenerextensionv1alpha1.NodeTemplate{
				Capacity: map[corev1.ResourceName]resource.Quantity{
					corev1.ResourceCPU: resource.MustParse("100m"),
				},
			},
			ProviderConfig: &apiruntime.RawExtension{
				Raw: workerConfigJSON,
			},
		}
		cloudProfileConfig = &apiv1alpha1.CloudProfileConfig{
			TypeMeta: metav1.TypeMeta{
				APIVersion: apiv1alpha1.SchemeGroupVersion.String(),
				Kind:       "CloudProfileConfig",
			},
			MachineTypes: []apiv1alpha1.MachineType{
				{
					Name: "large",
					ServerLabels: map[string]string{
						"foo": "bar",
					},
				},
			},
			MachineImages: []apiv1alpha1.MachineImages{
				{
					Name: "my-os",
					Versions: []apiv1alpha1.MachineImageVersion{
						{
							Version:      "1.0",
							Image:        "registry/my-os",
							Architecture: ptr.To[string]("amd64"),
						},
					},
				},
			},
		}
		shootVersionMajorMinor = "1.2"
		shootVersion = shootVersionMajorMinor + ".3"
		clusterWithoutImages = &extensionscontroller.Cluster{
			Shoot: &gardencorev1beta1.Shoot{
				Spec: gardencorev1beta1.ShootSpec{
					Resources: []gardencorev1beta1.NamedResourceReference{
						{
							Name: strings.TrimPrefix(ign.Name, v1beta1constants.ReferencedResourcesPrefix),
							ResourceRef: autoscalingv1.CrossVersionObjectReference{
								Kind:       "Secret",
								Name:       strings.TrimPrefix(ign.Name, v1beta1constants.ReferencedResourcesPrefix),
								APIVersion: "v1",
							},
						},
					},
					Kubernetes: gardencorev1beta1.Kubernetes{
						Version: shootVersion,
					},
					Provider: gardencorev1beta1.Provider{
						InfrastructureConfig: &apiruntime.RawExtension{
							Raw: []byte("{}"),
						},
					},
				},
				Status: gardencorev1beta1.ShootStatus{
					TechnicalID: technicalID,
				},
			},
		}
		cloudProfileConfigJSON, _ = json.Marshal(cloudProfileConfig)
		testCluster = &extensionscontroller.Cluster{
			CloudProfile: &gardencorev1beta1.CloudProfile{
				ObjectMeta: metav1.ObjectMeta{
					Name: "ironcore-metal",
				},
				Spec: gardencorev1beta1.CloudProfileSpec{
					ProviderConfig: &apiruntime.RawExtension{
						Raw: cloudProfileConfigJSON,
					},
				},
			},
			Shoot: clusterWithoutImages.Shoot,
		}
		w = &gardenerextensionv1alpha1.Worker{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "pool",
				Namespace: ns.Name,
			},
			Spec: gardenerextensionv1alpha1.WorkerSpec{
				DefaultSpec: gardenerextensionv1alpha1.DefaultSpec{},
				Region:      "foo",
				SecretRef: corev1.SecretReference{
					Name: "my-secret",
				},
				SSHPublicKey: nil,
				Pools: []gardenerextensionv1alpha1.WorkerPool{
					pool,
				},
			},
		}
	})

	return ns, &chartApplier
}

func mapToString(m map[string]interface{}) (string, error) {
	yamlData, err := yaml.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(yamlData), nil
}

var _ = AfterSuite(func() {
	Expect(testEnv.Stop()).To(Succeed())
})
