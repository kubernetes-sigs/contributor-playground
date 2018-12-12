/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cloud_provider

import (
	"k8s.io/api/core/v1"
	"github.com/golang/glog"
)

const (
	ServiceAnnotationLoadBalancerPrefix      = "service.beta.kubernetes.io/cce-load-balancer-"
	ServiceAnnotationLoadBalancerId          = ServiceAnnotationLoadBalancerPrefix + "id"
	ServiceAnnotationLoadBalancerInternalVpc = ServiceAnnotationLoadBalancerPrefix + "internal-vpc"
)

func ExtractAnnotationRequest(service *v1.Service) (*AnnotationRequest, *AnnotationRequest) {
	glog.V(4).Infof("start to ExtractAnnotationRequest: %v", service.Annotations)
	defaulted, request := &AnnotationRequest{}, &AnnotationRequest{}
	annotation := make(map[string]string)
	for k, v := range service.Annotations {
		annotation[k] = v
	}

	loadBalancerId, ok := annotation[ServiceAnnotationLoadBalancerId]
	if ok {
		defaulted.LoadBalancerId = loadBalancerId
		request.LoadBalancerId = loadBalancerId
	}

	loadBalancerInternalVpc, ok := annotation[ServiceAnnotationLoadBalancerInternalVpc]
	if ok {
		defaulted.LoadBalancerInternalVpc = loadBalancerInternalVpc
		request.LoadBalancerInternalVpc = loadBalancerInternalVpc
	}

	return defaulted, request
}
