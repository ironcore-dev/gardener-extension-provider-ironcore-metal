// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package controlplane

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"path/filepath"
	"strings"

	extensionscontroller "github.com/gardener/gardener/extensions/pkg/controller"
	"github.com/gardener/gardener/extensions/pkg/controller/controlplane/genericactuator"
	extensionssecretsmanager "github.com/gardener/gardener/extensions/pkg/util/secret/manager"
	gardencorev1beta1 "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	v1beta1constants "github.com/gardener/gardener/pkg/apis/core/v1beta1/constants"
	extensionsv1alpha1 "github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"
	"github.com/gardener/gardener/pkg/utils/chart"
	gutil "github.com/gardener/gardener/pkg/utils/gardener"
	kutil "github.com/gardener/gardener/pkg/utils/kubernetes"
	secretutils "github.com/gardener/gardener/pkg/utils/secrets"
	secretsmanager "github.com/gardener/gardener/pkg/utils/secrets/manager"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	autoscalingv1 "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/apis/autoscaling.k8s.io/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/ironcore-dev/gardener-extension-provider-ironcore-metal/charts"
	metalapi "github.com/ironcore-dev/gardener-extension-provider-ironcore-metal/pkg/apis/metal"
	"github.com/ironcore-dev/gardener-extension-provider-ironcore-metal/pkg/internal"
	"github.com/ironcore-dev/gardener-extension-provider-ironcore-metal/pkg/metal"
)

const (
	caNameControlPlane                        = "ca-" + metal.ProviderName + "-controlplane"
	cloudControllerManagerDeploymentName      = "cloud-controller-manager"
	cloudControllerManagerServerName          = "cloud-controller-manager-server"
	metalLoadBalancerControllerDeploymentName = "metal-load-balancer-controller-manager"
)

func secretConfigsFunc(namespace string) []extensionssecretsmanager.SecretConfigWithOptions {
	return []extensionssecretsmanager.SecretConfigWithOptions{
		{
			Config: &secretutils.CertificateSecretConfig{
				Name:       caNameControlPlane,
				CommonName: caNameControlPlane,
				CertType:   secretutils.CACert,
			},
			Options: []secretsmanager.GenerateOption{secretsmanager.Persist()},
		},
		{
			Config: &secretutils.CertificateSecretConfig{
				Name:                        cloudControllerManagerServerName,
				CommonName:                  metal.CloudControllerManagerName,
				DNSNames:                    kutil.DNSNamesForService(metal.CloudControllerManagerName, namespace),
				CertType:                    secretutils.ServerCert,
				SkipPublishingCACertificate: true,
			},
			Options: []secretsmanager.GenerateOption{secretsmanager.SignedByCA(caNameControlPlane)},
		},
	}
}

func shootAccessSecretsFunc(namespace string) []*gutil.AccessSecret {
	return []*gutil.AccessSecret{
		gutil.NewShootAccessSecret(cloudControllerManagerDeploymentName, namespace),
		gutil.NewShootAccessSecret(metalLoadBalancerControllerDeploymentName, namespace),
	}
}

