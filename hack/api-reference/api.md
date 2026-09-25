<p>Packages:</p>
<ul>
<li>
<a href="#ironcore-metal.provider.extensions.gardener.cloud%2fv1alpha1">ironcore-metal.provider.extensions.gardener.cloud/v1alpha1</a>
</li>
</ul>

<h2 id="ironcore-metal.provider.extensions.gardener.cloud/v1alpha1">ironcore-metal.provider.extensions.gardener.cloud/v1alpha1</h2>
<p>

</p>
Resource Types:
<ul>
<li>
<a href="#workerconfig">WorkerConfig</a>
</li>
</ul>

<h3 id="bgpfilter">BGPFilter
</h3>


<p>
(<em>Appears on:</em><a href="#calicobgpconfig">CalicoBgpConfig</a>)
</p>

<p>
BGPFilter contains configuration for BGPFilter resource.
</p>

<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>

<tr>
<td>
<code>name</code></br>
<em>
string
</em>
</td>
<td>
<p>Name is the name of the BGPFilter resource.</p>
</td>
</tr>
<tr>
<td>
<code>exportV4</code></br>
<em>
<a href="#bgpfilterrule">BGPFilterRule</a> array
</em>
</td>
<td>
<em>(Optional)</em>
<p>The ordered set of IPv4 BGPFilter rules acting on exporting routes to a peer.</p>
</td>
</tr>
<tr>
<td>
<code>importV4</code></br>
<em>
<a href="#bgpfilterrule">BGPFilterRule</a> array
</em>
</td>
<td>
<em>(Optional)</em>
<p>The ordered set of IPv4 BGPFilter rules acting on importing routes from a peer.</p>
</td>
</tr>
<tr>
<td>
<code>exportV6</code></br>
<em>
<a href="#bgpfilterrule">BGPFilterRule</a> array
</em>
</td>
<td>
<em>(Optional)</em>
<p>The ordered set of IPv6 BGPFilter rules acting on exporting routes to a peer.</p>
</td>
</tr>
<tr>
<td>
<code>importV6</code></br>
<em>
<a href="#bgpfilterrule">BGPFilterRule</a> array
</em>
</td>
<td>
<em>(Optional)</em>
<p>The ordered set of IPv6 BGPFilter rules acting on importing routes from a peer.</p>
</td>
</tr>

</tbody>
</table>


<h3 id="bgpfilterrule">BGPFilterRule
</h3>
<p><em>Underlying type: <a href="#struct{cidr-string-"json:\"cidr\"";-matchoperator-string-"json:\"matchoperator\"";-action-string-"json:\"action\""}">struct{CIDR string "json:\"cidr\""; MatchOperator string "json:\"matchOperator\""; Action string "json:\"action\""}</a></em></p>


<p>
(<em>Appears on:</em><a href="#bgpfilter">BGPFilter</a>)
</p>

<p>
BGPFilterRule defines a BGP filter rule consisting a single CIDR block and a filter action for this CIDR.
</p>


<h3 id="bgppeer">BgpPeer
</h3>


<p>
(<em>Appears on:</em><a href="#calicobgpconfig">CalicoBgpConfig</a>)
</p>

<p>
BgpPeer contains configuration for BGPPeer resource.
</p>

<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>

<tr>
<td>
<code>peerIP</code></br>
<em>
string
</em>
</td>
<td>
<p>PeerIP contains IP address of BGP peer followed by an optional port number to peer with.</p>
</td>
</tr>
<tr>
<td>
<code>asNumber</code></br>
<em>
<a href="#asnumber">ASNumber</a>
</em>
</td>
<td>
<p>ASNumber contains the AS number of the BGP peer.</p>
</td>
</tr>
<tr>
<td>
<code>nodeSelector</code></br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>NodeSelector is a key-value pair to select nodes that should have this peering.</p>
</td>
</tr>
<tr>
<td>
<code>filters</code></br>
<em>
string array
</em>
</td>
<td>
<em>(Optional)</em>
<p>Filters contains the filters for the BGP peer.</p>
</td>
</tr>

