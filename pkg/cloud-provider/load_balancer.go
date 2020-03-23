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
	"strings"
	"time"

	"k8s.io/klog"
	v1 "k8s.io/api/core/v1"
	cloudprovider "k8s.io/cloud-provider"
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
	ctx = context.WithValue(ctx, RequestID, GetRandom())
	// workaround to support old version, can be removed if not support old version
	lb, exist, err := bc.getServiceAssociatedBLB(ctx, clusterName, service)
	if err != nil {
		return nil, false, err
	}
	if !exist {
		return nil, exist, nil
	}

	var ip string
	if internalVpc, ok := service.Annotations[ServiceAnnotationLoadBalancerInternalVpc]; ok && internalVpc == "true" {
		ip = lb.Address // internal vpc ip
	} else {
		ip = lb.PublicIp // EIP
	}
	klog.V(3).Infof("[%v %v] GetLoadBalancer ip: %s", service.Namespace, service.Name, ip)

	return &v1.LoadBalancerStatus{Ingress: []v1.LoadBalancerIngress{{IP: ip}}}, true, nil
}

// GetLoadBalancerName returns the name of the load balancer. Implementations must treat the
// *v1.Service parameter as read-only and not modify it.
func (bc *Baiducloud) GetLoadBalancerName(ctx context.Context, clusterName string, service *v1.Service) string {
	//GCE requires that the name of a load balancer starts with a lower case letter.
	ret := "a" + string(service.UID)
	ret = strings.Replace(ret, "-", "", -1)
	//AWS requires that the name of a load balancer is shorter than 32 bytes.
	if len(ret) > 32 {
		ret = ret[:32]
	}
	return ret
}

// EnsureLoadBalancer creates a new load balancer 'name', or updates the existing one. Returns the status of the balancer
// Implementations must treat the *v1.Service and *v1.Node
// parameters as read-only and not modify them.
// Parameter 'clusterName' is the name of the cluster as presented to kube-controller-manager
func (bc *Baiducloud) EnsureLoadBalancer(ctx context.Context, clusterName string, service *v1.Service, nodes []*v1.Node) (*v1.LoadBalancerStatus, error) {
	ctx = context.WithValue(ctx, RequestID, GetRandom())
	serviceKey := fmt.Sprintf("%s/%s", service.Namespace, service.Name)
	klog.Infof(Message(ctx, fmt.Sprintf("EnsureLoadBalancer for service %s", serviceKey)))
	err := bc.validateService(service)
	if err != nil {
		return nil, err
	}

	// ensure BLB
	lb, err := bc.ensureBLB(ctx, clusterName, service)
	if err != nil {
		return nil, err
	}

	err = bc.reconcileListeners(ctx, clusterName, service)
	if err != nil {
		return nil, err
	}

	err = bc.reconcileBackendServers(ctx, clusterName, service, nodes)
	if err != nil {
		return nil, err
	}

	var pubIP string
	if internalIP, ok := service.Annotations[ServiceAnnotationLoadBalancerInternalVpc]; ok && internalIP == "true" {
		pubIP = lb.Address
		klog.Infof(Message(ctx, fmt.Sprintf("EnsureLoadBalancer for service %s/%s: use LoadBalancerInternalVpc ip %s", service.Namespace, service.Name, pubIP)))
	} else {
		// ensure EIP
		pubIP, err = bc.ensureEIP(ctx, clusterName, service)
		if err != nil {
			return nil, err
		}
		if service.Annotations == nil {
			service.Annotations = make(map[string]string, 0)
		}
		service.Annotations[ServiceAnnotationCceAutoAddEip] = pubIP
		klog.Infof(Message(ctx, fmt.Sprintf("EnsureLoadBalancer for service %s/%s: use EIP %s", service.Namespace, service.Name, pubIP)))
	}
	return &v1.LoadBalancerStatus{Ingress: []v1.LoadBalancerIngress{{IP: pubIP}}}, nil
}

// UpdateLoadBalancer updates hosts under the specified load balancer.
// Implementations must treat the *v1.Service and *v1.Node
// parameters as read-only and not modify them.
// Parameter 'clusterName' is the name of the cluster as presented to kube-controller-manager
func (bc *Baiducloud) UpdateLoadBalancer(ctx context.Context, clusterName string, service *v1.Service, nodes []*v1.Node) error {
	ctx = context.WithValue(ctx, RequestID, GetRandom())
	startTime := time.Now()
	serviceKey := fmt.Sprintf("%s/%s", service.Namespace, service.Name)
	defer func() {
		klog.V(4).Infof(Message(ctx, fmt.Sprintf("Finished UpdateLoadBalancer for service %q (%v)", serviceKey, time.Since(startTime))))
	}()
	klog.Infof(Message(ctx, fmt.Sprintf("UpdateLoadBalancer for service %s", serviceKey)))
	err := bc.reconcileBackendServers(ctx, clusterName, service, nodes)
	if err != nil {
		return err
	}
	klog.Infof(Message(ctx, fmt.Sprintf("UpdateLoadBalancer for service %s success", serviceKey)))
	return nil
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
	ctx = context.WithValue(ctx, RequestID, GetRandom())
	lb, exist, err := bc.getServiceAssociatedBLB(ctx, clusterName, service)
	if err != nil {
		return err
	}
	if !exist {
		msg := fmt.Sprintf("BLB for service %s already deleted", service.Name)
		klog.Info(Message(ctx, msg))
	}

	if internalIP, ok := service.Annotations[ServiceAnnotationLoadBalancerInternalVpc]; !ok || internalIP != "true" {
		err = bc.ensureEipDeleted(ctx, service, lb)
		if err != nil {
			return err
		}
	}

	if reserveLB, ok := service.Annotations[ServiceAnnotationLoadBalancerReserveLB]; !ok || reserveLB != "true" {
		err = bc.ensureBLBDeleted(ctx, lb)
		if err != nil {
			return err
		}
	}

	return nil
}
