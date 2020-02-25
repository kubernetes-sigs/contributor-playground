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

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"

	"icode.baidu.com/baidu/jpaas-caas/bce-sdk-go/blb"
)

const blbMaxRSNum int = 50
const defaultBLBRSWeight int = 100

func (bc *Baiducloud) reconcileBackendServers(ctx context.Context, clusterName string, service *v1.Service, nodes []*v1.Node) error {
	startTime := time.Now()
	serviceKey := fmt.Sprintf("%s/%s", service.Namespace, service.Name)
	defer func() {
		klog.Infof(Message(ctx, fmt.Sprintf("Finished reconcileBackendServers for service %q (%v)", serviceKey, time.Since(startTime))))
	}()
	lb, exist, err := bc.getServiceAssociatedBLB(ctx, clusterName, service)
	if err != nil {
		return err
	}
	if !exist {
		return fmt.Errorf("failed to reconcileBackendServers: lb not exist")
	}

	if service.Spec.ExternalTrafficPolicy == v1.ServiceExternalTrafficPolicyTypeLocal {
		nodes, err = bc.getServiceAssociatedNodes(ctx, service)
		if err != nil {
			return err
		}
		if len(nodes) == 0 {
			klog.Infof(Message(ctx, fmt.Sprintf("service %s has no nodes to add to lb, maybe has no pod, do nothing", serviceKey)))
			return nil
		}
		klog.Infof(Message(ctx, fmt.Sprintf("externalTrafficPolicy of service %s is Local, nodes is %+v", serviceKey, nodes)))
	}

	// extract annotation
	anno, err := ExtractServiceAnnotation(service)
	if err != nil {
		return fmt.Errorf("failed to ExtractServiceAnnotation %s, err: %v", service.Name, err)
	}
	// default rs num of a blb is 50
	targetRsNum := blbMaxRSNum
	if anno.LoadBalancerRsMaxNum > 0 {
		targetRsNum = anno.LoadBalancerRsMaxNum
	}
	// turn kube nodes list to backend list
	var candidateBackends []blb.BackendServer
	for _, node := range nodes {
		splitted := strings.Split(node.Spec.ProviderID, "//")
		if len(splitted) != 2 {
			msg := fmt.Sprintf("node %s has no spec.providerId", node.Name)
			bc.eventRecorder.Eventf(node, v1.EventTypeNormal, "Node has no providerID", msg)
			klog.Warningf(Message(ctx, msg))
			continue
		}
		name := splitted[1]
		candidateBackends = append(candidateBackends, blb.BackendServer{
			InstanceId: name,
		})
	}
	if len(candidateBackends) < targetRsNum {
		targetRsNum = len(candidateBackends)
	}
	klog.Infof(Message(ctx, fmt.Sprintf("nodes num is %d, target rs num is %d", len(candidateBackends), targetRsNum)))
	// get all existing rs from lb and change to map
	existingBackends, err := bc.getAllBackendServer(ctx, lb)
	if err != nil {
		return err
	}

	rsToAdd, rsToDel, err := mergeBackend(candidateBackends, existingBackends, targetRsNum)
	if err != nil {
		return err
	}
	klog.Infof(Message(ctx, fmt.Sprintf("find nodes %v to add to BLB %s for service %s", rsToAdd, lb.BlbId, serviceKey)))
	klog.Infof(Message(ctx, fmt.Sprintf("find nodes %v to del from BLB %s for service %s", rsToDel, lb.BlbId, serviceKey)))

	if len(rsToAdd) > 0 {
		args := blb.AddBackendServersArgs{
			LoadBalancerId:    lb.BlbId,
			BackendServerList: rsToAdd,
		}
		err = bc.clientSet.BLBClient.AddBackendServers(ctx, &args, bc.getSignOption(ctx))
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
		err = bc.clientSet.BLBClient.RemoveBackendServers(ctx, &args, bc.getSignOption(ctx))
		if err != nil {
			return err
		}
	}

	return nil
}

func (bc *Baiducloud) getServiceAssociatedNodes(ctx context.Context, service *v1.Service) ([]*v1.Node, error) {
	ep, err := bc.kubeClient.CoreV1().Endpoints(service.Namespace).Get(service.Name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	if len(ep.Subsets) == 0 {
		klog.Infof(Message(ctx, fmt.Sprintf("Endpoints %s/%s has no subsets", ep.Namespace, ep.Name)))
		return nil, nil
	}
	nodeMap := make(map[string]string, 0)
	for _, addr := range ep.Subsets[0].Addresses {
		nodeMap[*addr.NodeName] = *addr.NodeName
	}

	for _, addr := range ep.Subsets[0].NotReadyAddresses {
		nodeMap[*addr.NodeName] = *addr.NodeName
	}

	all_nodes, err := bc.kubeClient.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	result := make([]*v1.Node, 0)
	for _, node := range all_nodes.Items {
		if _, exist := nodeMap[node.Name]; exist {
			n := node.DeepCopy()
			klog.Infof(Message(ctx, fmt.Sprintf("Node is %s", n.Name)))
			result = append(result, n)
		}
	}

	return result, nil
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
	for insID := range existingBackendsMap {
		_, exist := candidateBackendsMap[insID]
		if !exist {
			rsToDel = append(rsToDel, blb.BackendServer{
				InstanceId: insID,
			})
			delete(existingBackendsMap, insID)
		}
	}

	// then, if number of rs in BLB still > targetBackendsNum, random choose rs in blb to delete
	numToDel := len(existingBackendsMap) - targetBackendsNum
	for insID := range existingBackendsMap {
		if numToDel > 0 {
			rsToDel = append(rsToDel, blb.BackendServer{InstanceId: insID})
			delete(existingBackendsMap, insID)
			numToDel = numToDel - 1
		}
	}

	// find rs to add
	if len(existingBackendsMap) < targetBackendsNum {
		numToAdd := targetBackendsNum - len(existingBackendsMap)
		for insID := range candidateBackendsMap {
			if numToAdd == 0 {
				break
			}
			if _, exist := existingBackendsMap[insID]; !exist {
				rsToAdd = append(rsToAdd, blb.BackendServer{
					InstanceId: insID,
					Weight:     defaultBLBRSWeight,
				})
				numToAdd = numToAdd - 1
			}
		}
	}
	return rsToAdd, rsToDel, nil
}

func (bc *Baiducloud) getAllBackendServer(ctx context.Context, lb *blb.LoadBalancer) ([]blb.BackendServer, error) {
	args := blb.DescribeBackendServersArgs{
		LoadBalancerId: lb.BlbId,
	}
	bs, err := bc.clientSet.BLBClient.DescribeBackendServers(ctx, &args, bc.getSignOption(ctx))
	if err != nil {
		return nil, err
	}
	return bs, nil
}

func (bc *Baiducloud) deleteAllBackendServers(ctx context.Context, lb *blb.LoadBalancer) error {
	allServers, err := bc.getAllBackendServer(ctx, lb)
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
		err = bc.clientSet.BLBClient.RemoveBackendServers(ctx, &args, bc.getSignOption(ctx))
		if err != nil {
			return err
		}
	}
	return nil
}
