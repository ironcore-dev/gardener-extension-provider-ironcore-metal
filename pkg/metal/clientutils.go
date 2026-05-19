// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package metal

import (
	extensionsv1alpha1 "github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
)

var metalScheme = runtime.NewScheme()

func init() {
	utilruntime.Must(corev1.AddToScheme(metalScheme))
	utilruntime.Must(extensionsv1alpha1.AddToScheme(metalScheme))
}