var (
	configChart = &chart.Chart{
		Name:       "cloud-provider-config",
		EmbeddedFS: charts.InternalChart,
		Path:       filepath.Join(charts.InternalChartsPath, "cloud-provider-config"),
		Objects: []*chart.Object{
			{Type: &corev1.ConfigMap{}, Name: internal.CloudProviderConfigMapName},
		},
	}

	controlPlaneChart = &chart.Chart{
		Name:       "seed-controlplane",
		EmbeddedFS: charts.InternalChart,
		Path:       filepath.Join(charts.InternalChartsPath, "seed-controlplane"),
		SubCharts: []*chart.Chart{
			{
				Name:   metal.CloudControllerManagerName,
				Images: []string{metal.CloudControllerManagerImageName},
				Objects: []*chart.Object{
					{Type: &corev1.Service{}, Name: "cloud-controller-manager"},
					{Type: &appsv1.Deployment{}, Name: "cloud-controller-manager"},
					{Type: &corev1.ConfigMap{}, Name: "cloud-controller-manager-observability-config"},
					{Type: &autoscalingv1.VerticalPodAutoscaler{}, Name: "cloud-controller-manager-vpa"},
				},
			},
			{
				Name:   "metal-load-balancer-controller-manager",
				Path:   filepath.Join(charts.InternalChartsPath, "metal-load-balancer-controller-manager"),
				Images: []string{metal.MetalLoadBalancerControllerManagerImageName},
				Objects: []*chart.Object{
					{Type: &appsv1.Deployment{}, Name: "metal-load-balancer-controller-manager"},
				},
			},
		},
	}

	controlPlaneShootChart = &chart.Chart{
		Name:       "shoot-system-components",
		EmbeddedFS: charts.InternalChart,
		Path:       filepath.Join(charts.InternalChartsPath, "shoot-system-components"),
		SubCharts: []*chart.Chart{
			{
				Name: "cloud-controller-manager",
				Path: filepath.Join(charts.InternalChartsPath, "cloud-controller-manager"),
				Objects: []*chart.Object{
					{Type: &rbacv1.ClusterRole{}, Name: "system:controller:cloud-node-controller"},
					{Type: &rbacv1.ClusterRoleBinding{}, Name: "system:controller:cloud-node-controller"},
					{Type: &rbacv1.ClusterRoleBinding{}, Name: "metal:cloud-provider"},
				},
			},
			{
				Name:   "metallb",
				Path:   filepath.Join(charts.InternalChartsPath, "metallb"),
				Images: []string{metal.MetallbControllerImageName, metal.MetallbSpeakerImageName},
				Objects: []*chart.Object{
					{Type: &rbacv1.ClusterRole{}, Name: "metallb:controller"},
					{Type: &rbacv1.ClusterRole{}, Name: "metallb:speaker"},
					{Type: &rbacv1.ClusterRoleBinding{}, Name: "metallb:controller"},
					{Type: &rbacv1.ClusterRoleBinding{}, Name: "metallb:speaker"},
					{Type: &corev1.ConfigMap{}, Name: "metallb-excludel2"},
					{Type: &appsv1.DaemonSet{}, Name: "metallb-speaker"},
					{Type: &appsv1.Deployment{}, Name: "metallb-controller"},
					{Type: &rbacv1.Role{}, Name: "metallb-controller"},
					{Type: &rbacv1.Role{}, Name: "metallb-pod-lister"},
					{Type: &rbacv1.RoleBinding{}, Name: "metallb-controller"},
					{Type: &rbacv1.RoleBinding{}, Name: "metallb-pod-lister"},
					{Type: &corev1.Secret{}, Name: "metallb-webhook-cert"},
					{Type: &corev1.Service{}, Name: "metallb-webhook-service"},
					{Type: &corev1.ServiceAccount{}, Name: "metallb-controller"},
					{Type: &corev1.ServiceAccount{}, Name: "metallb-speaker"},
				},
			},
			{
				Name:   "metal-load-balancer-controller-speaker",
				Path:   filepath.Join(charts.InternalChartsPath, "metal-load-balancer-controller-speaker"),
				Images: []string{metal.MetalLoadBalancerControllerSpeakerImageName},
				Objects: []*chart.Object{
					{Type: &rbacv1.ClusterRole{}, Name: "metal-load-balancer-controller:speaker"},
					{Type: &rbacv1.ClusterRoleBinding{}, Name: "metal-load-balancer-controller:speaker"},
					{Type: &appsv1.DaemonSet{}, Name: "metal-load-balancer-controller-speaker"},
					{Type: &corev1.ServiceAccount{}, Name: "metal-load-balancer-controller-speaker"},
					{Type: &rbacv1.ClusterRole{}, Name: "metal-load-balancer-controller:manager"},
					{Type: &rbacv1.ClusterRoleBinding{}, Name: "metal-load-balancer-controller:manager"},
				},
			},
		},
	}
)

