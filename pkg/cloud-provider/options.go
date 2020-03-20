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
	"strconv"

	"k8s.io/klog"
	v1 "k8s.io/api/core/v1"
)

const (
	// ServiceAnnotationLoadBalancerPrefix is the annotation prefix of LoadBalancer
	ServiceAnnotationLoadBalancerPrefix = "service.beta.kubernetes.io/cce-load-balancer-"
	// ServiceAnnotationCceAutoAddLoadBalancerID is the annotation of CCE adding LoadBalancerID
	ServiceAnnotationCceAutoAddLoadBalancerID = ServiceAnnotationLoadBalancerPrefix + "cce-add-id"
	// ServiceAnnotationCceAutoAddEip is the annotation of CCE adding Eip
	ServiceAnnotationCceAutoAddEip = ServiceAnnotationLoadBalancerPrefix + "cce-add-eip"
	// ServiceAnnotationLoadBalancerExistID is the annotation of user assign blbid
	ServiceAnnotationLoadBalancerExistID = ServiceAnnotationLoadBalancerPrefix + "exist-id"

	ServiceAnnotationLoadBalancerId = ServiceAnnotationLoadBalancerPrefix + "id"

	// ServiceAnnotationLoadBalancerInternalVpc is the annotation of LoadBalancerInternalVpc
	ServiceAnnotationLoadBalancerInternalVpc = ServiceAnnotationLoadBalancerPrefix + "internal-vpc"
	// ServiceAnnotationLoadBalancerAllocateVip is the annotation which indicates BLB with a VIP
	ServiceAnnotationLoadBalancerAllocateVip = ServiceAnnotationLoadBalancerPrefix + "allocate-vip"
	//ServiceAnnotationLoadBalancerSubnetID is the annotation which indicates the BCC type subnet the BLB will use
	ServiceAnnotationLoadBalancerSubnetID = ServiceAnnotationLoadBalancerPrefix + "subnet-id"
	// ServiceAnnotationLoadBalancerRsMaxNum is the annotation which set max num of rs of the BLB
	ServiceAnnotationLoadBalancerRsMaxNum = ServiceAnnotationLoadBalancerPrefix + "rs-max-num"
	// ServiceAnnotationLoadBalancerReserveBLB is the annotation which not delete BLB when delete service
	ServiceAnnotationLoadBalancerReserveLB = ServiceAnnotationLoadBalancerPrefix + "reserve-lb"

	ServiceAnnotationLoadBalancerBLBName = ServiceAnnotationLoadBalancerPrefix + "lb-name"

	// TODO:
	// ServiceAnnotationLoadBalancerScheduler is the annotation of load balancer which can be "RoundRobin"/"LeastConnection"/"Hash"
	ServiceAnnotationLoadBalancerScheduler = ServiceAnnotationLoadBalancerPrefix + "scheduler"
	// TODO:
	// ServiceAnnotationLoadBalancerHealthCheckTimeoutInSecond is the annotation of health check timeout, default 3s, [1, 60]
	ServiceAnnotationLoadBalancerHealthCheckTimeoutInSecond = ServiceAnnotationLoadBalancerPrefix + "health-check-timeout-in-second"
	// TODO:
	// ServiceAnnotationLoadBalancerHealthCheckInterval is the annotation of health check interval, default 3s, [1, 10]
	ServiceAnnotationLoadBalancerHealthCheckInterval = ServiceAnnotationLoadBalancerPrefix + "health-check-interval"
	// TODO:
	// ServiceAnnotationLoadBalancerUnhealthyThreshold is the annotation of unhealthy threshold, default 3, [2, 5]
	ServiceAnnotationLoadBalancerUnhealthyThreshold = ServiceAnnotationLoadBalancerPrefix + "unhealthy-threshold"
	// TODO:
	// ServiceAnnotationLoadBalancerHealthyThreshold is the annotation of healthy threshold, default 3, [2, 5]
	ServiceAnnotationLoadBalancerHealthyThreshold = ServiceAnnotationLoadBalancerPrefix + "healthy-threshold"
	// TODO:
	// ServiceAnnotationLoadBalancerHealthCheckString is the annotation of health check string
	ServiceAnnotationLoadBalancerHealthCheckString = ServiceAnnotationLoadBalancerPrefix + "health-check-string"

	// ServiceAnnotationElasticIPPrefix is the annotation prefix of ElasticIP
	ServiceAnnotationElasticIPPrefix = "service.beta.kubernetes.io/cce-elastic-ip-"
	// ServiceAnnotationElasticIPName is the annotation of ElasticIPName
	ServiceAnnotationElasticIPName = ServiceAnnotationElasticIPPrefix + "name"
	// ServiceAnnotationElasticIPPaymentTiming is the annotation of ElasticIPPaymentTiming
	ServiceAnnotationElasticIPPaymentTiming = ServiceAnnotationElasticIPPrefix + "payment-timing"
	// ServiceAnnotationElasticIPBillingMethod is the annotation of ElasticIPBillingMethod
	ServiceAnnotationElasticIPBillingMethod = ServiceAnnotationElasticIPPrefix + "billing-method"
	// ServiceAnnotationElasticIPBandwidthInMbps is the annotation of ElasticIPBandwidthInMbps
	ServiceAnnotationElasticIPBandwidthInMbps = ServiceAnnotationElasticIPPrefix + "bandwidth-in-mbps"
	// ServiceAnnotationElasticIPReservationLength is the annotation of ElasticIPReservationLength
	ServiceAnnotationElasticIPReservationLength = ServiceAnnotationElasticIPPrefix + "reservation-length"
)

