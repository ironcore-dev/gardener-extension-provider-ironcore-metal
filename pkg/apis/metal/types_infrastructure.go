// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package metal

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// InfrastructureConfig infrastructure configuration resource
type InfrastructureConfig struct {
	metav1.TypeMeta

	// Networks is the metal specific network configuration.
	Networks []Networks
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// InfrastructureStatus contains information about created infrastructure resources.
type InfrastructureStatus struct {
	metav1.TypeMeta
}

// Networks holds information about the Kubernetes and infrastructure networks.
type Networks struct {
	// Name is the name for this network.
	Name string
	// CIDR is the workers subnet range to create.
	CIDR string
	// ID is the ID for the workers' subnet.
	ID string
}