// valuesProvider is a ValuesProvider that provides metal-specific values for the 2 charts applied by the generic actuator.
type valuesProvider struct {
	client  client.Client
	decoder runtime.Decoder
}

// NewValuesProvider creates a new ValuesProvider for the generic actuator.
func NewValuesProvider(mgr manager.Manager) genericactuator.ValuesProvider {
	return &valuesProvider{
		client:  mgr.GetClient(),
		decoder: serializer.NewCodecFactory(mgr.GetScheme(), serializer.EnableStrict).UniversalDecoder(),
	}
}

func (vp *valuesProvider) GetControlPlaneExposureChartValues(ctx context.Context,
	cp *extensionsv1alpha1.ControlPlane,
	cluster *extensionscontroller.Cluster,
	secretsReader secretsmanager.Reader,
	checksums map[string]string) (map[string]any, error) {
	return map[string]any{}, nil
}

// GetConfigChartValues returns the values for the config chart applied by the generic actuator.
func (vp *valuesProvider) GetConfigChartValues(
	ctx context.Context,
	cp *extensionsv1alpha1.ControlPlane,
	cluster *extensionscontroller.Cluster,
) (map[string]any, error) {
	cpConfig := &metalapi.ControlPlaneConfig{}
	if cp.Spec.ProviderConfig != nil {
		if _, _, err := vp.decoder.Decode(cp.Spec.ProviderConfig.Raw, nil, cpConfig); err != nil {
			return nil, fmt.Errorf("could not decode providerConfig of controlplane '%s': %w", client.ObjectKeyFromObject(cp), err)
		}
	}
	return vp.getConfigChartValues(cluster, cpConfig)
}

// GetControlPlaneChartValues returns the values for the control plane chart applied by the generic actuator.
func (vp *valuesProvider) GetControlPlaneChartValues(
	ctx context.Context,
	cp *extensionsv1alpha1.ControlPlane,
	cluster *extensionscontroller.Cluster,
	secretsReader secretsmanager.Reader,
	checksums map[string]string,
	scaledDown bool,
) (
	map[string]any,
	error,
) {
	cpConfig := &metalapi.ControlPlaneConfig{}
	if cp.Spec.ProviderConfig != nil {
		if _, _, err := vp.decoder.Decode(cp.Spec.ProviderConfig.Raw, nil, cpConfig); err != nil {
			return nil, fmt.Errorf("could not decode providerConfig of controlplane '%s': %w", client.ObjectKeyFromObject(cp), err)
		}
	}

	return getControlPlaneChartValues(cpConfig, cp, cluster, secretsReader, checksums, scaledDown)
}

// GetControlPlaneShootChartValues returns the values for the control plane shoot chart applied by the generic actuator.
func (vp *valuesProvider) GetControlPlaneShootChartValues(
	_ context.Context,
	cp *extensionsv1alpha1.ControlPlane,
	cluster *extensionscontroller.Cluster,
	_ secretsmanager.Reader,
	_ map[string]string,
) (
	map[string]any,
	error,
) {
	cpConfig := &metalapi.ControlPlaneConfig{}
	if cp.Spec.ProviderConfig != nil {
		if _, _, err := vp.decoder.Decode(cp.Spec.ProviderConfig.Raw, nil, cpConfig); err != nil {
			return nil, fmt.Errorf("could not decode providerConfig of controlplane '%s': %w", client.ObjectKeyFromObject(cp), err)
		}
	}
	return vp.getControlPlaneShootChartValues(cluster, cpConfig)
}

// GetControlPlaneShootCRDsChartValues returns the values for the control plane shoot CRDs chart applied by the generic actuator.
// Currently, the provider extension does not specify a control plane shoot CRDs chart. That's why we simply return empty values.
func (vp *valuesProvider) GetControlPlaneShootCRDsChartValues(
	_ context.Context,
	_ *extensionsv1alpha1.ControlPlane,
	_ *extensionscontroller.Cluster,
) (map[string]any, error) {
	return map[string]any{}, nil
}