</tbody>
</table>


<h3 id="calicobgpconfig">CalicoBgpConfig
</h3>


<p>
(<em>Appears on:</em><a href="#calicoconfig">CalicoConfig</a>)
</p>

<p>
CalicoBgpConfig contains BGP configuration settings for calico.
</p>

<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>

<tr>
<td>
<code>asNumber</code></br>
<em>
<a href="#asnumber">ASNumber</a>
</em>
</td>
<td>
<p>ASNumber is the default AS number used by a node.</p>
</td>
</tr>
<tr>
<td>
<code>nodeToNodeMeshEnabled</code></br>
<em>
boolean
</em>
</td>
<td>
<em>(Optional)</em>
<p>nodeToNodeMeshEnabled enables the node-to-node mesh.</p>
</td>
</tr>
<tr>
<td>
<code>serviceLoadBalancerIPs</code></br>
<em>
string array
</em>
</td>
<td>
<em>(Optional)</em>
<p>ServiceLoadBalancerIPs are the CIDR blocks for Kubernetes Service LoadBalancer IPs.</p>
</td>
</tr>
<tr>
<td>
<code>serviceExternalIPs</code></br>
<em>
string array
</em>
</td>
<td>
<em>(Optional)</em>
<p>ServiceExternalIPs are the CIDR blocks for Kubernetes Service External IPs.</p>
</td>
</tr>
<tr>
<td>
<code>serviceClusterIPs</code></br>
<em>
string array
</em>
</td>
<td>
<em>(Optional)</em>
<p>ServiceClusterIPs are the CIDR blocks from which service cluster IPs are allocated.</p>
</td>
</tr>
<tr>
<td>
<code>bgpPeer</code></br>
<em>
<a href="#bgppeer">BgpPeer</a> array
</em>
</td>
<td>
<em>(Optional)</em>
<p>BGPPeer contains configuration for BGPPeer resource.</p>
</td>
</tr>
<tr>
<td>
<code>bgpFilter</code></br>
<em>
<a href="#bgpfilter">BGPFilter</a> array
</em>
</td>
<td>
<em>(Optional)</em>
<p>BGPFilter contains configuration for BGPFilter resource.</p>
</td>
</tr>

</tbody>
</table>


<h3 id="calicoconfig">CalicoConfig
</h3>


<p>
(<em>Appears on:</em><a href="#loadbalancerconfig">LoadBalancerConfig</a>)
</p>

<p>
CalicoConfig contains configuration settings for Calico.
</p>

<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>

<tr>
<td>
<code>calicoBgpConfig</code></br>
<em>
<a href="#calicobgpconfig">CalicoBgpConfig</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>CalicoBgpConfig contains BGP configuration settings for calico.</p>
</td>
</tr>
<tr>
<td>
<code>IPPools</code></br>
<em>
<a href="#calicoippool">CalicoIPPool</a> array
</em>
</td>
<td>
<em>(Optional)</em>
<p>CalicoIPPools are the CIDR blocks for LoadBalancer IPs.</p>
</td>
</tr>

</tbody>
</table>


<h3 id="calicoippool">CalicoIPPool
</h3>


<p>
(<em>Appears on:</em><a href="#calicoconfig">CalicoConfig</a>)
</p>

<p>
CalicoIPPool contains configuration for a Calico IP pool.
</p>

<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>

<tr>
<td>
<code>allowedUses</code></br>
<em>
<a href="#calicoippoolalloweduse">CalicoIPPoolAllowedUse</a> array
</em>
</td>
<td>
<em>(Optional)</em>
<p>CalicoIPPoolAllowedUses controls what the IP pool will be used for.</p>
</td>
</tr>
<tr>
<td>
<code>assignmentMode</code></br>
<em>
<a href="#calicoippoolassignmentmode">CalicoIPPoolAssignmentMode</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>CalicoIPPoolAssignmentMode determines the mode how IP addresses should be assigned from this pool.</p>
</td>
</tr>
<tr>
<td>
<code>cidr</code></br>
<em>
string
</em>
</td>
<td>
<p>CIDR is the CIDR block for the IP pool.</p>
</td>
</tr>
<tr>
<td>
<code>Disabled</code></br>
<em>
boolean
</em>
</td>
<td>
<em>(Optional)</em>
<p>When disabled is true, Calico IPAM will not assign addresses from this pool.<br />Default is false.</p>
</td>
</tr>

