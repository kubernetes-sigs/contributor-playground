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
	"net"
	"strings"
	"time"

	"github.com/golang/glog"
	"k8s.io/api/core/v1"
	"k8s.io/kubernetes/pkg/cloudprovider"

	"k8s.io/cloud-provider-baiducloud/pkg/sdk/blb"
	"k8s.io/cloud-provider-baiducloud/pkg/sdk/eip"
	"k8s.io/cloud-provider-baiducloud/pkg/sdk/util"
	"k8s.io/cloud-provider-baiducloud/pkg/sdk/vpc"
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
	result := ExtractServiceAnnotation(service)

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
	result := ExtractServiceAnnotation(service)
	err := bc.validateService(service)
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
		glog.V(3).Infof("[%v %v] EnsureLoadBalancer: ensureEIP failed, so delete BLB. ensureEIP error: %s", service.Namespace, service.Name, err)
		args := blb.DeleteLoadBalancerArgs{
			LoadBalancerId: lb.BlbId,
		}
		deleteLoadBalancerErr := bc.clientSet.Blb().DeleteLoadBalancer(&args)
		if deleteLoadBalancerErr != nil {
			glog.V(3).Infof("[%v %v] EnsureLoadBalancer: delete BLB error: %s", service.Namespace, service.Name, deleteLoadBalancerErr)
		}
		if service.Annotations != nil {
			delete(service.Annotations, ServiceAnnotationLoadBalancerId)
		}
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
	result := ExtractServiceAnnotation(service)
	serviceName := getServiceName(service)
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
	if lb.PublicIp != "" {
		if len(service.Spec.LoadBalancerIP) != 0 {
			glog.V(3).Infof("[%v %v] EnsureLoadBalancerDeleted: LoadBalancerIP is set, not delete EIP.", serviceName, clusterName)
			glog.V(2).Infof("[%v %v] EnsureLoadBalancerDeleted: delete %v FINISH", serviceName, clusterName, serviceName)
			return nil
		}
		glog.V(3).Infof("[%v %v] EnsureLoadBalancerDeleted: Start delete EIP: %s", serviceName, clusterName, lb.PublicIp)
		err = bc.deleteEIP(lb.PublicIp)
		if err != nil {
			return err
		}
	}
	glog.V(2).Infof("[%v %v] EnsureLoadBalancerDeleted: delete %v FINISH", serviceName, clusterName, serviceName)
	return nil
}

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
	service.Annotations[ServiceAnnotationLoadBalancerId] = lb.BlbId
	glog.V(2).Infof("[%v %v] WorkAround for old version, lb: %v", service.Namespace, service.Name, lb)
}

func (bc *Baiducloud) validateService(service *v1.Service) error {
	if len(service.Spec.Ports) == 0 {
		return fmt.Errorf("requested load balancer with no ports")
	}
	for _, port := range service.Spec.Ports {
		if port.Protocol != v1.ProtocolTCP {
			return fmt.Errorf("only TCP LoadBalancer is supported for Baidu CCE")
		}
	}
	return nil
}

