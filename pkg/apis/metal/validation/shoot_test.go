// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package validation

import (
	"github.com/gardener/gardener/pkg/apis/core"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/utils/ptr"
)

var _ = Describe("ShootConfig validation", func() {

	Describe("#ValidateWorkers", func() {

		var (
			workerConfig []core.Worker
			fldPath      *field.Path
		)

		BeforeEach(func() {
			workerConfig = []core.Worker{}
		})

		It("should return no errors for an empty configuration", func() {
			Expect(ValidateWorkers(workerConfig, fldPath, &core.ShootSpec{})).To(BeEmpty())
		})

		It("should return an error if the extra ignition secretRef is not in the shoot's resources", func() {
			workerConfig = []core.Worker{
				{
					ProviderConfig: &runtime.RawExtension{
						Raw: []byte(`{"extraIgnition": {"secretRef": "some-secret"}}`),
					},
					Zones: []string{"zone"},
					Volume: &core.Volume{
						Name:       ptr.To("volume"),
						Type:       ptr.To("persistentDisk"),
						VolumeSize: "10Gi",
					},
				},
			}
			Expect(ValidateWorkers(workerConfig, fldPath, &core.ShootSpec{})).To(ConsistOf(
				PointTo(MatchFields(IgnoreExtras, Fields{
					"Type":  Equal(field.ErrorTypeInvalid),
					"Field": Equal("[0].providerConfig.extraIgnition.secretRef"),
				})),
			))
		})

	})

})
