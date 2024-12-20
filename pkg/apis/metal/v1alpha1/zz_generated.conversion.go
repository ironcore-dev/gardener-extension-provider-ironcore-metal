//go:build !ignore_autogenerated
// +build !ignore_autogenerated

// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

// Code generated by conversion-gen. DO NOT EDIT.

package v1alpha1

import (
	unsafe "unsafe"

	metal "github.com/ironcore-dev/gardener-extension-provider-metal/pkg/apis/metal"
	v1 "k8s.io/api/core/v1"
	conversion "k8s.io/apimachinery/pkg/conversion"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

func init() {
	localSchemeBuilder.Register(RegisterConversions)
}

// RegisterConversions adds conversion functions to the given scheme.
// Public to allow building arbitrary schemes.
func RegisterConversions(s *runtime.Scheme) error {
	if err := s.AddGeneratedConversionFunc((*BgpPeer)(nil), (*metal.BgpPeer)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_BgpPeer_To_metal_BgpPeer(a.(*BgpPeer), b.(*metal.BgpPeer), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*metal.BgpPeer)(nil), (*BgpPeer)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_metal_BgpPeer_To_v1alpha1_BgpPeer(a.(*metal.BgpPeer), b.(*BgpPeer), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*CalicoBgpConfig)(nil), (*metal.CalicoBgpConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_CalicoBgpConfig_To_metal_CalicoBgpConfig(a.(*CalicoBgpConfig), b.(*metal.CalicoBgpConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*metal.CalicoBgpConfig)(nil), (*CalicoBgpConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_metal_CalicoBgpConfig_To_v1alpha1_CalicoBgpConfig(a.(*metal.CalicoBgpConfig), b.(*CalicoBgpConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*CloudControllerManagerConfig)(nil), (*metal.CloudControllerManagerConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_CloudControllerManagerConfig_To_metal_CloudControllerManagerConfig(a.(*CloudControllerManagerConfig), b.(*metal.CloudControllerManagerConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*metal.CloudControllerManagerConfig)(nil), (*CloudControllerManagerConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_metal_CloudControllerManagerConfig_To_v1alpha1_CloudControllerManagerConfig(a.(*metal.CloudControllerManagerConfig), b.(*CloudControllerManagerConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*CloudControllerNetworking)(nil), (*metal.CloudControllerNetworking)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_CloudControllerNetworking_To_metal_CloudControllerNetworking(a.(*CloudControllerNetworking), b.(*metal.CloudControllerNetworking), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*metal.CloudControllerNetworking)(nil), (*CloudControllerNetworking)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_metal_CloudControllerNetworking_To_v1alpha1_CloudControllerNetworking(a.(*metal.CloudControllerNetworking), b.(*CloudControllerNetworking), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*CloudProfileConfig)(nil), (*metal.CloudProfileConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_CloudProfileConfig_To_metal_CloudProfileConfig(a.(*CloudProfileConfig), b.(*metal.CloudProfileConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*metal.CloudProfileConfig)(nil), (*CloudProfileConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_metal_CloudProfileConfig_To_v1alpha1_CloudProfileConfig(a.(*metal.CloudProfileConfig), b.(*CloudProfileConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*ControlPlaneConfig)(nil), (*metal.ControlPlaneConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_ControlPlaneConfig_To_metal_ControlPlaneConfig(a.(*ControlPlaneConfig), b.(*metal.ControlPlaneConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*metal.ControlPlaneConfig)(nil), (*ControlPlaneConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_metal_ControlPlaneConfig_To_v1alpha1_ControlPlaneConfig(a.(*metal.ControlPlaneConfig), b.(*ControlPlaneConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*IgnitionConfig)(nil), (*metal.IgnitionConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_IgnitionConfig_To_metal_IgnitionConfig(a.(*IgnitionConfig), b.(*metal.IgnitionConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*metal.IgnitionConfig)(nil), (*IgnitionConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_metal_IgnitionConfig_To_v1alpha1_IgnitionConfig(a.(*metal.IgnitionConfig), b.(*IgnitionConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*InfrastructureConfig)(nil), (*metal.InfrastructureConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_InfrastructureConfig_To_metal_InfrastructureConfig(a.(*InfrastructureConfig), b.(*metal.InfrastructureConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*metal.InfrastructureConfig)(nil), (*InfrastructureConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_metal_InfrastructureConfig_To_v1alpha1_InfrastructureConfig(a.(*metal.InfrastructureConfig), b.(*InfrastructureConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*InfrastructureStatus)(nil), (*metal.InfrastructureStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_InfrastructureStatus_To_metal_InfrastructureStatus(a.(*InfrastructureStatus), b.(*metal.InfrastructureStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*metal.InfrastructureStatus)(nil), (*InfrastructureStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_metal_InfrastructureStatus_To_v1alpha1_InfrastructureStatus(a.(*metal.InfrastructureStatus), b.(*InfrastructureStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*LoadBalancerConfig)(nil), (*metal.LoadBalancerConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_LoadBalancerConfig_To_metal_LoadBalancerConfig(a.(*LoadBalancerConfig), b.(*metal.LoadBalancerConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*metal.LoadBalancerConfig)(nil), (*LoadBalancerConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_metal_LoadBalancerConfig_To_v1alpha1_LoadBalancerConfig(a.(*metal.LoadBalancerConfig), b.(*LoadBalancerConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*MachineImage)(nil), (*metal.MachineImage)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_MachineImage_To_metal_MachineImage(a.(*MachineImage), b.(*metal.MachineImage), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*metal.MachineImage)(nil), (*MachineImage)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_metal_MachineImage_To_v1alpha1_MachineImage(a.(*metal.MachineImage), b.(*MachineImage), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*MachineImageVersion)(nil), (*metal.MachineImageVersion)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_MachineImageVersion_To_metal_MachineImageVersion(a.(*MachineImageVersion), b.(*metal.MachineImageVersion), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*metal.MachineImageVersion)(nil), (*MachineImageVersion)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_metal_MachineImageVersion_To_v1alpha1_MachineImageVersion(a.(*metal.MachineImageVersion), b.(*MachineImageVersion), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*MachineImages)(nil), (*metal.MachineImages)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_MachineImages_To_metal_MachineImages(a.(*MachineImages), b.(*metal.MachineImages), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*metal.MachineImages)(nil), (*MachineImages)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_metal_MachineImages_To_v1alpha1_MachineImages(a.(*metal.MachineImages), b.(*MachineImages), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*MachineType)(nil), (*metal.MachineType)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_MachineType_To_metal_MachineType(a.(*MachineType), b.(*metal.MachineType), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*metal.MachineType)(nil), (*MachineType)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_metal_MachineType_To_v1alpha1_MachineType(a.(*metal.MachineType), b.(*MachineType), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*MetallbConfig)(nil), (*metal.MetallbConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_MetallbConfig_To_metal_MetallbConfig(a.(*MetallbConfig), b.(*metal.MetallbConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*metal.MetallbConfig)(nil), (*MetallbConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_metal_MetallbConfig_To_v1alpha1_MetallbConfig(a.(*metal.MetallbConfig), b.(*MetallbConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*RegionConfig)(nil), (*metal.RegionConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_RegionConfig_To_metal_RegionConfig(a.(*RegionConfig), b.(*metal.RegionConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*metal.RegionConfig)(nil), (*RegionConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_metal_RegionConfig_To_v1alpha1_RegionConfig(a.(*metal.RegionConfig), b.(*RegionConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*WorkerConfig)(nil), (*metal.WorkerConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_WorkerConfig_To_metal_WorkerConfig(a.(*WorkerConfig), b.(*metal.WorkerConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*metal.WorkerConfig)(nil), (*WorkerConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_metal_WorkerConfig_To_v1alpha1_WorkerConfig(a.(*metal.WorkerConfig), b.(*WorkerConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*WorkerStatus)(nil), (*metal.WorkerStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_WorkerStatus_To_metal_WorkerStatus(a.(*WorkerStatus), b.(*metal.WorkerStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*metal.WorkerStatus)(nil), (*WorkerStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_metal_WorkerStatus_To_v1alpha1_WorkerStatus(a.(*metal.WorkerStatus), b.(*WorkerStatus), scope)
	}); err != nil {
		return err
	}
	return nil
}

func autoConvert_v1alpha1_BgpPeer_To_metal_BgpPeer(in *BgpPeer, out *metal.BgpPeer, s conversion.Scope) error {
	out.PeerIP = in.PeerIP
	out.ASNumber = in.ASNumber
	out.NodeSelector = in.NodeSelector
	return nil
}

// Convert_v1alpha1_BgpPeer_To_metal_BgpPeer is an autogenerated conversion function.
func Convert_v1alpha1_BgpPeer_To_metal_BgpPeer(in *BgpPeer, out *metal.BgpPeer, s conversion.Scope) error {
	return autoConvert_v1alpha1_BgpPeer_To_metal_BgpPeer(in, out, s)
}

func autoConvert_metal_BgpPeer_To_v1alpha1_BgpPeer(in *metal.BgpPeer, out *BgpPeer, s conversion.Scope) error {
	out.PeerIP = in.PeerIP
	out.ASNumber = in.ASNumber
	out.NodeSelector = in.NodeSelector
	return nil
}

// Convert_metal_BgpPeer_To_v1alpha1_BgpPeer is an autogenerated conversion function.
func Convert_metal_BgpPeer_To_v1alpha1_BgpPeer(in *metal.BgpPeer, out *BgpPeer, s conversion.Scope) error {
	return autoConvert_metal_BgpPeer_To_v1alpha1_BgpPeer(in, out, s)
}

func autoConvert_v1alpha1_CalicoBgpConfig_To_metal_CalicoBgpConfig(in *CalicoBgpConfig, out *metal.CalicoBgpConfig, s conversion.Scope) error {
	out.ASNumber = in.ASNumber
	out.ServiceLoadBalancerIPs = *(*[]string)(unsafe.Pointer(&in.ServiceLoadBalancerIPs))
	out.ServiceExternalIPs = *(*[]string)(unsafe.Pointer(&in.ServiceExternalIPs))
	out.ServiceClusterIPs = *(*[]string)(unsafe.Pointer(&in.ServiceClusterIPs))
	out.BgpPeer = *(*[]metal.BgpPeer)(unsafe.Pointer(&in.BgpPeer))
	return nil
}

// Convert_v1alpha1_CalicoBgpConfig_To_metal_CalicoBgpConfig is an autogenerated conversion function.
func Convert_v1alpha1_CalicoBgpConfig_To_metal_CalicoBgpConfig(in *CalicoBgpConfig, out *metal.CalicoBgpConfig, s conversion.Scope) error {
	return autoConvert_v1alpha1_CalicoBgpConfig_To_metal_CalicoBgpConfig(in, out, s)
}

func autoConvert_metal_CalicoBgpConfig_To_v1alpha1_CalicoBgpConfig(in *metal.CalicoBgpConfig, out *CalicoBgpConfig, s conversion.Scope) error {
	out.ASNumber = in.ASNumber
	out.ServiceLoadBalancerIPs = *(*[]string)(unsafe.Pointer(&in.ServiceLoadBalancerIPs))
	out.ServiceExternalIPs = *(*[]string)(unsafe.Pointer(&in.ServiceExternalIPs))
	out.ServiceClusterIPs = *(*[]string)(unsafe.Pointer(&in.ServiceClusterIPs))
	out.BgpPeer = *(*[]BgpPeer)(unsafe.Pointer(&in.BgpPeer))
	return nil
}

// Convert_metal_CalicoBgpConfig_To_v1alpha1_CalicoBgpConfig is an autogenerated conversion function.
func Convert_metal_CalicoBgpConfig_To_v1alpha1_CalicoBgpConfig(in *metal.CalicoBgpConfig, out *CalicoBgpConfig, s conversion.Scope) error {
	return autoConvert_metal_CalicoBgpConfig_To_v1alpha1_CalicoBgpConfig(in, out, s)
}

func autoConvert_v1alpha1_CloudControllerManagerConfig_To_metal_CloudControllerManagerConfig(in *CloudControllerManagerConfig, out *metal.CloudControllerManagerConfig, s conversion.Scope) error {
	out.FeatureGates = *(*map[string]bool)(unsafe.Pointer(&in.FeatureGates))
	out.Networking = (*metal.CloudControllerNetworking)(unsafe.Pointer(in.Networking))
	return nil
}

// Convert_v1alpha1_CloudControllerManagerConfig_To_metal_CloudControllerManagerConfig is an autogenerated conversion function.
func Convert_v1alpha1_CloudControllerManagerConfig_To_metal_CloudControllerManagerConfig(in *CloudControllerManagerConfig, out *metal.CloudControllerManagerConfig, s conversion.Scope) error {
	return autoConvert_v1alpha1_CloudControllerManagerConfig_To_metal_CloudControllerManagerConfig(in, out, s)
}

func autoConvert_metal_CloudControllerManagerConfig_To_v1alpha1_CloudControllerManagerConfig(in *metal.CloudControllerManagerConfig, out *CloudControllerManagerConfig, s conversion.Scope) error {
	out.FeatureGates = *(*map[string]bool)(unsafe.Pointer(&in.FeatureGates))
	out.Networking = (*CloudControllerNetworking)(unsafe.Pointer(in.Networking))
	return nil
}

// Convert_metal_CloudControllerManagerConfig_To_v1alpha1_CloudControllerManagerConfig is an autogenerated conversion function.
func Convert_metal_CloudControllerManagerConfig_To_v1alpha1_CloudControllerManagerConfig(in *metal.CloudControllerManagerConfig, out *CloudControllerManagerConfig, s conversion.Scope) error {
	return autoConvert_metal_CloudControllerManagerConfig_To_v1alpha1_CloudControllerManagerConfig(in, out, s)
}

func autoConvert_v1alpha1_CloudControllerNetworking_To_metal_CloudControllerNetworking(in *CloudControllerNetworking, out *metal.CloudControllerNetworking, s conversion.Scope) error {
	out.ConfigureNodeAddresses = in.ConfigureNodeAddresses
	return nil
}

// Convert_v1alpha1_CloudControllerNetworking_To_metal_CloudControllerNetworking is an autogenerated conversion function.
func Convert_v1alpha1_CloudControllerNetworking_To_metal_CloudControllerNetworking(in *CloudControllerNetworking, out *metal.CloudControllerNetworking, s conversion.Scope) error {
	return autoConvert_v1alpha1_CloudControllerNetworking_To_metal_CloudControllerNetworking(in, out, s)
}

func autoConvert_metal_CloudControllerNetworking_To_v1alpha1_CloudControllerNetworking(in *metal.CloudControllerNetworking, out *CloudControllerNetworking, s conversion.Scope) error {
	out.ConfigureNodeAddresses = in.ConfigureNodeAddresses
	return nil
}

// Convert_metal_CloudControllerNetworking_To_v1alpha1_CloudControllerNetworking is an autogenerated conversion function.
func Convert_metal_CloudControllerNetworking_To_v1alpha1_CloudControllerNetworking(in *metal.CloudControllerNetworking, out *CloudControllerNetworking, s conversion.Scope) error {
	return autoConvert_metal_CloudControllerNetworking_To_v1alpha1_CloudControllerNetworking(in, out, s)
}

func autoConvert_v1alpha1_CloudProfileConfig_To_metal_CloudProfileConfig(in *CloudProfileConfig, out *metal.CloudProfileConfig, s conversion.Scope) error {
	out.MachineImages = *(*[]metal.MachineImages)(unsafe.Pointer(&in.MachineImages))
	out.RegionConfigs = *(*[]metal.RegionConfig)(unsafe.Pointer(&in.RegionConfigs))
	out.MachineTypes = *(*[]metal.MachineType)(unsafe.Pointer(&in.MachineTypes))
	return nil
}

// Convert_v1alpha1_CloudProfileConfig_To_metal_CloudProfileConfig is an autogenerated conversion function.
func Convert_v1alpha1_CloudProfileConfig_To_metal_CloudProfileConfig(in *CloudProfileConfig, out *metal.CloudProfileConfig, s conversion.Scope) error {
	return autoConvert_v1alpha1_CloudProfileConfig_To_metal_CloudProfileConfig(in, out, s)
}

func autoConvert_metal_CloudProfileConfig_To_v1alpha1_CloudProfileConfig(in *metal.CloudProfileConfig, out *CloudProfileConfig, s conversion.Scope) error {
	out.MachineImages = *(*[]MachineImages)(unsafe.Pointer(&in.MachineImages))
	out.RegionConfigs = *(*[]RegionConfig)(unsafe.Pointer(&in.RegionConfigs))
	out.MachineTypes = *(*[]MachineType)(unsafe.Pointer(&in.MachineTypes))
	return nil
}

// Convert_metal_CloudProfileConfig_To_v1alpha1_CloudProfileConfig is an autogenerated conversion function.
func Convert_metal_CloudProfileConfig_To_v1alpha1_CloudProfileConfig(in *metal.CloudProfileConfig, out *CloudProfileConfig, s conversion.Scope) error {
	return autoConvert_metal_CloudProfileConfig_To_v1alpha1_CloudProfileConfig(in, out, s)
}

func autoConvert_v1alpha1_ControlPlaneConfig_To_metal_ControlPlaneConfig(in *ControlPlaneConfig, out *metal.ControlPlaneConfig, s conversion.Scope) error {
	out.CloudControllerManager = (*metal.CloudControllerManagerConfig)(unsafe.Pointer(in.CloudControllerManager))
	out.LoadBalancerConfig = (*metal.LoadBalancerConfig)(unsafe.Pointer(in.LoadBalancerConfig))
	return nil
}

// Convert_v1alpha1_ControlPlaneConfig_To_metal_ControlPlaneConfig is an autogenerated conversion function.
func Convert_v1alpha1_ControlPlaneConfig_To_metal_ControlPlaneConfig(in *ControlPlaneConfig, out *metal.ControlPlaneConfig, s conversion.Scope) error {
	return autoConvert_v1alpha1_ControlPlaneConfig_To_metal_ControlPlaneConfig(in, out, s)
}

func autoConvert_metal_ControlPlaneConfig_To_v1alpha1_ControlPlaneConfig(in *metal.ControlPlaneConfig, out *ControlPlaneConfig, s conversion.Scope) error {
	out.CloudControllerManager = (*CloudControllerManagerConfig)(unsafe.Pointer(in.CloudControllerManager))
	out.LoadBalancerConfig = (*LoadBalancerConfig)(unsafe.Pointer(in.LoadBalancerConfig))
	return nil
}

// Convert_metal_ControlPlaneConfig_To_v1alpha1_ControlPlaneConfig is an autogenerated conversion function.
func Convert_metal_ControlPlaneConfig_To_v1alpha1_ControlPlaneConfig(in *metal.ControlPlaneConfig, out *ControlPlaneConfig, s conversion.Scope) error {
	return autoConvert_metal_ControlPlaneConfig_To_v1alpha1_ControlPlaneConfig(in, out, s)
}

func autoConvert_v1alpha1_IgnitionConfig_To_metal_IgnitionConfig(in *IgnitionConfig, out *metal.IgnitionConfig, s conversion.Scope) error {
	out.Raw = in.Raw
	out.SecretRef = (*v1.LocalObjectReference)(unsafe.Pointer(in.SecretRef))
	out.Override = in.Override
	return nil
}

// Convert_v1alpha1_IgnitionConfig_To_metal_IgnitionConfig is an autogenerated conversion function.
func Convert_v1alpha1_IgnitionConfig_To_metal_IgnitionConfig(in *IgnitionConfig, out *metal.IgnitionConfig, s conversion.Scope) error {
	return autoConvert_v1alpha1_IgnitionConfig_To_metal_IgnitionConfig(in, out, s)
}

func autoConvert_metal_IgnitionConfig_To_v1alpha1_IgnitionConfig(in *metal.IgnitionConfig, out *IgnitionConfig, s conversion.Scope) error {
	out.Raw = in.Raw
	out.SecretRef = (*v1.LocalObjectReference)(unsafe.Pointer(in.SecretRef))
	out.Override = in.Override
	return nil
}

// Convert_metal_IgnitionConfig_To_v1alpha1_IgnitionConfig is an autogenerated conversion function.
func Convert_metal_IgnitionConfig_To_v1alpha1_IgnitionConfig(in *metal.IgnitionConfig, out *IgnitionConfig, s conversion.Scope) error {
	return autoConvert_metal_IgnitionConfig_To_v1alpha1_IgnitionConfig(in, out, s)
}

func autoConvert_v1alpha1_InfrastructureConfig_To_metal_InfrastructureConfig(in *InfrastructureConfig, out *metal.InfrastructureConfig, s conversion.Scope) error {
	return nil
}

// Convert_v1alpha1_InfrastructureConfig_To_metal_InfrastructureConfig is an autogenerated conversion function.
func Convert_v1alpha1_InfrastructureConfig_To_metal_InfrastructureConfig(in *InfrastructureConfig, out *metal.InfrastructureConfig, s conversion.Scope) error {
	return autoConvert_v1alpha1_InfrastructureConfig_To_metal_InfrastructureConfig(in, out, s)
}

func autoConvert_metal_InfrastructureConfig_To_v1alpha1_InfrastructureConfig(in *metal.InfrastructureConfig, out *InfrastructureConfig, s conversion.Scope) error {
	return nil
}

// Convert_metal_InfrastructureConfig_To_v1alpha1_InfrastructureConfig is an autogenerated conversion function.
func Convert_metal_InfrastructureConfig_To_v1alpha1_InfrastructureConfig(in *metal.InfrastructureConfig, out *InfrastructureConfig, s conversion.Scope) error {
	return autoConvert_metal_InfrastructureConfig_To_v1alpha1_InfrastructureConfig(in, out, s)
}

func autoConvert_v1alpha1_InfrastructureStatus_To_metal_InfrastructureStatus(in *InfrastructureStatus, out *metal.InfrastructureStatus, s conversion.Scope) error {
	return nil
}

// Convert_v1alpha1_InfrastructureStatus_To_metal_InfrastructureStatus is an autogenerated conversion function.
func Convert_v1alpha1_InfrastructureStatus_To_metal_InfrastructureStatus(in *InfrastructureStatus, out *metal.InfrastructureStatus, s conversion.Scope) error {
	return autoConvert_v1alpha1_InfrastructureStatus_To_metal_InfrastructureStatus(in, out, s)
}

func autoConvert_metal_InfrastructureStatus_To_v1alpha1_InfrastructureStatus(in *metal.InfrastructureStatus, out *InfrastructureStatus, s conversion.Scope) error {
	return nil
}

// Convert_metal_InfrastructureStatus_To_v1alpha1_InfrastructureStatus is an autogenerated conversion function.
func Convert_metal_InfrastructureStatus_To_v1alpha1_InfrastructureStatus(in *metal.InfrastructureStatus, out *InfrastructureStatus, s conversion.Scope) error {
	return autoConvert_metal_InfrastructureStatus_To_v1alpha1_InfrastructureStatus(in, out, s)
}

func autoConvert_v1alpha1_LoadBalancerConfig_To_metal_LoadBalancerConfig(in *LoadBalancerConfig, out *metal.LoadBalancerConfig, s conversion.Scope) error {
	out.MetallbConfig = (*metal.MetallbConfig)(unsafe.Pointer(in.MetallbConfig))
	out.CalicoBgpConfig = (*metal.CalicoBgpConfig)(unsafe.Pointer(in.CalicoBgpConfig))
	return nil
}

// Convert_v1alpha1_LoadBalancerConfig_To_metal_LoadBalancerConfig is an autogenerated conversion function.
func Convert_v1alpha1_LoadBalancerConfig_To_metal_LoadBalancerConfig(in *LoadBalancerConfig, out *metal.LoadBalancerConfig, s conversion.Scope) error {
	return autoConvert_v1alpha1_LoadBalancerConfig_To_metal_LoadBalancerConfig(in, out, s)
}

func autoConvert_metal_LoadBalancerConfig_To_v1alpha1_LoadBalancerConfig(in *metal.LoadBalancerConfig, out *LoadBalancerConfig, s conversion.Scope) error {
	out.MetallbConfig = (*MetallbConfig)(unsafe.Pointer(in.MetallbConfig))
	out.CalicoBgpConfig = (*CalicoBgpConfig)(unsafe.Pointer(in.CalicoBgpConfig))
	return nil
}

// Convert_metal_LoadBalancerConfig_To_v1alpha1_LoadBalancerConfig is an autogenerated conversion function.
func Convert_metal_LoadBalancerConfig_To_v1alpha1_LoadBalancerConfig(in *metal.LoadBalancerConfig, out *LoadBalancerConfig, s conversion.Scope) error {
	return autoConvert_metal_LoadBalancerConfig_To_v1alpha1_LoadBalancerConfig(in, out, s)
}

func autoConvert_v1alpha1_MachineImage_To_metal_MachineImage(in *MachineImage, out *metal.MachineImage, s conversion.Scope) error {
	out.Name = in.Name
	out.Version = in.Version
	out.Image = in.Image
	out.Architecture = (*string)(unsafe.Pointer(in.Architecture))
	return nil
}

// Convert_v1alpha1_MachineImage_To_metal_MachineImage is an autogenerated conversion function.
func Convert_v1alpha1_MachineImage_To_metal_MachineImage(in *MachineImage, out *metal.MachineImage, s conversion.Scope) error {
	return autoConvert_v1alpha1_MachineImage_To_metal_MachineImage(in, out, s)
}

func autoConvert_metal_MachineImage_To_v1alpha1_MachineImage(in *metal.MachineImage, out *MachineImage, s conversion.Scope) error {
	out.Name = in.Name
	out.Version = in.Version
	out.Image = in.Image
	out.Architecture = (*string)(unsafe.Pointer(in.Architecture))
	return nil
}

// Convert_metal_MachineImage_To_v1alpha1_MachineImage is an autogenerated conversion function.
func Convert_metal_MachineImage_To_v1alpha1_MachineImage(in *metal.MachineImage, out *MachineImage, s conversion.Scope) error {
	return autoConvert_metal_MachineImage_To_v1alpha1_MachineImage(in, out, s)
}

func autoConvert_v1alpha1_MachineImageVersion_To_metal_MachineImageVersion(in *MachineImageVersion, out *metal.MachineImageVersion, s conversion.Scope) error {
	out.Version = in.Version
	out.Image = in.Image
	out.Architecture = (*string)(unsafe.Pointer(in.Architecture))
	return nil
}

// Convert_v1alpha1_MachineImageVersion_To_metal_MachineImageVersion is an autogenerated conversion function.
func Convert_v1alpha1_MachineImageVersion_To_metal_MachineImageVersion(in *MachineImageVersion, out *metal.MachineImageVersion, s conversion.Scope) error {
	return autoConvert_v1alpha1_MachineImageVersion_To_metal_MachineImageVersion(in, out, s)
}

func autoConvert_metal_MachineImageVersion_To_v1alpha1_MachineImageVersion(in *metal.MachineImageVersion, out *MachineImageVersion, s conversion.Scope) error {
	out.Version = in.Version
	out.Image = in.Image
	out.Architecture = (*string)(unsafe.Pointer(in.Architecture))
	return nil
}

// Convert_metal_MachineImageVersion_To_v1alpha1_MachineImageVersion is an autogenerated conversion function.
func Convert_metal_MachineImageVersion_To_v1alpha1_MachineImageVersion(in *metal.MachineImageVersion, out *MachineImageVersion, s conversion.Scope) error {
	return autoConvert_metal_MachineImageVersion_To_v1alpha1_MachineImageVersion(in, out, s)
}

func autoConvert_v1alpha1_MachineImages_To_metal_MachineImages(in *MachineImages, out *metal.MachineImages, s conversion.Scope) error {
	out.Name = in.Name
	out.Versions = *(*[]metal.MachineImageVersion)(unsafe.Pointer(&in.Versions))
	return nil
}

// Convert_v1alpha1_MachineImages_To_metal_MachineImages is an autogenerated conversion function.
func Convert_v1alpha1_MachineImages_To_metal_MachineImages(in *MachineImages, out *metal.MachineImages, s conversion.Scope) error {
	return autoConvert_v1alpha1_MachineImages_To_metal_MachineImages(in, out, s)
}

func autoConvert_metal_MachineImages_To_v1alpha1_MachineImages(in *metal.MachineImages, out *MachineImages, s conversion.Scope) error {
	out.Name = in.Name
	out.Versions = *(*[]MachineImageVersion)(unsafe.Pointer(&in.Versions))
	return nil
}

// Convert_metal_MachineImages_To_v1alpha1_MachineImages is an autogenerated conversion function.
func Convert_metal_MachineImages_To_v1alpha1_MachineImages(in *metal.MachineImages, out *MachineImages, s conversion.Scope) error {
	return autoConvert_metal_MachineImages_To_v1alpha1_MachineImages(in, out, s)
}

func autoConvert_v1alpha1_MachineType_To_metal_MachineType(in *MachineType, out *metal.MachineType, s conversion.Scope) error {
	out.Name = in.Name
	out.ServerLabels = *(*map[string]string)(unsafe.Pointer(&in.ServerLabels))
	return nil
}

// Convert_v1alpha1_MachineType_To_metal_MachineType is an autogenerated conversion function.
func Convert_v1alpha1_MachineType_To_metal_MachineType(in *MachineType, out *metal.MachineType, s conversion.Scope) error {
	return autoConvert_v1alpha1_MachineType_To_metal_MachineType(in, out, s)
}

func autoConvert_metal_MachineType_To_v1alpha1_MachineType(in *metal.MachineType, out *MachineType, s conversion.Scope) error {
	out.Name = in.Name
	out.ServerLabels = *(*map[string]string)(unsafe.Pointer(&in.ServerLabels))
	return nil
}

// Convert_metal_MachineType_To_v1alpha1_MachineType is an autogenerated conversion function.
func Convert_metal_MachineType_To_v1alpha1_MachineType(in *metal.MachineType, out *MachineType, s conversion.Scope) error {
	return autoConvert_metal_MachineType_To_v1alpha1_MachineType(in, out, s)
}

func autoConvert_v1alpha1_MetallbConfig_To_metal_MetallbConfig(in *MetallbConfig, out *metal.MetallbConfig, s conversion.Scope) error {
	out.IPAddressPool = *(*[]string)(unsafe.Pointer(&in.IPAddressPool))
	out.EnableSpeaker = in.EnableSpeaker
	out.EnableL2Advertisement = in.EnableL2Advertisement
	return nil
}

// Convert_v1alpha1_MetallbConfig_To_metal_MetallbConfig is an autogenerated conversion function.
func Convert_v1alpha1_MetallbConfig_To_metal_MetallbConfig(in *MetallbConfig, out *metal.MetallbConfig, s conversion.Scope) error {
	return autoConvert_v1alpha1_MetallbConfig_To_metal_MetallbConfig(in, out, s)
}

func autoConvert_metal_MetallbConfig_To_v1alpha1_MetallbConfig(in *metal.MetallbConfig, out *MetallbConfig, s conversion.Scope) error {
	out.IPAddressPool = *(*[]string)(unsafe.Pointer(&in.IPAddressPool))
	out.EnableSpeaker = in.EnableSpeaker
	out.EnableL2Advertisement = in.EnableL2Advertisement
	return nil
}

// Convert_metal_MetallbConfig_To_v1alpha1_MetallbConfig is an autogenerated conversion function.
func Convert_metal_MetallbConfig_To_v1alpha1_MetallbConfig(in *metal.MetallbConfig, out *MetallbConfig, s conversion.Scope) error {
	return autoConvert_metal_MetallbConfig_To_v1alpha1_MetallbConfig(in, out, s)
}

func autoConvert_v1alpha1_RegionConfig_To_metal_RegionConfig(in *RegionConfig, out *metal.RegionConfig, s conversion.Scope) error {
	out.Name = in.Name
	out.Server = in.Server
	out.CertificateAuthorityData = *(*[]byte)(unsafe.Pointer(&in.CertificateAuthorityData))
	return nil
}

// Convert_v1alpha1_RegionConfig_To_metal_RegionConfig is an autogenerated conversion function.
func Convert_v1alpha1_RegionConfig_To_metal_RegionConfig(in *RegionConfig, out *metal.RegionConfig, s conversion.Scope) error {
	return autoConvert_v1alpha1_RegionConfig_To_metal_RegionConfig(in, out, s)
}

func autoConvert_metal_RegionConfig_To_v1alpha1_RegionConfig(in *metal.RegionConfig, out *RegionConfig, s conversion.Scope) error {
	out.Name = in.Name
	out.Server = in.Server
	out.CertificateAuthorityData = *(*[]byte)(unsafe.Pointer(&in.CertificateAuthorityData))
	return nil
}

// Convert_metal_RegionConfig_To_v1alpha1_RegionConfig is an autogenerated conversion function.
func Convert_metal_RegionConfig_To_v1alpha1_RegionConfig(in *metal.RegionConfig, out *RegionConfig, s conversion.Scope) error {
	return autoConvert_metal_RegionConfig_To_v1alpha1_RegionConfig(in, out, s)
}

func autoConvert_v1alpha1_WorkerConfig_To_metal_WorkerConfig(in *WorkerConfig, out *metal.WorkerConfig, s conversion.Scope) error {
	out.ExtraIgnition = (*metal.IgnitionConfig)(unsafe.Pointer(in.ExtraIgnition))
	out.ExtraServerLabels = *(*map[string]string)(unsafe.Pointer(&in.ExtraServerLabels))
	return nil
}

// Convert_v1alpha1_WorkerConfig_To_metal_WorkerConfig is an autogenerated conversion function.
func Convert_v1alpha1_WorkerConfig_To_metal_WorkerConfig(in *WorkerConfig, out *metal.WorkerConfig, s conversion.Scope) error {
	return autoConvert_v1alpha1_WorkerConfig_To_metal_WorkerConfig(in, out, s)
}

func autoConvert_metal_WorkerConfig_To_v1alpha1_WorkerConfig(in *metal.WorkerConfig, out *WorkerConfig, s conversion.Scope) error {
	out.ExtraIgnition = (*IgnitionConfig)(unsafe.Pointer(in.ExtraIgnition))
	out.ExtraServerLabels = *(*map[string]string)(unsafe.Pointer(&in.ExtraServerLabels))
	return nil
}

// Convert_metal_WorkerConfig_To_v1alpha1_WorkerConfig is an autogenerated conversion function.
func Convert_metal_WorkerConfig_To_v1alpha1_WorkerConfig(in *metal.WorkerConfig, out *WorkerConfig, s conversion.Scope) error {
	return autoConvert_metal_WorkerConfig_To_v1alpha1_WorkerConfig(in, out, s)
}

func autoConvert_v1alpha1_WorkerStatus_To_metal_WorkerStatus(in *WorkerStatus, out *metal.WorkerStatus, s conversion.Scope) error {
	out.MachineImages = *(*[]metal.MachineImage)(unsafe.Pointer(&in.MachineImages))
	return nil
}

// Convert_v1alpha1_WorkerStatus_To_metal_WorkerStatus is an autogenerated conversion function.
func Convert_v1alpha1_WorkerStatus_To_metal_WorkerStatus(in *WorkerStatus, out *metal.WorkerStatus, s conversion.Scope) error {
	return autoConvert_v1alpha1_WorkerStatus_To_metal_WorkerStatus(in, out, s)
}

func autoConvert_metal_WorkerStatus_To_v1alpha1_WorkerStatus(in *metal.WorkerStatus, out *WorkerStatus, s conversion.Scope) error {
	out.MachineImages = *(*[]MachineImage)(unsafe.Pointer(&in.MachineImages))
	return nil
}

// Convert_metal_WorkerStatus_To_v1alpha1_WorkerStatus is an autogenerated conversion function.
func Convert_metal_WorkerStatus_To_v1alpha1_WorkerStatus(in *metal.WorkerStatus, out *WorkerStatus, s conversion.Scope) error {
	return autoConvert_metal_WorkerStatus_To_v1alpha1_WorkerStatus(in, out, s)
}
