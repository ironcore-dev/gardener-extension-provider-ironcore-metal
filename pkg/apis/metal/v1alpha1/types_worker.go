// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// WorkerConfig contains configuration settings for the worker nodes.
type WorkerConfig struct {
	metav1.TypeMeta
	// ExtraIgnition contains additional Ignition for Worker nodes.
	// +optional
	ExtraIgnition *IgnitionConfig `json:"extraIgnition,omitempty"`
	// ExtraServerLabels is a map of additional labels that are applied to the ServerClaim for Server selection.
	// +optional
	ExtraServerLabels map[string]string `json:"extraServerLabels,omitempty"`
	// IPAMConfig is a list of references to Network resources that should be used to assign IP addresses to the worker nodes.
	// +optional
	IPAMConfig []IPAMConfig `json:"ipamConfig,omitempty"`
	// Metadata is a key-value map of additional data which should be passed to the Machine.
	// +optional
	Metadata map[string]string `json:"metadata,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// WorkerStatus contains information about created worker resources.
type WorkerStatus struct {
	metav1.TypeMeta `json:",inline"`

	// MachineImages is a list of machine images that have been used in this worker. Usually, the extension controller
	// gets the mapping from name/version to the provider-specific machine image data in its componentconfig. However, if
	// a version that is still in use gets removed from this componentconfig it cannot reconcile anymore existing `Worker`
	// resources that are still using this version. Hence, it stores the used versions in the provider status to ensure
	// reconciliation is possible.
	// +optional
	MachineImages []MachineImage `json:"machineImages,omitempty"`
}

// MachineImage is a mapping from logical names and versions to metal-specific identifiers.
type MachineImage struct {
	// Name is the logical name of the machine image.
	Name string `json:"name"`
	// Version is the logical version of the machine image.
	Version string `json:"version"`
	// Image is the path to the image.
	Image string `json:"image"`
	// Architecture is the CPU architecture of the machine image.
	// +optional
	Architecture *string `json:"architecture,omitempty"`
}

// IgnitionConfig contains ignition settings.
type IgnitionConfig struct {
	// Raw contains an inline ignition config, which is merged with the config from the os extension.
	// +optional
	Raw string `json:"raw,omitempty"`

	// SecretRef is a reference to a resource in the shoot spec referencing a secret containing the ignition config.
	// +optional
	SecretRef string `json:"secretRef,omitempty"`

	// Override configures, if ignition keys set by the os-extension are overridden
	// by extra ignition.
	// +optional
	Override bool `json:"override,omitempty"`
}

// IPAMObjectReference is a reference to the IPAM object, which will be used for IP allocation.
type IPAMObjectReference struct {
	// Name is the name of resource being referenced.
	Name string `json:"name"`
	// APIGroup is the group for the resource being referenced.
	APIGroup string `json:"apiGroup"`
	// Kind is the type of resource being referenced.
	Kind string `json:"kind"`
}

// IPAMConfig is a reference to an IPAM resource.
type IPAMConfig struct {
	// MetadataKey is the name of metadata key for the network.
	MetadataKey string `json:"metadataKey"`
	// IPAMRef is a reference to the IPAM object, which will be used for IP allocation.
	IPAMRef *IPAMObjectReference `json:"ipamRef"`
}