const (
	// NodeAnnotationPrefix is the annotation prefix of Node
	NodeAnnotationPrefix = "node.alpha.kubernetes.io/"
	// NodeAnnotationVpcID is the annotation of VpcId on node
	NodeAnnotationVpcID = NodeAnnotationPrefix + "vpc-id"
	// NodeAnnotationVpcRouteTableID is the annotation of VpcRouteTableId on node
	NodeAnnotationVpcRouteTableID = NodeAnnotationPrefix + "vpc-route-table-id"
	// NodeAnnotationVpcRouteRuleID is the annotation of VpcRouteRuleId on node
	NodeAnnotationVpcRouteRuleID = NodeAnnotationPrefix + "vpc-route-rule-id"

	// NodeAnnotationCCMVersion is the version of CCM
	NodeAnnotationCCMVersion = NodeAnnotationPrefix + "ccm-version"

	// NodeAnnotationAdvertiseRoute indicates whether to advertise route to vpc route table
	NodeAnnotationAdvertiseRoute = NodeAnnotationPrefix + "advertise-route"
)

// ServiceAnnotation contains annotations from service
type ServiceAnnotation struct {
	/* BLB */
	LoadBalancerName         string
	LoadBalancerID           string
	CceAutoAddLoadBalancerID string
	CceAutoAddEip            string
	LoadBalancerExistID      string
	LoadBalancerInternalVpc  string
	LoadBalancerAllocateVip  string
	LoadBalancerSubnetID     string
	LoadBalancerScheduler    string
	LoadBalancerRsMaxNum     int
	LoadBalancerReserveLB    string

	LoadBalancerHealthCheckTimeoutInSecond int
	LoadBalancerHealthCheckInterval        int
	LoadBalancerUnhealthyThreshold         int
	LoadBalancerHealthyThreshold           int
	LoadBalancerHealthCheckString          string

	/* EIP */
	ElasticIPName              string
	ElasticIPPaymentTiming     string
	ElasticIPBillingMethod     string
	ElasticIPBandwidthInMbps   int
	ElasticIPReservationLength int
}

// NodeAnnotation contains annotations from node
type NodeAnnotation struct {
	VpcID           string
	VpcRouteTableID string
	VpcRouteRuleID  string
	CCMVersion      string
	AdvertiseRoute  bool
}

