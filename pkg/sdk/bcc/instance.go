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

package bcc

import (
	"encoding/json"

	"k8s.io/cloud-provider-baiducloud/pkg/sdk/bce"
)

const (
	InstanceStatusRunning            string = "Running"
	InstanceStatusStarting           string = "Starting"
	InstanceStatusStopping           string = "Stopping"
	InstanceStatusStopped            string = "Stopped"
	InstanceStatusDeleted            string = "Deleted"
	InstanceStatusScaling            string = "Scaling"
	InstanceStatusExpired            string = "Expired"
	InstanceStatusError              string = "Error"
	InstanceStatusSnapshotProcessing string = "SnapshotProcessing"
	InstanceStatusImageProcessing    string = "ImageProcessing"
)

// Instance define instance model
type Instance struct {
	InstanceId            string `json:"id"`
	InstanceName          string `json:"name"`
	Description           string `json:"desc"`
	Status                string `json:"status"`
	PaymentTiming         string `json:"paymentTiming"`
	CreationTime          string `json:"createTime"`
	ExpireTime            string `json:"expireTime"`
	PublicIP              string `json:"publicIp"`
	InternalIP            string `json:"internalIp"`
	CpuCount              int    `json:"cpuCount"`
	GpuCount              int    `json:"gpuCount"`
	MemoryCapacityInGB    int    `json:"memoryCapacityInGB"`
	localDiskSizeInGB     int    `json:"localDiskSizeInGB"`
	ImageId               string `json:"imageId"`
	NetworkCapacityInMbps int    `json:"networkCapacityInMbps"`
	PlacementPolicy       string `json:"placementPolicy"`
	ZoneName              string `json:"zoneName"`
	SubnetId              string `json:"subnetId"`
	VpcId                 string `json:"vpcId"`
}

type ListInstancesResponse struct {
	Marker      string     `json:"marker"`
	IsTruncated bool       `json:"isTruncated"`
	NextMarker  string     `json:"nextMarker"`
	MaxKeys     int        `json:"maxKeys"`
	Instances   []Instance `json:"instances"`
}

type GetInstanceResponse struct {
	Ins Instance `json:"instance"`
}

// ListInstances gets all Instances.
func (c *Client) ListInstances(option *bce.SignOption) ([]Instance, error) {

	req, err := bce.NewRequest("GET", c.GetURL("v2/instance", nil), nil)

	if err != nil {
		return nil, err
	}

	resp, err := c.SendRequest(req, option)

	if err != nil {
		return nil, err
	}

	bodyContent, err := resp.GetBodyContent()

	if err != nil {
		return nil, err
	}

	var insList *ListInstancesResponse
	err = json.Unmarshal(bodyContent, &insList)

	if err != nil {
		return nil, err
	}

	return insList.Instances, nil
}

// DescribeInstance describe a instance
func (c *Client) DescribeInstance(instanceID string, option *bce.SignOption) (*Instance, error) {

	req, err := bce.NewRequest("GET", c.GetURL("v2/instance"+"/"+instanceID, nil), nil)

	if err != nil {
		return nil, err
	}

	resp, err := c.SendRequest(req, option)

	if err != nil {
		return nil, err
	}

	bodyContent, err := resp.GetBodyContent()

	if err != nil {
		return nil, err
	}

	var ins GetInstanceResponse
	err = json.Unmarshal(bodyContent, &ins)

	if err != nil {
		return nil, err
	}

	return &ins.Ins, nil
}