// GetStorageClassesChartValues returns the values for the storage classes chart applied by the generic actuator.
func (vp *valuesProvider) GetStorageClassesChartValues(
	ctx context.Context,
	controlPlane *extensionsv1alpha1.ControlPlane,
	cluster *extensionscontroller.Cluster,
) (map[string]any, error) {
	values := make(map[string]any)
	return values, nil
}

// getControlPlaneChartValues collects and returns the control plane chart values.
func getControlPlaneChartValues(
	cpConfig *metalapi.ControlPlaneConfig,
	cp *extensionsv1alpha1.ControlPlane,
	cluster *extensionscontroller.Cluster,
	secretsReader secretsmanager.Reader,
	checksums map[string]string,
	scaledDown bool,
) (
	map[string]any,
	error,
) {
	ccm, err := getCCMChartValues(cpConfig, cp, cluster, secretsReader, checksums, scaledDown)
	if err != nil {
		return nil, err
	}

	metalLoadBalancerControllerManager, err := getMetalLoadBalancerControllerManagerChartValues(cpConfig)
	if err != nil {
		return nil, err
	}

	return map[string]any{
		"global": map[string]any{
			"genericTokenKubeconfigSecretName": extensionscontroller.GenericTokenKubeconfigSecretNameFromCluster(cluster),
		},
		metal.CloudControllerManagerName:             ccm,
		metal.MetalLoadBalancerControllerManagerName: metalLoadBalancerControllerManager,
	}, nil
}

func getMetalLoadBalancerControllerManagerChartValues(config *metalapi.ControlPlaneConfig) (map[string]any, error) {
	if config.LoadBalancerConfig == nil || config.LoadBalancerConfig.MetalLoadBalancerConfig == nil {
		return map[string]any{
			"enabled": false,
		}, nil
	}

	return map[string]any{
		"enabled":           true,
		"nodeCIDRMask":      config.LoadBalancerConfig.MetalLoadBalancerConfig.NodeCIDRMask,
		"allocateNodeCIDRs": config.LoadBalancerConfig.MetalLoadBalancerConfig.AllocateNodeCIDRs,
	}, nil
}

// getCCMChartValues collects and returns the CCM chart values.
func getCCMChartValues(
	cpConfig *metalapi.ControlPlaneConfig,
	cp *extensionsv1alpha1.ControlPlane,
	cluster *extensionscontroller.Cluster,
	secretsReader secretsmanager.Reader,
	checksums map[string]string,
	scaledDown bool,
) (map[string]any, error) {
	serverSecret, found := secretsReader.Get(cloudControllerManagerServerName)
	if !found {
		return nil, fmt.Errorf("secret %q not found", cloudControllerManagerServerName)
	}

	podLabels := map[string]any{
		v1beta1constants.LabelPodMaintenanceRestart: "true",
	}
	localAPI, ok := cluster.Seed.Annotations[metal.LocalMetalAPIAnnotation]
	if ok && localAPI == "true" {
		podLabels[metal.AllowEgressToIstioIngressLabel] = "allowed"
	}

	values := map[string]any{
		"enabled":     true,
		"replicas":    extensionscontroller.GetControlPlaneReplicas(cluster, scaledDown, 1),
		"clusterName": cp.Namespace,
		"podNetwork":  strings.Join(extensionscontroller.GetPodNetwork(cluster), ","),
		"podAnnotations": map[string]any{
			"checksum/config-" + internal.CloudProviderConfigMapName:      checksums[internal.CloudProviderConfigMapName],
			"checksum/secret-" + v1beta1constants.SecretNameCloudProvider: checksums[v1beta1constants.SecretNameCloudProvider],
		},
		"podLabels":       podLabels,
		"tlsCipherSuites": kutil.TLSCipherSuites,
		"secrets": map[string]any{
			"server": serverSecret.Name,
		},
	}

	if cpConfig.CloudControllerManager != nil {
		values[metal.CloudControllerManagerFeatureGatesKeyName] = cpConfig.CloudControllerManager.FeatureGates
	}

	overlayEnabled, err := isOverlayEnabled(cluster.Shoot.Spec.Networking)
	if err != nil {
		return nil, fmt.Errorf("failed to determine if overlay is enabled: %w", err)
	}
	values["configureCloudRoutes"] = !overlayEnabled

	return values, nil
}

