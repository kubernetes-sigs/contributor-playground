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
	"fmt"

	"github.com/golang/glog"
	v1 "k8s.io/api/core/v1"
	"k8s.io/kubernetes/pkg/cloudprovider"

	"k8s.io/cloud-provider-baiducloud/pkg/cloud-sdk/blb"
)

// workaround to support old version, can be removed if not support old version
func (bc *Baiducloud) workAround(service *v1.Service) {
	lb, exists, err := bc.getBCELoadBalancer(cloudprovider.GetLoadBalancerName(service))
	if err != nil {
		return
	}
	if !exists {
		return
	}
	if service.Annotations == nil {
		service.Annotations = make(map[string]string)
	}
	// TODO: 不会更新最终Service的annotaion，因为ip没变
	service.Annotations[ServiceAnnotationCceAutoAddLoadBalancerId] = lb.BlbId
	glog.V(2).Infof("[%v %v] WorkAround for old version, lb: %v", service.Namespace, service.Name, lb)
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
			return fmt.Errorf("UDP is not supported")
		case "HTTPS":
			return fmt.Errorf("UDP is not supported")
		default:
			return fmt.Errorf("target protocol is not supported: %v", port.Protocol)
		}
	}
	return nil
}

func (bc *Baiducloud) getBCELoadBalancer(name string) (lb *blb.LoadBalancer, exists bool, err error) {
	if len(name) == 0 {
		return nil, false, fmt.Errorf("LoadBalancerName is empty")
	}
	args := blb.DescribeLoadBalancersArgs{
		LoadBalancerName: name,
		ExactlyMatch:     true,
	}
	lbs, err := bc.clientSet.Blb().DescribeLoadBalancers(&args)
	if err != nil {
		glog.V(2).Infof("getBCELoadBalancer  blb not exists ! %v", args)
		return &blb.LoadBalancer{}, false, err
	}
	if len(lbs) != 1 {
		return &blb.LoadBalancer{}, false, nil
	}

	return &lbs[0], true, nil
}

func (bc *Baiducloud) getBCELoadBalancerById(id string) (lb *blb.LoadBalancer, exists bool, err error) {
	if len(id) == 0 {
		return nil, false, fmt.Errorf("LoadBalancerId is empty")
	}
	args := blb.DescribeLoadBalancersArgs{
		LoadBalancerId: id,
		ExactlyMatch:   true,
	}
	lbs, err := bc.clientSet.Blb().DescribeLoadBalancers(&args)
	if err != nil {
		glog.V(2).Infof("getBCELoadBalancer blb %s not exists: %v", args.LoadBalancerId, err)
		return &blb.LoadBalancer{}, false, err
	}
	if len(lbs) != 1 {
		return &blb.LoadBalancer{}, false, nil
	}
	return &lbs[0], true, nil
}

// This returns a human-readable version of the Service used to tag some resources.
// This is only used for human-readable convenience, and not to filter.
func getServiceName(service *v1.Service) string {
	return fmt.Sprintf("%s/%s", service.Namespace, service.Name)
}
