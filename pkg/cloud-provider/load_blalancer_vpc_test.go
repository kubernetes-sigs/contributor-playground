package cloud_provider

import (
	"context"
	"testing"

	v1 "k8s.io/api/core/v1"
)

func TestGetVpcInfoForBLB(t *testing.T) {
	cloud, resp, err := newCluster()
	if err != nil {
		t.Errorf("create cluster error, %v", err)
	}
	ctx := context.Background()
	// SubnetID is nil
	svc := &v1.Service{}
	svc.Annotations = make(map[string]string, 0)
	svc.Annotations[ServiceAnnotationLoadBalancerSubnetID] = ""
	_, _, err = cloud.getVpcInfoForBLB(ctx, svc)
	if err != nil {
		t.Errorf("getVpcInfoForBLB err, err: %s", err)
	}

	// assign SubnetID
	svc.Annotations[ServiceAnnotationLoadBalancerSubnetID] = resp.Nodes[0].SubnetID
	vpcID, newSubnetID, err := cloud.getVpcInfoForBLB(ctx, svc)
	if err != nil {
		t.Errorf("getVpcInfoForBLB err, err: %s", err)
	}
	if vpcID != resp.Nodes[0].VPCID || len(newSubnetID) == 0 {
		t.Errorf("getVpcInfoForBLB err, get vpcID : %v or get newSubnetID : %v", vpcID, newSubnetID)
	}
}
