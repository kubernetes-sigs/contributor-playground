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
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	cloudprovider "k8s.io/cloud-provider"
	"k8s.io/klog"

	"icode.baidu.com/baidu/jpaas-caas/bce-sdk-go/vpc"
	cce "icode.baidu.com/baidu/jpaas-caas/cloud-provider-baiducloud/pkg/temp-cce"
)

// Routes returns a routes interface along with whether the interface is supported.
func (bc *Baiducloud) Routes() (cloudprovider.Routes, bool) {
	return bc, true
}

// ListRoutes lists all managed routes that belong to the specified clusterName
func (bc *Baiducloud) ListRoutes(ctx context.Context, clusterName string) (routes []*cloudprovider.Route, err error) {
	ctx = context.WithValue(ctx, RequestID, GetRandom())
	startTime := time.Now()
	defer func() {
		klog.Infof(Message(ctx, fmt.Sprintf("Finished ListRoutes (%v)", time.Since(startTime))))
	}()
	rs, err := bc.getVpcRouteTable(ctx)
	if err != nil {
		return nil, err
	}

	// routeTableConflictDetection
	go bc.routeTableConflictDetection(ctx, rs)

	instanceResponse, err := bc.clientSet.CCEClient.ListClusterNodes(ctx, bc.ClusterID, bc.getSignOption(ctx))
	if err != nil {
		return nil, err
	}
	inss := instanceResponse.Nodes
	// Deprecated: there is no need to check node annotaions every cycle
	//vpcID := inss[0].VPCID
	nodename := make(map[string]string)
	for _, ins := range inss {
		if len(ins.Hostname) == 0 {
			nodename[ins.InstanceID] = ins.IP
		} else {
			nodename[ins.InstanceID] = ins.Hostname
		}
	}

	var kubeRoutes []*cloudprovider.Route
	for _, r := range rs {
		// filter instance route
		if r.NexthopType != "custom" {
			continue
		}

		insName, ok := nodename[r.NexthopID]
		if !ok {
			continue
		}
		route := &cloudprovider.Route{
			Name:            r.RouteRuleID,
			DestinationCIDR: r.DestinationAddress,
			TargetNode:      types.NodeName(insName),
		}

		advertiseRoute, err := bc.advertiseRoute(insName)
		if err != nil {
			continue
		}
		// use route.Blackhole to mark this route to be deleted
		route.Blackhole = !advertiseRoute

		// no need to check err
		// Deprecated: there is no need to check node annotaions every cycle
		//_ = bc.ensureRouteInfoToNode(insName, vpcID, r.RouteTableID, r.RouteRuleID)

		kubeRoutes = append(kubeRoutes, route)
	}
	return kubeRoutes, nil
}

// CreateRoute creates the described managed route
// route.Name will be ignored, although the cloud-provider may use nameHint
// to create a more user-meaningful name.
func (bc *Baiducloud) CreateRoute(ctx context.Context, clusterName string, nameHint string, kubeRoute *cloudprovider.Route) error {
	startTime := time.Now()
	defer func() {
		klog.Infof(Message(ctx, fmt.Sprintf("Finished CreateRoutes %+v (%v)", kubeRoute, time.Since(startTime))))
	}()
	klog.Infof(Message(ctx, fmt.Sprintf("CreateRoute: creating route. instance=%v cidr=%v", kubeRoute.TargetNode, kubeRoute.DestinationCIDR)))

	advertiseRoute, err := bc.advertiseRoute(string(kubeRoute.TargetNode))
	if err != nil {
		return err
	}

	if !advertiseRoute {
		klog.V(3).Infof("Node %s has annotation not to advertise route", string(kubeRoute.TargetNode))
		return nil
	}

	insID, err := bc.checkClusterNode(ctx, kubeRoute)
	if err != nil {
		return err
	}

	routeRule, err := bc.ensureCreateRule(ctx, kubeRoute, insID)
	if err != nil {
		return err
	}

	vpcID, err := bc.getVpcID(ctx)
	if err != nil {
		return err
	}
	err = bc.ensureRouteInfoToNode(string(kubeRoute.TargetNode), vpcID, routeRule.RouteTableID, routeRule.RouteRuleID)
	if err != nil {
		return err
	}

	klog.Infof(Message(ctx, fmt.Sprintf("CreateRoute for cluster: %v node: %v success", clusterName, kubeRoute.TargetNode)))
	return nil
}

