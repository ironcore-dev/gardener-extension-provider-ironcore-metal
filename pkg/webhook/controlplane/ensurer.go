// SPDX-FileCopyrightText: 2024 SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package controlplane

import (
	"context"
	"fmt"
	"slices"
	"strconv"

	"github.com/Masterminds/semver/v3"
	"github.com/coreos/go-systemd/v22/unit"
	extensionswebhook "github.com/gardener/gardener/extensions/pkg/webhook"
	extensionscontextwebhook "github.com/gardener/gardener/extensions/pkg/webhook/context"
	"github.com/gardener/gardener/extensions/pkg/webhook/controlplane/genericmutator"
	"github.com/gardener/gardener/pkg/apis/core/v1beta1"
	extensionsv1alpha1 "github.com/gardener/gardener/pkg/apis/extensions/v1alpha1"
	"github.com/gardener/gardener/pkg/component/nodemanagement/machinecontrollermanager"
	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	vpaautoscalingv1 "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/apis/autoscaling.k8s.io/v1"
	kubeletconfigv1beta1 "k8s.io/kubelet/config/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/ironcore-dev/gardener-extension-provider-ironcore-metal/imagevector"
	metalapi "github.com/ironcore-dev/gardener-extension-provider-ironcore-metal/pkg/apis/metal"
	"github.com/ironcore-dev/gardener-extension-provider-ironcore-metal/pkg/metal"
)

// NewEnsurer creates a new controlplane ensurer.
func NewEnsurer(logger logr.Logger, scheme *runtime.Scheme) genericmutator.Ensurer {
	return &ensurer{
		logger:  logger.WithName("ironcore-metal-controlplane-ensurer"),
		decoder: serializer.NewCodecFactory(scheme, serializer.EnableStrict).UniversalDecoder(),
	}
}

type ensurer struct {
	genericmutator.NoopEnsurer
	logger  logr.Logger
	decoder runtime.Decoder
}

// ImageVector is exposed for testing.
var ImageVector = imagevector.ImageVector()

// EnsureMachineControllerManagerDeployment ensures that the machine-controller-manager deployment conforms to the provider requirements.
func (e *ensurer) EnsureMachineControllerManagerDeployment(ctx context.Context, gctx extensionscontextwebhook.GardenContext, newObj, _ *appsv1.Deployment) error {
	image, err := ImageVector.FindImage(metal.MachineControllerManagerProviderIroncoreImageName)
	if err != nil {
		return err
	}
	cluster, err := gctx.GetCluster(ctx)
	if err != nil {
		return fmt.Errorf("failed to get cluster: %w", err)
	}

	cpConfig := &metalapi.ControlPlaneConfig{}
	if cluster.Shoot != nil && cluster.Shoot.Spec.Provider.ControlPlaneConfig != nil {
		cp := cluster.Shoot.Spec.Provider.ControlPlaneConfig
		if _, _, err := e.decoder.Decode(cp.Raw, nil, cpConfig); err != nil {
			return fmt.Errorf("could not decode providerConfig of controlplane for cluster %s: %w", client.ObjectKeyFromObject(cluster.Shoot), err)
		}
	}

	template := &newObj.Spec.Template
	ps := &template.Spec

	localAPI, ok := cluster.Seed.Annotations[metal.LocalMetalAPIAnnotation]
	if ok && localAPI == "true" {
		template.Labels = extensionswebhook.EnsureAnnotationOrLabel(template.Labels, metal.AllowEgressToIstioIngressLabel, "allowed")
	}

	ps.Containers = extensionswebhook.EnsureContainerWithName(
		newObj.Spec.Template.Spec.Containers,
		machinecontrollermanager.ProviderSidecarContainer(cluster.Shoot, newObj.Namespace, metal.ProviderName, image.String()),
	)

	if c := extensionswebhook.ContainerWithName(ps.Containers, metal.MachineControllerManagerProviderIroncoreImageName); c != nil {
		ensureMCMCommandLineArgs(c, cpConfig)
		c.VolumeMounts = extensionswebhook.EnsureVolumeMountWithName(c.VolumeMounts, corev1.VolumeMount{
			Name:      "cloudprovider",
			MountPath: "/etc/metal",
			ReadOnly:  true,
		})
	}

	ps.Volumes = extensionswebhook.EnsureVolumeWithName(ps.Volumes, corev1.Volume{
		Name: "cloudprovider",
		VolumeSource: corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				SecretName: "cloudprovider",
			},
		},
	})

	return nil
}

// EnsureMachineControllerManagerVPA ensures that the machine-controller-manager VPA conforms to the provider requirements.
func (e *ensurer) EnsureMachineControllerManagerVPA(_ context.Context, _ extensionscontextwebhook.GardenContext, newObj, _ *vpaautoscalingv1.VerticalPodAutoscaler) error {
	if newObj.Spec.ResourcePolicy == nil {
		newObj.Spec.ResourcePolicy = &vpaautoscalingv1.PodResourcePolicy{}
	}

	newObj.Spec.ResourcePolicy.ContainerPolicies = extensionswebhook.EnsureVPAContainerResourcePolicyWithName(
		newObj.Spec.ResourcePolicy.ContainerPolicies,
		machinecontrollermanager.ProviderSidecarVPAContainerPolicy(metal.ProviderName),
	)
	return nil
}

// EnsureKubeAPIServerDeployment ensures that the kube-apiserver deployment conforms to the provider requirements.
func (e *ensurer) EnsureKubeAPIServerDeployment(_ context.Context, _ extensionscontextwebhook.GardenContext, new, _ *appsv1.Deployment) error {
	template := &new.Spec.Template
	ps := &template.Spec

	if c := extensionswebhook.ContainerWithName(ps.Containers, "kube-apiserver"); c != nil {
		ensureKubeAPIServerCommandLineArgs(c)
	}

	return nil
}

