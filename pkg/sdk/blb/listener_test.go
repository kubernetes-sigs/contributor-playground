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

package blb

import (
	"fmt"
	"testing"

	"k8s.io/cloud-provider-baiducloud/pkg/sdk/util"
)

func TestCreateTCPListener(t *testing.T) {
	args := &CreateTCPListenerArgs{
		LoadBalancerId: "lb-e5b33752",
		ListenerPort:   8088,
		BackendPort:    8080,
		Scheduler:      "LeastConnection",
	}
	err := blbClient.CreateTCPListener(args)
	if err != nil {
		t.Error(err)
	}
}

func TestCreateUDPListener(t *testing.T) {
	args := &CreateUDPListenerArgs{
		LoadBalancerId:    "lb-f5d263e5",
		ListenerPort:      8888,
		BackendPort:       8888,
		Scheduler:         "LeastConnection",
		HealthCheckString: "hello",
	}
	err := blbClient.CreateUDPListener(args)
	if err != nil {
		t.Error(err)
	}
}

func TestCreateHTTPListener(t *testing.T) {
	args := &CreateHTTPListenerArgs{
		LoadBalancerId: "lb-f5d263e5",
		ListenerPort:   8899,
		BackendPort:    8899,
		Scheduler:      "LeastConnection",
	}
	err := blbClient.CreateHTTPListener(args)
	if err != nil {
		t.Error(err)
	}
}

func TestDescribeTCPListener(t *testing.T) {
	args := &DescribeTCPListenerArgs{
		LoadBalancerId: "lb-e5b33752",
		ListenerPort:   8088,
	}
	list, err := blbClient.DescribeTCPListener(args)

	if err != nil {
		t.Error(util.FormatTest("ListInstances", err.Error(), "nil"))
	}
	for _, blb := range list {
		fmt.Println(blb.ListenerPort)
	}
}

func TestDescribeUDPListener(t *testing.T) {
	args := &DescribeUDPListenerArgs{
		LoadBalancerId: "lb-07ab7a1d",
		// ListenerPort:   80,
	}
	list, err := blbClient.DescribeUDPListener(args)

	if err != nil {
		t.Error(util.FormatTest("DescribeUDPListener", err.Error(), "nil"))
	}
	for _, blb := range list {
		fmt.Println(blb.ListenerPort)
	}
}
func TestUpdateTCPListener(t *testing.T) {
	args := &UpdateTCPListenerArgs{
		LoadBalancerId: "lb-e5b33752",
		ListenerPort:   8088,
		BackendPort:    999,
	}
	err := blbClient.UpdateTCPListener(args)
	if err != nil {
		t.Error(err)
	}
}

func TestUpdateUDPListener(t *testing.T) {
	args := &UpdateUDPListenerArgs{
		LoadBalancerId:    "lb-f5d263e5",
		ListenerPort:      8888,
		BackendPort:       8019,
		Scheduler:         "RoundRobin",
		HealthCheckString: "A",
	}
	err := blbClient.UpdateUDPListener(args)
	if err != nil {
		t.Error(err)
	}
}

func TestDeleteListeners(t *testing.T) {
	args := &DeleteListenersArgs{
		LoadBalancerId: "lb-e5b33752",
		PortList:       []int{8088},
	}
	err := blbClient.DeleteListeners(args)
	if err != nil {
		t.Error(err)
	}
}
