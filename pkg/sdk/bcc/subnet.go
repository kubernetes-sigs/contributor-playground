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

// Subnet deinfe subnet of vpc
type Subnet struct {
	SubnetID    string `json:"subnetId"`
	Name        string `json:"name"`
	ZoneName    string `json:"zoneName"`
	Cidr        string `json:"cidr"`
	VpcID       string `json:"vpcId"`
	SubnetType  string `json:"subnetType"`
	Description string `json:"description"`
}

// CreateSubnetArgs define args create a subnet
type CreateSubnetArgs struct {
	Name        string `json:"name"`
	ZoneName    string `json:"zoneName"`
	Cidr        string `json:"cidr"`
	VpcID       string `json:"vpcId"`
	SubnetType  string `json:"subnetType,omitempty"`
	Description string `json:"description,omitempty"`
}

// CreateSubnetResponse define response of creating a subnet
type CreateSubnetResponse struct {
	SubnetID string `json:"subnetId"`
}

// CreateSubnet create a subnet
// https://cloud.baidu.com/doc/VPC/API.html#.E5.88.9B.E5.BB.BA.E5.AD.90.E7.BD.91
// func (c *Client) CreateSubnet(args *CreateSubnetArgs) (string, error) {

// }
