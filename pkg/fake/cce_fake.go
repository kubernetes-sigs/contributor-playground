package fake

import (
	"context"
	"fmt"

	"icode.baidu.com/baidu/jpaas-caas/bce-sdk-go/bce"
	"icode.baidu.com/baidu/jpaas-caas/bce-sdk-go/util"
	cce "icode.baidu.com/baidu/jpaas-caas/cloud-provider-baiducloud/pkg/temp-cce"
)

// FakeClient for unit test
type CceFakeClient struct {
	ClusterMap map[string]*cce.Cluster
	NodeMap    map[string]*cce.Node
}

// NewFakeClient for AppBLB fake client
func NewCceFakeClient() *CceFakeClient {
	return &CceFakeClient{
		ClusterMap: map[string]*cce.Cluster{},
		NodeMap:    map[string]*cce.Node{},
	}
}

// CreateCluster to create CCE cluster
// TODO: just for cce-ingress-controller's unit test
func (f *CceFakeClient) CreateCluster(ctx context.Context, args *cce.CreateClusterArgs) (*cce.CreateClusterResponse, error) {
	if args == nil {
		return nil, fmt.Errorf("CreateCluster failed: args is nil")
	}
	cluster := &cce.Cluster{}
	// Generate ClusterID
	for {
		clusterID := util.GenerateBCEShortID("c")
		if _, ok := f.ClusterMap[clusterID]; !ok {
			cluster.ClusterID = clusterID
			f.ClusterMap[clusterID] = cluster
			break
		}
	}
	// Generate Cluster Node
	nodes := []*cce.Node{}
	for i := 0; i < args.VMCount; i++ {
		node := &cce.Node{
			ClusterID:    cluster.ClusterID,
			InstanceName: "test",
			Hostname:     "test",
			IP:           "0.0.0.0",
			VPCID:        args.VPCID,
			SubnetID:     args.SubnetID,
		}
		for {
			instanceID := util.GenerateBCEShortID("i")
			if _, ok := f.NodeMap[instanceID]; !ok {
				node.InstanceID = instanceID
				f.NodeMap[instanceID] = node
				break
			}
		}
		nodes = append(nodes, node)
	}
	return &cce.CreateClusterResponse{
		ClusterID: cluster.ClusterID,
	}, nil
}

// ListClusterNodes list cluster nodes
func (f *CceFakeClient) ListClusterNodes(ctx context.Context, clusterID string, option *bce.SignOption) (*cce.ListClusterNodesResponse, error) {
	nodes := []*cce.Node{}
	if _, ok := f.ClusterMap[clusterID]; ok == false {
		return nil, fmt.Errorf("ClusterID %s not exist: NoSuchObject", clusterID)
	}
	for _, node := range f.NodeMap {
		if node.ClusterID == clusterID {
			nodes = append(nodes, node)
		}
	}
	return &cce.ListClusterNodesResponse{
		Nodes: nodes,
	}, nil
}