// DeleteRoute deletes the specified managed route
// Route should be as returned by ListRoutes
func (bc *Baiducloud) DeleteRoute(ctx context.Context, clusterName string, kubeRoute *cloudprovider.Route) error {
	startTime := time.Now()
	defer func() {
		klog.Infof(Message(ctx, fmt.Sprintf("Finished DeleteRoutes %v (%v)", kubeRoute, time.Since(startTime))))
	}()
	klog.Infof(Message(ctx, fmt.Sprintf("DeleteRoute: instance=%q cidr=%q", kubeRoute.TargetNode, kubeRoute.DestinationCIDR)))
	vpcTable, err := bc.getVpcRouteTable(ctx)
	if err != nil {
		klog.V(3).Infof("getVpcRouteTable error %s", err.Error())
		return err
	}
	for _, vr := range vpcTable {
		if vr.DestinationAddress == kubeRoute.DestinationCIDR && vr.SourceAddress == "0.0.0.0/0" {
			klog.V(3).Infof("DeleteRoute: DestinationAddress is %s .", vr.DestinationAddress)
			err := bc.clientSet.VPCClient.DeleteRoute(ctx, vr.RouteRuleID, bc.getSignOption(ctx))
			if err != nil {
				klog.V(3).Infof("Delete VPC route error %s", err.Error())
				return err
			}
		}
	}

	klog.Infof(Message(ctx, fmt.Sprintf("DeleteRoute: success, instance=%q cidr=%q", kubeRoute.TargetNode, kubeRoute.DestinationCIDR)))

	return nil
}

func (bc *Baiducloud) getVpcRouteTable(ctx context.Context) ([]vpc.RouteRule, error) {
	vpcid, err := bc.getVpcID(ctx)
	if err != nil {
		return nil, err
	}
	args := vpc.ListRouteArgs{
		VpcID: vpcid,
	}
	rs, err := bc.clientSet.VPCClient.ListRouteTable(ctx, &args, bc.getSignOption(ctx))
	if err != nil {
		return nil, err
	}

	if len(rs) < 1 {
		return nil, fmt.Errorf("VPC route length error: length is : %d", len(rs))
	}
	return rs, nil
}

// ensureRouteInfoToNode add below annotation to node
// node.alpha.kubernetes.io/vpc-id: "vpc-xxx"
// node.alpha.kubernetes.io/vpc-route-table-id: "rt-xxx"
// node.alpha.kubernetes.io/vpc-route-rule-id: "rr-xxx"
func (bc *Baiducloud) ensureRouteInfoToNode(nodeName, vpcID, vpcRouteTableID, vpcRouteRuleID string) error {
	curNode, err := bc.kubeClient.CoreV1().Nodes().Get(nodeName, metav1.GetOptions{})
	if err != nil {
		// skip unreachable node
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}
	if curNode.Annotations == nil {
		curNode.Annotations = make(map[string]string)
	}
	nodeAnnotation, err := ExtractNodeAnnotation(curNode)
	if err != nil {
		return err
	}

	isChanged := false
	if nodeAnnotation.VpcID != vpcID {
		curNode.Annotations[NodeAnnotationVpcID] = vpcID
		isChanged = true
	}
	if nodeAnnotation.VpcRouteTableID != vpcRouteTableID {
		curNode.Annotations[NodeAnnotationVpcRouteTableID] = vpcRouteTableID
		isChanged = true
	}
	if nodeAnnotation.VpcRouteRuleID != vpcRouteRuleID {
		curNode.Annotations[NodeAnnotationVpcRouteRuleID] = vpcRouteRuleID
		isChanged = true
	}
	if nodeAnnotation.CCMVersion != CCMVersion {
		curNode.Annotations[NodeAnnotationCCMVersion] = CCMVersion
		isChanged = true
	}
	if isChanged {
		j, err := json.Marshal(curNode.Annotations)
		if err != nil {
			return err
		}
		data := []byte(fmt.Sprintf(`{"metadata":{"annotations":%s}}`, j))
		_, err = bc.kubeClient.CoreV1().Nodes().Patch(nodeName, types.StrategicMergePatchType, data)
		if err != nil {
			klog.V(4).Infof("Patch error!")
			return err
		}
	}
	return nil
}

