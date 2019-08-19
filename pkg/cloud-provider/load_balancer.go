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
	"strings"

	"github.com/golang/glog"
	v1 "k8s.io/api/core/v1"
	"k8s.io/kubernetes/pkg/cloudprovider"

	"k8s.io/cloud-provider-baiducloud/pkg/cloud-sdk/blb"
	"k8s.io/cloud-provider-baiducloud/pkg/cloud-sdk/eip"
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

	if len(result.CceAutoAddLoadBalancerId) == 0 {
		return nil, false, nil
	}
	lb, exists, err := bc.getBCELoadBalancerById(result.CceAutoAddLoadBalancerId)
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
	glog.V(3).Infof("[%v %v] EnsureLoadBalancer(%v, %v, %v, %v, %v)",
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
	service.Annotations[ServiceAnnotationCceAutoAddEip] = pubIP
	glog.V(3).Infof("[%v %v] EnsureLoadBalancer: EXTERNAL-IP is %s", service.Namespace, service.Name, pubIP)
	return &v1.LoadBalancerStatus{Ingress: []v1.LoadBalancerIngress{{IP: pubIP}}}, nil
}

// UpdateLoadBalancer updates hosts under the specified load balancer.
// Implementations must treat the *v1.Service and *v1.Node
// parameters as read-only and not modify them.
// Parameter 'clusterName' is the name of the cluster as presented to kube-controller-manager
func (bc *Baiducloud) UpdateLoadBalancer(ctx context.Context, clusterName string, service *v1.Service, nodes []*v1.Node) error {
	glog.V(3).Infof("[%v %v] UpdateLoadBalancer(%v, %v, %v, %v, %v)",
		clusterName, service.Namespace, service.Name, bc.Region, service.Spec.LoadBalancerIP, service.Spec.Ports, service.Annotations)
	result, err := ExtractServiceAnnotation(service)
	if err != nil {
		return err
	}
	err = bc.validateService(service)
	if err != nil {
		return err
	}

	if len(result.CceAutoAddLoadBalancerId) != 0 {
		_, err = bc.EnsureLoadBalancer(ctx, clusterName, service, nodes)
		return err
	}
	glog.V(3).Infof("UpdateLoadBalancer is not necessary")
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
	// workaround to support old version, can be removed if not support old version
	bc.workAround(service)
	result, err := ExtractServiceAnnotation(service)
	if err != nil {
		// if annotation has error, then creation must be failed. So return nil to tell k8s lb has been deleted.
		return nil
	}
	serviceName := getServiceName(service)
	if len(result.CceAutoAddLoadBalancerId) == 0 {
		glog.V(1).Infof("[%v %v] EnsureLoadBalancerDeleted: target load balancer not create successful. So, no need to delete BLB and EIP", serviceName, clusterName)
		return nil
	}

	glog.V(2).Infof("[%v %v] EnsureLoadBalancerDeleted: START lbId=%q", serviceName, clusterName, result.CceAutoAddLoadBalancerId)

	// reconcile logic is capable of fully reconcile, so we can use this to delete
	service.Spec.Ports = []v1.ServicePort{}

	lb, existsLb, err := bc.getBCELoadBalancerById(result.CceAutoAddLoadBalancerId)
	glog.V(3).Infof("[%v %v] EnsureLoadBalancerDeleted: getBCELoadBalancer : %s", serviceName, clusterName, lb.BlbId)
	if err != nil {
		glog.V(3).Infof("[%v %v] EnsureLoadBalancerDeleted get error: %s", serviceName, clusterName, err.Error())
		return err
	}
	if !existsLb {
		glog.V(3).Infof("[%v %v] EnsureLoadBalancerDeleted: target blb not exist", serviceName, clusterName)
		err = bc.DeleteEipFinally(service, result, lb, serviceName, clusterName)
		if err != nil {
			return err
		}
		return nil
	}

	if len(result.LoadBalancerExistId) == 0 { //user does not assign the blbid in the annotation
		// start delete blb and eip, delete blb first
		glog.V(3).Infof("[%v %v] EnsureLoadBalancerDeleted: Start delete BLB: %s", serviceName, clusterName, lb.BlbId)
		args := blb.DeleteLoadBalancerArgs{
			LoadBalancerId: lb.BlbId,
		}
		err = bc.clientSet.Blb().DeleteLoadBalancer(&args)
		if err != nil {
			return err
		}
	} else if result.LoadBalancerExistId == "error_blb_has_been_used" {
		return nil
	} else {
		//get allListeners & delete Listeners
		allListeners, err := bc.getAllListeners(lb)
		if err != nil {
			return err
		}
		if len(allListeners) > 0 {
			err = bc.deleteListener(lb, allListeners)
			if err != nil {
				return err
			}
		}
		//get allServers & delete BackendServers
		allServers, err := bc.getAllBackendServer(lb)
		var removeList []string
		if err != nil {
			return err
		}
		for _, server := range allServers {
			removeList = append(removeList, server.InstanceId)
		}

		if len(removeList) > 0 {
			args := blb.RemoveBackendServersArgs{
				LoadBalancerId:    lb.BlbId,
				BackendServerList: removeList,
			}
			err = bc.clientSet.Blb().RemoveBackendServers(&args)
			if err != nil {
				return err
			}
		}

		// annotation "LoadBalancerInternalVpc" exists
		if result.LoadBalancerInternalVpc == "true" { //do not assign the eip
			if service.Annotations != nil {
				delete(service.Annotations, ServiceAnnotationCceAutoAddLoadBalancerId)
			}
			glog.V(3).Infof("[%v %v] EnsureLoadBalancerDeleted: use LoadBalancerInternalVpc, no EIP to delete", service.Namespace, service.Name)
			//todo recover eip for blb which has eip in the begin.
			glog.V(2).Infof("[%v %v] EnsureLoadBalancerDeleted: delete %v FINISH", serviceName, clusterName, serviceName)
			return nil
		}

		//annotation "LoadBalancerIP" exists
		//unbind eip & blb when user assigned the eip
		if len(service.Spec.LoadBalancerIP) != 0 { //use user’s eip, do not delete
			unbindArgs := eip.EipArgs{
				Ip: service.Spec.LoadBalancerIP,
			}
			// just unbind, not delete
			err := bc.clientSet.Eip().UnbindEip(&unbindArgs)
			if err != nil {
				glog.V(3).Infof("Unbind Eip error : %s", err.Error())
				return nil
			}
			return nil
		}

		//get targetEip
		var targetEip string
		if len(service.Status.LoadBalancer.Ingress) != 0 { // P0: use service EXTERNAL_IP
			targetEip = service.Status.LoadBalancer.Ingress[0].IP
		}
		if len(targetEip) == 0 { // P1: use BLB public ip
			targetEip = lb.PublicIp
		}
		//users may unbind eip manually
		if len(targetEip) == 0 { // get none EIP
			glog.V(3).Infof("Eip does not exist, Delete completed ")
			return nil
		}

		// blb if has eip in the begin
		if strings.Contains(lb.Desc, "cce_auto_create_eip") {
			glog.V(3).Infof("EnsureLoadBalancerDeleted: delete eip created by cce: %s", lb.Desc)
			unbindArgs := eip.EipArgs{
				Ip: targetEip,
			}
			lb.Desc = strings.TrimPrefix(lb.Desc, "cce_auto_create_eip")
			newLbArg := blb.UpdateLoadBalancerArgs{
				LoadBalancerId: lb.BlbId,
				Desc:           lb.Desc,
				Name:           lb.Name,
			}
			err = bc.clientSet.Blb().UpdateLoadBalancer(&newLbArg)
			if err != nil {
				return err
			}
			// unbind & delete
			err := bc.clientSet.Eip().UnbindEip(&unbindArgs)
			if err != nil {
				glog.V(3).Infof("Unbind Eip error : %s", err.Error())
				if strings.Contains(err.Error(), "EipNotFound") {
					return nil
				}
				return err
			}
			err = bc.deleteEIP(targetEip)
			if err != nil {
				return err
			}
		}
		if service.Annotations != nil {
			delete(service.Annotations, ServiceAnnotationCceAutoAddLoadBalancerId)
		}
		return nil
	}
	// does not assign blb, delete EIP
	err = bc.DeleteEipFinally(service, result, lb, serviceName, clusterName)
	if err != nil {
		return err
	}
	return nil
}
