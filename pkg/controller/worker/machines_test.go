// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package worker

import (
	"encoding/json"
	"fmt"

	"github.com/gardener/gardener/extensions/pkg/controller/worker"
	genericworkeractuator "github.com/gardener/gardener/extensions/pkg/controller/worker/genericactuator"
	gardencorev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	v1beta1constants "github.com/gardener/gardener/pkg/apis/core/v1beta1/constants"
	gardenerextensionv1alpha1 "github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"
	machinecontrollerv1alpha1 "github.com/gardener/machine-controller-manager/pkg/apis/machine/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/utils/ptr"
	. "sigs.k8s.io/controller-runtime/pkg/envtest/komega"

	"github.com/ironcore-dev/gardener-extension-provider-ironcore-metal/pkg/metal"
)

var _ = Describe("Machines", func() {
	ns, _ := SetupTest()

	When("deploying machine classes", func() {

		var (
			deploymentName     string
			className          string
			machineClass       *machinecontrollerv1alpha1.MachineClass
			machineClassSecret *corev1.Secret
			workerDelegate     genericworkeractuator.WorkerDelegate
		)

		dataYml := map[string]any{
			"a": map[string]any{
				"b": "foo",
				"c": "bar",
			},
		}
		yamlString, err := mapToString(dataYml)
		Expect(err).NotTo(HaveOccurred())

		BeforeEach(func(ctx SpecContext) {
			testCluster.CloudProfile.Spec.MachineCapabilities = []gardencorev1beta1.CapabilityDefinition{
				{Name: "architecture", Values: []string{"amd64", "arm64"}},
			}
			testCluster.CloudProfile.Spec.MachineTypes = []gardencorev1beta1.MachineType{
				{
					Name: pool.MachineType,
					Capabilities: gardencorev1beta1.Capabilities{
						"architecture": []string{"amd64"},
					},
				},
			}

			providerConfig := string(pool.ProviderConfig.Raw)
			workerPoolHash, err := worker.WorkerPoolHash(pool, testCluster, []string{providerConfig}, []string{providerConfig})
			Expect(err).NotTo(HaveOccurred())
			deploymentName = fmt.Sprintf("%s-%s-z%d", technicalID, pool.Name, 1)
			className = fmt.Sprintf("%s-%s", deploymentName, workerPoolHash)
			machineClass = &machinecontrollerv1alpha1.MachineClass{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: ns.Name,
					Name:      className,
				},
			}
			machineClassSecret = &corev1.Secret{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: ns.Name,
					Name:      className,
				},
			}
			By("deploying the machine class for a given multi zone cluster")
			decoder := serializer.NewCodecFactory(k8sClient.Scheme(), serializer.EnableStrict).UniversalDecoder()
			workerDelegate, err = NewWorkerDelegate(k8sClient, decoder, k8sClient.Scheme(), "", w, testCluster)
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func(ctx SpecContext) {
			Expect(k8sClient.Delete(ctx, machineClass)).To(Succeed())
			Expect(k8sClient.Delete(ctx, machineClassSecret)).To(Succeed())
		})

		It("should create the expected machine class for a multi zone cluster", func(ctx SpecContext) {
			Expect(workerDelegate.DeployMachineClasses(ctx)).To(Succeed())
			By("ensuring that the machine class for each pool has been deployed")
			machineClassProviderSpec := map[string]any{
				"image": "registry/my-os",
				"labels": map[string]any{
					metal.ClusterNameLabel: technicalID,
				},
				metal.ServerLabelsFieldName: map[string]string{
					"foo":  "bar",
					"foo1": "bar1",
				},
				metal.IgnitionFieldName:         yamlString,
				metal.IgnitionOverrideFieldName: true,
				metal.MetaDataFieldName: map[string]string{
					"foo": "bar",
					"baz": "100",
				},
			}

			Eventually(Object(machineClass)).Should(SatisfyAll(
				HaveField("CredentialsSecretRef", &corev1.SecretReference{
					Namespace: w.Spec.SecretRef.Namespace,
					Name:      w.Spec.SecretRef.Name,
				}),
				HaveField("SecretRef", &corev1.SecretReference{
					Namespace: ns.Name,
					Name:      className,
				}),
				HaveField("Provider", "ironcore-metal"),
				HaveField("NodeTemplate", &machinecontrollerv1alpha1.NodeTemplate{
					Architecture: pool.Architecture,

					Capacity:        pool.NodeTemplate.Capacity,
					VirtualCapacity: pool.NodeTemplate.VirtualCapacity,
					InstanceType:    pool.MachineType,
					Region:          w.Spec.Region,
					Zone:            "zone1",
				}),
				HaveField("ProviderSpec", runtime.RawExtension{
					Raw: encodeMap(machineClassProviderSpec),
				}),
			))

			By("ensuring that the machine class secret have been applied")

			Eventually(Object(machineClassSecret)).Should(SatisfyAll(
				HaveField("ObjectMeta.Labels", HaveKeyWithValue(v1beta1constants.GardenerPurpose, v1beta1constants.GardenPurposeMachineClass)),
				HaveField("Data", HaveKeyWithValue("userData", []byte("some-data"))),
			))
		})
	})

	It("should generate the machine deployments", func(ctx SpecContext) {
		By("creating a worker delegate")
		providerConfig := string(pool.ProviderConfig.Raw)
		workerPoolHash, err := worker.WorkerPoolHash(pool, testCluster, []string{providerConfig}, []string{providerConfig})
		Expect(err).NotTo(HaveOccurred())
		var (
			deploymentName1 = fmt.Sprintf("%s-%s-z%d", technicalID, pool.Name, 1)
			deploymentName2 = fmt.Sprintf("%s-%s-z%d", technicalID, pool.Name, 2)
			className1      = fmt.Sprintf("%s-%s", deploymentName1, workerPoolHash)
			className2      = fmt.Sprintf("%s-%s", deploymentName2, workerPoolHash)
		)
		decoder := serializer.NewCodecFactory(k8sClient.Scheme(), serializer.EnableStrict).UniversalDecoder()
		workerDelegate, err := NewWorkerDelegate(k8sClient, decoder, k8sClient.Scheme(), "", w, testCluster)
		Expect(err).NotTo(HaveOccurred())

		By("generating the machine deployments")
		machineDeployments, err := workerDelegate.GenerateMachineDeployments(ctx)
		Expect(err).NotTo(HaveOccurred())

		Expect(machineDeployments).To(Equal(worker.MachineDeployments{
			worker.MachineDeployment{
				Name:       deploymentName1,
				PoolName:   pool.Name,
				ClassName:  className1,
				SecretName: className1,
				Minimum:    worker.DistributeOverZones(0, pool.Minimum, 2),
				Maximum:    worker.DistributeOverZones(0, pool.Maximum, 2),
				Strategy: machinecontrollerv1alpha1.MachineDeploymentStrategy{
					Type: machinecontrollerv1alpha1.RollingUpdateMachineDeploymentStrategyType,
					RollingUpdate: &machinecontrollerv1alpha1.RollingUpdateMachineDeployment{
						UpdateConfiguration: machinecontrollerv1alpha1.UpdateConfiguration{
							MaxSurge:       ptr.To(worker.DistributePositiveIntOrPercent(0, pool.MaxSurge, 2, pool.Maximum)),
							MaxUnavailable: ptr.To(worker.DistributePositiveIntOrPercent(0, pool.MaxUnavailable, 2, pool.Minimum)),
						},
					},
				},
				Labels:               pool.Labels,
				Annotations:          pool.Annotations,
				Taints:               pool.Taints,
				MachineConfiguration: genericworkeractuator.ReadMachineConfiguration(pool),
				Priority:             ptr.To(int32(1)),
			},
			worker.MachineDeployment{
				Name:       deploymentName2,
				PoolName:   pool.Name,
				ClassName:  className2,
				SecretName: className2,
				Minimum:    worker.DistributeOverZones(1, pool.Minimum, 2),
				Maximum:    worker.DistributeOverZones(1, pool.Maximum, 2),
				Strategy: machinecontrollerv1alpha1.MachineDeploymentStrategy{
					Type: machinecontrollerv1alpha1.RollingUpdateMachineDeploymentStrategyType,
					RollingUpdate: &machinecontrollerv1alpha1.RollingUpdateMachineDeployment{
						UpdateConfiguration: machinecontrollerv1alpha1.UpdateConfiguration{
							MaxSurge:       ptr.To(worker.DistributePositiveIntOrPercent(1, pool.MaxSurge, 2, pool.Maximum)),
							MaxUnavailable: ptr.To(worker.DistributePositiveIntOrPercent(1, pool.MaxUnavailable, 2, pool.Minimum)),
						},
					},
				},
				Labels:               pool.Labels,
				Annotations:          pool.Annotations,
				Taints:               pool.Taints,
				MachineConfiguration: genericworkeractuator.ReadMachineConfiguration(pool),
				Priority:             ptr.To(int32(1)),
			},
		}))
	})

	It("should set PoolName on each MachineDeployment to the worker pool name", func(ctx SpecContext) {
		decoder := serializer.NewCodecFactory(k8sClient.Scheme(), serializer.EnableStrict).UniversalDecoder()
		workerDelegate, err := NewWorkerDelegate(k8sClient, decoder, k8sClient.Scheme(), "", w, testCluster)
		Expect(err).NotTo(HaveOccurred())

		machineDeployments, err := workerDelegate.GenerateMachineDeployments(ctx)
		Expect(err).NotTo(HaveOccurred())

		for _, md := range machineDeployments {
			Expect(md.PoolName).ToNot(BeEmpty(), "PoolName must not be empty for deployment %s", md.Name)
			Expect(md.PoolName).To(Equal(pool.Name), "PoolName must match the worker pool name for deployment %s", md.Name)
		}
	})

	Context("in-place update strategy", func() {
		var (
			inPlacePoolName       = "pool-inplace"
			manualInPlacePoolName = "pool-manual-inplace"
		)

		newInPlacePool := func(name string, strategy gardencorev1beta1.MachineUpdateStrategy) gardenerextensionv1alpha1.WorkerPool {
			return gardenerextensionv1alpha1.WorkerPool{
				Name:                name,
				MachineType:         pool.MachineType,
				Minimum:             pool.Minimum,
				Maximum:             pool.Maximum,
				MaxSurge:            intstr.FromInt32(3),
				MaxUnavailable:      intstr.FromInt32(1),
				Architecture:        pool.Architecture,
				MachineImage:        pool.MachineImage,
				Volume:              pool.Volume,
				Zones:               pool.Zones,
				Labels:              pool.Labels,
				Annotations:         pool.Annotations,
				Taints:              pool.Taints,
				NodeTemplate:        pool.NodeTemplate,
				ProviderConfig:      pool.ProviderConfig,
				UpdateStrategy:      ptr.To(strategy),
				KubernetesVersion:   ptr.To(shootVersion),
				NodeAgentSecretName: pool.NodeAgentSecretName,
				UserDataSecretRef:   pool.UserDataSecretRef,
			}
		}

		It("should generate in-place machine deployments with AutoInPlaceUpdate strategy", func(ctx SpecContext) {
			autoPool := newInPlacePool(inPlacePoolName, gardencorev1beta1.AutoInPlaceUpdate)
			w.Spec.Pools = []gardenerextensionv1alpha1.WorkerPool{autoPool}

			workerPoolHash, err := worker.WorkerPoolHash(autoPool, testCluster, nil, nil)
			Expect(err).NotTo(HaveOccurred())

			deploymentName1 := fmt.Sprintf("%s-%s-z1", technicalID, inPlacePoolName)
			deploymentName2 := fmt.Sprintf("%s-%s-z2", technicalID, inPlacePoolName)
			className1 := fmt.Sprintf("%s-%s", deploymentName1, workerPoolHash)
			className2 := fmt.Sprintf("%s-%s", deploymentName2, workerPoolHash)

			decoder := serializer.NewCodecFactory(k8sClient.Scheme(), serializer.EnableStrict).UniversalDecoder()
			workerDelegate, err := NewWorkerDelegate(k8sClient, decoder, k8sClient.Scheme(), "", w, testCluster)
			Expect(err).NotTo(HaveOccurred())

			machineDeployments, err := workerDelegate.GenerateMachineDeployments(ctx)
			Expect(err).NotTo(HaveOccurred())

			Expect(machineDeployments).To(Equal(worker.MachineDeployments{
				worker.MachineDeployment{
					Name:       deploymentName1,
					PoolName:   inPlacePoolName,
					ClassName:  className1,
					SecretName: className1,
					Minimum:    worker.DistributeOverZones(0, autoPool.Minimum, 2),
					Maximum:    worker.DistributeOverZones(0, autoPool.Maximum, 2),
					Strategy: machinecontrollerv1alpha1.MachineDeploymentStrategy{
						Type: machinecontrollerv1alpha1.InPlaceUpdateMachineDeploymentStrategyType,
						InPlaceUpdate: &machinecontrollerv1alpha1.InPlaceUpdateMachineDeployment{
							OrchestrationType: machinecontrollerv1alpha1.OrchestrationTypeAuto,
							UpdateConfiguration: machinecontrollerv1alpha1.UpdateConfiguration{
								MaxSurge:       ptr.To(worker.DistributePositiveIntOrPercent(0, autoPool.MaxSurge, 2, autoPool.Maximum)),
								MaxUnavailable: ptr.To(worker.DistributePositiveIntOrPercent(0, autoPool.MaxUnavailable, 2, autoPool.Minimum)),
							},
						},
					},
					Labels:               autoPool.Labels,
					Annotations:          autoPool.Annotations,
					Taints:               autoPool.Taints,
					MachineConfiguration: genericworkeractuator.ReadMachineConfiguration(autoPool),
					Priority:             ptr.To(int32(1)),
				},
				worker.MachineDeployment{
					Name:       deploymentName2,
					PoolName:   inPlacePoolName,
					ClassName:  className2,
					SecretName: className2,
					Minimum:    worker.DistributeOverZones(1, autoPool.Minimum, 2),
					Maximum:    worker.DistributeOverZones(1, autoPool.Maximum, 2),
					Strategy: machinecontrollerv1alpha1.MachineDeploymentStrategy{
						Type: machinecontrollerv1alpha1.InPlaceUpdateMachineDeploymentStrategyType,
						InPlaceUpdate: &machinecontrollerv1alpha1.InPlaceUpdateMachineDeployment{
							OrchestrationType: machinecontrollerv1alpha1.OrchestrationTypeAuto,
							UpdateConfiguration: machinecontrollerv1alpha1.UpdateConfiguration{
								MaxSurge:       ptr.To(worker.DistributePositiveIntOrPercent(1, autoPool.MaxSurge, 2, autoPool.Maximum)),
								MaxUnavailable: ptr.To(worker.DistributePositiveIntOrPercent(1, autoPool.MaxUnavailable, 2, autoPool.Minimum)),
							},
						},
					},
					Labels:               autoPool.Labels,
					Annotations:          autoPool.Annotations,
					Taints:               autoPool.Taints,
					MachineConfiguration: genericworkeractuator.ReadMachineConfiguration(autoPool),
					Priority:             ptr.To(int32(1)),
				},
			}))
		})

		It("should generate in-place machine deployments with ManualInPlaceUpdate strategy", func(ctx SpecContext) {
			manualPool := newInPlacePool(manualInPlacePoolName, gardencorev1beta1.ManualInPlaceUpdate)
			w.Spec.Pools = []gardenerextensionv1alpha1.WorkerPool{manualPool}

			workerPoolHash, err := worker.WorkerPoolHash(manualPool, testCluster, nil, nil)
			Expect(err).NotTo(HaveOccurred())

			deploymentName1 := fmt.Sprintf("%s-%s-z1", technicalID, manualInPlacePoolName)
			deploymentName2 := fmt.Sprintf("%s-%s-z2", technicalID, manualInPlacePoolName)
			className1 := fmt.Sprintf("%s-%s", deploymentName1, workerPoolHash)
			className2 := fmt.Sprintf("%s-%s", deploymentName2, workerPoolHash)

			decoder := serializer.NewCodecFactory(k8sClient.Scheme(), serializer.EnableStrict).UniversalDecoder()
			workerDelegate, err := NewWorkerDelegate(k8sClient, decoder, k8sClient.Scheme(), "", w, testCluster)
			Expect(err).NotTo(HaveOccurred())

			machineDeployments, err := workerDelegate.GenerateMachineDeployments(ctx)
			Expect(err).NotTo(HaveOccurred())

			Expect(machineDeployments).To(Equal(worker.MachineDeployments{
				worker.MachineDeployment{
					Name:       deploymentName1,
					PoolName:   manualInPlacePoolName,
					ClassName:  className1,
					SecretName: className1,
					Minimum:    worker.DistributeOverZones(0, manualPool.Minimum, 2),
					Maximum:    worker.DistributeOverZones(0, manualPool.Maximum, 2),
					Strategy: machinecontrollerv1alpha1.MachineDeploymentStrategy{
						Type: machinecontrollerv1alpha1.InPlaceUpdateMachineDeploymentStrategyType,
						InPlaceUpdate: &machinecontrollerv1alpha1.InPlaceUpdateMachineDeployment{
							OrchestrationType: machinecontrollerv1alpha1.OrchestrationTypeManual,
							UpdateConfiguration: machinecontrollerv1alpha1.UpdateConfiguration{
								MaxSurge:       ptr.To(worker.DistributePositiveIntOrPercent(0, manualPool.MaxSurge, 2, manualPool.Maximum)),
								MaxUnavailable: ptr.To(worker.DistributePositiveIntOrPercent(0, manualPool.MaxUnavailable, 2, manualPool.Minimum)),
							},
						},
					},
					Labels:               manualPool.Labels,
					Annotations:          manualPool.Annotations,
					Taints:               manualPool.Taints,
					MachineConfiguration: genericworkeractuator.ReadMachineConfiguration(manualPool),
					Priority:             ptr.To(int32(1)),
				},
				worker.MachineDeployment{
					Name:       deploymentName2,
					PoolName:   manualInPlacePoolName,
					ClassName:  className2,
					SecretName: className2,
					Minimum:    worker.DistributeOverZones(1, manualPool.Minimum, 2),
					Maximum:    worker.DistributeOverZones(1, manualPool.Maximum, 2),
					Strategy: machinecontrollerv1alpha1.MachineDeploymentStrategy{
						Type: machinecontrollerv1alpha1.InPlaceUpdateMachineDeploymentStrategyType,
						InPlaceUpdate: &machinecontrollerv1alpha1.InPlaceUpdateMachineDeployment{
							OrchestrationType: machinecontrollerv1alpha1.OrchestrationTypeManual,
							UpdateConfiguration: machinecontrollerv1alpha1.UpdateConfiguration{
								MaxSurge:       ptr.To(worker.DistributePositiveIntOrPercent(1, manualPool.MaxSurge, 2, manualPool.Maximum)),
								MaxUnavailable: ptr.To(worker.DistributePositiveIntOrPercent(1, manualPool.MaxUnavailable, 2, manualPool.Minimum)),
							},
						},
					},
					Labels:               manualPool.Labels,
					Annotations:          manualPool.Annotations,
					Taints:               manualPool.Taints,
					MachineConfiguration: genericworkeractuator.ReadMachineConfiguration(manualPool),
					Priority:             ptr.To(int32(1)),
				},
			}))
		})

		It("should generate mixed rolling and in-place deployments when pools have different strategies", func(ctx SpecContext) {
			autoPool := newInPlacePool(inPlacePoolName, gardencorev1beta1.AutoInPlaceUpdate)
			autoPool.Zones = []string{"zone1"}
			rollingPool := pool
			rollingPool.Zones = []string{"zone1"}

			w.Spec.Pools = []gardenerextensionv1alpha1.WorkerPool{rollingPool, autoPool}

			decoder := serializer.NewCodecFactory(k8sClient.Scheme(), serializer.EnableStrict).UniversalDecoder()
			workerDelegate, err := NewWorkerDelegate(k8sClient, decoder, k8sClient.Scheme(), "", w, testCluster)
			Expect(err).NotTo(HaveOccurred())

			machineDeployments, err := workerDelegate.GenerateMachineDeployments(ctx)
			Expect(err).NotTo(HaveOccurred())

			Expect(machineDeployments).To(HaveLen(2))
			Expect(machineDeployments[0].Strategy.Type).To(Equal(machinecontrollerv1alpha1.RollingUpdateMachineDeploymentStrategyType))
			Expect(machineDeployments[0].Strategy.RollingUpdate).NotTo(BeNil())
			Expect(machineDeployments[0].Strategy.InPlaceUpdate).To(BeNil())
			Expect(machineDeployments[1].Strategy.Type).To(Equal(machinecontrollerv1alpha1.InPlaceUpdateMachineDeploymentStrategyType))
			Expect(machineDeployments[1].Strategy.InPlaceUpdate).NotTo(BeNil())
			Expect(machineDeployments[1].Strategy.InPlaceUpdate.OrchestrationType).To(Equal(machinecontrollerv1alpha1.OrchestrationTypeAuto))
			Expect(machineDeployments[1].Strategy.RollingUpdate).To(BeNil())
		})

		It("should produce the same worker pool hash regardless of providerConfig for in-place pools", func(ctx SpecContext) {
			pool1 := newInPlacePool(inPlacePoolName, gardencorev1beta1.AutoInPlaceUpdate)
			pool1.ProviderConfig = &runtime.RawExtension{Raw: []byte(`{"config":"value-a"}`)}

			pool2 := newInPlacePool(inPlacePoolName, gardencorev1beta1.AutoInPlaceUpdate)
			pool2.ProviderConfig = &runtime.RawExtension{Raw: []byte(`{"config":"value-b"}`)}

			hash1, err := worker.WorkerPoolHash(pool1, testCluster, nil, nil)
			Expect(err).NotTo(HaveOccurred())
			hash2, err := worker.WorkerPoolHash(pool2, testCluster, nil, nil)
			Expect(err).NotTo(HaveOccurred())

			Expect(hash1).To(Equal(hash2))
		})
	})
})

func encodeMap(m map[string]any) []byte {
	data, err := json.Marshal(m)
	Expect(err).To(Succeed())
	return data
}
