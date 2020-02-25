package cloud_provider

import (
	"context"
	"strings"
	"testing"

	"icode.baidu.com/baidu/jpaas-caas/bce-sdk-go/bce"
	"icode.baidu.com/baidu/jpaas-caas/bce-sdk-go/vpc"
	cce "icode.baidu.com/baidu/jpaas-caas/cloud-provider-baiducloud/pkg/temp-cce"
	"k8s.io/apimachinery/pkg/types"
	cloudprovider "k8s.io/cloud-provider"
)

var routeruletableID string
var routeruleIDs []string

// beforeTest create vpc subnet cluster and routerule.
func beforeTestRoute() (cloud *Baiducloud, nodesResq *cce.ListClusterNodesResponse, err error) {
	cloud = NewFakeCloud("")
	ctx := context.Background()
	// Create VPC
	vpcID, err := cloud.clientSet.VPCClient.CreateVPC(ctx, &vpc.CreateVPCArgs{
		Name: "",
		CIDR: "10.0.0.0/8",
	}, nil)
	if err != nil {
		return nil, nil, err
	}
	temp := strings.Split(vpcID, "/")
	vpcID = temp[0]
	routeruletableID = temp[1]
	//  Create Subnet for test
	subnetID, err := cloud.clientSet.VPCClient.CreateSubnet(ctx, &vpc.CreateSubnetArgs{
		VPCID:      vpcID,
		SubnetType: vpc.SubnetTypeBCC,
		CIDR:       "10.0.0.0/16",
	}, nil)
	if err != nil {
		return nil, nil, err
	}
	// Create CCE Cluster
	resq, err := cloud.clientSet.CCEClient.CreateCluster(ctx, &cce.CreateClusterArgs{
		ClusterName: "",
		VPCID:       vpcID,
		SubnetID:    subnetID,
		VMCount:     2,
	})
	if err != nil || resq == nil {
		return nil, nil, err
	}
	cloud.CloudConfig.ClusterID = resq.ClusterID
	nodesResq, err = cloud.clientSet.CCEClient.ListClusterNodes(ctx, resq.ClusterID, nil)
	if err != nil || nodesResq == nil {
		return nil, nil, err
	}
	if len(nodesResq.Nodes) == 0 {
		return nil, nil, err
	}
	// create routerule
	args := []vpc.CreateRouteRuleArgs{
		{
			RouteTableID:       routeruletableID,
			SourceAddress:      "0.0.0.0/0",
			DestinationAddress: "test",
			NexthopID:          "test",
			NexthopType:        "test",
			Description:        "test",
		},
		{
			RouteTableID:       routeruletableID,
			SourceAddress:      "0.0.0.0/0",
			DestinationAddress: "test2",
			NexthopID:          "test2",
			NexthopType:        "test2",
			Description:        "test2",
		},
	}
	for _, arg := range args {
		rrid, err := cloud.clientSet.VPCClient.CreateRouteRule(ctx, &arg, &bce.SignOption{
			CustomSignFunc: CCEServiceSign,
		})
		if err != nil {
			return nil, nil, err
		}
		routeruleIDs = append(routeruleIDs, rrid)
	}
	return cloud, nodesResq, nil
}

