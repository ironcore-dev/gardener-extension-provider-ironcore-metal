// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// InfrastructureConfig infrastructure configuration resource
type InfrastructureConfig struct {
	metav1.TypeMeta `json:",inline"`

	// Networks is the metal specific network configuration.
	// +optional
	Networks []Networks `json:"networks,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// InfrastructureStatus contains information about created infrastructure resources.
type InfrastructureStatus struct {
	metav1.TypeMeta `json:",inline"`
}

// Networks holds information about the Kubernetes and infrastructure networks.
type Networks struct {
	// Name is the name for this CIDR.
	Name string `json:"name"`
	// CIDR is the workers subnet range to create.
	CIDR string `json:"cidr"`
	// ID is the ID for the workers' subnet.
	// +optional
	ID string `json:"id,omitempty"`
}