func (bc *Baiducloud) ensureBLB(ctx context.Context, clusterName string, service *v1.Service, nodes []*v1.Node, serviceAnnotation *ServiceAnnotation) (*blb.LoadBalancer, error) {
	var lb *blb.LoadBalancer
	var err error
	if len(serviceAnnotation.LoadBalancerId) == 0 { // blb not exist, create one and update annotation
		glog.V(3).Infof("[%v %v] EnsureLoadBalancer create blb!", service.Namespace, service.Name)
		vpcId, subnetId, err := bc.getVpcInfoForBLB()
		if err != nil {
			return nil, fmt.Errorf("Can't get VPC info for BLB: %v\n", err)
		}
		allocateVip := false
		if serviceAnnotation.LoadBalancerAllocateVip == "true" {
			allocateVip = true
		}
		args := blb.CreateLoadBalancerArgs{
			Name:        bc.ClusterID + "/" + getServiceName(service),
			VpcID:       vpcId,
			SubnetID:    subnetId,
			Desc:        "auto generated by cce:" + bc.ClusterID,
			AllocateVip: allocateVip,
		}
		glog.V(3).Infof("[%v %v] EnsureLoadBalancer create blb args: %v", service.Namespace, service.Name, args)
		resp, err := bc.clientSet.Blb().CreateLoadBalancer(&args)
		if err != nil {
			return nil, err
		}
		glog.V(3).Infof("[%v %v] EnsureLoadBalancer create blb success, BLB name: %s, BLB id: %s, BLB address: %s.", service.Namespace, service.Name, resp.Name, resp.LoadBalancerId, resp.Address)
		argsDesc := blb.DescribeLoadBalancersArgs{
			LoadBalancerId: resp.LoadBalancerId,
		}
		lbs, err := bc.clientSet.Blb().DescribeLoadBalancers(&argsDesc)
		if err != nil {
			return nil, err
		}
		if len(lbs) != 1 {
			tryCount := 0
			for {
				tryCount ++
				if tryCount > 10 {
					return nil, fmt.Errorf("EnsureLoadBalancer create blb success but query get none")
				}
				glog.V(3).Infof("[%v %v] EnsureLoadBalancer create blb success but query get none, tryCount: ", service.Namespace, service.Name, tryCount)
				lbs, err = bc.clientSet.Blb().DescribeLoadBalancers(&argsDesc)
				if err != nil {
					return nil, err
				}
				if len(lbs) == 1 {
					glog.V(3).Infof("[%v %v] EnsureLoadBalancer create blb success and query get one, tryCount: ", service.Namespace, service.Name, tryCount)
					break
				}
				time.Sleep(10 * time.Second)
			}
		}
		lb = &lbs[0]
		if service.Annotations == nil {
			service.Annotations = make(map[string]string)
		}
		service.Annotations[ServiceAnnotationLoadBalancerId] = lb.BlbId
	} else { // blb already exist, get info from cloud
		var exists bool
		lb, exists, err = bc.getBCELoadBalancerById(serviceAnnotation.LoadBalancerId)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, fmt.Errorf("EnsureLoadBalancer getBCELoadBalancerById failed, target blb not exist, blb id: %v", serviceAnnotation.LoadBalancerId)
		}
		glog.V(3).Infof("[%v %v] EnsureLoadBalancer: blb already exists: %v", service.Namespace, service.Name, lb)
	}
	lb, err = bc.waitForLoadBalancer(lb)
	if err != nil {
		return nil, err
	}

	// update listener
	glog.V(2).Infof("[%v %v] EnsureLoadBalancer: reconcileListeners!", service.Namespace, service.Name)
	err = bc.reconcileListeners(service, lb)
	if err != nil {
		return nil, err
	}
	lb, err = bc.waitForLoadBalancer(lb)
	if err != nil {
		return nil, err
	}

	// update backend server
	glog.V(2).Infof("[%v %v] EnsureLoadBalancer: reconcileBackendServers!", service.Namespace, service.Name)
	err = bc.reconcileBackendServers(nodes, lb)
	if err != nil {
		return nil, err
	}
	lb, err = bc.waitForLoadBalancer(lb)
	if err != nil {
		return nil, err
	}

	return lb, nil
}

func (bc *Baiducloud) ensureEIP(ctx context.Context, clusterName string, service *v1.Service, nodes []*v1.Node, serviceAnnotation *ServiceAnnotation, lb *blb.LoadBalancer) (string, error) {
	if lb == nil {
		return "", fmt.Errorf("[%v %v] EnsureLoadBalancer: ensureEIP need not nil lb", service.Namespace, service.Name)
	}
	if len(service.Spec.LoadBalancerIP) == 0 {
		// not set LoadBalancerIP
		return bc.ensureEIPWithNoSpecificIP(ctx, clusterName, service, nodes, serviceAnnotation, lb)
	} else {
		// set LoadBalancerIP
		return bc.ensureEIPWithSpecificIP(ctx, clusterName, service, nodes, serviceAnnotation, lb)
	}
}

