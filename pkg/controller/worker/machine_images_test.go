// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package worker

import (
	"github.com/gardener/gardener/pkg/apis/core/v1beta1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	apiv1alpha1 "github.com/ironcore-dev/gardener-extension-provider-ironcore-metal/pkg/apis/metal/v1alpha1"
)

var _ = Describe("MachinesImages", func() {
	ns, _ := SetupTest()

	It("should update the worker status with GEP-33 capabilities", func(ctx SpecContext) {
		By("defining and setting infrastructure status for worker")
		infraStatus := &apiv1alpha1.InfrastructureStatus{
			TypeMeta: metav1.TypeMeta{
				APIVersion: apiv1alpha1.SchemeGroupVersion.String(),
				Kind:       "InfrastructureStatus",
			},
		}
		w.Spec.InfrastructureProviderStatus = &runtime.RawExtension{Object: infraStatus}

		By("setting up GEP-33 capabilities and machine types in the test cluster")
		testCluster.CloudProfile.Spec.MachineCapabilities = []v1beta1.CapabilityDefinition{
			{Name: "architecture", Values: []string{"amd64", "arm64"}},
		}

		machineTypeName := w.Spec.Pools[0].MachineType
		testCluster.CloudProfile.Spec.MachineTypes = []v1beta1.MachineType{
			{
				Name: machineTypeName,
				Capabilities: v1beta1.Capabilities{
					"architecture": []string{"amd64"},
				},
			},
		}

		By("creating the worker")
		Expect(k8sClient.Create(ctx, w)).To(Succeed())

		By("creating a worker delegate")
		decoder := serializer.NewCodecFactory(k8sClient.Scheme(), serializer.EnableStrict).UniversalDecoder()
		workerDelegate, err := NewWorkerDelegate(k8sClient, decoder, k8sClient.Scheme(), "", w, testCluster)
		Expect(err).NotTo(HaveOccurred())

		By("calling the updating machine image status")
		err = workerDelegate.UpdateMachineImagesStatus(ctx)
		Expect(err).NotTo(HaveOccurred())

		By("ensuring that the worker provider status has been updated with capabilities")
		Eventually(func(g Gomega) {
			workerKey := types.NamespacedName{Namespace: ns.Name, Name: w.Name}
			err := k8sClient.Get(ctx, workerKey, w)
			Expect(client.IgnoreNotFound(err)).To(Succeed())
			g.Expect(err).NotTo(HaveOccurred())

			expectedWorkerStatus := &apiv1alpha1.WorkerStatus{
				TypeMeta: metav1.TypeMeta{
					APIVersion: apiv1alpha1.SchemeGroupVersion.String(),
					Kind:       "WorkerStatus",
				},
				MachineImages: []apiv1alpha1.MachineImage{
					{
						Name:    "my-os",
						Version: "1.0",
						Image:   "registry/my-os",
						Capabilities: v1beta1.Capabilities{
							"architecture": {"amd64"},
						},
					},
				},
			}

			workerStatus := &apiv1alpha1.WorkerStatus{}
			_, _, err = decoder.Decode(w.Status.ProviderStatus.Raw, nil, workerStatus)
			Expect(err).NotTo(HaveOccurred())
			g.Expect(workerStatus).To(Equal(expectedWorkerStatus))
		}).Should(Succeed())
	})

	It("should return an error if no matching image is found for the requested architecture", func(ctx SpecContext) {
		By("setting up GEP-33 capabilities with arm64 requirement")
		testCluster.CloudProfile.Spec.MachineCapabilities = []v1beta1.CapabilityDefinition{
			{Name: "architecture", Values: []string{"amd64", "arm64"}},
		}

		machineTypeName := w.Spec.Pools[0].MachineType
		testCluster.CloudProfile.Spec.MachineTypes = []v1beta1.MachineType{
			{
				Name: machineTypeName,
				Capabilities: v1beta1.Capabilities{
					"architecture": []string{"arm64"},
				},
			},
		}

		decoder := serializer.NewCodecFactory(k8sClient.Scheme(), serializer.EnableStrict).UniversalDecoder()
		workerDelegate, err := NewWorkerDelegate(k8sClient, decoder, k8sClient.Scheme(), "", w, testCluster)
		Expect(err).NotTo(HaveOccurred())

		By("calling the updating machine image status")
		err = workerDelegate.UpdateMachineImagesStatus(ctx)

		By("ensuring it fails because there is no arm64 image in the cloud profile")
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("could not find machine image"))
	})
})
