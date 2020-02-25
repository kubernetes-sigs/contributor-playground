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
	"context"
	"fmt"

	"icode.baidu.com/baidu/jpaas-caas/bce-sdk-go/blb"
	v1 "k8s.io/api/core/v1"
	"k8s.io/klog"
)

// workaround to support old version, can be removed if not support old version
func (bc *Baiducloud) workAround(ctx context.Context, clusterName string, service *v1.Service) {
	serviceKey := fmt.Sprintf("%s/%s", service.Namespace, service.Name)
	klog.Infof(Message(ctx, fmt.Sprintf("WorkAround for service %s begin", serviceKey)))
	defer klog.Infof(Message(ctx, fmt.Sprintf("WorkAround for service %s end", serviceKey)))
	lb, exists, err := bc.getBLBByName(ctx, bc.GetLoadBalancerName(ctx, clusterName, service))
	if err != nil {
		return
	}
	if !exists {
		return
	}
	if service.Annotations == nil {
		service.Annotations = make(map[string]string)
	}

	service.Annotations[ServiceAnnotationCceAutoAddLoadBalancerID] = lb.BlbId
}

func (bc *Baiducloud) validateService(service *v1.Service) error {
	if len(service.Spec.Ports) == 0 {
		return fmt.Errorf("requested load balancer with no ports")
	}
	for _, port := range service.Spec.Ports {
		switch port.Protocol {
		case "TCP":
			continue
		case "UDP":
			continue
		case "HTTP":
			return fmt.Errorf("HTTP is not supported")
		case "HTTPS":
			return fmt.Errorf("HTTPS is not supported")
		default:
			return fmt.Errorf("target protocol is not supported: %v", port.Protocol)
		}
	}
	return nil
}

func (bc *Baiducloud) getBLBByName(ctx context.Context, name string) (lb *blb.LoadBalancer, exists bool, err error) {
	if len(name) == 0 {
		return nil, false, fmt.Errorf("LoadBalancerName is empty")
	}
	args := blb.DescribeLoadBalancersArgs{
		LoadBalancerName: name,
		ExactlyMatch:     true,
	}
	lbs, err := bc.clientSet.BLBClient.DescribeLoadBalancers(ctx, &args, bc.getSignOption(ctx))
	if err != nil {
		klog.Errorf(Message(ctx, fmt.Sprintf("getBLBByName failed: %s", err)))
		return &blb.LoadBalancer{}, false, err
	}

	if len(lbs) == 0 {
		klog.Info(Message(ctx, fmt.Sprintf("BLB named %s not exist", name)))
		return &blb.LoadBalancer{}, false, nil
	}
	if len(lbs) > 1 {
		msg := fmt.Sprintf("multi BLB named %s exist", name)
		klog.Warning(Message(ctx, msg))
		return &blb.LoadBalancer{}, false, fmt.Errorf(msg)
	}

	return &lbs[0], true, nil
}

func (bc *Baiducloud) getBLBByID(ctx context.Context, id string) (lb *blb.LoadBalancer, exists bool, err error) {
	if len(id) == 0 {
		return nil, false, fmt.Errorf("LoadBalancerId is empty")
	}
	args := blb.DescribeLoadBalancersArgs{
		LoadBalancerId: id,
		ExactlyMatch:   true,
	}
	lbs, err := bc.clientSet.BLBClient.DescribeLoadBalancers(ctx, &args, bc.getSignOption(ctx))
	if err != nil {
		klog.Infof(Message(ctx, fmt.Sprintf("getBLBByID blb %s not exists: %v", args.LoadBalancerId, err)))
		return &blb.LoadBalancer{}, false, err
	}
	if len(lbs) == 0 {
		msg := fmt.Sprintf("BLB with id %s not exist", id)
		klog.Warning(Message(ctx, msg))
		return &blb.LoadBalancer{}, false, fmt.Errorf(msg)
	}
	return &lbs[0], true, nil
}

// This returns a human-readable version of the Service used to tag some resources.
// This is only used for human-readable convenience, and not to filter.
func getBlbName(clusterID string, service *v1.Service) string {
	blbName := fmt.Sprintf("CCE/SVC/%s/%s/%s", clusterID, service.Namespace, service.Name)
	if annotationName, ok := service.Annotations[ServiceAnnotationLoadBalancerBLBName]; ok {
		blbName = annotationName
	}
	return blbName
}
