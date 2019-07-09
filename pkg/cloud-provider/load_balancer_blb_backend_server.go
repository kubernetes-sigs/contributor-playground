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
	"strings"

	"github.com/golang/glog"
	v1 "k8s.io/api/core/v1"
	"k8s.io/cloud-provider-baiducloud/pkg/cloud-sdk/blb"
)

const BLBMaxRSNum int = 50
const DefaultBLBRSWeight int = 100

func (bc *Baiducloud) reconcileBackendServers(service *v1.Service, nodes []*v1.Node, lb *blb.LoadBalancer) error {
	// extract annotation
	anno, err := ExtractServiceAnnotation(service)
	if err != nil {
		return fmt.Errorf("failed to ExtractServiceAnnotation %s, err: %v", service.Name, err)
	}
	// default rs num of a blb is 50
	targetRsNum := BLBMaxRSNum
	if anno.LoadBalancerRsMaxNum > 0 {
		targetRsNum = anno.LoadBalancerRsMaxNum
	}
	if len(nodes) < targetRsNum {
		targetRsNum = len(nodes)
	}
	glog.Infof("nodes num is %d, target rs num is %d", len(nodes), targetRsNum)

	// turn kube nodes list to backend list
	var candidateBackends []blb.BackendServer
	for _, node := range nodes {
		splitted := strings.Split(node.Spec.ProviderID, "//")
		if len(splitted) != 2 {
			glog.Warningf("node %s has no spec.providerId", node.Name)
			continue
		}
		name := splitted[1]
		candidateBackends = append(candidateBackends, blb.BackendServer{
			InstanceId: name,
		})
	}

	// get all existing rs from lb and change to map
	existingBackends, err := bc.getAllBackendServer(lb)
	if err != nil {
		return err
	}

	rsToAdd, rsToDel, err := mergeBackend(candidateBackends, existingBackends, targetRsNum)
	if err != nil {
		return err
	}
	glog.Infof("find nodes %v to add to BLB %s", rsToAdd, lb.BlbId)
	glog.Infof("find nodes %v to del from BLB %s", rsToDel, lb.BlbId)

	if len(rsToAdd) > 0 {
		args := blb.AddBackendServersArgs{
			LoadBalancerId:    lb.BlbId,
			BackendServerList: rsToAdd,
		}
		err = bc.clientSet.Blb().AddBackendServers(&args)
		if err != nil {
			return err
		}
	}

	if len(rsToDel) > 0 {
		var delList []string
		for _, rs := range rsToDel {
			delList = append(delList, rs.InstanceId)
		}
		args := blb.RemoveBackendServersArgs{
			LoadBalancerId:    lb.BlbId,
			BackendServerList: delList,
		}
		err = bc.clientSet.Blb().RemoveBackendServers(&args)
		if err != nil {
			return err
		}
	}

	return nil
}

/*
case 1:
candidateBackends: ["1", "2", "3"] existingBackends: ["4", "5"] targetBackendsNum: 1
rsToAdd: ["1"]  rsToDel: ["4", "5"]

case 2:
candidateBackends: ["1", "2", "3", "4"] existingBackends: ["4", "5"] targetBackendsNum: 1
rsToAdd: []  rsToDel: ["5"]

case 3:
candidateBackends: ["1", "2", "3", "4"] existingBackends: ["4", "5"] targetBackendsNum: 3
rsToAdd: ["1", "2"]  rsToDel: ["5"]
*/
// candidateBackends contains all ready kubernetes nodes
// existingBackends is real rss(nodes) bound to BLB
func mergeBackend(candidateBackends, existingBackends []blb.BackendServer, targetBackendsNum int) (
	[]blb.BackendServer, []blb.BackendServer, error) {

	if targetBackendsNum > len(candidateBackends) || targetBackendsNum <= 0 {
		return nil, nil, fmt.Errorf("targetBackendsNum %d is invalid", targetBackendsNum)
	}

	// turn existingBackends to map
	existingBackendsMap := make(map[string]int)
	for _, backend := range existingBackends {
		existingBackendsMap[backend.InstanceId] = 0
	}

	// turn candidateBackends to map
	candidateBackendsMap := make(map[string]int)
	for _, backend := range candidateBackends {
		candidateBackendsMap[backend.InstanceId] = 0
	}

	// find rs to delete
	var rsToAdd, rsToDel []blb.BackendServer
	// first find rs that is not in kubernetes to delete from blb
	for insId := range existingBackendsMap {
		_, exist := candidateBackendsMap[insId]
		if !exist {
			rsToDel = append(rsToDel, blb.BackendServer{
				InstanceId: insId,
			})
			delete(existingBackendsMap, insId)
		}
	}

	// then, if number of rs in BLB still > targetBackendsNum, random choose rs in blb to delete
	numToDel := len(existingBackendsMap) - targetBackendsNum
	for insId := range existingBackendsMap {
		if numToDel > 0 {
			rsToDel = append(rsToDel, blb.BackendServer{InstanceId: insId})
			delete(existingBackendsMap, insId)
			numToDel = numToDel - 1
		}
	}

	// find rs to add
	if len(existingBackendsMap) < targetBackendsNum {
		numToAdd := targetBackendsNum - len(existingBackendsMap)
		for insId := range candidateBackendsMap {
			if numToAdd == 0 {
				break
			}
			if _, exist := existingBackendsMap[insId]; !exist {
				rsToAdd = append(rsToAdd, blb.BackendServer{
					InstanceId: insId,
					Weight:     DefaultBLBRSWeight,
				})
				numToAdd = numToAdd - 1
			}
		}
	}
	return rsToAdd, rsToDel, nil
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