func (bc *Baiducloud) ensureEIPWithNoSpecificIP(ctx context.Context, clusterName string, service *v1.Service, nodes []*v1.Node, serviceAnnotation *ServiceAnnotation, lb *blb.LoadBalancer) (string, error) {
	pubIP := lb.PublicIp
	if len(pubIP) == 0 { // blb not bind eip, mostly case ==>
		glog.V(2).Infof("[%v %v] EnsureLoadBalancer: createEIP!", service.Namespace, service.Name)
		args, err := bc.getEipCreateArgsFromAnnotation(serviceAnnotation)
		if err != nil {
			glog.Errorf("[%v %v] getEipCreateArgsFromAnnotation failed: %v", service.Namespace, service.Name, err)
			return "", err
		}
		if len(args.Name) == 0 {
			args.Name = lb.Name // default EIP name = lb name
		}
		pubIP, err = bc.createEIP(args, lb)
		if err != nil {
			if pubIP != "" {
				args := eip.EipArgs{
					Ip: pubIP,
				}
				bc.clientSet.Eip().DeleteEip(&args)
			}
			return "", err
		}
	} else { // blb already bind eip
		glog.V(3).Infof("[%v %v] EnsureLoadBalancer: blb's eip already exists, start to ensure...", service.Namespace, service.Name)
		argsGet := eip.GetEipsArgs{
			Ip: pubIP,
		}
		eips, err := bc.clientSet.Eip().GetEips(&argsGet)
		if err != nil {
			return "", err
		}
		if len(eips) == 0 {
			err = fmt.Errorf("[%v %v] EnsureLoadBalancer: EIP %s not Exist", service.Namespace, service.Name, pubIP)
			return "", err
		}
		targetEip := eips[0]
		if (len(serviceAnnotation.ElasticIPPaymentTiming) != 0 && serviceAnnotation.ElasticIPPaymentTiming != targetEip.PaymentTiming) ||
			(len(serviceAnnotation.ElasticIPBillingMethod) != 0 && serviceAnnotation.ElasticIPBillingMethod != targetEip.BillingMethod) {
			glog.V(3).Infof("[%v %v] EnsureLoadBalancer: EIP config change, need delete old eip and create new one", service.Namespace, service.Name)
			// TODO
			//pubIP, err = bc.deleteOldAndCreateNewEip(service, serviceAnnotation, pubIP, lb)
			//if err != nil {
			//	return "", err
			//}
			return "", fmt.Errorf("not support change ElasticIP PaymentTiming or ElasticIP BillingMethod, you can delete old and create a new one")
		}
		if serviceAnnotation.ElasticIPBandwidthInMbps != 0 && serviceAnnotation.ElasticIPBandwidthInMbps != targetEip.BandwidthInMbps {
			glog.V(3).Infof("[%v %v] EnsureLoadBalancer: EIP config change, need change ElasticIPBandwidthInMbps", service.Namespace, service.Name)
			// just validate args
			_, err := bc.getEipCreateArgsFromAnnotation(serviceAnnotation)
			if err != nil {
				glog.Errorf("[%v %v] Eip Args error: %v", service.Namespace, service.Name, err)
				return "", err
			}
			err = bc.resizeEip(service, serviceAnnotation, pubIP)
			if err != nil {
				return "", err
			}
		}
	}
	return pubIP, nil
}