// ExtractServiceAnnotation extract annotations from service
func ExtractServiceAnnotation(service *v1.Service) (*ServiceAnnotation, error) {
	klog.V(4).Infof("start to ExtractServiceAnnotation: %v", service.Annotations)
	result := &ServiceAnnotation{}
	annotation := make(map[string]string)
	for k, v := range service.Annotations {
		annotation[k] = v
	}

	loadBalancerName, exist := annotation[ServiceAnnotationLoadBalancerBLBName]
	if exist {
		result.LoadBalancerName = loadBalancerName
	}

	loadBalancerID, exist := annotation[ServiceAnnotationLoadBalancerId]
	if exist {
		result.LoadBalancerID = loadBalancerID
	}

	autoAddLoadBalancerID, exist := annotation[ServiceAnnotationCceAutoAddLoadBalancerID]
	if exist {
		result.CceAutoAddLoadBalancerID = autoAddLoadBalancerID
	}
	cceAddEip, exist := annotation[ServiceAnnotationCceAutoAddEip]
	if exist {
		result.CceAutoAddEip = cceAddEip

	}
	LoadBalancerExistID, exist := annotation[ServiceAnnotationLoadBalancerExistID]
	if exist {
		result.LoadBalancerExistID = LoadBalancerExistID
	}

	loadBalancerInternalVpc, exist := annotation[ServiceAnnotationLoadBalancerInternalVpc]
	if exist {
		result.LoadBalancerInternalVpc = loadBalancerInternalVpc
	}

	loadBalancerAllocateVip, ok := annotation[ServiceAnnotationLoadBalancerAllocateVip]
	if ok {
		result.LoadBalancerAllocateVip = loadBalancerAllocateVip
	}

	loadBalancerSubnetID, ok := annotation[ServiceAnnotationLoadBalancerSubnetID]
	if ok {
		result.LoadBalancerSubnetID = loadBalancerSubnetID
	}

	loadBalancerRsNum, ok := annotation[ServiceAnnotationLoadBalancerRsMaxNum]
	if ok {
		i, err := strconv.Atoi(loadBalancerRsNum)
		if err != nil {
			return nil, fmt.Errorf("ServiceAnnotationLoadBalancerRsMaxNum must be int, err: %v", err)
		} else if i <= 0 || i > blbMaxRSNum {
			return nil, fmt.Errorf("ServiceAnnotationLoadBalancerRsMaxNum must be in (0, 50)")
		} else {
			result.LoadBalancerRsMaxNum = i
		}
	}

	loadBalancerScheduler, ok := annotation[ServiceAnnotationLoadBalancerScheduler]
	if ok {
		result.LoadBalancerScheduler = loadBalancerScheduler
	}

	loadBalancerReserveLB, ok := annotation[ServiceAnnotationLoadBalancerReserveLB]
	if ok {
		result.LoadBalancerReserveLB = loadBalancerReserveLB
	}

	loadBalancerHealthCheckTimeoutInSecond, exist := annotation[ServiceAnnotationLoadBalancerHealthCheckTimeoutInSecond]
	if exist {
		i, err := strconv.Atoi(loadBalancerHealthCheckTimeoutInSecond)
		if err != nil {
			return nil, fmt.Errorf("ServiceAnnotationLoadBalancerHealthCheckTimeoutInSecond must be int")
		} else {
			result.LoadBalancerHealthCheckTimeoutInSecond = i
		}
	}

	loadBalancerHealthCheckInterval, exist := annotation[ServiceAnnotationLoadBalancerHealthCheckInterval]
	if exist {
		i, err := strconv.Atoi(loadBalancerHealthCheckInterval)
		if err != nil {
			return nil, fmt.Errorf("ServiceAnnotationLoadBalancerHealthCheckInterval must be int")
		} else {
			result.LoadBalancerHealthCheckInterval = i
		}
	}

	loadBalancerUnhealthyThreshold, exist := annotation[ServiceAnnotationLoadBalancerUnhealthyThreshold]
	if exist {
		i, err := strconv.Atoi(loadBalancerUnhealthyThreshold)
		if err != nil {
			return nil, fmt.Errorf("ServiceAnnotationLoadBalancerUnhealthyThreshold must be int")
		} else {
			result.LoadBalancerUnhealthyThreshold = i
		}
	}

	loadBalancerHealthyThreshold, exist := annotation[ServiceAnnotationLoadBalancerHealthyThreshold]
	if exist {
		i, err := strconv.Atoi(loadBalancerHealthyThreshold)
		if err != nil {
			return nil, fmt.Errorf("ServiceAnnotationLoadBalancerHealthyThreshold must be int")
		} else {
			result.LoadBalancerHealthyThreshold = i
		}
	}

	loadBalancerHealthCheckString, exist := annotation[ServiceAnnotationLoadBalancerHealthCheckString]
	if exist {
		result.LoadBalancerHealthCheckString = loadBalancerHealthCheckString
	}

	elasticIPName, exist := annotation[ServiceAnnotationElasticIPName]
	if exist {
		result.ElasticIPName = elasticIPName
	}

	elasticIPPaymentTiming, exist := annotation[ServiceAnnotationElasticIPPaymentTiming]
	if exist {
		result.ElasticIPPaymentTiming = elasticIPPaymentTiming
	}

	elasticIPBillingMethod, exist := annotation[ServiceAnnotationElasticIPBillingMethod]
	if exist {
		result.ElasticIPBillingMethod = elasticIPBillingMethod
	}

	elasticIPBandwidthInMbps, exist := annotation[ServiceAnnotationElasticIPBandwidthInMbps]
	if exist {
		i, err := strconv.Atoi(elasticIPBandwidthInMbps)
		if err != nil {
			return nil, fmt.Errorf("ServiceAnnotationElasticIPBandwidthInMbps must be int")
		} else {
			result.ElasticIPBandwidthInMbps = i
		}
	}

	elasticIPReservationLength, exist := annotation[ServiceAnnotationElasticIPReservationLength]
	if exist {
		i, err := strconv.Atoi(elasticIPReservationLength)
		if err != nil {
			return nil, fmt.Errorf("ServiceAnnotationElasticIPReservationLength must be int")
		} else {
			result.ElasticIPReservationLength = i
		}
	}

	return result, nil
}

