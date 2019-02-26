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
	"strings"

	"k8s.io/api/core/v1"

	"k8s.io/cloud-provider-baiducloud/pkg/cloud-sdk/blb"
)

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
