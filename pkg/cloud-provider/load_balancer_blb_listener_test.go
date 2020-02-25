package cloud_provider

import (
	"context"
	"testing"

	"icode.baidu.com/baidu/jpaas-caas/bce-sdk-go/bce"
	"icode.baidu.com/baidu/jpaas-caas/bce-sdk-go/blb"
	api "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func beforeTestListener() (*Baiducloud, *blb.CreateLoadBalancerResponse, error) {
	cloud, _, resp, err := beforeTestBlb()
	if err != nil {
		return nil, nil, err
	}
	ctx := context.Background()
	// create tcplistener
	argsTcp := blb.CreateTCPListenerArgs{
		LoadBalancerId: resp.LoadBalancerId,
		ListenerPort:   11,
		BackendPort:    12,
		Scheduler:      "RoundRobin",
	}
	err = cloud.clientSet.BLBClient.CreateTCPListener(ctx, &argsTcp, &bce.SignOption{
		CustomSignFunc: CCEServiceSign,
	})
	if err != nil {
		return nil, nil, err
	}
	// create udplistener
	argsUdp := blb.CreateUDPListenerArgs{
		LoadBalancerId:    resp.LoadBalancerId,
		ListenerPort:      13,
		BackendPort:       14,
		Scheduler:         "RoundRobin",
		HealthCheckString: "HealthCheck",
	}
	err = cloud.clientSet.BLBClient.CreateUDPListener(ctx, &argsUdp, &bce.SignOption{
		CustomSignFunc: CCEServiceSign,
	})
	if err != nil {
		return nil, nil, err
	}
	return cloud, resp, nil
}

func TestGetAllListeners(t *testing.T) {
	cloud, resp, err := beforeTestListener()
	if err != nil {
		t.Errorf("beforeTestListener err, err: %v", err)
	}
	ctx := context.Background()
	case1 := &blb.LoadBalancer{
		BlbId: resp.LoadBalancerId,
	}
	pl, err := cloud.getAllListeners(ctx, case1)
	if err != nil {
		t.Errorf("getAllListeners err, err: %v", err)
	}
	if len(pl) != 2 || pl[0].Protocol != "TCP" ||
		pl[1].Protocol != "UDP" {
		t.Errorf("getAllListeners err, get pl: %v", pl)
	}
	case2 := &blb.LoadBalancer{
		BlbId: "test",
	}
	_, err = cloud.getAllListeners(ctx, case2)
	if err == nil {
		t.Errorf("getAllListeners err, there should be an error but get nil")
	}
}

func TestCreateListener(t *testing.T) {
	cloud, resp, err := beforeTestListener()
	if err != nil {
		t.Errorf("beforeTestListener err, err: %v", err)
	}
	ctx := context.Background()
	lb := &blb.LoadBalancer{
		BlbId: resp.LoadBalancerId,
	}
	pls := []PortListener{
		{
			Port:     12,
			Protocol: "TCP",
			NodePort: 11,
		},
		{
			Port:     12,
			Protocol: "UDP",
			NodePort: 11,
		},
	}
	for _, pl := range pls {
		err = cloud.createListener(ctx, lb, pl)
		if err != nil {
			t.Errorf("createListener err, err: %v", err)
		}
	}
	pls = []PortListener{
		{
			Port:     12,
			Protocol: "HTTPS",
			NodePort: 11,
		},
		{
			Port:     12,
			Protocol: "HTTP",
			NodePort: 11,
		},
		{
			Port:     12,
			Protocol: "test",
			NodePort: 11,
		},
	}
	for _, pl := range pls {
		err = cloud.createListener(ctx, lb, pl)
		if err == nil {
			t.Errorf("createListener err, there should be  but get nil")
		}
	}
	lb1 := &blb.LoadBalancer{
		BlbId: "11",
	}
	pls = []PortListener{
		{
			Port:     12,
			Protocol: "TCP",
			NodePort: 11,
		},
		{
			Port:     12,
			Protocol: "UDP",
			NodePort: 11,
		},
	}
	for _, pl := range pls {
		err = cloud.createListener(ctx, lb1, pl)
		if err == nil {
			t.Errorf("createListener err, there should be  but get nil")
		}
	}
}

