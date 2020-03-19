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
	"time"

	"icode.baidu.com/baidu/jpaas-caas/bce-sdk-go/blb"
	v1 "k8s.io/api/core/v1"
	"k8s.io/klog"
)

// PortListener describe listener port
type PortListener struct {
	Port     int
	Protocol string
	NodePort int32
}

func (bc *Baiducloud) reconcileListeners(ctx context.Context, clusterName string, service *v1.Service) error {
	startTime := time.Now()
	serviceKey := fmt.Sprintf("%s/%s", service.Namespace, service.Name)
	defer func() {
		klog.V(4).Infof(Message(ctx, fmt.Sprintf("Finished reconcileListeners for service %q (%v)", serviceKey, time.Since(startTime))))
	}()
	// add expected ports
	expected := make(map[int]PortListener)
	for _, servicePort := range service.Spec.Ports {
		expected[int(servicePort.Port)] = PortListener{
			Port:     int(servicePort.Port),
			Protocol: string(servicePort.Protocol),
			NodePort: servicePort.NodePort,
		}
	}

	lb, exist, err := bc.getServiceAssociatedBLB(ctx, clusterName, service)
	if err != nil {
		return err
	}
	if !exist {
		return fmt.Errorf("failed to reconcileListeners: lb not exist")
	}

	// delete or update unexpected ports
	all, err := bc.getAllListeners(ctx, lb)
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
				klog.Infof(Message(ctx, fmt.Sprintf("reconcileListeners for service %s: update listener with new config: %v", serviceKey, port)))
				err := bc.updateListener(ctx, lb, port)
				if err != nil {
					return err
				}
			}
			delete(expected, l.Port)
		}
	}
	// delete listener
	if len(deleteList) > 0 {
		klog.Infof(Message(ctx, fmt.Sprintf("reconcileListeners for service %s: delete unexpected listener: %v", serviceKey, deleteList)))
		err = bc.deleteListener(ctx, lb, deleteList)
		if err != nil {
			return err
		}
	}

	// create expected listener
	klog.Infof(Message(ctx, fmt.Sprintf("reconcileListeners for service %s: create expected listener: %v", serviceKey, expected)))
	for _, pl := range expected {
		err := bc.createListener(ctx, lb, pl)
		if err != nil {
			return err
		}
	}

	return nil
}

func (bc *Baiducloud) createListener(ctx context.Context, lb *blb.LoadBalancer, pl PortListener) error {
	switch pl.Protocol {
	case "UDP":
		args := blb.CreateUDPListenerArgs{
			LoadBalancerId:    lb.BlbId,
			ListenerPort:      pl.Port,
			BackendPort:       int(pl.NodePort),
			Scheduler:         "RoundRobin",
			HealthCheckString: "HealthCheck",
		}
		err := bc.clientSet.BLBClient.CreateUDPListener(ctx, &args, bc.getSignOption(ctx))
		if err != nil {
			return err
		}
		return nil
	case "TCP":
		args := blb.CreateTCPListenerArgs{
			LoadBalancerId: lb.BlbId,
			ListenerPort:   pl.Port,
			BackendPort:    int(pl.NodePort),
			Scheduler:      "RoundRobin",
		}
		err := bc.clientSet.BLBClient.CreateTCPListener(ctx, &args, bc.getSignOption(ctx))
		if err != nil {
			return err
		}
		return nil
	case "HTTP":
		// TODO
	case "HTTPS":
		// TODO
	}
	return fmt.Errorf("CreateListener protocol not match: %s", pl.Protocol)
}

func (bc *Baiducloud) updateListener(ctx context.Context, lb *blb.LoadBalancer, pl PortListener) error {
	switch pl.Protocol {
	case "UDP":
		args := blb.UpdateUDPListenerArgs{
			LoadBalancerId:    lb.BlbId,
			ListenerPort:      pl.Port,
			BackendPort:       int(pl.NodePort),
			Scheduler:         "RoundRobin",
			HealthCheckString: "HealthCheck",
		}
		err := bc.clientSet.BLBClient.UpdateUDPListener(ctx, &args, bc.getSignOption(ctx))
		if err != nil {
			return err
		}
		return nil
	case "TCP":
		args := blb.UpdateTCPListenerArgs{
			LoadBalancerId: lb.BlbId,
			ListenerPort:   pl.Port,
			BackendPort:    int(pl.NodePort),
			Scheduler:      "RoundRobin",
		}
		err := bc.clientSet.BLBClient.UpdateTCPListener(ctx, &args, bc.getSignOption(ctx))
		if err != nil {
			return err
		}
		return nil
	case "HTTP":
		// TODO
	case "HTTPS":
		// TODO
	}
	return fmt.Errorf("updateListener protocol not match: %s", pl.Protocol)
}

func (bc *Baiducloud) getAllListeners(ctx context.Context, lb *blb.LoadBalancer) ([]PortListener, error) {
	var allListeners []PortListener

	// add TCPlisteners
	describeTCPListenerArgs := blb.DescribeTCPListenerArgs{
		LoadBalancerId: lb.BlbId,
	}
	tcpListeners, err := bc.clientSet.BLBClient.DescribeTCPListener(ctx, &describeTCPListenerArgs, bc.getSignOption(ctx))
	if err != nil {
		return nil, err
	}
	for _, listener := range tcpListeners {
		allListeners = append(allListeners, PortListener{
			Port:     listener.ListenerPort,
			Protocol: "TCP",
			NodePort: int32(listener.BackendPort),
		})
	}

	// add UDPlisteners
	describeUDPListenerArgs := blb.DescribeUDPListenerArgs{
		LoadBalancerId: lb.BlbId,
	}
	udpListeners, err := bc.clientSet.BLBClient.DescribeUDPListener(ctx, &describeUDPListenerArgs, bc.getSignOption(ctx))
	if err != nil {
		return nil, err
	}
	for _, listener := range udpListeners {
		allListeners = append(allListeners, PortListener{
			Port:     listener.ListenerPort,
			Protocol: "UDP",
			NodePort: int32(listener.BackendPort),
		})
	}

	// TODO: add HTTP,HTTPS
	return allListeners, nil
}

func (bc *Baiducloud) deleteListener(ctx context.Context, lb *blb.LoadBalancer, pl []PortListener) error {
	var portList []int
	for _, l := range pl {
		portList = append(portList, l.Port)
	}
	args := blb.DeleteListenersArgs{
		LoadBalancerId: lb.BlbId,
		PortList:       portList,
	}
	err := bc.clientSet.BLBClient.DeleteListeners(ctx, &args, bc.getSignOption(ctx))
	if err != nil {
		return err
	}
	return nil
}