// EnsureKubeControllerManagerDeployment ensures that the kube-controller-manager deployment conforms to the provider requirements.
func (e *ensurer) EnsureKubeControllerManagerDeployment(ctx context.Context, gctx extensionscontextwebhook.GardenContext, new, _ *appsv1.Deployment) error {
	template := &new.Spec.Template
	ps := &template.Spec

	cluster, err := gctx.GetCluster(ctx)
	if err != nil {
		return err
	}

	allocateNodeCIDRs := true
	if networkingConfig := cluster.Shoot.Spec.Networking; networkingConfig != nil && slices.Contains(networkingConfig.IPFamilies, v1beta1.IPFamilyIPv6) {
		allocateNodeCIDRs = false
	}

	if c := extensionswebhook.ContainerWithName(ps.Containers, "kube-controller-manager"); c != nil {
		ensureKubeControllerManagerCommandLineArgs(c, allocateNodeCIDRs)
	}

	return nil
}

// EnsureKubeSchedulerDeployment ensures that the kube-scheduler deployment conforms to the provider requirements.
func (e *ensurer) EnsureKubeSchedulerDeployment(_ context.Context, _ extensionscontextwebhook.GardenContext, _, _ *appsv1.Deployment) error {
	return nil
}

// EnsureClusterAutoscalerDeployment ensures that the cluster-autoscaler deployment conforms to the provider requirements.
func (e *ensurer) EnsureClusterAutoscalerDeployment(_ context.Context, _ extensionscontextwebhook.GardenContext, _, _ *appsv1.Deployment) error {
	return nil
}

func ensureMCMCommandLineArgs(c *corev1.Container, cp *metalapi.ControlPlaneConfig) {
	c.Args = extensionswebhook.EnsureStringWithPrefix(c.Args, "--metal-kubeconfig=", "/etc/metal/kubeconfig")
	if cp.NodeNamePolicy == "" {
		return
	}
	switch cp.NodeNamePolicy {
	case metalapi.NodeNamePolicyBMCName, metalapi.NodeNamePolicyServerName, metalapi.NodeNamePolicyServerClaimName:
		c.Args = extensionswebhook.EnsureStringWithPrefix(c.Args, "--node-name-policy=", string(cp.NodeNamePolicy))
	}
}

func ensureKubeAPIServerCommandLineArgs(c *corev1.Container) {
	c.Command = extensionswebhook.EnsureNoStringWithPrefix(c.Command, "--cloud-provider=")
	c.Command = extensionswebhook.EnsureNoStringWithPrefix(c.Command, "--cloud-config=")
}

func ensureKubeControllerManagerCommandLineArgs(c *corev1.Container, allocateNodeCIDRs bool) {
	c.Command = extensionswebhook.EnsureStringWithPrefix(c.Command, "--cloud-provider=", "external")
	c.Command = extensionswebhook.EnsureNoStringWithPrefix(c.Command, "--cloud-config=")

	// allocate-node-cidrs is a boolean flag and could be enabled by name without an explicit value passed. Therefore, we delete all prefixes (without including "=" in the prefix)
	c.Command = extensionswebhook.EnsureNoStringWithPrefix(c.Command, "--allocate-node-cidrs")
	c.Command = extensionswebhook.EnsureStringWithPrefix(c.Command, "--allocate-node-cidrs=", strconv.FormatBool(allocateNodeCIDRs))
}

// EnsureKubeletServiceUnitOptions ensures that the kubelet.service unit options conform to the provider requirements.
func (e *ensurer) EnsureKubeletServiceUnitOptions(_ context.Context, _ extensionscontextwebhook.GardenContext, _ *semver.Version, new, _ []*unit.UnitOption) ([]*unit.UnitOption, error) {
	if opt := extensionswebhook.UnitOptionWithSectionAndName(new, "Service", "ExecStart"); opt != nil {
		command := extensionswebhook.DeserializeCommandLine(opt.Value)
		command = ensureKubeletCommandLineArgs(command)
		opt.Value = extensionswebhook.SerializeCommandLine(command, 1, " \\\n    ")
	}

	new = extensionswebhook.EnsureUnitOption(new, &unit.UnitOption{
		Section: "Service",
		Name:    "ExecStartPre",
		Value:   `/bin/sh -c 'hostnamectl set-hostname $(hostname -f)'`,
	})

	return new, nil
}

func ensureKubeletCommandLineArgs(command []string) []string {
	command = extensionswebhook.EnsureStringWithPrefix(command, "--cloud-provider=", "external")
	return command
}

// EnsureKubeletConfiguration ensures that the kubelet configuration conforms to the provider requirements.
func (e *ensurer) EnsureKubeletConfiguration(_ context.Context, _ extensionscontextwebhook.GardenContext, _ *semver.Version, _, _ *kubeletconfigv1beta1.KubeletConfiguration) error {
	return nil
}

// EnsureKubernetesGeneralConfiguration ensures that the kubernetes general configuration conforms to the provider requirements.
func (e *ensurer) EnsureKubernetesGeneralConfiguration(_ context.Context, _ extensionscontextwebhook.GardenContext, _, _ *string) error {
	return nil
}

// EnsureAdditionalUnits ensures that additional required system units are added.
func (e *ensurer) EnsureAdditionalUnits(_ context.Context, _ extensionscontextwebhook.GardenContext, _, _ *[]extensionsv1alpha1.Unit) error {
	return nil
}

// EnsureAdditionalFiles ensures that additional required system files are added.
func (e *ensurer) EnsureAdditionalFiles(_ context.Context, _ extensionscontextwebhook.GardenContext, _, _ *[]extensionsv1alpha1.File) error {
	return nil
}