func isOverlayEnabled(networking *gardencorev1beta1.Networking) (bool, error) {
	if networking == nil || networking.ProviderConfig == nil {
		return false, nil
	}

	obj, err := runtime.Decode(unstructured.UnstructuredJSONScheme, networking.ProviderConfig.Raw)
	if err != nil {
		return false, err
	}

	u, ok := obj.(*unstructured.Unstructured)
	if !ok {
		return false, fmt.Errorf("object %T is not an unstructured.Unstructured", obj)
	}

	enabled, ok, err := unstructured.NestedBool(u.UnstructuredContent(), "overlay", "enabled")
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}

	return enabled, nil
}

// getControlPlaneShootChartValues collects and returns the control plane shoot chart values.
func (vp *valuesProvider) getControlPlaneShootChartValues(cluster *extensionscontroller.Cluster, cp *metalapi.ControlPlaneConfig) (map[string]any, error) {
	if cluster.Shoot == nil {
		return nil, fmt.Errorf("cluster %s does not contain a shoot object", cluster.ObjectMeta.Name)
	}

	metallb, err := getMetallbChartValues(cp)
	if err != nil {
		return nil, err
	}

	calicoBgp, err := getCalicoBgpChartValues(cp, cluster)
	if err != nil {
		return nil, err
	}

	metalLoadBalancerControllerSpeaker, err := getMetalLoadBalancerControllerSpeakerChartValues(cp)
	if err != nil {
		return nil, fmt.Errorf("failed to get metal load balancer controller chart values: %w", err)
	}

	return map[string]any{
		metal.CloudControllerManagerName:             map[string]any{"enabled": true},
		metal.MetallbName:                            metallb,
		metal.CalicoBgpName:                          calicoBgp,
		metal.MetalLoadBalancerControllerSpeakerName: metalLoadBalancerControllerSpeaker,
	}, nil
}

// getMetalLoadBalancerControllerSpeakerChartValues collects and returns the Metal Load Balancer Controller chart values.
func getMetalLoadBalancerControllerSpeakerChartValues(cpConfig *metalapi.ControlPlaneConfig) (map[string]any, error) {
	if cpConfig.LoadBalancerConfig == nil || cpConfig.LoadBalancerConfig.MetalLoadBalancerConfig == nil {
		return map[string]any{
			"enabled": false,
		}, nil
	}

	return map[string]any{
		"enabled":         true,
		"vni":             cpConfig.LoadBalancerConfig.MetalLoadBalancerConfig.VNI,
		"metalBondServer": cpConfig.LoadBalancerConfig.MetalLoadBalancerConfig.MetalBondServer,
	}, nil
}

// getConfigChartValues collects and returns the config chart values.
func (vp *valuesProvider) getConfigChartValues(cluster *extensionscontroller.Cluster, cpConfig *metalapi.ControlPlaneConfig) (map[string]any, error) {
	values := map[string]any{
		metal.ClusterFieldName: cluster.ObjectMeta.Name,
	}

	if cpConfig.CloudControllerManager != nil {
		if cpConfig.CloudControllerManager.Networking != nil {
			values[metal.CloudControllerManagerNetworkingKeyName] = map[string]any{
				metal.CloudControllerManagerNodeAddressesConfigKeyName: cpConfig.CloudControllerManager.Networking.ConfigureNodeAddresses,
			}
			ipamKind := cpConfig.CloudControllerManager.Networking.IPAMKind
			if ipamKind != nil {
				values[metal.CloudControllerManagerNetworkingKeyName] = map[string]any{
					metal.CloudControllerManagerNodeIPAMKindKeyName: map[string]any{
						"apiGroup": ipamKind.APIGroup,
						"kind":     ipamKind.Kind,
					},
				}
			}
		}
	}
	return values, nil
}

