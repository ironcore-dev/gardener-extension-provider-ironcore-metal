// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package worker

import (
	"context"

	v1beta1constants "github.com/gardener/gardener/pkg/apis/core/v1beta1/constants"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/ironcore-dev/gardener-extension-provider-ironcore-metal/pkg/metal"
)

func (w *workerDelegate) GetMachineControllerManagerChartValues(ctx context.Context) (map[string]any, error) {
	namespace := &corev1.Namespace{}
	if err := w.client.Get(ctx, client.ObjectKey{Name: w.worker.Namespace}, namespace); err != nil {
		return nil, err
	}

	podLabels := map[string]any{
		v1beta1constants.LabelPodMaintenanceRestart: "true",
	}
	localAPI, ok := w.cluster.Seed.Annotations[metal.LocalMetalAPIAnnotation]
	if ok && localAPI == "true" {
		podLabels[metal.AllowEgressToIstioIngressLabel] = "allowed"
	}

	return map[string]any{
		"providerName": metal.ProviderName,
		"namespace": map[string]any{
			"uid": namespace.UID,
		},
		"podLabels": podLabels,
	}, nil
}

func (w *workerDelegate) GetMachineControllerManagerShootChartValues(ctx context.Context) (map[string]any, error) {
	return map[string]any{
		"providerName": metal.ProviderName,
	}, nil
}
