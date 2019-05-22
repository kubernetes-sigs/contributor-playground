package cloud_provider

import (
	api "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func buildService() *api.Service {
	return &api.Service{
		ObjectMeta: meta_v1.ObjectMeta{
			Name:      "foo",
			Namespace: api.NamespaceDefault,
		},
		Spec: api.ServiceSpec{},
	}
}

func TestExtractServiceAnnotation(t *testing.T) {
	svc := buildService()

	result, err := ExtractServiceAnnotation(svc)
	if err != nil {
		t.Errorf("failed to extract service annotation: %v", err)
	}
	if result.LoadBalancerSubnetId != "" {
		t.Errorf("failed")
	}

	data := map[string]string{}
	data[ServiceAnnotationLoadBalancerSubnetId] = ""

	svc.SetAnnotations(data)

	result, err = ExtractServiceAnnotation(svc)
	if err != nil {
		t.Errorf("failed to extract service annotation: %v", err)
	}
	if result.LoadBalancerSubnetId != "" {
		t.Errorf("failed")
	}

	data[ServiceAnnotationLoadBalancerSubnetId] = "abc"
	svc.SetAnnotations(data)

	result, err = ExtractServiceAnnotation(svc)
	if err != nil {
		t.Errorf("failed to extract service annotation: %v", err)
	}
	if result.LoadBalancerSubnetId != "abc" {
		t.Errorf("extrac service annotation wrong")
	}
}