</tbody>
</table>


<h3 id="calicoippoolalloweduse">CalicoIPPoolAllowedUse
</h3>
<p><em>Underlying type: string</em></p>


<p>
(<em>Appears on:</em><a href="#calicoippool">CalicoIPPool</a>)
</p>

<p>
CalicoIPPoolAllowedUse controls what the IP pool will be used for.
</p>


<h3 id="calicoippoolassignmentmode">CalicoIPPoolAssignmentMode
</h3>
<p><em>Underlying type: string</em></p>


<p>
(<em>Appears on:</em><a href="#calicoippool">CalicoIPPool</a>)
</p>

<p>
CalicoIPPoolAssignmentMode determines the mode how IP addresses should be assigned from this pool.
</p>


<h3 id="cloudcontrollermanagerconfig">CloudControllerManagerConfig
</h3>


<p>
(<em>Appears on:</em><a href="#controlplaneconfig">ControlPlaneConfig</a>)
</p>

<p>
CloudControllerManagerConfig contains configuration settings for the cloud-controller-manager.
</p>

<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>

<tr>
<td>
<code>featureGates</code></br>
<em>
object (keys:string, values:boolean)
</em>
</td>
<td>
<em>(Optional)</em>
<p>FeatureGates contains information about enabled feature gates.</p>
</td>
</tr>
<tr>
<td>
<code>networking</code></br>
<em>
<a href="#cloudcontrollernetworking">CloudControllerNetworking</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>Networking contains configuration settings for CCM networking.</p>
</td>
</tr>
<tr>
<td>
<code>podPrefixSize</code></br>
<em>
integer
</em>
</td>
<td>
<em>(Optional)</em>
<p>PodPrefixSize is the prefix size for pod CIDRs assigned to nodes.<br />When non-zero, passed as --pod-prefix-size to the CCM.</p>
</td>
</tr>

</tbody>
</table>


<h3 id="cloudcontrollernetworking">CloudControllerNetworking
</h3>


<p>
(<em>Appears on:</em><a href="#cloudcontrollermanagerconfig">CloudControllerManagerConfig</a>)
</p>

<p>
CloudControllerNetworking contains configuration settings for CCM networking.
</p>

<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>

<tr>
<td>
<code>configureNodeAddresses</code></br>
<em>
boolean
</em>
</td>
<td>
<em>(Optional)</em>
<p>ConfigureNodeAddresses enables the configuration of node addresses.</p>
</td>
</tr>
<tr>
<td>
<code>ipamKind</code></br>
<em>
<a href="#ipamkind">IPAMKind</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>IPAMKind enables the IPAM integration.</p>
</td>
</tr>

</tbody>
</table>


<h3 id="cloudprofileconfig">CloudProfileConfig
</h3>


<p>
CloudProfileConfig contains provider-specific configuration that is embedded into Gardener's `CloudProfile`
resource.
</p>

<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>

<tr>
<td>
<code>machineImages</code></br>
<em>
<a href="#machineimages">MachineImages</a> array
</em>
</td>
<td>
<p>MachineImages is the list of machine images that are understood by the controller. It maps<br />logical names and versions to provider-specific identifiers.</p>
</td>
</tr>
<tr>
<td>
<code>regionConfigs</code></br>
<em>
<a href="#regionconfig">RegionConfig</a> array
</em>
</td>
<td>
<p>RegionConfigs is the list of supported regions.</p>
</td>
</tr>
<tr>
<td>
<code>machineTypes</code></br>
<em>
<a href="#machinetype">MachineType</a> array
</em>
</td>
<td>
<p></p>
</td>
</tr>

