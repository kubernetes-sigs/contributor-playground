package temp_cce

import (
	"context"
	"time"

	"icode.baidu.com/baidu/jpaas-caas/bce-sdk-go/bce"
	"icode.baidu.com/baidu/jpaas-caas/bce-sdk-go/vpc"
)

const (
	InstanceStatusRunning      InstanceStatus = "RUNNING"
	InstanceStatusCreating     InstanceStatus = "CREATING"
	InstanceStatusDeleting     InstanceStatus = "DELETING"
	InstanceStatusDeleted      InstanceStatus = "DELETED"
	InstanceStatusCreateFailed InstanceStatus = "CREATE_FAILED"
	InstanceStatusError        InstanceStatus = "ERROR"
	InstanceStatusReady        InstanceStatus = "READY"
)

// Interface defines the interface of CCE Client
type Interface interface {
	CreateCluster(ctxd context.Context, args *CreateClusterArgs) (*CreateClusterResponse, error)

	ListClusterNodes(ctx context.Context, clusterID string, option *bce.SignOption) (*ListClusterNodesResponse, error)

	// TODO: Add more
}

// CreateClusterArgs createCluster's args
// TODO: just for cce-ingress-controller's unit test
type CreateClusterArgs struct {
	ClusterName string `json:"clusterName"`
	VPCID       string `json:"vpcId"`
	SubnetID    string `json:"subnetId"`
	VMCount     int    `json:"vmCount"`
}

// CreateClusterResponse createCluster's response
type CreateClusterResponse struct {
	ClusterID string   `json:"clusterUuid"`
	OrderID   []string `json:"orderId"`
}

// Cluster for CCE Cluster
// TODO: just for cce-ingress-controller's unit test
type Cluster struct {
	ClusterID    string  `json:"clusterUuid"`
	InstanceList []*Node `json:"instanceList"`
}

// Node fot CCE Node
type Node struct {
	InstanceID   string       `json:"instanceShortId"`
	InstanceName string       `json:"instanceName"`
	Hostname     string       `json:"hostname"`
	InstanceType InstanceType `json:"instanceType"`

	Status InstanceStatus `json:"status"`

	IP           string `json:"fixIp"`
	EIP          string `json:"eip"`
	EIPBandwidth int    `json:"eipBandwidth"`

	VPCID      string         `json:"vpcId"`
	VPCCIDR    string         `json:"vpcCidr"`
	SubnetID   string         `json:"subnetId"`
	SubnetType vpc.SubnetType `json:"subnetType"`

	AvailableZone string `json:"availableZone"`

	ClusterID string `json:"clusterUuid"`

	CPU         int        `json:"cpu"`     // unit = core
	Memory      int        `json:"memory"`  // unit = GB
	SysDiskSize int        `json:"sysDisk"` // unit = GB
	CDSDisk     []*CDSDisk `json:"cdsDisk,omitempty"`

	RuntimeVersion string `json:"runtimeVersion"`

	PaymentMethod PaymentType `json:"paymentMethod"`

	CreateTime time.Time `json:"createTime"`
	DeleteTime time.Time `json:"deleteTime"`
	ExpireTime time.Time `json:"expireTime"`
}

// InstanceType node instance type
type InstanceType string

// InstanceStatus node instance status
type InstanceStatus string

// PaymentType node payment type
type PaymentType string

// CDSDisk node's cds disk
type CDSDisk struct {
}

// ListClusterNodesResponse the return of ListClusterNodes
type ListClusterNodesResponse struct {
	Marker      string  `json:"marker"`
	IsTruncated bool    `json:"isTruncated"`
	NextMarker  string  `json:"nextMarker"`
	MaxKeys     int     `json:"maxKeys"`
	Nodes       []*Node `json:"nodes"`
}
