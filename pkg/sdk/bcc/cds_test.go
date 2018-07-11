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
	"fmt"
	"testing"

	"k8s.io/cloud-provider-baiducloud/pkg/sdk/bce"
	"k8s.io/cloud-provider-baiducloud/pkg/sdk/util"
)

func TestDeleteVolume(t *testing.T) {
	// ts := httptest.NewServer(InstancesHandler())
	// defer ts.Close()
	// bccClient.SetDebug(true)
	// bccClient.Endpoint = ts.URL
	err := bccClient.DeleteVolume("v-MK288vVC")
	if err != nil {
		t.Error(util.FormatTest("DeleteVolume", err.Error(), "nil"))
	}
}

var expectBill = &bce.Billing{
	PaymentTiming: "Postpaid",
	BillingMethod: "ByTraffic",
}
var expectCreateVolumeArgs = &CreateVolumeArgs{
	PurchaseCount: 1,
	Billing:       expectBill,
	StorageType:   STORAGE_TYPE_STD1,
	CdsSizeInGB:   10,
}

func TestCreateVolumes(t *testing.T) {
	// ts := httptest.NewServer(EipHandler())
	// defer ts.Close()
	// eipClient.Endpoint = ts.URL
	_, err := bccClient.CreateVolumes(expectCreateVolumeArgs)
	if err != nil {
		t.Error(err)
	}
}

func TestGetVolumeList(t *testing.T) {
	// ts := httptest.NewServer(EipHandler())
	// defer ts.Close()
	// eipClient.Endpoint = ts.URL
	list, err := bccClient.GetVolumeList(nil)
	if err != nil {
		t.Error(err)
	}
	for _, v := range list {
		fmt.Println(v.Id)
	}
}

func TestDescribeVolume(t *testing.T) {
	// ts := httptest.NewServer(InstancesHandler())
	// defer ts.Close()
	// bccClient.Endpoint = ts.URL
	ins, err := bccClient.DescribeVolume("v-31wjHWIU")
	if err != nil {
		t.Error(util.FormatTest("TestDescribeVolume", err.Error(), "nil"))
	}
	fmt.Println(ins.Id)
}

var expectAttach = &AttachCDSVolumeArgs{
	VolumeId:   "v-JCvK3cpI",
	InstanceId: "i-NN0KeMyw",
}

func TestAttachCDSVolume(t *testing.T) {
	// ts := httptest.NewServer(EipHandler())
	// defer ts.Close()
	// eipClient.Endpoint = ts.URL
	bccClient.SetDebug(true)
	att, err := bccClient.AttachCDSVolume(expectAttach)
	if err != nil {
		t.Error(err)
	} else {
		fmt.Println(att.Device)
	}

}

func TestDetachCDSVolume(t *testing.T) {
	bccClient.SetDebug(true)
	err := bccClient.DetachCDSVolume(expectAttach)
	if err != nil {
		t.Error(err)
	}
}

func TestDeleteCDSVolume(t *testing.T) {
	bccClient.SetDebug(true)
	err := bccClient.DeleteCDS("v-JCvK3cpI")
	if err != nil {
		t.Error(err)
	}
}
