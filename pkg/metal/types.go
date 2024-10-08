// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package metal

import (
	extensionsv1alpha1 "github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"
)

const (
	// ProviderName is the name of the metal provider.
	ProviderName = "provider-metal"

	// CloudControllerManagerImageName is the name of the cloud-controller-manager image.
	CloudControllerManagerImageName = "cloud-controller-manager"
	// MachineControllerManagerImageName is the name of the MachineControllerManager image.
	MachineControllerManagerImageName = "machine-controller-manager"
	// MachineControllerManagerProviderIroncoreImageName is the name of the MachineController metal image.
	MachineControllerManagerProviderIroncoreImageName = "machine-controller-manager-provider-metal"
	// MetallbSpeakerImageName is the name of the metallb speaker to deploy to the shoot.
	MetallbSpeakerImageName = "metallb-speaker"
	// MetallbControllerImageName is the name of the metallb controller to deploy to the shoot.
	MetallbControllerImageName = "metallb-controller"

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
	// ServerLabels is the name of the server labels field
	ServerLabels = "serverLabels"
	// ClusterNameLabel is the name is the label key of the cluster name
	ClusterNameLabel = "extension.metal.dev/cluster-name"

	// CloudProviderConfigName is the name of the secret containing the cloud provider config.
	CloudProviderConfigName = "cloud-provider-config"
	// CloudControllerManagerName is a constant for the name of the CloudController deployed by the worker controller.
	CloudControllerManagerName = "cloud-controller-manager"
	// CalicoBgpName is a constant for the name of the Calico BGP deployed by the worker controller.
	CalicoBgpName = "calico-bgp"
	// MetallbName is a constant for the name of the MetalLB deployed by the worker controller.
	MetallbName = "metallb"
	// MachineControllerManagerName is a constant for the name of the machine-controller-manager.
	MachineControllerManagerName = "machine-controller-manager"
	// ShootCalicoNetworkType is the network type for calico in a shoot.
	ShootCalicoNetworkType = "calico"
	// MachineControllerManagerVpaName is the name of the VerticalPodAutoscaler of the machine-controller-manager deployment.
	MachineControllerManagerVpaName = "machine-controller-manager-vpa"
	// MachineControllerManagerMonitoringConfigName is the name of the ConfigMap containing monitoring stack configurations for machine-controller-manager.
	MachineControllerManagerMonitoringConfigName = "machine-controller-manager-monitoring-config"
)

var (
	// UsernamePrefix is a constant for the username prefix of components deployed by metal.
	UsernamePrefix = extensionsv1alpha1.SchemeGroupVersion.Group + ":" + ProviderName + ":"
)
