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

package cce

import (
	"bytes"
	"encoding/json"
	"fmt"

	"k8s.io/cloud-provider-baiducloud/pkg/cloud-sdk/bcc"
	"k8s.io/cloud-provider-baiducloud/pkg/cloud-sdk/bce"
)

const (
	InstanceStatusRunning      string = "RUNNING"
	InstanceStatusCreating     string = "CREATING"
	InstanceStatusDeleting     string = "DELETING"
	InstanceStatusDeleted      string = "DELETED"
	InstanceStatusCreateFailed string = "CREATE_FAILED"
	InstanceStatusError        string = "ERROR"
)

// CceInstance define instance of cce
type CceInstance struct {
	InstanceId            string `json:"id"`
	InstanceName          string `json:"name"`
	InstanceType          string `json:"instancetype"`
	Description           string `json:"desc"`
	Status                string `json:"status"`
	PaymentTiming         string `json:"paymentTiming"`
	CreationTime          string `json:"createTime"`
	ExpireTime            string `json:"expireTime"`
	PublicIP              string `json:"publicIp"`
	InternalIP            string `json:"internalIp"`
	CpuCount              int    `json:"cpu"`
	GpuCount              int    `json:"gpu"`
	MemoryCapacityInGB    int    `json:"memory"`
	LocalDiskSizeInGB     int    `json:"localDiskSizeInGB"`
	ImageId               string `json:"imageId"`
	NetworkCapacityInMbps int    `json:"networkCapacityInMbps"`
	PlacementPolicy       string `json:"placementPolicy"`
	ZoneName              string `json:"zoneName"`
	SubnetId              string `json:"subnetId"`
	VpcId                 string `json:"vpcId"`
}

// ListInstancesResponse define response of cce list
type ListInstancesResponse struct {
	Instances []CceInstance `json:"instances"`
}


// NodeConfig is the config for node
type NodeConfig struct {
	InstanceType int    `json:"instanceType"`
	CPU          int    `json:"cpu,omitempty"`
	Memory       int    `json:"memory,omitempty"`
	GpuCount     int    `json:"gpuCount,omitempty"`
	GpuCard      string `json:"gpuCard,omitempty"`
	DiskSize     int    `json:"diskSize,omitempty"`
}

// CceCluster define cluster of cce
type CceCluster struct {
	ClusterUuid string     `json:"clusterUuid"`
	NodeConfig  NodeConfig `json:"nodeConfig"`
}

