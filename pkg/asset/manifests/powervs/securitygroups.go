package powervs

import (
	"fmt"

	"k8s.io/utils/ptr"
	capibmcloud "sigs.k8s.io/cluster-api-provider-ibmcloud/api/v1beta2"

	"github.com/openshift/installer/pkg/types"
)

const (
	clusterWideSGNameSuffix  = "sg-cluster-wide"
	openshiftNetSGNameSuffix = "sg-openshift-net"
	kubeAPILBSGNameSuffix    = "sg-kube-api-lb"
	controlPlaneSGNameSuffix = "sg-control-plane"
	cpInternalSGNameSuffix   = "sg-cp-internal"
)

func buildControlPlaneSecurityGroup(infraID string) capibmcloud.VPCSecurityGroup {
	controlPlaneSGNamePtr := ptr.To(fmt.Sprintf("%s-%s", infraID, controlPlaneSGNameSuffix))
	clusterWideSGNamePtr := ptr.To(fmt.Sprintf("%s-%s", infraID, clusterWideSGNameSuffix))
	kubeAPILBSGNamePtr := ptr.To(fmt.Sprintf("%s-%s", infraID, kubeAPILBSGNameSuffix))

	return capibmcloud.VPCSecurityGroup{
		Name: controlPlaneSGNamePtr,
		Rules: []*capibmcloud.VPCSecurityGroupRule{
			{
				// Kubernetes API - inbound via cluster
				Action:    capibmcloud.VPCSecurityGroupRuleActionAllow,
				Direction: capibmcloud.VPCSecurityGroupRuleDirectionInbound,
				Source: &capibmcloud.VPCSecurityGroupRulePrototype{
					PortRange: &capibmcloud.VPCSecurityGroupPortRange{
						MaximumPort: 6443,
						MinimumPort: 6443,
					},
					Protocol: capibmcloud.VPCSecurityGroupRuleProtocolTCP,
					Remotes: []capibmcloud.VPCSecurityGroupRuleRemote{
						{
							RemoteType:        capibmcloud.VPCSecurityGroupRuleRemoteTypeSG,
							SecurityGroupName: clusterWideSGNamePtr,
						},
					},
				},
			},
			{
				// Kubernetes API - inbound via LB
				Action:    capibmcloud.VPCSecurityGroupRuleActionAllow,
				Direction: capibmcloud.VPCSecurityGroupRuleDirectionInbound,
				Source: &capibmcloud.VPCSecurityGroupRulePrototype{
					PortRange: &capibmcloud.VPCSecurityGroupPortRange{
						MaximumPort: 6443,
						MinimumPort: 6443,
					},
					Protocol: capibmcloud.VPCSecurityGroupRuleProtocolTCP,
					Remotes: []capibmcloud.VPCSecurityGroupRuleRemote{
						{
							RemoteType:        capibmcloud.VPCSecurityGroupRuleRemoteTypeSG,
							SecurityGroupName: kubeAPILBSGNamePtr,
						},
					},
				},
			},
			{
				// Machine Config Server - inbound via LB
				Action:    capibmcloud.VPCSecurityGroupRuleActionAllow,
				Direction: capibmcloud.VPCSecurityGroupRuleDirectionInbound,
				Source: &capibmcloud.VPCSecurityGroupRulePrototype{
					PortRange: &capibmcloud.VPCSecurityGroupPortRange{
						MaximumPort: 22623,
						MinimumPort: 22623,
					},
					Protocol: capibmcloud.VPCSecurityGroupRuleProtocolTCP,
					Remotes: []capibmcloud.VPCSecurityGroupRuleRemote{
						{
							RemoteType:        capibmcloud.VPCSecurityGroupRuleRemoteTypeSG,
							SecurityGroupName: kubeAPILBSGNamePtr,
						},
					},
				},
			},
			{
				// Kubernetes default ports
				Action:    capibmcloud.VPCSecurityGroupRuleActionAllow,
				Direction: capibmcloud.VPCSecurityGroupRuleDirectionInbound,
				Source: &capibmcloud.VPCSecurityGroupRulePrototype{
					PortRange: &capibmcloud.VPCSecurityGroupPortRange{
						MaximumPort: 10259,
						MinimumPort: 10257,
					},
					Protocol: capibmcloud.VPCSecurityGroupRuleProtocolTCP,
					Remotes: []capibmcloud.VPCSecurityGroupRuleRemote{
						{
							RemoteType:        capibmcloud.VPCSecurityGroupRuleRemoteTypeSG,
							SecurityGroupName: clusterWideSGNamePtr,
						},
					},
				},
			},
		},
	}
}

func buildCPInternalSecurityGroup(infraID string) capibmcloud.VPCSecurityGroup {
	cpInternalSGNamePtr := ptr.To(fmt.Sprintf("%s-%s", infraID, cpInternalSGNameSuffix))

	return capibmcloud.VPCSecurityGroup{
		Name: cpInternalSGNamePtr,
		Rules: []*capibmcloud.VPCSecurityGroupRule{
			{
				// etcd internal traffic
				Action:    capibmcloud.VPCSecurityGroupRuleActionAllow,
				Direction: capibmcloud.VPCSecurityGroupRuleDirectionInbound,
				Source: &capibmcloud.VPCSecurityGroupRulePrototype{
					PortRange: &capibmcloud.VPCSecurityGroupPortRange{
						MaximumPort: 2380,
						MinimumPort: 2379,
					},
					Protocol: capibmcloud.VPCSecurityGroupRuleProtocolTCP,
					Remotes: []capibmcloud.VPCSecurityGroupRuleRemote{
						{
							RemoteType:        capibmcloud.VPCSecurityGroupRuleRemoteTypeSG,
							SecurityGroupName: cpInternalSGNamePtr,
						},
					},
				},
			},
		},
	}
}

func getVPCSecurityGroups(infraID string, publishStrategy types.PublishingStrategy) []capibmcloud.VPCSecurityGroup {
	// IBM Cloud currently relies on 5 SecurityGroups to manage traffic and 1 SecurityGroup for bootstrapping.
	securityGroups := make([]capibmcloud.VPCSecurityGroup, 0, 6)
	// Generate the Cluster's primary SG's.
	// securityGroups = append(securityGroups, buildClusterWideSecurityGroup(infraID, allSubnets))
	// securityGroups = append(securityGroups, buildOpenshiftNetSecurityGroup(infraID, allSubnets))
	// securityGroups = append(securityGroups, buildKubeAPILBSecurityGroup(infraID))
	// securityGroups = append(securityGroups, buildControlPlaneSecurityGroup(infraID))
	securityGroups = append(securityGroups, buildCPInternalSecurityGroup(infraID))

	// Generate the bootstrap SG.
	// securityGroups = append(securityGroups, buildBootstrapSecurityGroup(infraID, allSubnets, publishStrategy))

	return securityGroups
}
