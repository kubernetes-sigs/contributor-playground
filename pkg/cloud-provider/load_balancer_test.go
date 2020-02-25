package cloud_provider

import (
	"context"
	"testing"

	api "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestGetLoadBalancer(t *testing.T) {
	cloud, _, lbResp, err := beforeTestBlb()
	if err != nil {
		t.Errorf("beforeTestBlb err, err: %s", err)
	}
	ctx := context.Background()
	// case is nil
	clusterName := ""
	svc := buildService()
	_, exist, err := cloud.GetLoadBalancer(ctx, clusterName, svc)
	if err != nil {
		t.Errorf("GetLoadBalancer err, err: %s", err)
	}
	if !exist {
		t.Errorf("GetLoadBalancer err, blb should exist")
	}
	// right case1
	clusterName = "test"
	data := map[string]string{}
	data[ServiceAnnotationCceAutoAddLoadBalancerID] = lbResp.LoadBalancerId
	svc.SetAnnotations(data)
	result, exist, err := cloud.GetLoadBalancer(ctx, clusterName, svc)
	if err != nil {
		t.Errorf("GetLoadBalancer err, err: %s", err)
	}
	if !exist {
		t.Errorf("GetLoadBalancer err, blb should exist")
	}
	if len(result.Ingress[0].IP) == 0 {
		t.Errorf("GetLoadBalancer err, get result %v", result)
	}
	// right case2
	svc = buildService()
	clusterName = "test"
	data = map[string]string{}
	data[ServiceAnnotationCceAutoAddLoadBalancerID] = lbResp.LoadBalancerId
	data[ServiceAnnotationLoadBalancerInternalVpc] = "true"
	svc.SetAnnotations(data)
	result2, exist2, err := cloud.GetLoadBalancer(ctx, clusterName, svc)
	if err != nil {
		t.Errorf("GetLoadBalancer err, err: %s", err)
	}
	if !exist2 {
		t.Errorf("GetLoadBalancer err, blb should exist")
	}
	if len(result.Ingress[0].IP) == 0 {
		t.Errorf("GetLoadBalancer err, get result %v", result2)
	}
}

func TestEnsureLoadBalancer(t *testing.T) {
	cloud, _, _, err := beforeTestBlb()
	if err != nil {
		t.Errorf("beforeTestBlb err, err: %s", err)
	}
	ctx := context.Background()
	// case is nil
	clusterName := ""
	svc := &api.Service{
		ObjectMeta: meta_v1.ObjectMeta{
			Name:      "foo",
			Namespace: api.NamespaceDefault,
		},
		Spec: api.ServiceSpec{
			Ports: []api.ServicePort{
				api.ServicePort{
					Protocol: "TCP",
				},
			},
		},
	}
	nodes := []*v1.Node{
		&v1.Node{},
	}
	_, err = cloud.EnsureLoadBalancer(ctx, clusterName, svc, nodes)
	if err == nil {
		t.Errorf("EnsureLoadBalancer err， should be an error here")
	}
	// todo right case1
	clusterName = "test"
	svc.Spec.Ports[0].Protocol = "TCP"
}