</tbody>
</table>


<h3 id="controlplaneconfig">ControlPlaneConfig
</h3>


<p>
ControlPlaneConfig contains configuration settings for the control plane.
</p>

<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>

<tr>
<td>
<code>cloudControllerManager</code></br>
<em>
<a href="#cloudcontrollermanagerconfig">CloudControllerManagerConfig</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>CloudControllerManager contains configuration settings for the cloud-controller-manager.</p>
</td>
</tr>
<tr>
<td>
<code>loadBalancerConfig</code></br>
<em>
<a href="#loadbalancerconfig">LoadBalancerConfig</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>LoadBalancerConfig contains configuration settings for the shoot loadbalancing.</p>
</td>
</tr>
<tr>
<td>
<code>nodeNamePolicy</code></br>
<em>
<a href="#nodenamepolicy">NodeNamePolicy</a>
</em>
</td>
<td>
<p>NodeNamePolicy is a policy for generating hostnames for the worker nodes.</p>
</td>
</tr>

</tbody>
</table>


<h3 id="ipamconfig">IPAMConfig
</h3>


<p>
(<em>Appears on:</em><a href="#workerconfig">WorkerConfig</a>)
</p>

<p>
IPAMConfig is a reference to an IPAM resource.
</p>

<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>

<tr>
<td>
<code>metadataKey</code></br>
<em>
string
</em>
</td>
<td>
<p>MetadataKey is the name of metadata key for the network.</p>
</td>
</tr>
<tr>
<td>
<code>ipamRef</code></br>
<em>
<a href="#ipamobjectreference">IPAMObjectReference</a>
</em>
</td>
<td>
<p>IPAMRef is a reference to the IPAM object, which will be used for IP allocation.</p>
</td>
</tr>

</tbody>
</table>


<h3 id="ipamkind">IPAMKind
</h3>


<p>
(<em>Appears on:</em><a href="#cloudcontrollernetworking">CloudControllerNetworking</a>)
</p>

<p>
IPAMKind specifiers the IPAM objects in-use.
</p>

<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>

<tr>
<td>
<code>apiGroup</code></br>
<em>
string
</em>
</td>
<td>
<p>APIGroup is the resource group.</p>
</td>
</tr>
<tr>
<td>
<code>kind</code></br>
<em>
string
</em>
</td>
<td>
<p>Kind is the resource type.</p>
</td>
</tr>

</tbody>
</table>


<h3 id="ipamobjectreference">IPAMObjectReference
</h3>


<p>
(<em>Appears on:</em><a href="#ipamconfig">IPAMConfig</a>)
</p>

<p>
IPAMObjectReference is a reference to the IPAM object, which will be used for IP allocation.
</p>

<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>

<tr>
<td>
<code>name</code></br>
<em>
string
</em>
</td>
<td>
<p>Name is the name of resource being referenced.</p>
</td>
</tr>
<tr>
<td>
<code>apiGroup</code></br>
<em>
string
</em>
</td>
<td>
<p>APIGroup is the group for the resource being referenced.</p>
</td>
</tr>
<tr>
<td>
<code>kind</code></br>
<em>
string
</em>
</td>
<td>
<p>Kind is the type of resource being referenced.</p>
</td>
</tr>

</tbody>
</table>


<h3 id="ignitionconfig">IgnitionConfig
</h3>


<p>
(<em>Appears on:</em><a href="#workerconfig">WorkerConfig</a>)
</p>

<p>
IgnitionConfig contains ignition settings.
</p>

<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>

<tr>
<td>
<code>raw</code></br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Raw contains an inline ignition config, which is merged with the config from the os extension.</p>
</td>
</tr>
<tr>
<td>
<code>secretRef</code></br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>SecretRef is a reference to a resource in the shoot spec referencing a secret containing the ignition config.</p>
</td>
</tr>
<tr>
<td>
<code>override</code></br>
<em>
boolean
</em>
</td>
<td>
<em>(Optional)</em>
<p>Override configures, if ignition keys set by the os-extension are overridden<br />by extra ignition.</p>
</td>
</tr>

