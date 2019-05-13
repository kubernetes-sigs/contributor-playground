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
	"math/rand"
	"net"
	"time"

	"github.com/golang/glog"

	"k8s.io/cloud-provider-baiducloud/pkg/cloud-sdk/vpc"
)

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
func (bc *Baiducloud) getSubnetForBLB(serviceAnnotation *ServiceAnnotation) (string, string, error) {
	// get VPC id
	vpcId, err := bc.getVpcID()
	if err != nil {
		return "", "", fmt.Errorf("Can't get VPC for BLB: %v\n", err)
	}
	// user set subnet id in annotation
	subnetId := serviceAnnotation.LoadBalancerSubnetId
	if subnetId != "" {
		glog.V(3).Infof("Find subnetId %v in annotation for BLB", subnetId)
		subnet, err := bc.clientSet.Vpc().DescribeSubnet(subnetId)
		if err != nil {
			return "", "", fmt.Errorf("Can't get subnet with subnetId %v in annotation: %v\n", subnetId, err)
		}
		if subnet.SubnetType != "BCC" {
			return "", "", fmt.Errorf("Can't use subnet with subnetId %v in annotation: subnet type is not BCC\n", subnetId)
		}
		glog.V(3).Infof("Use subnet with id %v in annotation for BLB", subnetId)
		return vpcId, subnetId, nil
	}
	// get subnet id from instance
	ins, err := bc.clientSet.Cce().ListInstances(bc.ClusterID)
	if err != nil {
		return "", "", err
	}
	if len(ins) == 0 {
		return "", "", fmt.Errorf("getSubnetForBLB failed since instance num is zero")
	}
	// random select a VM to choose subnet
	randomVM := ins[rand.Intn(len(ins))]
	subnetId = randomVM.SubnetId

	// check subnet
	subnet, err := bc.clientSet.Vpc().DescribeSubnet(subnetId)
	if err != nil {
		return "", "", fmt.Errorf("DescribeSubnet failed: %v", err)
	}
	if subnet.SubnetType == "BCC" {
		return subnet.VpcID, subnetId, nil
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
	tryCount := 0
	for { // loop
		tryCount++
		if tryCount > 10 {
			return "", "", fmt.Errorf("CreateSubnet failed after 10 retries")
		}
		_, cidr, err := net.ParseCIDR(currentCidr)
		if err != nil {
			return "", "", fmt.Errorf("ParseCIDR failed: %v", err)
		}
		mask, _ := cidr.Mask.Size()
		nextCidr, notExist := NextSubnet(cidr, mask)
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