// getMetallbChartValues collects and returns the MetalLB chart values.
func getMetallbChartValues(cpConfig *metalapi.ControlPlaneConfig) (map[string]any, error) {
	if cpConfig.LoadBalancerConfig == nil || cpConfig.LoadBalancerConfig.MetallbConfig == nil {
		return map[string]any{
			"enabled": false,
		}, nil
	}

	for _, cidr := range cpConfig.LoadBalancerConfig.MetallbConfig.IPAddressPool {
		if err := parseAddressPool(cidr); err != nil {
			return nil, fmt.Errorf("invalid CIDR %q in pool: %w", cidr, err)
		}
	}

	return map[string]any{
		"enabled": true,
		"speaker": map[string]any{
			"enabled": cpConfig.LoadBalancerConfig.MetallbConfig.EnableSpeaker,
		},
		"l2Advertisement": map[string]any{
			"enabled": cpConfig.LoadBalancerConfig.MetallbConfig.EnableL2Advertisement,
		},
		"ipAddressPool": cpConfig.LoadBalancerConfig.MetallbConfig.IPAddressPool,
	}, nil
}

// getCalicoBgpChartValues collects and returns the Calico BGP chart values.
func getCalicoBgpChartValues(
	cpConfig *metalapi.ControlPlaneConfig,
	cluster *extensionscontroller.Cluster,
) (map[string]any, error) {
	if cpConfig.LoadBalancerConfig == nil || cpConfig.LoadBalancerConfig.CalicoBgpConfig == nil {
		return map[string]any{
			"enabled": false,
			"bgp": map[string]any{
				"enabled": false,
			},
		}, nil
	}

	var serviceLbIPs, serviceExtIPs, serviceClusterIPs []string
	var peers []map[string]any
	var filters []map[string]any
	if cpConfig.LoadBalancerConfig.CalicoBgpConfig != nil &&
		*cluster.Shoot.Spec.Networking.Type == metal.ShootCalicoNetworkType {
		if cpConfig.LoadBalancerConfig.CalicoBgpConfig.ServiceLoadBalancerIPs != nil {
			for _, cidr := range cpConfig.LoadBalancerConfig.CalicoBgpConfig.ServiceLoadBalancerIPs {
				if err := parseAddressPool(cidr); err != nil {
					return nil, fmt.Errorf("invalid CIDR %q in pool: %w", cidr, err)
				}
				serviceLbIPs = append(serviceLbIPs, cidr)
			}
		}

		if cpConfig.LoadBalancerConfig.CalicoBgpConfig.ServiceExternalIPs != nil {
			for _, cidr := range cpConfig.LoadBalancerConfig.CalicoBgpConfig.ServiceExternalIPs {
				if err := parseAddressPool(cidr); err != nil {
					return nil, fmt.Errorf("invalid CIDR %q in pool: %w", cidr, err)
				}
				serviceExtIPs = append(serviceExtIPs, cidr)
			}
		}

		if cpConfig.LoadBalancerConfig.CalicoBgpConfig.ServiceClusterIPs != nil {
			for _, cidr := range cpConfig.LoadBalancerConfig.CalicoBgpConfig.ServiceClusterIPs {
				if err := parseAddressPool(cidr); err != nil {
					return nil, fmt.Errorf("invalid CIDR %q in pool: %w", cidr, err)
				}
				serviceClusterIPs = append(serviceClusterIPs, cidr)
			}
		}

		if cpConfig.LoadBalancerConfig.CalicoBgpConfig.BGPFilter != nil {
			for _, bgpFilter := range cpConfig.LoadBalancerConfig.CalicoBgpConfig.BGPFilter {
				var exportV4Filters, exportV6Filters, importV4Filters, importV6Filters []map[string]any
				var err error

				if bgpFilter.ExportV4 != nil {
					exportV4Filters, err = processFilters(bgpFilter.ExportV4)
					if err != nil {
						return nil, err
					}
				}
				if bgpFilter.ImportV4 != nil {
					importV4Filters, err = processFilters(bgpFilter.ImportV4)
					if err != nil {
						return nil, err
					}
				}
				if bgpFilter.ExportV6 != nil {
					exportV6Filters, err = processFilters(bgpFilter.ExportV6)
					if err != nil {
						return nil, err
					}
				}
				if bgpFilter.ImportV6 != nil {
					importV6Filters, err = processFilters(bgpFilter.ImportV6)
					if err != nil {
						return nil, err
					}
				}

				var filterMap = map[string]any{
					"name": bgpFilter.Name,
				}
				if len(exportV4Filters) > 0 {
					filterMap["exportV4"] = exportV4Filters
				}
				if len(importV4Filters) > 0 {
					filterMap["importV4"] = importV4Filters
				}
				if len(exportV6Filters) > 0 {
					filterMap["exportV6"] = exportV6Filters
				}
				if len(importV6Filters) > 0 {
					filterMap["importV6"] = importV6Filters
				}
				filters = append(filters, filterMap)
			}
		}

		if cpConfig.LoadBalancerConfig.CalicoBgpConfig.BgpPeer != nil {
			for _, peer := range cpConfig.LoadBalancerConfig.CalicoBgpConfig.BgpPeer {
				peerMap := map[string]any{
					"peerIP":       peer.PeerIP,
					"asNumber":     peer.ASNumber,
					"nodeSelector": peer.NodeSelector,
				}
				if len(peer.Filters) > 0 {
					peerMap["filters"] = peer.Filters
				}
				peers = append(peers, peerMap)
			}
		}
	}

	bgpValues := map[string]any{
		"enabled":                true,
		"asNumber":               cpConfig.LoadBalancerConfig.CalicoBgpConfig.ASNumber,
		"serviceLoadBalancerIPs": serviceLbIPs,
		"serviceExternalIPs":     serviceExtIPs,
		"serviceClusterIPs":      serviceClusterIPs,
		"nodeToNodeMeshEnabled":  cpConfig.LoadBalancerConfig.CalicoBgpConfig.NodeToNodeMeshEnabled,
		"bgpPeer":                peers,
	}

	if len(filters) > 0 {
		bgpValues["bgpFilter"] = filters
	}

	return map[string]any{
		"enabled": true,
		"bgp":     bgpValues,
	}, nil
}

