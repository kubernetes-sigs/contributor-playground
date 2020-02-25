package temp_cce

import (
	"context"
	"encoding/json"
	"fmt"

	"icode.baidu.com/baidu/jpaas-caas/bce-sdk-go/bce"
)

// CreateCluster create CCE cluster
func (c *Client) CreateCluster(ctx context.Context, args *CreateClusterArgs) (*CreateClusterResponse, error) {
	return &CreateClusterResponse{}, nil
}

// ListClusterNodes gets all Instances of a cluster.
func (c *Client) ListClusterNodes(ctx context.Context, clusterID string, option *bce.SignOption) (*ListClusterNodesResponse, error) {
	if clusterID == "" {
		return nil, fmt.Errorf("clusterID is nil")
	}

	params := map[string]string{
		"clusterUuid": clusterID,
	}

	req, err := bce.NewRequest("GET", c.GetURL("v1/node", params), nil)

	if err != nil {
		return nil, err
	}

	resp, err := c.SendRequest(ctx, req, option)

	if err != nil {
		return nil, err
	}

	bodyContent, err := resp.GetBodyContent()

	if err != nil {
		return nil, err
	}

	var nodesResq ListClusterNodesResponse
	err = json.Unmarshal(bodyContent, &nodesResq)

	if err != nil {
		return nil, err
	}

	return &nodesResq, nil
}
