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

	"github.com/golang/glog"
	"k8s.io/api/core/v1"
	"k8s.io/kubernetes/pkg/cloudprovider"

	"k8s.io/cloud-provider-baiducloud/pkg/cloud-sdk/blb"
)

// LoadBalancer returns a balancer interface. Also returns true if the interface is supported, false otherwise.
func (bc *Baiducloud) LoadBalancer() (cloudprovider.LoadBalancer, bool) {
	return bc, true
}

// GetLoadBalancer returns whether the specified load balancer exists, and
// if so, what its status is.
// Implementations must treat the *v1.Service parameter as read-only and not modify it.
// Parameter 'clusterName' is the name of the cluster as presented to kube-controller-manager
func (bc *Baiducloud) GetLoadBalancer(ctx context.Context, clusterName string, service *v1.Service) (status *v1.LoadBalancerStatus, exists bool, err error) {
	// workaround to support old version, can be removed if not support old version
	bc.workAround(service)
	result, err := ExtractServiceAnnotation(service)
	if err != nil {
		return nil, false, err
	}

	if len(result.LoadBalancerId) == 0 {
		return nil, false, nil
	}
	lb, exists, err := bc.getBCELoadBalancerById(result.LoadBalancerId)
	if err != nil {
		return nil, false, err
	}
	if !exists {
		return nil, false, nil
	}

	var ip string
	if result.LoadBalancerInternalVpc == "true" {
		ip = lb.Address // internal vpc ip
	} else {
		ip = lb.PublicIp // EIP
	}
	glog.V(3).Infof("[%v %v] GetLoadBalancer ip: %s", service.Namespace, service.Name, ip)

	return &v1.LoadBalancerStatus{Ingress: []v1.LoadBalancerIngress{{IP: ip}}}, true, nil
}

// EnsureLoadBalancer creates a new load balancer 'name', or updates the existing one. Returns the status of the balancer
// Implementations must treat the *v1.Service and *v1.Node
// parameters as read-only and not modify them.
// Parameter 'clusterName' is the name of the cluster as presented to kube-controller-manager
func (bc *Baiducloud) EnsureLoadBalancer(ctx context.Context, clusterName string, service *v1.Service, nodes []*v1.Node) (*v1.LoadBalancerStatus, error) {
	glog.V(3).Infof("[%v %v] EnsureLoadBalancer(%v, %v, %v, %v, %v, %v, %v)",
		clusterName, service.Namespace, service.Name, bc.Region, service.Spec.LoadBalancerIP, service.Spec.Ports, service.Annotations)
	// workaround to support old version, can be removed if not support old version
	bc.workAround(service)
	result, err := ExtractServiceAnnotation(service)
	if err != nil {
		return nil, err
	}
	err = bc.validateService(service)
	if err != nil {
		return nil, err
	}

	// ensure BLB
	lb, err := bc.ensureBLB(ctx, clusterName, service, nodes, result)
	if err != nil {
		return nil, err
	}
	if result.LoadBalancerInternalVpc == "true" {
		glog.V(3).Infof("[%v %v] EnsureLoadBalancer: use LoadBalancerInternalVpc, EXTERNAL-IP is %s", service.Namespace, service.Name, lb.Address)
		return &v1.LoadBalancerStatus{Ingress: []v1.LoadBalancerIngress{{IP: lb.Address}}}, nil
	}

	// ensure EIP
	pubIP, err := bc.ensureEIP(ctx, clusterName, service, nodes, result, lb)
	if err != nil {
		return nil, err
	}

	glog.V(3).Infof("[%v %v] EnsureLoadBalancer: EXTERNAL-IP is %s", service.Namespace, service.Name, pubIP)
	return &v1.LoadBalancerStatus{Ingress: []v1.LoadBalancerIngress{{IP: pubIP}}}, nil
}

// UpdateLoadBalancer updates hosts under the specified load balancer.
// Implementations must treat the *v1.Service and *v1.Node
// parameters as read-only and not modify them.
// Parameter 'clusterName' is the name of the cluster as presented to kube-controller-manager
func (bc *Baiducloud) UpdateLoadBalancer(ctx context.Context, clusterName string, service *v1.Service, nodes []*v1.Node) error {
	_, err := bc.EnsureLoadBalancer(ctx, clusterName, service, nodes)
	return err
}