func (bc *Baiducloud) ensureEIPWithSpecificIP(ctx context.Context, clusterName string, service *v1.Service, nodes []*v1.Node, serviceAnnotation *ServiceAnnotation, lb *blb.LoadBalancer) (string, error) {
	pubIP := lb.PublicIp
	loadBalancerIP := service.Spec.LoadBalancerIP
	glog.V(3).Infof("[%v %v] EnsureLoadBalancer: Try to bind Custom LoadBalancerIP %s to BLB %s.", service.Namespace, service.Name, loadBalancerIP, lb.BlbId)
	if len(pubIP) == 0 { // blb not bind target eip
		// check eip status
		argsGet := eip.GetEipsArgs{
			Ip: loadBalancerIP,
		}
		eips, err := bc.clientSet.Eip().GetEips(&argsGet)
		if err != nil {
			return "", err
		}
		if len(eips) == 0 {
			err = fmt.Errorf("[%v %v] EnsureLoadBalancer: EIP %s not Exist", service.Namespace, service.Name, loadBalancerIP)
			return "", err
		} else {
			eipStatus := eips[0].Status
			for index := 0; (index < 10) && (eipStatus != "available"); index++ {
				glog.V(3).Infof("[%v %v] Eip: %s is not available, retry:  %d", service.Namespace, service.Name, loadBalancerIP, index)
				time.Sleep(10 * time.Second)
				eips, err := bc.clientSet.Eip().GetEips(&argsGet)
				if err != nil {
					return "", err
				}
				eipStatus = eips[0].Status
			}
			glog.V(3).Infof("Eip final status is: %s", eipStatus)
			if eipStatus != "available" {
				return "", fmt.Errorf("[%v %v] EnsureLoadBalancer: target eip not available", service.Namespace, service.Name)
			}
		}

		// bind
		lb.Status = "unknown" // add here to do loop
		for index := 0; (index < 10) && (lb.Status != "available"); index++ {
			glog.V(3).Infof("[%v %v] BLB: %s is not available, retry:  %d", service.Namespace, service.Name, lb.BlbId, index)
			time.Sleep(10 * time.Second)
			newlb, exist, err := bc.getBCELoadBalancerById(lb.BlbId)
			if err != nil {
				glog.V(3).Infof("getBCELoadBalancer error: %s", lb.BlbId)
				return "", err
			}
			if !exist {
				glog.V(3).Infof("getBCELoadBalancer not exist: %s", lb.BlbId)
				return "", fmt.Errorf("BLB not exists:%s", lb.BlbId)
			}
			lb = newlb
			glog.V(3).Infof("[%v %v] BLB status is : %s", service.Namespace, service.Name, lb.Status)
		}
		argsBind := &eip.BindEipArgs{
			Ip:           loadBalancerIP,
			InstanceId:   lb.BlbId,
			InstanceType: eip.BLB,
		}
		glog.V(3).Infof("[%v %v] Bind EIP: %v", service.Namespace, service.Name, argsBind)
		glog.V(3).Infof("[%v %v] Bind BLB: %v", service.Namespace, service.Name, lb)
		err = bc.clientSet.Eip().BindEip(argsBind)
		if err != nil {
			glog.V(3).Infof("BindEip error: %v", err)
			return "", err
		}
		lb.PublicIp = loadBalancerIP
		pubIP = loadBalancerIP
		glog.V(3).Infof("[%v %v] EnsureLoadBalancer: Bind EIP to BLB success.", service.Namespace, service.Name)
	} else { // blb already bind eip
		if pubIP == loadBalancerIP { // blb bind correct LoadBalancerIP
			glog.V(3).Infof("[%v %v] EnsureLoadBalancer: BLB %s already bind EIP %s.", service.Namespace, service.Name, lb.BlbId, pubIP)
		} else { // blb not bind correct LoadBalancerIP, need update
			glog.V(3).Infof("[%v %v] EnsureLoadBalancer: BLB %s already bind EIP %s, but need updating to %s.", service.Namespace, service.Name, lb.BlbId, pubIP, loadBalancerIP)
			// check eip status
			argsGet := eip.GetEipsArgs{
				Ip: pubIP,
			}
			eips, err := bc.clientSet.Eip().GetEips(&argsGet)
			if err != nil {
				return "", err
			}
			if len(eips) > 0 {
				unbindArgs := eip.EipArgs{
					Ip: pubIP,
				}
				// just unbind, not delete
				err := bc.clientSet.Eip().UnbindEip(&unbindArgs)
				if err != nil {
					glog.V(3).Infof("Unbind Eip error : %s", err.Error())
					return "", err
				}
			}
			// bind
			lb.Status = "unknown" // add here to do loop
			for index := 0; (index < 10) && (lb.Status != "available"); index++ {
				glog.V(3).Infof("[%v %v] BLB: %s is not available, retry:  %d", service.Namespace, service.Name, lb.BlbId, index)
				time.Sleep(10 * time.Second)
				newlb, exist, err := bc.getBCELoadBalancerById(lb.BlbId)
				if err != nil {
					glog.V(3).Infof("getBCELoadBalancer error: %s", lb.BlbId)
					return "", err
				}
				if !exist {
					glog.V(3).Infof("getBCELoadBalancer not exist: %s", lb.BlbId)
					return "", fmt.Errorf("BLB not exists:%s", lb.BlbId)
				}
				lb = newlb
				glog.V(3).Infof("[%v %v] BLB status is : %s", service.Namespace, service.Name, lb.Status)
			}
			argsBind := &eip.BindEipArgs{
				Ip:           loadBalancerIP,
				InstanceId:   lb.BlbId,
				InstanceType: eip.BLB,
			}
			glog.V(3).Infof("[%v %v] Bind EIP: %v", service.Namespace, service.Name, argsBind)
			glog.V(3).Infof("[%v %v] Bind BLB: %v", service.Namespace, service.Name, lb)
			err = bc.clientSet.Eip().BindEip(argsBind)
			if err != nil {
				glog.V(3).Infof("BindEip error: %v", err)
				return "", err
			}
			lb.PublicIp = loadBalancerIP
			pubIP = loadBalancerIP
			glog.V(3).Infof("[%v %v] EnsureLoadBalancer: Bind EIP to BLB success.", service.Namespace, service.Name)
		}
	}
	return pubIP, nil
}

