// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package infrastructure

import (
	"encoding/json"

	extensionscontroller "github.com/gardener/gardener/extensions/pkg/controller"
	gardencorev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	extensionsv1alpha1 "github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"
	"github.com/go-logr/logr"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/ptr"

	metalv1alpha1 "github.com/ironcore-dev/gardener-extension-provider-ironcore-metal/pkg/apis/metal/v1alpha1"
)

var (
	shootVersionMajorMinor = "1.2"
	shootVersion           = shootVersionMajorMinor + ".3"
)

var _ = Describe("Actuator Reconcile", func() {
	var (
		log     logr.Logger
		infra   *extensionsv1alpha1.Infrastructure
		cluster *extensionscontroller.Cluster
		act     *actuator
	)

	BeforeEach(func(ctx SpecContext) {
		log = logr.Discard()

		infrastructureConfig := metalv1alpha1.InfrastructureConfig{
			Networks: []metalv1alpha1.Networks{
				{Name: "worker-network-1", CIDR: "10.10.10.0/24", ID: "1"},
				{Name: "worker-network-2", CIDR: "10.10.20.0/24", ID: "2"},
			},
		}
		infrastructureConfigRaw, err := json.Marshal(infrastructureConfig)
		Expect(err).To(Succeed())

		infra = &extensionsv1alpha1.Infrastructure{}
		infra.Name = "some-infra"
		infra.Namespace = metav1.NamespaceDefault
		infra.Spec.ProviderConfig = &runtime.RawExtension{
			Raw: infrastructureConfigRaw,
		}

		Expect(k8sClient.Create(ctx, infra)).To(Succeed())
		DeferCleanup(k8sClient.Delete, infra)

		cluster = &extensionscontroller.Cluster{
			Shoot: &gardencorev1beta1.Shoot{
				Spec: gardencorev1beta1.ShootSpec{
					Kubernetes: gardencorev1beta1.Kubernetes{
						Version: shootVersion,
					},
					Networking: &gardencorev1beta1.Networking{
						Pods:     ptr.To("100.12.12.0/8"),
						Services: ptr.To("100.12.13/8"),
					},
				},
			},
		}

		act = &actuator{client: k8sClient}
	})

	Describe("#Reconcile", func() {
		It("should update infra.Status.Networking.Nodes", func(ctx SpecContext) {
			err := act.Reconcile(ctx, log, infra, cluster)
			Expect(err).NotTo(HaveOccurred())

			// Verify that infra.Status.Networking.Nodes is updated
			expectedNodes := []string{"10.10.10.0/24", "10.10.20.0/24"}
			Expect(infra.Status.Networking.Nodes).To(Equal(expectedNodes))
		})

		It("should copy the Pod and Service CIDRs from the Shoot spec to infra.Status.Networking", func(ctx SpecContext) {
			err := act.Reconcile(ctx, log, infra, cluster)
			Expect(err).NotTo(HaveOccurred())

			// Verify that infra.Status.Networking.Pods and Services are updated
			Expect(infra.Status.Networking.Pods).To(Equal([]string{"100.12.12.0/8"}))
			Expect(infra.Status.Networking.Services).To(Equal([]string{"100.12.13/8"}))
		})
	})
})