// ExtractNodeAnnotation extract annotations from node
func ExtractNodeAnnotation(node *v1.Node) (*NodeAnnotation, error) {
	klog.V(4).Infof("start to ExtractNodeAnnotation: %v", node.Annotations)
	result := &NodeAnnotation{}
	annotation := make(map[string]string)
	for k, v := range node.Annotations {
		annotation[k] = v
	}

	vpcID, ok := annotation[NodeAnnotationVpcID]
	if ok {
		result.VpcID = vpcID
	}

	vpcRouteTableID, ok := annotation[NodeAnnotationVpcRouteTableID]
	if ok {
		result.VpcRouteTableID = vpcRouteTableID
	}

	vpcRouteRuleID, ok := annotation[NodeAnnotationVpcRouteRuleID]
	if ok {
		result.VpcRouteRuleID = vpcRouteRuleID
	}

	ccmVersion, ok := annotation[NodeAnnotationCCMVersion]
	if ok {
		result.CCMVersion = ccmVersion
	}

	advertiseRoute, ok := annotation[NodeAnnotationAdvertiseRoute]
	if ok {
		advertise, err := strconv.ParseBool(advertiseRoute)
		if err != nil {
			return nil, fmt.Errorf("NodeAnnotationAdvertiseRoute syntex error: %v", err)
		}
		result.AdvertiseRoute = advertise
	} else {
		result.AdvertiseRoute = true
	}

	return result, nil
}
