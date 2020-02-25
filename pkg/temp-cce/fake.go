package temp_cce

import (
	"context"
	"fmt"

	"icode.baidu.com/baidu/jpaas-caas/bce-sdk-go/bce"
	"icode.baidu.com/baidu/jpaas-caas/bce-sdk-go/util"
)

// FakeClient for unit test
type FakeClient struct {
	ClusterMap map[string]*Cluster
	NodeMap    map[string]*Node
}

// NewFakeClient for AppBLB fake client
func NewFakeClient() *FakeClient {
	return &FakeClient{
		ClusterMap: map[string]*Cluster{},
		NodeMap:    map[string]*Node{},
	}
}

// CreateCluster to create CCE cluster
// TODO: just for cce-ingress-controller's unit test
func (f *FakeClient) CreateCluster(ctx context.Context, args *CreateClusterArgs) (*CreateClusterResponse, error) {
	if args == nil {
		return nil, fmt.Errorf("CreateCluster failed: args is nil")
	}

	cluster := &Cluster{}

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
	nodes := []*Node{}
	for i := 0; i < args.VMCount; i++ {
		node := &Node{
			ClusterID: cluster.ClusterID,
			VPCID:     args.VPCID,
			SubnetID:  args.SubnetID,
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

	return &CreateClusterResponse{
		ClusterID: cluster.ClusterID,
	}, nil
}

// ListClusterNodes list cluster nodes
func (f *FakeClient) ListClusterNodes(ctx context.Context, clusterID string, option *bce.SignOption) (*ListClusterNodesResponse, error) {
	nodes := []*Node{}

	if _, ok := f.ClusterMap[clusterID]; ok == false {
		return nil, fmt.Errorf("ClusterID %s not exist: NoSuchObject", clusterID)
	}

	for _, node := range f.NodeMap {
		if node.ClusterID == clusterID {
			nodes = append(nodes, node)
		}
	}

	return &ListClusterNodesResponse{
		Nodes: nodes,
	}, nil
}
