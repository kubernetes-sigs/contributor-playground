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
	"math/rand"

	v1 "k8s.io/api/core/v1"
	"k8s.io/klog"

	"icode.baidu.com/baidu/jpaas-caas/bce-sdk-go/vpc"
)

func (bc *Baiducloud) getVpcInfoForBLB(ctx context.Context, service *v1.Service) (string, string, error) {
	// get VPC id
	vpcID, err := bc.getVpcID(ctx)
	if err != nil {
		return "", "", fmt.Errorf("Can't get VPC for BLB: %v\n ", err)
	}
	// user set subnet id in annotation
	subnetID, ok := service.Annotations[ServiceAnnotationLoadBalancerSubnetID]
	if ok {
		if subnetID != "" {
			klog.V(3).Infof("Find subnetId %v in annotation for BLB", subnetID)
			subnetIsTypeBCC, err := bc.subnetIsTypeBCC(ctx, subnetID)
			if err != nil {
				return "", "", err
			}
			if !subnetIsTypeBCC {
				return "", "", fmt.Errorf("SubnetId %v in annotation is not type BCC", subnetID)
			}
			klog.V(3).Infof("Use subnet with id %v in annotation for BLB", subnetID)
			return vpcID, subnetID, nil
		}
	}

	// get subnet id from instance
	instanceResponse, err := bc.clientSet.CCEClient.ListClusterNodes(ctx, bc.ClusterID, bc.getSignOption(ctx))
	if err != nil {
		return "", "", err
	}
	ins := instanceResponse.Nodes
	if len(ins) == 0 {
		return "", "", fmt.Errorf("getVpcInfoForBLB failed since instance num is zero")
	}
	// random select a VM to choose subnet
	randomVM := ins[rand.Intn(len(ins))]
	subnetID = randomVM.SubnetID

	// check subnet
	subnetIsTypeBCC, err := bc.subnetIsTypeBCC(ctx, subnetID)
	if err != nil {
		return "", "", fmt.Errorf("DescribeSubnet failed: %v", err)
	}
	if subnetIsTypeBCC {
		return vpcID, subnetID, nil
	}

	// get subnet list and choose preferred one
	listSubnetArgs := &vpc.ListSubnetArgs{
		VPCID: vpcID,
	}
	subnets, err := bc.clientSet.VPCClient.ListSubnet(ctx, listSubnetArgs, bc.getSignOption(ctx))
	if err != nil {
		return "", "", fmt.Errorf("ListSubnet failed: %v", err)
	}
	for _, subnet := range subnets {
		if subnet.Name == "系统预定义子网" {
			return subnet.VPCID, subnet.SubnetID, nil
		}
		if subnet.Name == "CCE-Reserve" {
			return subnet.VPCID, subnet.SubnetID, nil
		}
	}

	return "", "", fmt.Errorf("no suitable subnet found for BLB")
}

func (bc *Baiducloud) subnetIsTypeBCC(ctx context.Context, subnetID string) (bool, error) {
	subnet, err := bc.clientSet.VPCClient.DescribeSubnet(ctx, subnetID, bc.getSignOption(ctx))
	if err != nil {
		return false, fmt.Errorf("DescribeSubnet failed: %v", err)
	}
	return subnet.SubnetType == "BCC", nil
}