</tbody>
</table>


<h3 id="infrastructureconfig">InfrastructureConfig
</h3>


<p>
InfrastructureConfig infrastructure configuration resource
</p>

<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>

<tr>
<td>
<code>networks</code></br>
<em>
<a href="#networks">Networks</a> array
</em>
</td>
<td>
<em>(Optional)</em>
<p>Networks is the metal specific network configuration.</p>
</td>
</tr>

</tbody>
</table>


<h3 id="infrastructurestatus">InfrastructureStatus
</h3>


<p>
InfrastructureStatus contains information about created infrastructure resources.
</p>


<h3 id="loadbalancerconfig">LoadBalancerConfig
</h3>


<p>
(<em>Appears on:</em><a href="#controlplaneconfig">ControlPlaneConfig</a>)
</p>

<p>
LoadBalancerConfig contains configuration settings for the shoot loadbalancing.
</p>

<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>

<tr>
<td>
<code>metallbConfig</code></br>
<em>
<a href="#metallbconfig">MetallbConfig</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>MetallbConfig contains configuration settings for metallb.</p>
</td>
</tr>
<tr>
<td>
<code>calicoConfig</code></br>
<em>
<a href="#calicoconfig">CalicoConfig</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>CalicoConfig contains configuration settings for calico.</p>
</td>
</tr>
<tr>
<td>
<code>metalLoadBalancerConfig</code></br>
<em>
<a href="#metalloadbalancerconfig">MetalLoadBalancerConfig</a>
</em>
</td>
<td>
<p>MetalLoadBalancerConfig contains configuration settings for the metal load balancer.</p>
</td>
</tr>

</tbody>
</table>


<h3 id="machineimage">MachineImage
</h3>


<p>
(<em>Appears on:</em><a href="#workerstatus">WorkerStatus</a>)
</p>

<p>
MachineImage is a mapping from logical names and versions to metal-specific identifiers.
</p>

<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>

<tr>
<td>
<code>name</code></br>
<em>
string
</em>
</td>
<td>
<p>Name is the logical name of the machine image.</p>
</td>
</tr>
<tr>
<td>
<code>version</code></br>
<em>
string
</em>
</td>
<td>
<p>Version is the logical version of the machine image.</p>
</td>
</tr>
<tr>
<td>
<code>image</code></br>
<em>
string
</em>
</td>
<td>
<p>Image is the path to the image.</p>
</td>
</tr>
<tr>
<td>
<code>architecture</code></br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Architecture is the CPU architecture of the machine image.</p>
</td>
</tr>
<tr>
<td>
<code>capabilities</code></br>
<em>
<a href="#capabilities">Capabilities</a>
</em>
</td>
<td>
<p>Capabilities of the machine image.</p>
</td>
</tr>

</tbody>
</table>


<h3 id="machineimageflavor">MachineImageFlavor
</h3>


<p>
(<em>Appears on:</em><a href="#machineimageversion">MachineImageVersion</a>)
</p>

<p>
MachineImageFlavor is a flavor of the machine image version that supports a specific set of capabilities.
</p>

<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>

<tr>
<td>
<code>capabilities</code></br>
<em>
<a href="#capabilities">Capabilities</a>
</em>
</td>
<td>
<p>Capabilities is the set of capabilities that are supported by the image in this flavor.</p>
</td>
</tr>
<tr>
<td>
<code>image</code></br>
<em>
string
</em>
</td>
<td>
<p>Image is the path to the image.</p>
</td>
</tr>

</tbody>
</table>


<h3 id="machineimageversion">MachineImageVersion
</h3>


<p>
(<em>Appears on:</em><a href="#machineimages">MachineImages</a>)
</p>

<p>
MachineImageVersion contains a version and a provider-specific identifier.
</p>

<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>