func TestUpdateListener(t *testing.T) {
	cloud, resp, err := beforeTestListener()
	if err != nil {
		t.Errorf("beforeTestListener err, err: %v", err)
	}
	ctx := context.Background()
	lb := &blb.LoadBalancer{
		BlbId: resp.LoadBalancerId,
	}
	pls := []PortListener{
		{
			Port:     12,
			Protocol: "TCP",
			NodePort: 11,
		},
		{
			Port:     12,
			Protocol: "UDP",
			NodePort: 11,
		},
	}
	for _, pl := range pls {
		err = cloud.updateListener(ctx, lb, pl)
		if err != nil {
			t.Errorf("updateListener err, err: %v", err)
		}
	}
	pls = []PortListener{
		{
			Port:     12,
			Protocol: "HTTP",
			NodePort: 11,
		},
		{
			Port:     12,
			Protocol: "HTTPS",
			NodePort: 11,
		},
		{
			Port:     12,
			Protocol: "TEST",
			NodePort: 11,
		},
	}
	for _, pl := range pls {
		err = cloud.updateListener(ctx, lb, pl)
		if err == nil {
			t.Errorf("updateListener err, there should be an errorbut get nil")
		}
	}
	lb1 := &blb.LoadBalancer{
		BlbId: "11",
	}
	pls = []PortListener{
		{
			Port:     12,
			Protocol: "TCP",
			NodePort: 11,
		},
		{
			Port:     12,
			Protocol: "UDP",
			NodePort: 11,
		},
	}
	for _, pl := range pls {
		err = cloud.updateListener(ctx, lb1, pl)
		if err == nil {
			t.Errorf("updateListener err, there should be an error but get nil")
		}
	}
}

func TestDeleteListener(t *testing.T) {
	cloud, resp, err := beforeTestListener()
	if err != nil {
		t.Errorf("beforeTestListener err, err: %v", err)
	}
	ctx := context.Background()
	// right case
	lb := &blb.LoadBalancer{
		BlbId: resp.LoadBalancerId,
	}
	pls, err := cloud.getAllListeners(ctx, lb)
	if err != nil {
		t.Errorf("getAllListeners err, err: %v", err)
	}
	if len(pls) == 0 {
		t.Errorf("getAllListeners err, pls : %v ", pls)
	}
	err = cloud.deleteListener(ctx, lb, pls)
	if err != nil {
		t.Errorf("deleteListener err, err: %v", err)
	}
	pls, err = cloud.getAllListeners(ctx, lb)
	if err != nil {
		t.Errorf("getAllListeners err, err: %v", err)
	}
	if len(pls) != 0 {
		t.Errorf("deleteListener err, pls: %v left", pls)
	}
}

func TestReconcileListeners(t *testing.T) {
	cloud, _, err := beforeTestListener()
	if err != nil {
		t.Errorf("beforeTestListener err, err: %v", err)
	}
	ctx := context.Background()
	//service
	case1 := &api.Service{
		ObjectMeta: meta_v1.ObjectMeta{
			Name:      "foo",
			Namespace: api.NamespaceDefault,
		},
		Spec: api.ServiceSpec{
			Ports: []api.ServicePort{
				{
					Name:     "test",
					Port:     11,
					Protocol: "TCP",
					NodePort: 13,
				},
				{
					Name:     "test",
					Port:     12,
					Protocol: "UDP",
					NodePort: 14,
				},
			},
		},
	}
	// dirty blb
	err = cloud.reconcileListeners(ctx, cloud.ClusterName, case1)
	if err != nil {
		t.Errorf("reconcileListeners err, err %v", err)
	}
	err = cloud.reconcileListeners(ctx, cloud.ClusterName,case1)
	if err != nil {
		t.Errorf("reconcileListeners err, err %v", err)
	}
	// to complate...
}
