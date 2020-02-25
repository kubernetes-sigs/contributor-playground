package cloud_provider

import (
	"context"
	"testing"

	"icode.baidu.com/baidu/jpaas-caas/bce-sdk-go/bce"
	"icode.baidu.com/baidu/jpaas-caas/bce-sdk-go/blb"
	cce "icode.baidu.com/baidu/jpaas-caas/cloud-provider-baiducloud/pkg/temp-cce"
	api "k8s.io/api/core/v1"
)

func beforeTestBackend() (*Baiducloud, *cce.ListClusterNodesResponse, *blb.CreateLoadBalancerResponse, error) {
	cloud, nodesRes, blbRes, err := beforeTestBlb()
	if err != nil {
		return nil, nil, nil, err
	}
	return cloud, nodesRes, blbRes, nil
}

// case1: bs is nil
// case2: add bs, then get bs
func TestGetAllBackendServer(t *testing.T) {
	cloud, _, resp, err := beforeTestBlb()
	if err != nil {
		t.Errorf("beforeTestBlb err, err: %v", err)
	}
	ctx := context.Background()
	// bs is nil
	lb := &blb.LoadBalancer{
		BlbId: resp.LoadBalancerId,
	}
	bs, err := cloud.getAllBackendServer(ctx, lb)
	if err != nil {
		t.Errorf("getAllBackendServer err, err: %v", err)
	}
	if len(bs) != 0 {
		t.Errorf("getAllBackendServer err, bs  should be nil but get : %v", bs)
	}
	// add bs
	bsAdd := []blb.BackendServer{
		{
			InstanceId: "1",
		},
		{
			InstanceId: "2",
		},
	}
	args := blb.AddBackendServersArgs{
		LoadBalancerId:    lb.BlbId,
		BackendServerList: bsAdd,
	}
	err = cloud.clientSet.BLBClient.AddBackendServers(ctx, &args, &bce.SignOption{
		CustomSignFunc: CCEServiceSign,
	})
	if err != nil {
		t.Errorf("AddBackendServers err, err: %v", err)
	}
	// get bs
	bs, err = cloud.getAllBackendServer(ctx, lb)
	if err != nil {
		t.Errorf("getAllBackendServer err, err: %v", err)
	}
	if len(bs) != 2 {
		t.Errorf("getAllBackendServer err, bs  should be nil but get : %v", bs)
	}
}

// case1: bs is nil
// case2: add bs, then delete bs
func TestDeleteAllBackendServers(t *testing.T) {
	cloud, _, resp, err := beforeTestBlb()
	if err != nil {
		t.Errorf("beforeTestBlb err, err: %v", err)
	}
	ctx := context.Background()
	lb := &blb.LoadBalancer{
		BlbId: resp.LoadBalancerId,
	}
	// bs is nil
	err = cloud.deleteAllBackendServers(ctx, lb)
	if err != nil {
		t.Errorf("deleteAllBackendServers err, err : %v", err)
	}
	// add bs
	bsAdd := []blb.BackendServer{
		{
			InstanceId: "1",
		},
		{
			InstanceId: "2",
		},
	}
	args := blb.AddBackendServersArgs{
		LoadBalancerId:    lb.BlbId,
		BackendServerList: bsAdd,
	}
	err = cloud.clientSet.BLBClient.AddBackendServers(ctx, &args, &bce.SignOption{
		CustomSignFunc: CCEServiceSign,
	})
	if err != nil {
		t.Errorf("AddBackendServers err, err: %v", err)
	}
	err = cloud.deleteAllBackendServers(ctx, lb)
	if err != nil {
		t.Errorf("deleteAllBackendServers err, err : %v", err)
	}
	bs, err := cloud.getAllBackendServer(ctx, lb)
	if err != nil {
		t.Errorf("getAllBackendServer err, err: %v", err)
	}
	if len(bs) != 0 {
		t.Errorf("getAllBackendServer err, bs should be nil but get : %v", bs)
	}
}
func arraysHelp(des []blb.BackendServer, ss ...string) bool {
	if len(ss) == 0 {
		return false
	}
	count := 0
	for _, d := range des {
		for i, s := range ss {
			if d.InstanceId == s {
				count++
				ss[i] = "-1&&-|-&&1-"
			}
		}
	}
	return count == len(ss)
}