func (bc *Baiducloud) getBCELoadBalancer(name string) (lb *blb.LoadBalancer, exists bool, err error) {
	args := blb.DescribeLoadBalancersArgs{
		LoadBalancerName: name,
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
	args := blb.DescribeLoadBalancersArgs{
		LoadBalancerId: id,
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

// PortListener describe listener port
type PortListener struct {
	Port     int
	Protocol string
	NodePort int32
}

func (bc *Baiducloud) reconcileListeners(service *v1.Service, lb *blb.LoadBalancer) error {
	expected := make(map[int]PortListener)
	// add expected ports
	for _, servicePort := range service.Spec.Ports {
		expected[int(servicePort.Port)] = PortListener{
			Port:     int(servicePort.Port),
			Protocol: string(servicePort.Protocol),
			NodePort: servicePort.NodePort,
		}
	}
	// delete or update unexpected ports
	all, err := bc.getAllListeners(lb)
	if err != nil {
		return err
	}
	var deleteList []PortListener
	for _, l := range all {
		port, ok := expected[l.Port]
		if !ok {
			// delete listener port
			// add to deleteList
			deleteList = append(deleteList, l)
		} else {
			if l != port {
				// update listener port
				err := bc.updateListener(lb, port)
				if err != nil {
					return err
				}
			}
			delete(expected, l.Port)
		}
	}
	// delete listener
	if len(deleteList) > 0 {
		err = bc.deleteListener(lb, deleteList)
		if err != nil {
			return err
		}
	}

	// create expected listener
	for _, pl := range expected {
		err := bc.createListener(lb, pl)
		if err != nil {
			return err
		}
	}
	return nil
}

func (bc *Baiducloud) findPortListener(lb *blb.LoadBalancer, port int, proto string) (PortListener, error) {
	switch proto {
	case "HTTP":
		// TODO
	case "TCP":
		args := blb.DescribeTCPListenerArgs{
			LoadBalancerId: lb.BlbId,
			ListenerPort:   port,
		}
		ls, err := bc.clientSet.Blb().DescribeTCPListener(&args)
		if err != nil {
			return PortListener{}, err
		}
		if len(ls) < 1 {
			return PortListener{}, fmt.Errorf("there is no tcp listener blb:%s  port:%d", lb.BlbId, port)
		}
		return PortListener{
			Port:     ls[0].ListenerPort,
			NodePort: int32(ls[0].BackendPort),
			Protocol: proto,
		}, nil
	case "HTTPS":
		// TODO
	case "UDP":
		// TODO
	}
	return PortListener{}, fmt.Errorf("protocol not match: %s", proto)
}

func (bc *Baiducloud) getAllListeners(lb *blb.LoadBalancer) ([]PortListener, error) {
	allListeners := []PortListener{}
	// add TCPlisteners
	args := blb.DescribeTCPListenerArgs{
		LoadBalancerId: lb.BlbId,
	}
	ls, err := bc.clientSet.Blb().DescribeTCPListener(&args)
	if err != nil {
		return nil, err
	}
	for _, listener := range ls {
		allListeners = append(allListeners, PortListener{
			Port:     listener.ListenerPort,
			Protocol: "TCP",
			NodePort: int32(listener.BackendPort),
		})
	}

	// add HTTPlisteners HTTPS UDP
	// TODO
	return allListeners, nil
}

func (bc *Baiducloud) createListener(lb *blb.LoadBalancer, pl PortListener) error {
	switch pl.Protocol {
	case "HTTP":
		// TODO
	case "TCP":
		args := blb.CreateTCPListenerArgs{
			LoadBalancerId: lb.BlbId,
			ListenerPort:   pl.Port,
			BackendPort:    int(pl.NodePort),
			Scheduler:      "RoundRobin",
		}
		err := bc.clientSet.Blb().CreateTCPListener(&args)
		if err != nil {
			return err
		}
		return nil
	case "HTTPS":
		// TODO
	case "UDP":
		// TODO
	}
	return fmt.Errorf("CreateListener protocol not match: %s", pl.Protocol)
}

func (bc *Baiducloud) updateListener(lb *blb.LoadBalancer, pl PortListener) error {
	switch pl.Protocol {
	case "HTTP":
		// TODO
	case "TCP":
		args := blb.UpdateTCPListenerArgs{
			LoadBalancerId: lb.BlbId,
			ListenerPort:   pl.Port,
			BackendPort:    int(pl.NodePort),
			Scheduler:      "RoundRobin",
		}
		err := bc.clientSet.Blb().UpdateTCPListener(&args)
		if err != nil {
			return err
		}
		return nil
	case "HTTPS":
		// TODO
	case "UDP":
		// TODO
	}
	return fmt.Errorf("updateListener protocol not match: %s", pl.Protocol)
}

func (bc *Baiducloud) deleteListener(lb *blb.LoadBalancer, pl []PortListener) error {
	portList := []int{}
	for _, l := range pl {
		portList = append(portList, l.Port)
	}
	args := blb.DeleteListenersArgs{
		LoadBalancerId: lb.BlbId,
		PortList:       portList,
	}
	err := bc.clientSet.Blb().DeleteListeners(&args)
	if err != nil {
		return err
	}
	return nil
}

func (bc *Baiducloud) getAllBackendServer(lb *blb.LoadBalancer) ([]blb.BackendServer, error) {
	args := blb.DescribeBackendServersArgs{
		LoadBalancerId: lb.BlbId,
	}
	bs, err := bc.clientSet.Blb().DescribeBackendServers(&args)
	if err != nil {
		return nil, err
	}
	return bs, nil
}

func (bc *Baiducloud) reconcileBackendServers(nodes []*v1.Node, lb *blb.LoadBalancer) error {
	expectedServer := make(map[string]string)
	for _, node := range nodes {
		splitted := strings.Split(node.Spec.ProviderID, "//")
		name := splitted[1]
		expectedServer[name] = node.ObjectMeta.Name
	}
	allBS, err := bc.getAllBackendServer(lb)
	if err != nil {
		return err
	}
	var removeList []string
	// remove unexpected servers
	for _, bs := range allBS {
		_, exists := expectedServer[bs.InstanceId]
		if !exists {
			removeList = append(removeList, bs.InstanceId)
		}
		delete(expectedServer, bs.InstanceId)
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
	var addList []blb.BackendServer
	// add expected servers
	for insID, _ := range expectedServer {
		addList = append(addList, blb.BackendServer{
			InstanceId: insID,
			Weight:     100,
		})
	}
	if len(addList) > 0 {
		args := blb.AddBackendServersArgs{
			LoadBalancerId:    lb.BlbId,
			BackendServerList: addList,
		}
		err = bc.clientSet.Blb().AddBackendServers(&args)
		if err != nil {
			return err
		}
	}
	return nil
}

func (bc *Baiducloud) createEIP(args *eip.CreateEipArgs, lb *blb.LoadBalancer) (string, error) {
	glog.V(3).Infof("CreateEip:  %v", args)
	ip, err := bc.clientSet.Eip().CreateEip(args)
	if err != nil {
		return "", err
	}
	argsGet := eip.GetEipsArgs{
		Ip: ip,
	}
	eips, err := bc.clientSet.Eip().GetEips(&argsGet)
	if err != nil {
		return "", err
	}
	if len(eips) > 0 {
		eipStatus := eips[0].Status
		for index := 0; (index < 10) && (eipStatus != "available"); index++ {
			glog.V(3).Infof("Eip: %s is not available, retry:  %d", ip, index)
			time.Sleep(10 * time.Second)
			eips, err := bc.clientSet.Eip().GetEips(&argsGet)
			if err != nil {
				return "", err
			}
			if len(eips) == 0 {
				return "", fmt.Errorf("createEIP failed: CreateEip success but query failed")
			}
			eipStatus = eips[0].Status
		}
		glog.V(3).Infof("Eip status is: %s", eipStatus)
	}
	lb.Status = "unknown" // add here to do loop
	for index := 0; (index < 10) && (lb.Status != "available"); index++ {
		glog.V(3).Infof("BLB: %s is not available, retry:  %d", lb.BlbId, index)
		time.Sleep(10 * time.Second)
		newlb, exist, err := bc.getBCELoadBalancerById(lb.BlbId)
		if err != nil {
			glog.V(3).Infof("getBCELoadBalancer error: %s", lb.BlbId)
			return "", err
		}
		if !exist {
			glog.V(3).Infof("getBCELoadBalancer not exist: %s", lb.BlbId)
			return "", fmt.Errorf("BLB not exists:%s", lb.BlbId)
		}
		lb = newlb
		glog.V(3).Infof("BLB status is : %s", lb.Status)
	}
	argsBind := &eip.BindEipArgs{
		Ip:           ip,
		InstanceId:   lb.BlbId,
		InstanceType: eip.BLB,
	}
	glog.V(3).Infof("BindEip:  %v", argsBind)
	glog.V(3).Infof("Bind BLB: %v", lb)
	err = bc.clientSet.Eip().BindEip(argsBind)
	if err != nil {
		glog.V(3).Infof("BindEip error: %v", err)
		return ip, err
	}
	lb.PublicIp = ip
	glog.V(3).Infof("createEIP: lb.PublicIp is %s", lb.PublicIp)
	return ip, nil
}

func (bc *Baiducloud) deleteEIP(ip string) error {
	argsGet := eip.GetEipsArgs{
		Ip: ip,
	}
	eips, err := bc.clientSet.Eip().GetEips(&argsGet)
	if err != nil {
		return err
	}
	if len(eips) > 0 {
		eipStatus := eips[0].Status
		for index := 0; (index < 10) && (eipStatus != "available"); index++ {
			glog.V(3).Infof("Eip: %s is not available, retry:  %d", ip, index)
			time.Sleep(10 * time.Second)
			eips, err := bc.clientSet.Eip().GetEips(&argsGet)
			if err != nil {
				return err
			}
			eipStatus = eips[0].Status
		}
	}
	args := eip.EipArgs{
		Ip: ip,
	}
	err = bc.clientSet.Eip().DeleteEip(&args)
	if err != nil {
		return err
	}
	return nil
}

func (bc *Baiducloud) waitForLoadBalancer(lb *blb.LoadBalancer) (*blb.LoadBalancer, error) {
	lb.Status = "unknown" // add here to do loop
	for index := 0; (index < 10) && (lb.Status != "available"); index++ {
		glog.V(3).Infof("BLB: %s is not available, retry:  %d", lb.BlbId, index)
		time.Sleep(10 * time.Second)
		newlb, exist, err := bc.getBCELoadBalancerById(lb.BlbId)
		if err != nil {
			glog.V(3).Infof("getBCELoadBalancer error: %s", lb.BlbId)
			return newlb, err
		}
		if !exist {
			glog.V(3).Infof("getBCELoadBalancer not exist: %s", lb.BlbId)
			return newlb, fmt.Errorf("BLB not exists:%s", lb.BlbId)
		}
		lb = newlb
		glog.V(3).Infof("BLB status is : %s", lb.Status)
		if index == 9 && lb.Status != "available" {
			return nil, fmt.Errorf("waitForLoadBalancer failed after retry")
		}
	}

	return lb, nil
}

// TODO: 存在很大的优化空间
// 背景：BLB与NAT子网存在冲突，当用户的集群在NAT子网内时，CCE会创建一个保留子网，类型是通用型，名字是CCE-Reserve，给BLB
// CCE-Reserve 参数：
// (1) 名字: CCE-Reserve
// (2) 可用区：第一台虚机所在可用区
// (3) CIDR
//		   IP：第一台虚机所在子网的下一个子网
//		   Mask：第一台虚机所在子网的Mask
// (4) VPC：第一台虚机所在VPC
// (5) 类型：通用型
func (bc *Baiducloud) getVpcInfoForBLB() (string, string, error) {
	// get prefer vpc info
	ins, err := bc.clientSet.Cce().ListInstances(bc.ClusterID)
	if err != nil {
		return "", "", err
	}
	if len(ins) == 0 {
		return "", "", fmt.Errorf("getVpcInfoForBLB failed since instance num is zero")
	}
	vpcId := ins[0].VpcId
	subnetId := ins[0].SubnetId

	// check subnet
	subnet, err := bc.clientSet.Vpc().DescribeSubnet(subnetId)
	if err != nil {
		return "", "", fmt.Errorf("DescribeSubnet failed: %v", err)
	}
	if subnet.SubnetType == "BCC" {
		return vpcId, subnetId, nil
	}

	// get subnet list and choose preferred one
	params := make(map[string]string, 0)
	params["vpcId"] = subnet.VpcID
	subnets, err := bc.clientSet.Vpc().ListSubnet(params)
	if err != nil {
		return "", "", fmt.Errorf("ListSubnet failed: %v", err)
	}
	for _, subnet := range subnets {
		if subnet.Name == "系统预定义子网" {
			return subnet.VpcID, subnet.SubnetID, nil
		}
		if subnet.Name == "CCE-Reserve" {
			return subnet.VpcID, subnet.SubnetID, nil
		}
	}

	// create one
	currentCidr := subnet.Cidr
	for { // loop
		_, cidr, err := net.ParseCIDR(currentCidr)
		if err != nil {
			return "", "", fmt.Errorf("ParseCIDR failed: %v", err)
		}
		mask, _ := cidr.Mask.Size()
		nextCidr, notExist := util.NextSubnet(cidr, mask)
		if notExist {
			return "", "", fmt.Errorf("NextSubnet failed: %v", err)
		}
		currentCidr = nextCidr.String()
		createSubnetArgs := &vpc.CreateSubnetArgs{
			Name:       "CCE-Reserve",
			ZoneName:   subnet.ZoneName,
			Cidr:       nextCidr.String(),
			VpcID:      subnet.VpcID,
			SubnetType: "BCC",
		}
		newSubnetId, err := bc.clientSet.Vpc().CreateSubnet(createSubnetArgs)
		if err != nil {
			glog.V(3).Infof("CreateSubnet failed: %v, will try again.", err)
			time.Sleep(3 * time.Second)
			continue
		}
		return subnet.VpcID, newSubnetId, nil
	}
}

func (bc *Baiducloud) getEipCreateArgsFromAnnotation(serviceAnnotation *ServiceAnnotation) (*eip.CreateEipArgs, error) {
	var args *eip.CreateEipArgs

	paymentTiming := serviceAnnotation.ElasticIPPaymentTiming
	if len(paymentTiming) == 0 {
		paymentTiming = eip.PAYMENTTIMING_POSTPAID // default Postpaid
	}
	billingMethod := serviceAnnotation.ElasticIPBillingMethod
	if len(billingMethod) == 0 {
		billingMethod = eip.BILLINGMETHOD_BYTRAFFIC // default ByTraffic
	}
	bandwidthInMbps := serviceAnnotation.ElasticIPBandwidthInMbps
	reservationLength := serviceAnnotation.ElasticIPReservationLength
	switch paymentTiming {
	case eip.PAYMENTTIMING_PREPAID:
		if len(serviceAnnotation.ElasticIPBillingMethod) != 0 {
			return nil, fmt.Errorf("when using Prepaid EIP, do not need to set ElasticIPBillingMethod")
		}
		if bandwidthInMbps == 0 { // not set bandwidthInMbps
			bandwidthInMbps = 200
		} else {
			if bandwidthInMbps < 1 || bandwidthInMbps > 200 {
				return nil, fmt.Errorf("prepaid EIP bandwidthInMbps should in [1, 200]")
			}
		}
		reservationLengthAllowed := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36}
		rightReservationLength := false
		for _, length := range reservationLengthAllowed {
			if reservationLength == length {
				rightReservationLength = true
			}
		}
		if !rightReservationLength {
			return nil, fmt.Errorf("prepaid EIP reservationLength should in [1,2,3,4,5,6,7,8,9,12,24,36]")
		}
		args = &eip.CreateEipArgs{
			Name:            serviceAnnotation.ElasticIPName,
			BandwidthInMbps: bandwidthInMbps,
			Billing: &eip.Billing{
				PaymentTiming: paymentTiming,
				Reservation: &eip.Reservation{
					ReservationLength:   reservationLength,
					ReservationTimeUnit: "Month",
				},
			},
		}
	case eip.PAYMENTTIMING_POSTPAID:
		switch billingMethod {
		case eip.BILLINGMETHOD_BYTRAFFIC:
			if bandwidthInMbps == 0 { // not set bandwidthInMbps
				bandwidthInMbps = 1000
			} else {
				if bandwidthInMbps < 1 || bandwidthInMbps > 1000 {
					return nil, fmt.Errorf("postpaid ByTraffic EIP bandwidthInMbps should in [1, 1000]")
				}
			}
		case eip.BILLINGMETHOD_BYBANDWIDTH:
			if bandwidthInMbps == 0 { // not set bandwidthInMbps
				bandwidthInMbps = 200
			} else {
				if bandwidthInMbps < 1 || bandwidthInMbps > 200 {
					return nil, fmt.Errorf("postpaid ByBandwidth EIP bandwidthInMbps should in [1, 200]")
				}
			}
		default:
			return nil, fmt.Errorf("not support target ElasticIPBillingMethod: %v", billingMethod)
		}
		args = &eip.CreateEipArgs{
			Name:            serviceAnnotation.ElasticIPName,
			BandwidthInMbps: bandwidthInMbps,
			Billing: &eip.Billing{
				PaymentTiming: paymentTiming,
				BillingMethod: billingMethod,
			},
		}
	default:
		return nil, fmt.Errorf("not support target ElasticIPPaymentTiming: %v", paymentTiming)
	}

	return args, nil
}

func (bc *Baiducloud) deleteOldAndCreateNewEip(service *v1.Service, serviceAnnotation *ServiceAnnotation, oldEip string, lb *blb.LoadBalancer) (string, error) {
	err := bc.deleteEIP(oldEip)
	if err != nil {
		return "", err
	}
	glog.V(2).Infof("[%v %v] EnsureLoadBalancer: createEIP!", service.Namespace, service.Name)
	args, err := bc.getEipCreateArgsFromAnnotation(serviceAnnotation)
	if err != nil {
		glog.Errorf("[%v %v] getEipCreateArgsFromAnnotation failed: %v", service.Namespace, service.Name, err)
		return "", err
	}
	if len(args.Name) == 0 {
		args.Name = lb.Name // default EIP name = lb name
	}
	pubIP, err := bc.createEIP(args, lb)
	if err != nil {
		if pubIP != "" {
			args := eip.EipArgs{
				Ip: pubIP,
			}
			bc.clientSet.Eip().DeleteEip(&args)
		}
		return "", err
	}
	return pubIP, nil
}

func (bc *Baiducloud) resizeEip(service *v1.Service, serviceAnnotation *ServiceAnnotation, targetEip string) error {
	return bc.clientSet.Eip().ResizeEip(&eip.ResizeEipArgs{
		BandwidthInMbps: serviceAnnotation.ElasticIPBandwidthInMbps,
		Ip:              targetEip,
	})
}
