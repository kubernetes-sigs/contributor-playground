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

	"k8s.io/api/core/v1"

	"k8s.io/cloud-provider-baiducloud/pkg/cloud-sdk/blb"
)

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