func TestGetVpcID(t *testing.T) {
	ctx := context.Background()
	cloud := NewFakeCloud("")
	cases := []struct {
		name      string
		clusterID string
		cidr      string
	}{
		{
			name: "CCE Cluster Nodes subnetType is BCC",
			clusterID: func() string {
				// Create VPC
				vpcID, err := cloud.clientSet.VPCClient.CreateVPC(ctx, &vpc.CreateVPCArgs{
					Name: "",
					CIDR: "10.0.0.0/8",
				}, nil)
				if err != nil {
					t.Errorf("CreateVPC failed: %v", err)
					return ""
				}
				//  Create Subnet for test
				subnetID, err := cloud.clientSet.VPCClient.CreateSubnet(ctx, &vpc.CreateSubnetArgs{
					VPCID:      vpcID,
					SubnetType: vpc.SubnetTypeBCC,
					CIDR:       "10.0.0.0/16",
				}, nil)
				if err != nil {
					t.Errorf("CreateSubnet failed: %v", err)
					return ""
				}
				// Create CCE Cluster
				resq, err := cloud.clientSet.CCEClient.CreateCluster(ctx, &cce.CreateClusterArgs{
					ClusterName: "",
					VPCID:       vpcID,
					SubnetID:    subnetID,
					VMCount:     2,
				})
				if err != nil || resq == nil {
					t.Errorf("CreateCluster failed: %v", err)
					return ""
				}
				return resq.ClusterID
			}(),
			cidr: "10.0.0.0/16",
		},
		{
			name: "CCE Cluster Nodes not exist",
			clusterID: func() string {
				// Create VPC
				vpcID, err := cloud.clientSet.VPCClient.CreateVPC(ctx, &vpc.CreateVPCArgs{
					Name: "",
					CIDR: "10.0.0.0/8",
				}, nil)
				if err != nil {
					t.Errorf("CreateVPC failed: %v", err)
					return ""
				}
				//  Create Subnet for test
				subnetID, err := cloud.clientSet.VPCClient.CreateSubnet(ctx, &vpc.CreateSubnetArgs{
					VPCID:      vpcID,
					SubnetType: vpc.SubnetTypeBCC,
					CIDR:       "10.0.0.0/16",
				}, nil)
				if err != nil {
					t.Errorf("CreateSubnet failed: %v", err)
					return ""
				}
				// Create CCE Cluster
				resq, err := cloud.clientSet.CCEClient.CreateCluster(ctx, &cce.CreateClusterArgs{
					ClusterName: "",
					VPCID:       vpcID,
					SubnetID:    subnetID,
					VMCount:     0,
				})
				if err != nil || resq == nil {
					t.Errorf("CreateCluster failed: %v", err)
					return ""
				}
				return resq.ClusterID
			}(),
		},
		{
			name: `CCE Cluster Nodes subnetType not BCC, but contains "系统预定义子网"`,
			clusterID: func() string {
				// Create VPC
				vpcID, err := cloud.clientSet.VPCClient.CreateVPC(ctx, &vpc.CreateVPCArgs{
					Name: "",
					CIDR: "10.0.0.0/8",
				}, nil)
				if err != nil {
					t.Errorf("CreateVPC failed: %v", err)
					return ""
				}
				//  Create Subnet for test
				subnetID, err := cloud.clientSet.VPCClient.CreateSubnet(ctx, &vpc.CreateSubnetArgs{
					VPCID:      vpcID,
					SubnetType: vpc.SubnetTypeBCCNAT,
					CIDR:       "10.0.0.0/16",
				}, nil)
				if err != nil {
					t.Errorf("CreateSubnet failed: %v", err)
					return ""
				}
				if _, err := cloud.clientSet.VPCClient.CreateSubnet(ctx, &vpc.CreateSubnetArgs{
					Name:       "系统预定义子网",
					VPCID:      vpcID,
					SubnetType: vpc.SubnetTypeBCC,
					CIDR:       "10.10.0.0/16",
				}, nil); err != nil {
					t.Errorf("CreateSubnet failed: %v", err)
					return ""
				}
				// Create CCE Cluster
				resq, err := cloud.clientSet.CCEClient.CreateCluster(ctx, &cce.CreateClusterArgs{
					ClusterName: "",
					VPCID:       vpcID,
					SubnetID:    subnetID,
					VMCount:     2,
				})
				if err != nil || resq == nil {
					t.Errorf("CreateCluster failed: %v", err)
					return ""
				}
				return resq.ClusterID
			}(),
			cidr: "10.10.0.0/16",
		},
		{
			name: `CCE Cluster Nodes subnetType not BCC, but contains "CCE-Reserver"`,
			clusterID: func() string {
				// Create VPC
				vpcID, err := cloud.clientSet.VPCClient.CreateVPC(ctx, &vpc.CreateVPCArgs{
					Name: "",
					CIDR: "10.0.0.0/8",
				}, nil)
				if err != nil {
					t.Errorf("CreateVPC failed: %v", err)
					return ""
				}
				//  Create Subnet for test
				subnetID, err := cloud.clientSet.VPCClient.CreateSubnet(ctx, &vpc.CreateSubnetArgs{
					VPCID:      vpcID,
					SubnetType: vpc.SubnetTypeBCCNAT,
					CIDR:       "10.0.0.0/16",
				}, nil)
				if err != nil {
					t.Errorf("CreateSubnet failed: %v", err)
					return ""
				}
				if _, err := cloud.clientSet.VPCClient.CreateSubnet(ctx, &vpc.CreateSubnetArgs{
					Name:       "CCE-Reserve",
					VPCID:      vpcID,
					SubnetType: vpc.SubnetTypeBCC,
					CIDR:       "10.30.0.0/16",
				}, nil); err != nil {
					t.Errorf("CreateSubnet failed: %v", err)
					return ""
				}
				// Create CCE Cluster
				resq, err := cloud.clientSet.CCEClient.CreateCluster(ctx, &cce.CreateClusterArgs{
					ClusterName: "",
					VPCID:       vpcID,
					SubnetID:    subnetID,
					VMCount:     2,
				})
				if err != nil || resq == nil {
					t.Errorf("CreateCluster failed: %v", err)
					return ""
				}
				return resq.ClusterID
			}(),
			cidr: "10.30.0.0/16",
		},
		{
			name: `CCE Cluster Nodes VPC subnetType not BCC, and not contains "系统预定义子网" or "CCE-Reserve"`,
			clusterID: func() string {
				// Create VPC
				vpcID, err := cloud.clientSet.VPCClient.CreateVPC(ctx, &vpc.CreateVPCArgs{
					Name: "",
					CIDR: "10.0.0.0/8",
				}, nil)
				if err != nil {
					t.Errorf("CreateVPC failed: %v", err)
					return ""
				}
				//  Create Subnet for test
				subnetID, err := cloud.clientSet.VPCClient.CreateSubnet(ctx, &vpc.CreateSubnetArgs{
					VPCID:      vpcID,
					SubnetType: vpc.SubnetTypeBCCNAT,
					CIDR:       "10.1.5.0/24",
				}, nil)
				if err != nil {
					t.Errorf("CreateSubnet failed: %v", err)
					return ""
				}
				// Create CCE Cluster
				resq, err := cloud.clientSet.CCEClient.CreateCluster(ctx, &cce.CreateClusterArgs{
					ClusterName: "",
					VPCID:       vpcID,
					SubnetID:    subnetID,
					VMCount:     2,
				})
				if err != nil || resq == nil {
					t.Errorf("CreateCluster failed: %v", err)
					return ""
				}
				return resq.ClusterID
			}(),
			cidr: "10.1.6.0/24",
		},
		{
			name: `CCE Cluster Nodes VPC subnetType not BCC, and not contains "系统预定义子网" or "CCE-Reserve"`,
			clusterID: func() string {
				// Create VPC
				vpcID, err := cloud.clientSet.VPCClient.CreateVPC(ctx, &vpc.CreateVPCArgs{
					Name: "",
					CIDR: "10.0.0.0/8",
				}, nil)
				if err != nil {
					t.Errorf("CreateVPC failed: %v", err)
					return ""
				}
				//  Create Subnet for test
				subnetID, err := cloud.clientSet.VPCClient.CreateSubnet(ctx, &vpc.CreateSubnetArgs{
					VPCID:      vpcID,
					SubnetType: vpc.SubnetTypeBCCNAT,
					CIDR:       "10.0.0.0/16",
				}, nil)
				if err != nil {
					t.Errorf("CreateSubnet failed: %v", err)
					return ""
				}
				// Create CCE Cluster
				resq, err := cloud.clientSet.CCEClient.CreateCluster(ctx, &cce.CreateClusterArgs{
					ClusterName: "",
					VPCID:       vpcID,
					SubnetID:    subnetID,
					VMCount:     2,
				})
				if err != nil || resq == nil {
					t.Errorf("CreateCluster failed: %v", err)
					return ""
				}
				return resq.ClusterID
			}(),
			cidr: "10.1.0.0/16",
		},
	}
	for _, c := range cases {
		cloud.CloudConfig.ClusterID = c.clusterID
		vpcID, err := cloud.getVpcID(ctx)
		if err != nil {
			if c.name == "CCE Cluster Nodes not exist" {
				continue
			}
			t.Errorf("getVPCInfo '%s' failed: %v", c.name, err)
			return
		}
		// Check If VPCID match
		nodesResq, err := cloud.clientSet.CCEClient.ListClusterNodes(ctx, c.clusterID, nil)
		if err != nil || nodesResq == nil {
			t.Errorf("getVPCInfo %s failed: %v", c.name, err)
			return
		}
		if len(nodesResq.Nodes) == 0 {
			if c.name == "CCE Cluster Nodes not exist" {
				return
			}
			t.Errorf("getVPC %s failed: cluster's node cannot be 0", c.name)
			return
		}
		if nodesResq.Nodes[0].VPCID != vpcID {
			t.Errorf("getVPCInfo %s failed: gotVPCID=%s wantVPCID=%s", c.name, nodesResq.Nodes[0].VPCID, vpcID)
			return
		}
	}
}

