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
)

func TestListRouteTable(t *testing.T) {
	// ts := httptest.NewServer(EipHandler())
	// defer ts.Close()
	// eipClient.Endpoint = ts.URL
	// eips, err := eipClient.GetEips(nil)
	bccClient.Endpoint = "bcc.bce-api.baidu.com"
	bccClient.SetDebug(true)
	args := ListRouteArgs{
		VpcID: "vpc-4yprmi7pjscp",
	}
	rs, err := bccClient.ListRouteTable(&args)
	if err != nil {
		t.Error(err)
	}
	for _, r := range rs {
		// fmt.Printf("%+v", r)
		fmt.Println(r.RouteRuleID)
		// fmt.Println(r.NexthopID)
	}
}

func TestCreateRouteRule(t *testing.T) {
	// bccClient.Endpoint = "bcc.bce-api.baidu.com"
	args := CreateRouteRuleArgs{
		RouteTableID:       "rt-wc5rd05e8fzs",
		SourceAddress:      "0.0.0.0/0",
		DestinationAddress: "172.17.112.0/24",
		NexthopID:          "i-ddUE7vVn",
		NexthopType:        "custom",
		Description:        "a",
	}
	id, err := bccClient.CreateRouteRule(&args)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(id)
}

func TestDeleteRoute(t *testing.T) {
	// bccClient.Endpoint = "bcc.bce-api.baidu.com"
	err := bccClient.DeleteRoute("rr-p9dbxrxdcsrh")
	if err != nil {
		t.Error(err)
	}
}
