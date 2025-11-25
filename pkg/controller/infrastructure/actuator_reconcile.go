// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gardener/gardener/extensions/pkg/controller"
	extensionsv1alpha1 "github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"
	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"

	metalv1alpha1 "github.com/ironcore-dev/gardener-extension-provider-ironcore-metal/pkg/apis/metal/v1alpha1"
)

// Reconcile implements infrastructure actuator reconciliation
func (a *actuator) Reconcile(ctx context.Context, log logr.Logger, infra *extensionsv1alpha1.Infrastructure, cluster *controller.Cluster) error {
	var infrastructureConfig metalv1alpha1.InfrastructureConfig
	err := json.Unmarshal(infra.Spec.ProviderConfig.Raw, &infrastructureConfig)
	if err != nil {
		return fmt.Errorf("failed to unmarshal infrastructure config: %w", err)
	}

	originalInfra := infra.DeepCopy()

	var newNodes []string
	if infrastructureConfig.Networks != nil {
		for _, network := range infrastructureConfig.Networks {
			if network.Name == "" {
				return fmt.Errorf("network name is required")
			}
			newNodes = append(newNodes, network.CIDR)
		}
		if infra.Status.Networking == nil {
			infra.Status.Networking = &extensionsv1alpha1.InfrastructureStatusNetworking{}
		}
		infra.Status.Networking.Nodes = newNodes
	}

	if infra.Status.Networking != nil {
		if cluster.Shoot.Spec.Networking.Pods != nil {
			infra.Status.Networking.Pods = []string{*cluster.Shoot.Spec.Networking.Pods}
		}
		if cluster.Shoot.Spec.Networking.Services != nil {
			infra.Status.Networking.Services = []string{*cluster.Shoot.Spec.Networking.Services}
		}
	}

	if err := a.client.Status().Patch(ctx, infra, client.MergeFrom(originalInfra)); err != nil {
		return fmt.Errorf("failed to patch infrastructure status: %w", err)
	}

	return a.reconcile(ctx, log, infra, cluster)
}

func (a *actuator) reconcile(ctx context.Context, log logr.Logger, infra *extensionsv1alpha1.Infrastructure, cluster *controller.Cluster) error {
	return nil
}
