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
	"fmt"
	"net/http/httptest"
	"testing"

	"k8s.io/cloud-provider-baiducloud/pkg/sdk/util"
)

func TestListInstances(t *testing.T) {
	ts := httptest.NewServer(ClusterHandler())
	defer ts.Close()
	cceClient.Endpoint = ts.URL
	cceClient.SetDebug(true)
	list, err := cceClient.ListInstances("a")

	if err != nil {
		t.Error(util.FormatTest("ListInstances", err.Error(), "nil"))
	}
	for _, ins := range list {
		fmt.Println(ins.VpcId)
	}
}

func TestScaleUp(t *testing.T) {
	ts := httptest.NewServer(ClusterHandler())
	defer ts.Close()
	cceClient.Endpoint = ts.URL
	cceClient.SetDebug(true)
	args := &ScaleUpClusterArgs{
		ClusterID: "c-NqYwWEhu",
		OrderContent: OrderContent{
			Items: []OrderItem{
				OrderItem{
					Config: BccOrderConfig{
						CPU: 100,
					},
				},
			},
		},
	}
	res, err := cceClient.ScaleUpCluster(args)

	if err != nil {
		t.Fatalf("ScaleUpCluster fail: %v", err)
	}
	if res.ClusterID != "c-NqYwWEhu" {
		t.Fatalf("ScaleUpCluster ClusterID fail")
	}
}

func TestScaleDown(t *testing.T) {
	ts := httptest.NewServer(ClusterHandler())
	defer ts.Close()
	cceClient.Endpoint = ts.URL
	cceClient.SetDebug(true)
	args := &ScaleDownClusterArgs{
		ClusterID: "c-NqYwWEhu",
		AuthCode:  "123456",
	}

	err := cceClient.ScaleDownCluster(args)

	if err != nil {
		t.Fatalf("ScaleDownCluster fail: %v", err)
	}

}