func (bc *Baiducloud) getVpcID(ctx context.Context) (string, error) {
	if bc.VpcID == "" {
		instanceResponse, err := bc.clientSet.CCEClient.ListClusterNodes(ctx, bc.ClusterID, bc.getSignOption(ctx))
		if err != nil {
			return "", err
		}
		ins := instanceResponse.Nodes
		if len(ins) > 0 {
			bc.VpcID = ins[0].VPCID
			bc.SubnetID = ins[0].SubnetID
		} else {
			return "", fmt.Errorf("Get vpcid error\n ")
		}
	}
	return bc.VpcID, nil
}

func (bc *Baiducloud) routeTableConflictDetection(ctx context.Context, rs []vpc.RouteRule) {
	klog.Infof(Message(ctx, fmt.Sprintf("start routeTable conflict detection.")))
	if len(rs) < 2 {
		return
	}
	var cceRR []vpc.RouteRule
	var otherRR []vpc.RouteRule
	for i := 0; i < len(rs); i++ {
		if strings.Contains(rs[i].Description, "auto generated by cce") {
			cceRR = append(cceRR, rs[i])
		} else {
			otherRR = append(otherRR, rs[i])
		}
	}
	if len(cceRR) == 0 || len(otherRR) == 0 {
		return
	}
	for i := 0; i < len(otherRR); i++ {
		for j := 0; j < len(cceRR); j++ {
			if bc.isConflict(otherRR[i], cceRR[j]) {
				klog.V(4).Infof("RouteTable conflict detected, custom routeRule %v may conflict with cce routeRule %v", otherRR[i], cceRR[j])
				if bc.eventRecorder != nil {
					bc.eventRecorder.Eventf(&v1.ObjectReference{
						Kind: "VPC",
						Name: "RouteTableConflict",
					}, v1.EventTypeWarning, "RouteTableConflictDetection", "RouteTable conflict detected, custom routeRule %v may conflict with cce routeRule %v", otherRR[i], cceRR[j])
				}
			}
		}
	}
}

func (bc *Baiducloud) isConflict(otherRR vpc.RouteRule, cceRR vpc.RouteRule) bool {
	// rule 1: 用户路由的目标网段 是 CCE实例路由的目标网段 的子网
	{
		_, cidrBlock, err := net.ParseCIDR("0.0.0.0/0")
		if err != nil {
			klog.Errorf("cidrBlock net.ParseCIDR failed: %v", err)
			return false
		}
		_, cceCidr, err := net.ParseCIDR(cceRR.DestinationAddress)
		if err != nil {
			klog.Errorf("cceRR %v net.ParseCIDR failed: %v", cceRR, err)
			return false
		}
		_, otherCidr, err := net.ParseCIDR(otherRR.DestinationAddress)
		if err != nil {
			klog.Errorf("otherRR %v net.ParseCIDR failed: %v", otherRR, err)
			return false
		}
		err = VerifyNoOverlap([]*net.IPNet{cceCidr, otherCidr}, cidrBlock)
		if err != nil {
			klog.Errorf("VerifyNoOverlap: %v", err)
			return true
		}
		return false
	}

	// rule 2: TODO
}

