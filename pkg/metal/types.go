// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package metal

import (
	extensionsv1alpha1 "github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	// ProviderName is the name of the metal provider.
	ProviderName = "provider-ironcore-metal"

	// CloudControllerManagerImageName is the name of the cloud-controller-manager image.
	CloudControllerManagerImageName = "cloud-controller-manager"
	// MachineControllerManagerImageName is the name of the MachineControllerManager image.
	MachineControllerManagerImageName = "machine-controller-manager"
	// MachineControllerManagerProviderIroncoreImageName is the name of the MachineController metal image.
	MachineControllerManagerProviderIroncoreImageName = "machine-controller-manager-provider-ironcore-metal"
	// MetallbSpeakerImageName is the name of the metallb speaker to deploy to the shoot.
	MetallbSpeakerImageName = "metallb-speaker"
	// MetallbControllerImageName is the name of the metallb controller to deploy to the shoot.
	MetallbControllerImageName = "metallb-controller"
	// MetalLoadBalancerControllerSpeakerImageName is the name of the metal load balancer controller to deploy to the shoot.
	MetalLoadBalancerControllerSpeakerImageName = "metal-load-balancer-controller-speaker"
	// MetalLoadBalancerControllerManagerImageName is the name of the metal load balancer controller manager to deploy to the seed.
	MetalLoadBalancerControllerManagerImageName = "metal-load-balancer-controller-manager"

	// UsernameFieldName is the field in a secret where the namespace is stored at.
	UsernameFieldName = "username"
	// NamespaceFieldName is the field in a secret where the namespace is stored at.
	NamespaceFieldName = "namespace"
	// KubeConfigFieldName is containing the effective kubeconfig to access an metal cluster.
	KubeConfigFieldName = "kubeconfig"
	// TokenFieldName is containing the token to access an metal cluster.
	TokenFieldName = "token"
	// ClusterFieldName is the name of the cluster field
	ClusterFieldName = "clusterName"
	// LabelsFieldName is the name of the labels field
	LabelsFieldName = "labels"
	// UserDataFieldName is the name of the user data field
	UserDataFieldName = "userData"
	// ImageFieldName is the name of the image field
	ImageFieldName = "image"
	// ServerLabelsFieldName is the name of the server labels field
	ServerLabelsFieldName = "serverLabels"
	// IgnitionFieldName is the name of the ignition field
	IgnitionFieldName = "ignition"
	// IgnitionOverrideFieldName is the name of the ignitionOverride field
	IgnitionOverrideFieldName = "ignitionOverride"
	// MetaDataFieldName is the name of the metadata field
	MetaDataFieldName = "metaData"
	// IPAMConfigFieldName is the name of the ipamConfig field
	IPAMConfigFieldName = "ipamConfig"
	// ClusterNameLabel is the name is the label key of the cluster name
	ClusterNameLabel = "extension.metal.dev/cluster-name"
	// LocalMetalAPIAnnotation is the name of the annotation to mark a seed, which contains a local metal API shoot
	LocalMetalAPIAnnotation = "metal.ironcore.dev/local-metal-api"
	// AllowEgressToIstioIngressLabel is the label key to allow egress to the istio ingress gateway
	AllowEgressToIstioIngressLabel = "networking.resources.gardener.cloud/to-all-istio-ingresses-istio-ingressgateway-tcp-9443"

	// CloudProviderConfigName is the name of the secret containing the cloud provider config.
	CloudProviderConfigName = "cloud-provider-config"
	// CloudControllerManagerName is a constant for the name of the CloudController deployed by the worker controller.
	CloudControllerManagerName = "cloud-controller-manager"
	// CloudControllerManagerFeatureGatesKeyName is the key name for the feature gates key in CCM configuration
	CloudControllerManagerFeatureGatesKeyName = "featureGates"
	// CloudControllerManagerNetworkingKeyName is the key name for the networking key in CCM configuration
	CloudControllerManagerNetworkingKeyName = "networking"
	// CloudControllerManagerNodeAddressesConfigKeyName is the key name for the networking key in CCM configuration
	CloudControllerManagerNodeAddressesConfigKeyName = "configureNodeAddresses"
	// CloudControllerManagerNodeIPAMKindKeyName is the key name for the networking ipamKind key in CCM configuration
	CloudControllerManagerNodeIPAMKindKeyName = "ipamKind"
	// CalicoBgpName is a constant for the name of the Calico BGP deployed by the worker controller.
	CalicoBgpName = "calico-bgp"
	// CalicoIPPoolName is a constant for the name of the Calico IP Pools deployed by the worker controller.
	CalicoIPPoolName = "calico-ipam"
	// MetallbName is a constant for the name of the MetalLB deployed by the worker controller.
	MetallbName = "metallb"
	// MetalLoadBalancerControllerSpeakerName is a constant for the name of the metal load balancer controller.
	MetalLoadBalancerControllerSpeakerName = "metal-load-balancer-controller-speaker"
	// MetalLoadBalancerControllerManagerName is a constant for the name of the metal load balancer controller manager.
	MetalLoadBalancerControllerManagerName = "metal-load-balancer-controller-manager"
	// MachineControllerManagerName is a constant for the name of the machine-controller-manager.
	MachineControllerManagerName = "machine-controller-manager"
	// ShootCalicoNetworkType is the network type for calico in a shoot.
	ShootCalicoNetworkType = "calico"
	// MachineControllerManagerVpaName is the name of the VerticalPodAutoscaler of the machine-controller-manager deployment.
	MachineControllerManagerVpaName = "machine-controller-manager-vpa"
	// MachineControllerManagerMonitoringConfigName is the name of the ConfigMap containing monitoring stack configurations for machine-controller-manager.
	MachineControllerManagerMonitoringConfigName = "machine-controller-manager-monitoring-config"

	// FieldOwner for server side apply
	FieldOwner client.FieldOwner = ProviderName
)

var (
	// UsernamePrefix is a constant for the username prefix of components deployed by metal.
	UsernamePrefix = extensionsv1alpha1.SchemeGroupVersion.Group + ":" + ProviderName + ":"
)
