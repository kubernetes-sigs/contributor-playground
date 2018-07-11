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

func TestAddBackendServers(t *testing.T) {
	args := &AddBackendServersArgs{
		LoadBalancerId: "lb-e5b33752",
		BackendServerList: []BackendServer{
			BackendServer{
				InstanceId: "i-YWIy3FQx",
				Weight:     50,
			},
			BackendServer{
				InstanceId: "i-vfBlsqNG",
				Weight:     50,
			},
		},
	}
	err := blbClient.AddBackendServers(args)
	if err != nil {
		t.Error(err)
	}
}

func TestDescribeBackendServers(t *testing.T) {
	args := &DescribeBackendServersArgs{
		LoadBalancerId: "lb-e5b33752",
	}
	list, err := blbClient.DescribeBackendServers(args)
	if err != nil {
		fmt.Println(err)
		t.Error(util.FormatTest("DescribeBackendServers", err.Error(), "nil"))
	}
	for _, blb := range list {
		fmt.Println(blb)
	}
}

func TestUpdateBackendServers(t *testing.T) {
	args := &UpdateBackendServersArgs{
		LoadBalancerId: "lb-e5b33752",
		BackendServerList: []BackendServer{BackendServer{
			InstanceId: "i-vfBlsqNG",
			Weight:     99,
		}},
	}
	err := blbClient.UpdateBackendServers(args)
	if err != nil {
		t.Error(err)
	}
}

func TestRemoveBackendServers(t *testing.T) {
	args := &RemoveBackendServersArgs{
		LoadBalancerId:    "lb-e5b33752",
		BackendServerList: []string{"i-vfBlsqNG", "i-vfBlsqNG"},
	}

	err := blbClient.RemoveBackendServers(args)
	if err != nil {
		t.Error(err)
	}
}