func TestIsConflict(t *testing.T) {
	cloud, _, err := newCluster()
	if err != nil {
		t.Errorf("create cluster error, %v", err)
	}
	cases1 := []vpc.RouteRule{
		{
			RouteRuleID:        "rr-test1",
			RouteTableID:       "rt-test",
			SourceAddress:      "",
			DestinationAddress: "10.0.0.0/8",
			NexthopID:          "",
			NexthopType:        "",
			Description:        "",
		},
		{
			RouteRuleID:        "rr-test2",
			RouteTableID:       "rt-test",
			SourceAddress:      "",
			DestinationAddress: "10.0.0.0/4",
			NexthopID:          "",
			NexthopType:        "",
			Description:        "",
		},
		{
			RouteRuleID:        "rr-test3",
			RouteTableID:       "rt-test",
			SourceAddress:      "",
			DestinationAddress: "10.0.0.0/16",
			NexthopID:          "",
			NexthopType:        "",
			Description:        "",
		},
	}
	cases2 := []vpc.RouteRule{
		{
			RouteRuleID:        "rr-test4",
			RouteTableID:       "rt-test",
			SourceAddress:      "",
			DestinationAddress: "10.0.0.0/8",
			NexthopID:          "",
			NexthopType:        "",
			Description:        "",
		},
		{
			RouteRuleID:        "rr-test5",
			RouteTableID:       "rt-test",
			SourceAddress:      "",
			DestinationAddress: "10.0.0.0/4",
			NexthopID:          "",
			NexthopType:        "",
			Description:        "",
		},
		{
			RouteRuleID:        "rr-test6",
			RouteTableID:       "rt-test",
			SourceAddress:      "",
			DestinationAddress: "10.0.0.0/16",
			NexthopID:          "",
			NexthopType:        "",
			Description:        "",
		},
	}
	for _, c1 := range cases1 {
		for _, c2 := range cases2 {
			re := cloud.isConflict(c1, c2)
			if !re {
				t.Errorf("isConflict error, c1: %v , c2: %v", c1, c2)
			}
		}
	}
	cases3 := []vpc.RouteRule{
		{
			RouteRuleID:        "rr-test1",
			RouteTableID:       "rt-test",
			SourceAddress:      "",
			DestinationAddress: "100.0.0/8",
			NexthopID:          "",
			NexthopType:        "",
			Description:        "",
		},
		{
			RouteRuleID:        "rr-test2",
			RouteTableID:       "rt-test",
			SourceAddress:      "",
			DestinationAddress: "192.168.0.0/4",
			NexthopID:          "",
			NexthopType:        "",
			Description:        "",
		},
		{
			RouteRuleID:        "rr-test3",
			RouteTableID:       "rt-test",
			SourceAddress:      "",
			DestinationAddress: "192.168.0.0/16",
			NexthopID:          "",
			NexthopType:        "",
			Description:        "",
		},
	}
	cases4 := []vpc.RouteRule{
		{
			RouteRuleID:        "rr-test4",
			RouteTableID:       "rt-test",
			SourceAddress:      "",
			DestinationAddress: "11.0.0.0/8",
			NexthopID:          "",
			NexthopType:        "",
			Description:        "",
		},
		{
			RouteRuleID:        "rr-test5",
			RouteTableID:       "rt-test",
			SourceAddress:      "",
			DestinationAddress: "10.0.0.0/4",
			NexthopID:          "",
			NexthopType:        "",
			Description:        "",
		},
		{
			RouteRuleID:        "rr-test6",
			RouteTableID:       "rt-test",
			SourceAddress:      "",
			DestinationAddress: "10.0.0.0/16",
			NexthopID:          "",
			NexthopType:        "",
			Description:        "",
		},
	}
	for _, c3 := range cases3 {
		for _, c4 := range cases4 {
			re := cloud.isConflict(c3, c4)
			if re {
				t.Errorf("isConflict error, c1: %v , c2: %v", c3, c4)
			}
		}
	}
}