<tr>
<td>
<code>version</code></br>
<em>
string
</em>
</td>
<td>
<p>Version is the version of the image.</p>
</td>
</tr>
<tr>
<td>
<code>image</code></br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Image is the path to the image.</p>
</td>
</tr>
<tr>
<td>
<code>architecture</code></br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>Architecture is the CPU architecture of the machine image.</p>
</td>
</tr>
<tr>
<td>
<code>capabilityFlavors</code></br>
<em>
<a href="#machineimageflavor">MachineImageFlavor</a> array
</em>
</td>
<td>
<p>CapabilityFlavors is a collection of all images for that version with capabilities.</p>
</td>
</tr>

</tbody>
</table>


<h3 id="machineimages">MachineImages
</h3>


<p>
(<em>Appears on:</em><a href="#cloudprofileconfig">CloudProfileConfig</a>)
</p>

<p>
MachineImages is a mapping from logical names and versions to provider-specific identifiers.
</p>

<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>

<tr>
<td>
<code>name</code></br>
<em>
string
</em>
</td>
<td>
<p>Name is the logical name of the machine image.</p>
</td>
</tr>
<tr>
<td>
<code>versions</code></br>
<em>
<a href="#machineimageversion">MachineImageVersion</a> array
</em>
</td>
<td>
<p>Versions contains versions and a provider-specific identifier.</p>
</td>
</tr>

</tbody>
</table>


<h3 id="machinetype">MachineType
</h3>


<p>
(<em>Appears on:</em><a href="#cloudprofileconfig">CloudProfileConfig</a>)
</p>

<p>

</p>

<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>

<tr>
<td>
<code>name</code></br>
<em>
string
</em>
</td>
<td>
<p></p>
</td>
</tr>
<tr>
<td>
<code>serverLabels</code></br>
<em>
object (keys:string, values:string)
</em>
</td>
<td>
<p></p>
</td>
</tr>

</tbody>
</table>


<h3 id="metalloadbalancerconfig">MetalLoadBalancerConfig
</h3>


<p>
(<em>Appears on:</em><a href="#loadbalancerconfig">LoadBalancerConfig</a>)
</p>

<p>
MetalLoadBalancerConfig contains configuration settings for the metal load balancer.
</p>

<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>

<tr>
<td>
<code>nodeCIDRMask</code></br>
<em>
integer
</em>
</td>
<td>
<p>NodeCIDRMask is the mask for the node CIDR.</p>
</td>
</tr>
<tr>
<td>
<code>allocateNodeCIDRs</code></br>
<em>
boolean
</em>
</td>
<td>
<p>AllocateNodeCIDRs enables the allocation of node CIDRs.</p>
</td>
</tr>
<tr>
<td>
<code>vni</code></br>
<em>
integer
</em>
</td>
<td>
<p>VNI is the VNI used for IP announcements.</p>
</td>
</tr>
<tr>
<td>
<code>metalBondServer</code></br>
<em>
string
</em>
</td>
<td>
<p>MetalBondServer is the URL of the metal bond server.</p>
</td>
</tr>

</tbody>
</table>


<h3 id="metallbconfig">MetallbConfig
</h3>


<p>
(<em>Appears on:</em><a href="#loadbalancerconfig">LoadBalancerConfig</a>)
</p>

<p>
MetallbConfig contains configuration settings for metallb.
</p>

<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>

<tr>
<td>
<code>ipAddressPool</code></br>
<em>
string array
</em>
</td>
<td>
<em>(Optional)</em>
<p>IPAddressPool contains IP address pools for metallb.</p>
</td>
</tr>
<tr>
<td>
<code>enableSpeaker</code></br>
<em>
boolean
</em>
</td>
<td>
<em>(Optional)</em>
<p>EnableSpeaker enables the metallb speaker.</p>
</td>
</tr>
<tr>
<td>
<code>enableL2Advertisement</code></br>
<em>
boolean
</em>
</td>
<td>
<em>(Optional)</em>
<p>EnableL2Advertisement enables L2 advertisement.</p>
</td>
</tr>

</tbody>
</table>


<h3 id="networks">Networks
</h3>


<p>
(<em>Appears on:</em><a href="#infrastructureconfig">InfrastructureConfig</a>)
</p>

