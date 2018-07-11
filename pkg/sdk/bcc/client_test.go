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
	"net/http"
	"os"
	"path"

	"k8s.io/cloud-provider-baiducloud/pkg/sdk/bce"
)

var credentials, _ = bce.NewCredentialsFromFile("../aksk-test.json")

//var bceConfig = bce.NewConfig(credentials)
var bceConfig = &bce.Config{
	Credentials: credentials,
	Checksum:    true,
	Region:      os.Getenv("BOS_REGION"),
}
var bccConfig = NewConfig(bceConfig)
var bccClient = NewClient(bccConfig)

func InstancesHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/v2/instance", func(w http.ResponseWriter, r *http.Request) {
		handleInstanceList(w, r)
	})
	mux.HandleFunc("/v2/instance/", func(w http.ResponseWriter, r *http.Request) {
		handleDescribeInstance(w, r)
	})
	return mux
}

func handleInstanceList(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		response := `{
    "nextMarker": "i-xktdeMSf",
    "marker": "i-IyWRtII7",
    "maxKeys": 1,
    "isTruncated": true,
    "instances": [
        {
            "id": "i-IyWRtII7",
            "createTime": "2015-08-06T13:23:13Z",
            "name": "instance-j93wzbn1",
            "status": "Running",
            "desc": "console",
            "paymentTiming":"Postpaid",
            "expireTime": null,
            "internalIp": "192.168.6.15",
            "publicIp": "-",
            "cpuCount": 1,
            "memoryCapacityInGB": 1,
            "localDiskSizeInGB": 0,
            "networkCapacityInMbps": 1,
            "imageId": "m-3zfBY1Ku",
              "placementPolicy": "default",
              "zoneName": "cn-bj-a"
        }
    ]
}`
		fmt.Fprint(w, response)
	}
}

func handleDescribeInstance(w http.ResponseWriter, r *http.Request) {
	_, id := path.Split(r.URL.Path)
	if id != "i-YufwpQAe" {
		return
	}
	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		response := `{
    "instance": {
        "id": "i-YufwpQAe",
        "createTime": "2015-07-09T10:27:15Z",
        "name": "instance-luz2ef4l-1",
        "status": "Stopped",  
        "desc": "console",
        "paymentTiming":"Postpaid",
        "expireTime": null,
        "internalIp": "192.168.0.25",
        "publicIp": "-",
        "cpuCount": 1,
        "memoryCapacityInGB": 1,
        "localDiskSizeInGB": 0,
        "networkCapacityInMbps": 5,
        "imageId": "m-nky7qeom",
        "placementPolicy": "default",
        "zoneName": "cn-bj-a",
        "subnetId": "sbn-oioiadda",
        "vpcId": "vpc-i80sab3o"
    }
}`
		fmt.Fprint(w, response)
	}
}
