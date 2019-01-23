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
	"github.com/golang/glog"
	"k8s.io/api/core/v1"
)

const (
	ServiceAnnotationLoadBalancerPrefix      = "service.beta.kubernetes.io/cce-load-balancer-"
	ServiceAnnotationLoadBalancerId          = ServiceAnnotationLoadBalancerPrefix + "id"
	ServiceAnnotationLoadBalancerInternalVpc = ServiceAnnotationLoadBalancerPrefix + "internal-vpc"
	ServiceAnnotationLoadBalancerAllocateVip = ServiceAnnotationLoadBalancerPrefix + "allocate-vip"
)

const (
	NodeAnnotationPrefix          = "node.alpha.kubernetes.io/"
	NodeAnnotationVpcId           = NodeAnnotationPrefix + "vpc-id"
	NodeAnnotationVpcRouteTableId = NodeAnnotationPrefix + "vpc-route-table-id"
	NodeAnnotationVpcRouteRuleId  = NodeAnnotationPrefix + "vpc-route-rule-id"
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

	loadBalancerAllocateVip, ok := annotation[ServiceAnnotationLoadBalancerAllocateVip]
	if ok {
		defaulted.LoadBalancerAllocateVip = loadBalancerAllocateVip
		request.LoadBalancerAllocateVip = loadBalancerAllocateVip
	}

	return defaulted, request
}

func ExtractNodeAnnotation(node *v1.Node) *NodeAnnotation {
	glog.V(4).Infof("start to ExtractNodeAnnotation: %v", node.Annotations)
	result := &NodeAnnotation{}
	annotation := make(map[string]string)
	for k, v := range node.Annotations {
		annotation[k] = v
	}

	vpcId, ok := annotation[NodeAnnotationVpcId]
	if ok {
		result.VpcId = vpcId
	}

	vpcRouteTableId, ok := annotation[NodeAnnotationVpcRouteTableId]
	if ok {
		result.VpcRouteTableId = vpcRouteTableId
	}

	vpcRouteRuleId, ok := annotation[NodeAnnotationVpcRouteRuleId]
	if ok {
		result.VpcRouteRuleId = vpcRouteRuleId
	}

	return result
}