<p>
Networks holds information about the Kubernetes and infrastructure networks.
</p>

<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>

<tr>
<td>
<code>name</code></br>
<em>
string
</em>
</td>
<td>
<p>Name is the name for this CIDR.</p>
</td>
</tr>
<tr>
<td>
<code>cidr</code></br>
<em>
string
</em>
</td>
<td>
<p>CIDR is the workers subnet range to create.</p>
</td>
</tr>
<tr>
<td>
<code>id</code></br>
<em>
string
</em>
</td>
<td>
<em>(Optional)</em>
<p>ID is the ID for the workers' subnet.</p>
</td>
</tr>

</tbody>
</table>


<h3 id="nodenamepolicy">NodeNamePolicy
</h3>
<p><em>Underlying type: string</em></p>


<p>
(<em>Appears on:</em><a href="#controlplaneconfig">ControlPlaneConfig</a>)
</p>

<p>
NodeNamePolicy is a policy for generating hostnames for the worker nodes.
</p>


<h3 id="regionconfig">RegionConfig
</h3>


<p>
(<em>Appears on:</em><a href="#cloudprofileconfig">CloudProfileConfig</a>)
</p>

<p>
RegionConfig is the definition of a region.
</p>

<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>

<tr>
<td>
<code>name</code></br>
<em>
string
</em>
</td>
<td>
<p>Name is the name of a region.</p>
</td>
</tr>
<tr>
<td>
<code>server</code></br>
<em>
string
</em>
</td>
<td>
<p>Server is the server endpoint of this region.</p>
</td>
</tr>
<tr>
<td>
<code>certificateAuthorityData</code></br>
<em>
integer array
</em>
</td>
<td>
<p>CertificateAuthorityData is the CA data of the region server.</p>
</td>
</tr>

</tbody>
</table>


<h3 id="workerconfig">WorkerConfig
</h3>


<p>
WorkerConfig contains configuration settings for the worker nodes.
</p>

<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>

<tr>
<td>
<code>extraIgnition</code></br>
<em>
<a href="#ignitionconfig">IgnitionConfig</a>
</em>
</td>
<td>
<em>(Optional)</em>
<p>ExtraIgnition contains additional Ignition for Worker nodes.</p>
</td>
</tr>
<tr>
<td>
<code>extraServerLabels</code></br>
<em>
object (keys:string, values:string)
</em>
</td>
<td>
<em>(Optional)</em>
<p>ExtraServerLabels is a map of additional labels that are applied to the ServerClaim for Server selection.</p>
</td>
</tr>
<tr>
<td>
<code>ipamConfig</code></br>
<em>
<a href="#ipamconfig">IPAMConfig</a> array
</em>
</td>
<td>
<em>(Optional)</em>
<p>IPAMConfig is a list of references to Network resources that should be used to assign IP addresses to the worker nodes.</p>
</td>
</tr>
<tr>
<td>
<code>metadata</code></br>
<em>
object (keys:string, values:string)
</em>
</td>
<td>
Refer to the Kubernetes API documentation for the fields of the <code>metadata</code> field.
</td>
</tr>

</tbody>
</table>


<h3 id="workerstatus">WorkerStatus
</h3>


<p>
WorkerStatus contains information about created worker resources.
</p>

<table>
<thead>
<tr>
<th>Field</th>
<th>Description</th>
</tr>
</thead>
<tbody>

<tr>
<td>
<code>machineImages</code></br>
<em>
<a href="#machineimage">MachineImage</a> array
</em>
</td>
<td>
<em>(Optional)</em>
<p>MachineImages is a list of machine images that have been used in this worker. Usually, the extension controller<br />gets the mapping from name/version to the provider-specific machine image data in its componentconfig. However, if<br />a version that is still in use gets removed from this componentconfig it cannot reconcile anymore existing `Worker`<br />resources that are still using this version. Hence, it stores the used versions in the provider status to ensure<br />reconciliation is possible.</p>
</td>
</tr>

</tbody>
</table>