func (bc *Baiducloud) advertiseRoute(nodename string) (bool, error) {

	// check node resource in k8s has advertise route annotation, if is false, not create route
	curNode, err := bc.kubeClient.CoreV1().Nodes().Get(nodename, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			return false, nil
		}
		return true, err
	}

	if curNode.Annotations == nil {
		curNode.Annotations = make(map[string]string)
	}
	nodeAnnotation, err := ExtractNodeAnnotation(curNode)
	if err != nil {
		return true, err
	}
	return nodeAnnotation.AdvertiseRoute, nil
}

func (bc *Baiducloud) checkClusterNode(ctx context.Context, kubeRoute *cloudprovider.Route) (string, error) {
	var node *cce.Node
	instanceResponse, err := bc.clientSet.CCEClient.ListClusterNodes(ctx, bc.ClusterID, bc.getSignOption(ctx))
	if err != nil {
		return "", err
	}

	for _, ins := range instanceResponse.Nodes {
		if ins.Hostname == string(kubeRoute.TargetNode) || ins.IP == string(kubeRoute.TargetNode) {
			node = ins
			break
		}
	}

	if node == nil {
		klog.Errorf(Message(ctx, fmt.Sprintf("InstanceId not found for k8s node %s, not create route", string(kubeRoute.TargetNode))))
		return "", fmt.Errorf("InstanceId not found for k8s node %s, create route failed", string(kubeRoute.TargetNode))
	}

	if node.Status == cce.InstanceStatusCreateFailed || node.Status == cce.InstanceStatusDeleted ||
		node.Status == cce.InstanceStatusDeleting || node.Status == cce.InstanceStatusError {
		klog.V(3).Infof("No need to create route, instance has a wrong status: %s", node.Status)
		return "", nil
	}

	return node.InstanceID, nil
}

func (bc *Baiducloud) ensureCreateRule(ctx context.Context, kubeRoute *cloudprovider.Route, insID string) (vpc.RouteRule, error) {
	vpcRoutes, err := bc.getVpcRouteTable(ctx)
	if err != nil {
		return vpc.RouteRule{}, err
	}

	for _, vr := range vpcRoutes {
		if vr.DestinationAddress == kubeRoute.DestinationCIDR && vr.SourceAddress == "0.0.0.0/0" && vr.NexthopID == insID {
			klog.Infof(Message(ctx, fmt.Sprintf("route rule %+v already exist", vr)))
			return vr, nil
		}
		if vr.DestinationAddress == kubeRoute.DestinationCIDR && vr.SourceAddress == "0.0.0.0/0" {
			err := bc.clientSet.VPCClient.DeleteRoute(ctx, vr.RouteRuleID, bc.getSignOption(ctx))
			if err != nil {
				klog.Infof("Delete VPC route error %s", err)
				return vpc.RouteRule{}, err
			}
		}
	}

	args := vpc.CreateRouteRuleArgs{
		RouteTableID:       vpcRoutes[0].RouteTableID,
		NexthopType:        "custom",
		Description:        fmt.Sprintf("auto generated by cce:%s", bc.ClusterID),
		DestinationAddress: kubeRoute.DestinationCIDR,
		SourceAddress:      "0.0.0.0/0",
		NexthopID:          insID,
	}
	klog.Infof(Message(ctx, fmt.Sprintf("CreateRoute: create args %v", args)))
	routeRuleID, err := bc.clientSet.VPCClient.CreateRouteRule(ctx, &args, bc.getSignOption(ctx))
	if err != nil {
		return vpc.RouteRule{}, err
	}
	return vpc.RouteRule{
		RouteTableID:       args.RouteTableID,
		NexthopType:        args.NexthopType,
		Description:        args.Description,
		DestinationAddress: args.DestinationAddress,
		SourceAddress:      args.SourceAddress,
		NexthopID:          args.NexthopID,
		RouteRuleID:        routeRuleID,
	}, nil
}