//  case 1:
// 	candidateBackends: ["1", "2", "3"] existingBackends: ["4", "5"] targetBackendsNum: 1
// 	rsToAdd: ["1"] or ["2"] or ["3"]  rsToDel: ["4", "5"]
// 	case 2:
// 	candidateBackends: ["1", "2", "3", "4"] existingBackends: ["4", "5"] targetBackendsNum: 1
// 	rsToAdd: []  rsToDel: ["5"]
// 	case 3:
// 	candidateBackends: ["1", "2", "3", "4"] existingBackends: ["4", "5"] targetBackendsNum: 3
// 	rsToAdd: ["1", "2"] or ["1", "3"] or ["3", "2"] rsToDel: ["5"]
func TestMergeBackend(t *testing.T) {
	// case1
	candidateBackends := []blb.BackendServer{
		{
			InstanceId: "1",
		},
		{
			InstanceId: "2",
		},
		{
			InstanceId: "3",
		},
	}
	existingBackends := []blb.BackendServer{
		{
			InstanceId: "4",
		},
		{
			InstanceId: "5",
		},
	}
	rsToAdd, rsToDel, err := mergeBackend(candidateBackends, existingBackends, 1)
	if err != nil {
		t.Errorf("mergeBackend err, err: %v", err)
	}
	if !arraysHelp(rsToAdd, "1") && !arraysHelp(rsToAdd, "2") && !arraysHelp(rsToAdd, "3") ||
		!arraysHelp(rsToDel, "4") || !arraysHelp(rsToDel, "5") {
		t.Errorf("mergeBackend err, want 1 or 2 or 3| 4 ,5 get %v | %v, %v", rsToAdd[0].InstanceId,
			rsToDel[0].InstanceId, rsToDel[1].InstanceId)
	}
	// case2
	candidateBackends = []blb.BackendServer{
		{
			InstanceId: "1",
		},
		{
			InstanceId: "2",
		},
		{
			InstanceId: "3",
		},
		{
			InstanceId: "4",
		},
	}
	existingBackends = []blb.BackendServer{
		{
			InstanceId: "4",
		},
		{
			InstanceId: "5",
		},
	}
	rsToAdd, rsToDel, err = mergeBackend(candidateBackends, existingBackends, 1)
	if err != nil {
		t.Errorf("mergeBackend err, err: %v", err)
	}
	if len(rsToAdd) != 0 || !arraysHelp(rsToDel, "5") {
		t.Errorf("mergeBackend err, want nil |  5 get %v | %v", rsToAdd[0].InstanceId,
			rsToDel[0].InstanceId)
	}
	//case3
	candidateBackends = []blb.BackendServer{
		{
			InstanceId: "1",
		},
		{
			InstanceId: "2",
		},
		{
			InstanceId: "3",
		},
		{
			InstanceId: "4",
		},
	}
	existingBackends = []blb.BackendServer{
		{
			InstanceId: "4",
		},
		{
			InstanceId: "5",
		},
	}
	rsToAdd, rsToDel, err = mergeBackend(candidateBackends, existingBackends, 3)
	if err != nil {
		t.Errorf("mergeBackend err, err: %v", err)
	}
	if !arraysHelp(rsToAdd, "1", "2") && !arraysHelp(rsToAdd, "2", "3") &&
		!arraysHelp(rsToAdd, "1", "3") || !arraysHelp(rsToDel, "5") {
		t.Errorf("mergeBackend err, want 1，2 or 2，3 or 1，3| 5 get %v, %v | %v", rsToAdd[0].InstanceId,
			rsToAdd[1].InstanceId, rsToDel[0].InstanceId)
	}
}

// case1: expected input, get the right output
// case2: nodes or lb is nil
func TestReconcileBackendServers(t *testing.T) {
	cloud, nodesRes, _, err := beforeTestBackend()
	if err != nil {
		t.Errorf("beforeTestBackend err, err: %v", err)
	}
	ctx := context.Background()
	svc := buildService()
	// case1
	nodes := []*api.Node{
		&api.Node{
			Spec: api.NodeSpec{
				ProviderID: "test//" + nodesRes.Nodes[0].InstanceID,
			},
		},
		&api.Node{
			Spec: api.NodeSpec{
				ProviderID: "test//" + nodesRes.Nodes[1].InstanceID,
			},
		},
	}
	err = cloud.reconcileBackendServers(ctx, cloud.ClusterName, svc, nodes)
	if err != nil {
		t.Errorf("reconcileBackendServers err, err: %v p", err)
	}
	// case2
	nodes = []*api.Node{
		&api.Node{
			Spec: api.NodeSpec{
				ProviderID: "test//dadad",
			},
		},
		&api.Node{
			Spec: api.NodeSpec{
				ProviderID: "test//sdasd",
			},
		},
	}
	err = cloud.reconcileBackendServers(ctx, cloud.ClusterName, svc, nodes)
	if err != nil {
		t.Errorf("reconcileBackendServers err, err: %v", err)
	}
}
