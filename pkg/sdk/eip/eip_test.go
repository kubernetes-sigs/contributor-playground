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

package eip

import (
	"testing"
)

func TestCreateEip(t *testing.T) {
	bill := &Billing{
		PaymentTiming: "Postpaid",
		BillingMethod: "ByTraffic",
	}
	args := &CreateEipArgs{
		BandwidthInMbps: 998,
		Billing:         bill,
		Name:            "k8stestcgy",
	}
	ip, err := eipClient.CreateEip(args)
	if err != nil {
		t.Error(err)
	}
	if ip != "180.181.3.133" {
		t.Error("ip error")
	}
}

var expectResizeEip = &ResizeEipArgs{
	BandwidthInMbps: 111,
	Ip:              "180.76.242.209",
}

func TestResizeEip(t *testing.T) {
	err := eipClient.ResizeEip(expectResizeEip)
	if err != nil {
		t.Error(err)
	}
}

var expectBindEip = &BindEipArgs{
	Ip:           "180.76.247.62",
	InstanceType: "BCC",
	InstanceId:   "i-VAEyKKTh",
}
var expectUnbindEip = &EipArgs{
	Ip: "180.76.154.83",
}

func TestBindEip(t *testing.T) {
	err := eipClient.BindEip(expectBindEip)
	if err != nil {
		t.Error(err)
	}
}

func TestUnbindEip(t *testing.T) {
	err := eipClient.UnbindEip(expectUnbindEip)
	if err != nil {
		t.Error(err)
	}
}

func TestDeleteEip(t *testing.T) {
	err := eipClient.DeleteEip(expectUnbindEip)
	if err != nil {
		t.Error(err)
	}
}
func TestGetEips(t *testing.T) {
	eips, err := eipClient.GetEips(nil)
	if err != nil {
		t.Error(err)
	}
	for _, eip := range eips {
		if eip.Eip != "180.181.3.133" && eip.Eip != "180.181.3.134" {
			t.Fatal("eip errpr")
		}
	}
}