// DescribeCluster describe the cluster
func (c *Client) DescribeCluster(clusterID string) (*CceCluster, error) {
	if clusterID == "" {
		return nil, fmt.Errorf("clusterID should not be nil")
	}
	req, err := bce.NewRequest("GET", c.GetURL("/v1/cluster/"+clusterID, nil), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.SendRequest(req, nil)
	if err != nil {
		return nil, err
	}

	bodyContent, err := resp.GetBodyContent()
	if err != nil {
		return nil, err
	}

	var cceCluster CceCluster
	err = json.Unmarshal(bodyContent, &cceCluster)
	if err != nil {
		return nil, err
	}

	return &cceCluster, nil
}

// ListInstances gets all Instances of a cluster.
func (c *Client) ListInstances(clusterID string) ([]CceInstance, error) {
	if clusterID == "" {
		return nil, fmt.Errorf("clusterID should not be nil")
	}
	params := map[string]string{
		"clusterid": clusterID,
	}
	req, err := bce.NewRequest("GET", c.GetURL("/v1/instance", params), nil)

	if err != nil {
		return nil, err
	}

	resp, err := c.SendRequest(req, nil)

	if err != nil {
		return nil, err
	}

	bodyContent, err := resp.GetBodyContent()

	if err != nil {
		return nil, err
	}

	var insList ListInstancesResponse
	err = json.Unmarshal(bodyContent, &insList)

	if err != nil {
		return nil, err
	}

	return insList.Instances, nil
}

// ScaleUpClusterArgs define  args
type ScaleUpClusterArgs struct {
	ClusterID       string              `json:"clusterUuid,omitempty"`
	CdsPreMountInfo bcc.CdsPreMountInfo `json:"cdsPreMountInfo,omitempty"`
	OrderContent    OrderContent        `json:"orderContent,omitempty"`
}

// ScaleUpClusterResponse define  args
type ScaleUpClusterResponse struct {
	ClusterID string   `json:"clusterUuid"`
	OrderID   []string `json:"orderId"`
}

// OrderContent define  bcc order content
type OrderContent struct {
	PaymentMethod []string    `json:"paymentMethod,omitempty"`
	Items         []OrderItem `json:"items,omitempty"`
}

// OrderItem define  bcc order content item
type OrderItem struct {
	Config        interface{} `json:"config,omitempty"`
	PaymentMethod []string    `json:"paymentMethod,omitempty"`
}

// BccOrderConfig define BCC order config
type BccOrderConfig struct {
	// 付费类型，一期只支持postpay
	ProductType string `json:"productType,omitempty"`
	Region      string `json:"region,omitempty"`
	LogicalZone string `json:"logicalZone,omitempty"`
	// 普通BCC
	InstanceType string `json:"instanceType,omitempty"`
	// 这些参数默认就行 容器产品用不到
	FpgaCard string `json:"fpgaCard,omitempty"`
	GpuCard  int    `json:"gpuCard,omitempty"`
	GpuCount int    `json:"gpuCount,omitempty"`

	CPU    int `json:"cpu,omitempty"`
	Memory int `json:"memory,omitempty"`
	// 就一个镜像 ubuntu1604
	ImageType string `json:"imageType,omitempty"`
	// 系统类型
	OsType string `json:"osType,omitempty"`
	// 系统版本
	OsVersion string `json:"osVersion,omitempty"`
	// 系统盘大小
	DiskSize int `json:"diskSize,omitempty"`
	// 暂时为空
	EbsSize []int `json:"ebsSize,omitempty"`
	// 是否需要购买EIP
	IfBuyEip int `json:"ifBuyEip,omitempty"`
	// eip名称
	EipName        string `json:"eipName,omitempty"`
	SubProductType string `json:"subProductType,omitempty"`
	// eip带宽
	BandwidthInMbps int `json:"bandwidthInMbps,omitempty"`

	SubnetUuiD      string `json:"subnetUuid,omitempty"`      // 子网uuid
	SecurityGroupID string `json:"securityGroupId,omitempty"` // 安全组id

	AdminPass        string `json:"adminPass,omitempty"`
	AdminPassConfirm string `json:"adminPassConfirm,omitempty"`
	PurchaseLength   int    `json:"purchaseLength,omitempty"`
	// 购买的虚机个数
	PurchaseNum int `json:"purchaseNum,omitempty"`

	AutoRenewTimeUnit   string                `json:"autoRenewTimeUnit,omitempty"`
	AutoRenewTime       int64                 `json:"autoRenewTime,omitempty"`
	CreateEphemeralList []CreateEphemeralList `json:"createEphemeralList,omitempty"`
	// 是否自动续费 默认即可 后付费不存在这个问题
	AutoRenew bool `json:"autoRenew,omitempty"`
	// 镜像id 用默认即可 固定是ubuntu1604
	ImageID           string `json:"imageId,omitempty"`
	OsName            string `json:"osName,omitempty"`
	SecurityGroupName string `json:"securityGroupName,omitempty"`
	// BCC
	ServiceType string `json:"serviceType,omitempty"`
}

// CreateEphemeralList define storage
type CreateEphemeralList struct {
	// 磁盘存储类型 从页面创建虚机时 看到请求 默认是ssd
	StorageType string `json:"storageType,omitempty"`
	// 磁盘大小
	SizeInGB int `json:"sizeInGB,omitempty"`
}

// CdsOrderConfig define CDS order config
type CdsOrderConfig struct {
	// 付费类型，一期只支持postpay
	productType string `json:"productType,omitempty"`
	// "zoneA"
	logicalZone    string `json:"logicalZone,omitempty"`
	region         string `json:"region,omitempty"`         // "bj"
	purchaseNum    int    `json:"purchaseNum,omitempty"`    // 1
	purchaseLength int    `json:"purchaseLength,omitempty"` // 1
	autoRenewTime  int    `json:"autoRenewTime,omitempty"`  // 0
	// "month"
	autoRenewTimeUnit string               `json:"autoRenewTimeUnit,omitempty"`
	cdsDiskSize       []bcc.DiskSizeConfig `json:"cdsDiskSize,omitempty"`
	// "CDS"
	serviceType string `json:"serviceType,omitempty"`
}

// EipOrderConfig define CDS order config
type EipOrderConfig struct {
	// 付费类型，一期只支持postpay
	ProductType     string `json:"productType,omitempty"`
	BandwidthInMbps int    `json:"bandwidthInMbps,omitempty"` // 1000
	Region          string `json:"region,omitempty"`          // "bj"
	SubProductType  string `json:"subProductType,omitempty"`  // "netraffic",
	// EIP购买数量应该是购买BCC数量的总和
	PurchaseNum       int    `json:"purchaseNum,omitempty"`
	PurchaseLength    int    `json:"purchaseLength,omitempty"`    // 1
	AutoRenewTime     int    `json:"autoRenewTime,omitempty"`     // 0
	AutoRenewTimeUnit string `json:"autoRenewTimeUnit,omitempty"` // "month",
	Name              string `json:"name,omitempty"`              // "kkk"
	ServiceType       string `json:"serviceType,omitempty"`       // "EIP"
}

// ScaleDownClusterArgs define  args
type ScaleDownClusterArgs struct {
	ClusterID string     `json:"clusterUuid"`
	AuthCode  string     `json:"authCode"`
	NodeInfos []NodeInfo `json:"nodeInfo"`
}

// NodeInfo define instanceid
type NodeInfo struct {
	InstanceID string `json:"instanceId"`
}

// ScaleDownClusterResponse define  args
type ScaleDownClusterResponse struct {
	ClusterID string   `json:"clusterUuid"`
	OrderID   []string `json:"orderId"`
}

// ScaleUpCluster scaleup a  cluster
func (c *Client) ScaleUpCluster(args *ScaleUpClusterArgs) (*ScaleUpClusterResponse, error) {
	var params map[string]string
	if args != nil {
		params = map[string]string{
			"clientToken": c.GenerateClientToken(),
			"scalingUp":   "",
		}
	}
	postContent, err := json.Marshal(args)
	if err != nil {
		return nil, err
	}
	req, err := bce.NewRequest("POST", c.GetURL("v1/cluster", params), bytes.NewBuffer(postContent))
	if err != nil {
		return nil, err
	}
	resp, err := c.SendRequest(req, nil)
	if err != nil {
		return nil, err
	}
	bodyContent, err := resp.GetBodyContent()
	if err != nil {
		return nil, err
	}
	var scResp *ScaleUpClusterResponse
	err = json.Unmarshal(bodyContent, &scResp)

	if err != nil {
		return nil, err
	}
	return scResp, nil
}

// ScaleDownCluster scale down a  cluster
func (c *Client) ScaleDownCluster(args *ScaleDownClusterArgs) error {
	var params map[string]string
	if args != nil {
		params = map[string]string{
			"clientToken": c.GenerateClientToken(),
			"scalingDown": "",
		}
	}
	postContent, err := json.Marshal(args)
	if err != nil {
		return err
	}
	req, err := bce.NewRequest("POST", c.GetURL("v1/cluster", params), bytes.NewBuffer(postContent))
	if err != nil {
		return err
	}
	_, err = c.SendRequest(req, nil)
	return err
}
