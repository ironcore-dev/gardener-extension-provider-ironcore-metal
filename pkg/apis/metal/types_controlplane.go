// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package metal

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ControlPlaneConfig contains configuration settings for the control plane.
type ControlPlaneConfig struct {
	metav1.TypeMeta

	// CloudControllerManager contains configuration settings for the cloud-controller-manager.
	CloudControllerManager *CloudControllerManagerConfig

	// LoadBalancerConfig contains configuration settings for the shoot loadbalancing.
	LoadBalancerConfig *LoadBalancerConfig

	// HostnamePolicy is a policy for generating hostnames for the worker nodes.
	HostnamePolicy NodeNamePolicy
}

// NodeNamePolicy is a policy for generating hostnames for the worker nodes.
type NodeNamePolicy string

const (
	// NodeNamePolicyServerName is a policy for generating hostnames based on the Server.
	NodeNamePolicyServerName NodeNamePolicy = "ServerName"
)

// CloudControllerNetworking contains configuration settings for CCM networking.
type CloudControllerNetworking struct {
	// ConfigureNodeAddresses enables the configuration of node addresses.
	ConfigureNodeAddresses bool
	// IPAMKind enables the IPAM integration.
	IPAMKind *IPAMKind
}

// IPAMKind specifiers the IPAM objects in-use.
type IPAMKind struct {
	// APIGroup is the resource group.
	APIGroup string
	// Kind is the resource type.
	Kind string
}

// CloudControllerManagerConfig contains configuration settings for the cloud-controller-manager.
type CloudControllerManagerConfig struct {
	// FeatureGates contains information about enabled feature gates.
	FeatureGates map[string]bool

	// Networking contains configuration settings for CCM networking.
	Networking *CloudControllerNetworking
}

// LoadBalancerConfig contains configuration settings for the shoot loadbalancing.
type LoadBalancerConfig struct {
	// MetallbConfig contains configuration settings for metallb.
	MetallbConfig *MetallbConfig

	// CalicoBgpConfig contains configuration settings for calico.
	CalicoBgpConfig *CalicoBgpConfig
}

// MetallbConfig contains configuration settings for metallb.
type MetallbConfig struct {
	// IPAddressPool contains IP address pools for metallb.
	IPAddressPool []string

	// EnableSpeaker enables the metallb speaker.
	EnableSpeaker bool

	// EnableL2Advertisement enables L2 advertisement.
	EnableL2Advertisement bool
}

// CalicoBgpConfig contains BGP configuration settings for calico.
type CalicoBgpConfig struct {
	// ASNumber is the default AS number used by a node.
	ASNumber int

	// nodeToNodeMeshEnabled enables the node-to-node mesh.
	NodeToNodeMeshEnabled bool

	// ServiceLoadBalancerIPs are the CIDR blocks for Kubernetes Service LoadBalancer IPs.
	ServiceLoadBalancerIPs []string

	// ServiceExternalIPs are the CIDR blocks for Kubernetes Service External IPs.
	ServiceExternalIPs []string

	// ServiceClusterIPs are the CIDR blocks from which service cluster IPs are allocated.
	ServiceClusterIPs []string

	// BGPPeer contains configuration for BGPPeer resource.
	BgpPeer []BgpPeer

	// BGPFilter contains configuration for BGPFilter resource.
	BGPFilter []BGPFilter
}

// BgpPeer contains configuration for BGPPeer resource.
type BgpPeer struct {
	// PeerIP contains IP address of BGP peer followed by an optional port number to peer with.
	PeerIP string

	// ASNumber contains the AS number of the BGP peer.
	ASNumber int

	// NodeSelector is a key-value pair to select nodes that should have this peering.
	NodeSelector string

	// Filters contain the filters for the BGP peer.
	Filters []string
}

// BGPFilter contains configuration for BGPFilter resource.
type BGPFilter struct {
	// Name is the name of the BGPFilter resource.
	Name string

	// The ordered set of IPv4 BGPFilter rules acting on exporting routes to a peer.
	ExportV4 []BGPFilterRule

	// The ordered set of IPv4 BGPFilter rules acting on importing routes from a peer.
	ImportV4 []BGPFilterRule

	// The ordered set of IPv6 BGPFilter rules acting on exporting routes to a peer.
	ExportV6 []BGPFilterRule

	// The ordered set of IPv6 BGPFilter rules acting on importing routes from a peer.
	ImportV6 []BGPFilterRule
}

// BGPFilterRule defines a BGP filter rule consisting a single CIDR block and a filter action for this CIDR.
type BGPFilterRule struct {
	CIDR string

	// +kubebuilder:validation:Enum=Equal;NotEqual;In;NotIn
	MatchOperator string

	// +kubebuilder:validation:Enum=Accept;Reject
	Action string
}
