// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package validation

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	gomegatypes "github.com/onsi/gomega/types"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/utils/ptr"

	apismetal "github.com/ironcore-dev/gardener-extension-provider-ironcore-metal/pkg/apis/metal"
)

var _ = Describe("Infrastructure validation", func() {
	fldPath := field.NewPath("spec", "providerConfig")

	DescribeTable("#ValidateInfrastructureConfig",
		func(infra *apismetal.InfrastructureConfig, workloadIdentity bool, matcher gomegatypes.GomegaMatcher) {
			errs := ValidateInfrastructureConfig(infra, nil, nil, nil, workloadIdentity, fldPath)
			Expect(errs).To(matcher)
		},

		// Legacy mode (workloadIdentity=false) — metalNamespace is optional
		Entry("should pass with nil infra in legacy mode",
			nil, false,
			BeEmpty()),
		Entry("should pass with empty infra in legacy mode",
			&apismetal.InfrastructureConfig{}, false,
			BeEmpty()),
		Entry("should pass with a valid metalNamespace in legacy mode",
			&apismetal.InfrastructureConfig{MetalNamespace: ptr.To("my-namespace")}, false,
			BeEmpty()),
		Entry("should reject an invalid metalNamespace in legacy mode",
			&apismetal.InfrastructureConfig{MetalNamespace: ptr.To("Invalid_NS")}, false,
			ContainElement(HaveField("Field", "spec.providerConfig.metalNamespace"))),

		// WI mode (workloadIdentity=true) — metalNamespace is required
		Entry("should require metalNamespace when workload identity is active and infra is nil",
			nil, true,
			ContainElement(HaveField("Field", "spec.providerConfig.metalNamespace"))),
		Entry("should require metalNamespace when workload identity is active and field is absent",
			&apismetal.InfrastructureConfig{}, true,
			ContainElement(HaveField("Field", "spec.providerConfig.metalNamespace"))),
		Entry("should pass when workload identity is active and metalNamespace is valid",
			&apismetal.InfrastructureConfig{MetalNamespace: ptr.To("my-namespace")}, true,
			BeEmpty()),
		Entry("should reject an invalid metalNamespace when workload identity is active",
			&apismetal.InfrastructureConfig{MetalNamespace: ptr.To("Invalid_NS")}, true,
			ContainElement(HaveField("Field", "spec.providerConfig.metalNamespace"))),
	)
})
