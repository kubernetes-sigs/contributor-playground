package cloud_provider

import (
	"context"
	"strings"
	"testing"

	"icode.baidu.com/baidu/jpaas-caas/bce-sdk-go/vpc"
	cce "icode.baidu.com/baidu/jpaas-caas/cloud-provider-baiducloud/pkg/temp-cce"
	"k8s.io/apimachinery/pkg/types"
)

// common func
func newCluster() (cloud *Baiducloud, nodesResq *cce.ListClusterNodesResponse, err error) {
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
		ClusterName: "test",
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

	return cloud, nodesResq, nil

}

// providerID = ""
// providerID = "test"
// providerID = "test/test"
// providerID = "test//test"
// right input
func TestGetInstanceByProviderID(t *testing.T) {
	ctx := context.Background()

	cloud, nodesResq, err := newCluster()
	if err != nil {
		t.Errorf("create cluster error, %v", err)
	}

	cases := []string{
		"",
		"test",
		"test/test",
		"test//test",
		//nodesResq.Nodes[0].InstanceID,
	}

	for _, c := range cases {
		_, err := cloud.getInstanceByProviderID(ctx, c)
		if err != nil {
			t.Logf("getInstanceByProviderID err %v", err)
		}
	}

	cases = []string{
		"test//" + nodesResq.Nodes[0].InstanceID,
	}

	for _, c := range cases {
		node, err := cloud.getInstanceByProviderID(ctx, c)
		if err == nil {
			t.Errorf("getInstanceByProviderID should be err")
			if !strings.Contains(c, node.InstanceID) {
				t.Errorf("getInstanceByProviderID err, providerID %s , instanceID %s", c, node.InstanceID)
			}
		}
	}

}

// name = ""
// name = "test"
// right input
func TestGetInstanceByNodeName(t *testing.T) {
	ctx := context.Background()

	cloud, nodesResq, err := newCluster()
	if err != nil {
		t.Errorf("create cluster error, %v", err)
	}

	cases := []types.NodeName{
		"",
		"test1",
	}

	for _, c := range cases {
		_, err := cloud.getInstanceByNodeName(ctx, c)
		if err == nil {
			t.Errorf("get instance name should get error")
		}
	}

	cases = []types.NodeName{
		types.NodeName(nodesResq.Nodes[0].Hostname),
		types.NodeName(nodesResq.Nodes[0].IP),
	}

	for _, c := range cases {
		node, err := cloud.getInstanceByNodeName(ctx, c)
		if err != nil {
			t.Errorf("get instance by name error, %v", err)
		}
		if node.Hostname != string(c) && node.IP != string(c) {
			t.Errorf("get instance by name error, want %s , get %s", string(c), node.IP)
		}
	}
}

// providerID = ""
// providerID = "test"
// providerID = "test/test"
// providerID = "test//test"
// right input
func TestNodeAddressesByProviderID(t *testing.T) {
	ctx := context.Background()

	cloud, nodesResq, err := newCluster()
	if err != nil {
		t.Errorf("create cluster error, %v", err)
	}

	cases := []string{
		"",
		"test",
		"test/test",
		"test//test",
		nodesResq.Nodes[0].InstanceID,
	}

	for _, c := range cases {
		_, err := cloud.NodeAddressesByProviderID(ctx, c)
		if err == nil {
			t.Errorf("NodeAddressesByProviderID err, should be error here!!!")
		}
	}

	cases = []string{
		"test//" + nodesResq.Nodes[0].InstanceID,
	}

	for _, c := range cases {
		address, err := cloud.NodeAddressesByProviderID(ctx, c)
		if err != nil {
			t.Errorf("NodeAddressesByProviderID err %v", err)
		}
		if address[0].Address != "0.0.0.0" || address[1].Address != "0.0.0.0" {
			t.Errorf("NodeAddressesByProviderID err, providerID %s , addresses %v", c, address)
		}
	}

}