func TestGetVpcRouteTable(t *testing.T) {
	ctx := context.Background()
	cloud, _, err := beforeTestRoute()
	if err != nil {
		t.Errorf("beforeTest failed err: %s", err)
	}
	rrs, err := cloud.getVpcRouteTable(ctx)
	if err != nil {
		t.Errorf("getVpcRouteTable error ,%s", err)
	}
	for _, rr := range rrs {
		if rr.RouteTableID != routeruletableID {
			t.Errorf("getVpcRouteTable error , want %s , get %s", routeruletableID, rr.RouteTableID)
		}
	}
}

func TestCheckClusterNode(t *testing.T) {
	ctx := context.Background()
	cloud, resp, err := beforeTestRoute()
	if err != nil {
		t.Errorf("beforeTest failed err: %s", err)
	}
	// right case
	kubeRoute := &cloudprovider.Route{
		Name:            "test1",
		TargetNode:      types.NodeName(resp.Nodes[0].Hostname),
		DestinationCIDR: "test",
		Blackhole:       false,
	}
	_, err = cloud.checkClusterNode(ctx, kubeRoute)
	if err != nil {
		t.Errorf("checkClusterNode err, err: %s", err)
	}
	kubeRoute = &cloudprovider.Route{
		Name:            "test1",
		TargetNode:      types.NodeName(resp.Nodes[0].IP),
		DestinationCIDR: "test",
		Blackhole:       false,
	}
	_, err = cloud.checkClusterNode(ctx, kubeRoute)
	if err != nil {
		t.Errorf("checkClusterNode err, err: %s", err)
	}

	kubeRoute = &cloudprovider.Route{
		Name:            "test",
		TargetNode:      "test1",
		DestinationCIDR: "test",
		Blackhole:       false,
	}

	_, err = cloud.checkClusterNode(ctx, kubeRoute)
	if err == nil {
		t.Errorf("checkClusterNode err, should be an error")
	}
	kubeRoute = &cloudprovider.Route{
		Name:            "",
		TargetNode:      "",
		DestinationCIDR: "",
		Blackhole:       false,
	}
	_, err = cloud.checkClusterNode(ctx, kubeRoute)
	if err == nil {
		t.Errorf("checkClusterNode err, should be an error")
	}
}