func processFilters(filtersConfig []metalapi.BGPFilterRule) ([]map[string]any, error) {
	var filters []map[string]any
	for _, filter := range filtersConfig {
		if err := parseAddressPool(filter.CIDR); err != nil {
			return nil, fmt.Errorf("invalid CIDR %q in pool: %w", filter.CIDR, err)
		}
		filterMap := map[string]any{
			"cidr":          filter.CIDR,
			"action":        filter.Action,
			"matchOperator": filter.MatchOperator,
		}
		filters = append(filters, filterMap)
	}
	return filters, nil
}

func parseAddressPool(cidr string) error {
	if !strings.Contains(cidr, "-") {
		_, _, err := net.ParseCIDR(cidr)
		if err != nil {
			return fmt.Errorf("invalid CIDR %q", cidr)
		}
		return nil
	}
	fs := strings.SplitN(cidr, "-", 2)
	if len(fs) != 2 {
		return fmt.Errorf("invalid IP range %q", cidr)
	}
	start := net.ParseIP(strings.TrimSpace(fs[0]))
	if start == nil {
		return fmt.Errorf("invalid IP range %q: invalid start IP %q", cidr, fs[0])
	}
	end := net.ParseIP(strings.TrimSpace(fs[1]))
	if end == nil {
		return fmt.Errorf("invalid IP range %q: invalid end IP %q", cidr, fs[1])
	}
	if bytes.Compare(start, end) > 0 {
		return fmt.Errorf("invalid IP range %q: start IP %q is after the end IP %q", cidr, start, end)
	}
	return nil
}
