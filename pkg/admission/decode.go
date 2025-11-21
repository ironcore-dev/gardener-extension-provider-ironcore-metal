// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package admission

import (
	"github.com/gardener/gardener/extensions/pkg/util"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/ironcore-dev/gardener-extension-provider-ironcore-metal/pkg/apis/metal"
)

// DecodeControlPlaneConfig decodes the `ControlPlaneConfig` from the given `RawExtension`.
func DecodeControlPlaneConfig(decoder runtime.Decoder, cp *runtime.RawExtension) (*metal.ControlPlaneConfig, error) {
	controlPlaneConfig := &metal.ControlPlaneConfig{}
	if err := util.Decode(decoder, cp.Raw, controlPlaneConfig); err != nil {
		return nil, err
	}

	return controlPlaneConfig, nil
}

// DecodeInfrastructureConfig decodes the `InfrastructureConfig` from the given `RawExtension`.
func DecodeInfrastructureConfig(decoder runtime.Decoder, infra *runtime.RawExtension) (*metal.InfrastructureConfig, error) {
	infraConfig := &metal.InfrastructureConfig{}
	if err := util.Decode(decoder, infra.Raw, infraConfig); err != nil {
		return nil, err
	}

	return infraConfig, nil
}