// EnsureLoadBalancerDeleted deletes the specified load balancer if it
// exists, returning nil if the load balancer specified either didn't exist or
// was successfully deleted.
// This construction is useful because many cloud providers' load balancers
// have multiple underlying components, meaning a Get could say that the LB
// doesn't exist even if some part of it is still laying around.
// Implementations must treat the *v1.Service parameter as read-only and not modify it.
// Parameter 'clusterName' is the name of the cluster as presented to kube-controller-manager
func (bc *Baiducloud) EnsureLoadBalancerDeleted(ctx context.Context, clusterName string, service *v1.Service) error {
	// workaround to support old version, can be removed if not support old version
	bc.workAround(service)
	result, err := ExtractServiceAnnotation(service)
	if err != nil {
		// if annotation has error, then creation must be failed. So return nil to tell k8s lb has been deleted.
		return nil
	}
	serviceName := getServiceName(service)
	if len(result.LoadBalancerId) == 0 {
		glog.V(1).Infof("[%v %v] EnsureLoadBalancerDeleted: target load balancer not create successful. So, no need to delete BLB and EIP", serviceName, clusterName)
		return nil
	}

	glog.V(2).Infof("[%v %v] EnsureLoadBalancerDeleted: START lbId=%q", serviceName, clusterName, result.LoadBalancerId)

	// reconcile logic is capable of fully reconcile, so we can use this to delete
	service.Spec.Ports = []v1.ServicePort{}

	lb, existsLb, err := bc.getBCELoadBalancerById(result.LoadBalancerId)
	glog.V(3).Infof("[%v %v] EnsureLoadBalancerDeleted: getBCELoadBalancer : %s", serviceName, clusterName, lb.BlbId)
	if err != nil {
		glog.V(3).Infof("[%v %v] EnsureLoadBalancerDeleted get error: %s", serviceName, clusterName, err.Error())
		return err
	}
	if !existsLb {
		glog.V(3).Infof("[%v %v] EnsureLoadBalancerDeleted: target blb not exist", serviceName, clusterName)
		return nil
	}

	// start delete blb and eip, delete blb first
	glog.V(3).Infof("[%v %v] EnsureLoadBalancerDeleted: Start delete BLB: %s", serviceName, clusterName, lb.BlbId)
	args := blb.DeleteLoadBalancerArgs{
		LoadBalancerId: lb.BlbId,
	}
	err = bc.clientSet.Blb().DeleteLoadBalancer(&args)
	if err != nil {
		return err
	}

	// delete EIP
	if result.LoadBalancerInternalVpc == "true" {
		if service.Annotations != nil {
			delete(service.Annotations, ServiceAnnotationLoadBalancerId)
		}
		glog.V(3).Infof("[%v %v] EnsureLoadBalancerDeleted: use LoadBalancerInternalVpc, no EIP to delete", service.Namespace, service.Name, lb.Address)
		glog.V(2).Infof("[%v %v] EnsureLoadBalancerDeleted: delete %v FINISH", serviceName, clusterName, serviceName)
		return nil
	}
	if len(service.Spec.LoadBalancerIP) != 0 {
		if service.Annotations != nil {
			delete(service.Annotations, ServiceAnnotationLoadBalancerId)
		}
		glog.V(3).Infof("[%v %v] EnsureLoadBalancerDeleted: LoadBalancerIP is set, not delete EIP.", serviceName, clusterName)
		glog.V(2).Infof("[%v %v] EnsureLoadBalancerDeleted: delete %v FINISH", serviceName, clusterName, serviceName)
		return nil
	}
	glog.V(3).Infof("[%v %v] EnsureLoadBalancerDeleted: Start delete EIP: %s", serviceName, clusterName, lb.PublicIp)
	var targetEip string
	if len(service.Status.LoadBalancer.Ingress) != 0 { // P0: use service EXTERNAL_IP
		targetEip = service.Status.LoadBalancer.Ingress[0].IP
	}
	if len(targetEip) == 0 { // P1: use BLB public ip
		targetEip = lb.PublicIp
	}
	if len(targetEip) == 0 { // get none EIP
		return fmt.Errorf("EnsureLoadBalancerDeleted failed: can not get a EIP to delete")
	}
	err = bc.deleteEIP(targetEip)
	if err != nil {
		return err
	}
	if service.Annotations != nil {
		delete(service.Annotations, ServiceAnnotationLoadBalancerId)
	}
	glog.V(2).Infof("[%v %v] EnsureLoadBalancerDeleted: delete %v FINISH", serviceName, clusterName, serviceName)
	return nil
}
