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
	"testing"

	"fmt"

	"k8s.io/cloud-provider-baiducloud/pkg/sdk/util"
)

func TestDescribeInstance(t *testing.T) {
	// ts := httptest.NewServer(InstancesHandler())
	// defer ts.Close()
	bccClient.SetDebug(true)
	// bccClient.Endpoint = ts.URL
	// ins, err := bccClient.DescribeInstance("i-YufwpQAe", nil)
	ins, err := bccClient.DescribeInstance("i-7VUJvwqR", nil)
	if err != nil {
		t.Error(util.FormatTest("ListInstances", err.Error(), "nil"))
	}
	if ins.InstanceName != "instance-luz2ef4l-1" {
		t.Error("name error!")
	}
}

func TestListInstances(t *testing.T) {
	// ts := httptest.NewServer(InstancesHandler())
	// defer ts.Close()
	// bccClient.Endpoint = ts.URL
	// bccClient.Endpoint = "bcc.bce-api.baidu.com"
	bccClient.SetDebug(true)
	list, err := bccClient.ListInstances(nil)

	if err != nil {
		t.Error(util.FormatTest("ListInstances", err.Error(), "nil"))
	}
	for _, ins := range list {
		fmt.Println(ins.VpcId)
		if ins.InstanceId != "i-IyWRtII7" {
			// t.Error("instanceId error")
		}
	}
}
