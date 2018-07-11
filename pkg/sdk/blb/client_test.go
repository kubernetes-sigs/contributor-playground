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

package blb

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"

	"k8s.io/cloud-provider-baiducloud/pkg/sdk/bce"
	"github.com/gorilla/mux"
)

var (
	testHTTPServer *httptest.Server
	blbClient      *Client
)

func init() {
	var credentials, _ = bce.NewCredentialsFromFile("../aksk-test.json")

	//var bceConfig = bce.NewConfig(credentials)
	var bceConfig = &bce.Config{
		Credentials: credentials,
		Checksum:    true,
		Region:      os.Getenv("BOS_REGION"),
	}
	var bccConfig = NewConfig(bceConfig)
	blbClient = NewBLBClient(bccConfig)
	r := mux.NewRouter()
	// loadbalancer
	r.HandleFunc("/v1/blb", handleGetBLB).Methods("GET")
	r.HandleFunc("/v1/blb/{blbid}", handleDeleteBLB).Methods("DELETE")
	r.HandleFunc("/v1/blb/{blbid}", handleUpdateBLB).Methods("PUT")
	r.HandleFunc("/v1/blb", handleCreateBLB).Methods("POST")

	// backendserver
	r.HandleFunc("/v1/blb/{blbid}/backendserver", handleCreateBackendServers).Methods("POST")
	r.HandleFunc("/v1/blb/{blbid}/backendserver", handleDescribeBackendServers).Methods("GET")
	r.HandleFunc("/v1/blb/{blbid}/backendserver", handleDescribeBackendServers).Methods("PUT")
	r.HandleFunc("/v1/blb/{blbid}/backendserver", handleDescribeBackendServers).Methods("DELETE")

	// listener
	r.HandleFunc("/v1/blb/{blbid}/TCPlistener", handleCreateTCPlistener).Methods("POST")
	r.HandleFunc("/v1/blb/{blbid}/UDPlistener", handleCreateUDPlistener).Methods("POST")
	r.HandleFunc("/v1/blb/{blbid}/HTTPlistener", handleCreateHTTPListener).Methods("POST")
	r.HandleFunc("/v1/blb/{blbid}/TCPlistener", handleDescribeTCPListeners).Methods("GET")
	r.HandleFunc("/v1/blb/{blbid}/UDPlistener", handleDescribeUDPListeners).Methods("GET")

	r.HandleFunc("/v1/blb/{blbid}/UDPlistener", handleUpdateUDPListener).Methods("PUT")
	r.HandleFunc("/v1/blb/{blbid}/TCPlistener", handleUpdateTCPListener).Methods("PUT")
	r.HandleFunc("/v1/blb/{blbid}/listener", handleDeleteListeners).Methods("PUT")
	// start
	testHTTPServer = httptest.NewServer(r)
	blbClient.Endpoint = testHTTPServer.URL
}

func handleGetBLB(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := `{
    "blbList":[
        {
            "blbId":"lb-a7e5zPPk",
            "status":"available",
            "name":"test-blb",
            "desc":"用于生产环境",
            "vpcId":"vpc-fy6vdykpwkqb",
            "address":"10.32.249.113"
        },
        {
            "blbId": "lb-gj5gVpeq",
            "status":"available",
            "name": "nametest",
            "desc": "用于测试环境",
            "vpcId":"vpc-a8n5p6kybbx4",
            "address": "10.32.251.4"
        }
    ],
    "marker": "blb-0A20F971",
    "nextMarker": "blb-0A20FB09",
    "isTruncated": true,
    "maxKeys": 2
}`
	fmt.Fprint(w, response)
}

func handleCreateBLB(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	args := &CreateLoadBalancerArgs{}
	json.Unmarshal(body, args)
	if args.Name != expectCreateBLB.Name {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	response := `{
    "blbId": "lb-BLuOPSLZ",
    "name": "blb-for-test",
    "desc": "",
    "address": "10.32.251.93"
}`
	fmt.Fprint(w, response)
}

func handleUpdateBLB(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	blbid := vars["blbid"]
	if blbid == "lb-e5b33752" {
		w.WriteHeader(200)
		return
	}
	w.WriteHeader(400)
}

func handleDeleteBLB(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	blbid := vars["blbid"]
	if blbid == "lb-426fad2b" {
		w.WriteHeader(200)
	} else {
		w.WriteHeader(400)
	}

}

func handleCreateBackendServers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	args := &AddBackendServersArgs{}
	json.Unmarshal(body, args)
	fmt.Println(args.LoadBalancerId)
	vars := mux.Vars(r)
	blbid := vars["blbid"]
	if blbid != "lb-e5b33752" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func handleDescribeBackendServers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	blbid := vars["blbid"]
	if blbid != "lb-e5b33752" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	response := `{
    "backendServerList": [
        {
            "instanceId": "i-YfAibl4A",
            "weight": 50
        }
    ],
    "marker": "rs-0A6BE9BB",
    "nextMarker": null,
    "isTruncated": false,
    "maxKeys": 1000
}`
	fmt.Fprint(w, response)
}

func handleUpdateBackendServers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	blbid := vars["blbid"]
	if blbid == "lb-e5b33752" {
		w.WriteHeader(200)
		return
	}
	w.WriteHeader(400)
}

func handleDeleteBackendServers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	blbid := vars["blbid"]
	if blbid == "lb-e5b33752" {
		w.WriteHeader(200)
	} else {
		w.WriteHeader(400)
	}

}

func handleCreateTCPlistener(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	args := &CreateTCPListenerArgs{}
	json.Unmarshal(body, args)
	vars := mux.Vars(r)
	blbid := vars["blbid"]
	if blbid != "lb-e5b33752" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if args.ListenerPort != 8088 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func handleCreateUDPlistener(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	args := &CreateUDPListenerArgs{}
	json.Unmarshal(body, args)
	vars := mux.Vars(r)
	blbid := vars["blbid"]
	if blbid != "lb-f5d263e5" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if args.ListenerPort != 8888 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func handleCreateHTTPListener(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	args := &CreateUDPListenerArgs{}
	json.Unmarshal(body, args)
	vars := mux.Vars(r)
	blbid := vars["blbid"]
	if blbid != "lb-f5d263e5" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if args.ListenerPort != 8899 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func handleDescribeTCPListeners(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	blbid := vars["blbid"]
	if blbid != "lb-e5b33752" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	response := `{
    "listenerList": [
        {
            "listenerPort": 80,
            "backendPort": 80,
            "scheduler": "RoundRobin",
            "healthCheckTimeoutInSecond": 3,
            "healthCheckInterval": 3,
            "unhealthyThreshold": 3,
            "healthyThreshold": 3
        },
        {
            "listenerPort": 88,
            "backendPort": 88,
            "scheduler": "RoundRobin",
            "healthCheckTimeoutInSecond": 2,
            "healthCheckInterval": 4,
            "unhealthyThreshold": 3,
            "healthyThreshold": 3
        }
    ],
    "marker": "listener-0050",
    "nextMarker": null,
    "isTruncated": false,
    "maxKeys": 2
}`
	fmt.Fprint(w, response)
}

func handleDescribeUDPListeners(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	blbid := vars["blbid"]
	if blbid != "lb-07ab7a1d" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	response := `{
    "listenerList": [
        {
            "listenerPort": 80,
            "backendPort": 80,
            "scheduler": "RoundRobin",
            "healthCheckTimeoutInSecond": 3,
            "healthCheckInterval": 3,
            "unhealthyThreshold": 3,
            "healthyThreshold": 3
        },
        {
            "listenerPort": 88,
            "backendPort": 88,
            "scheduler": "RoundRobin",
            "healthCheckTimeoutInSecond": 2,
            "healthCheckInterval": 4,
            "unhealthyThreshold": 3,
            "healthyThreshold": 3
        }
    ],
    "marker": "listener-0050",
    "nextMarker": null,
    "isTruncated": false,
    "maxKeys": 2
}`
	fmt.Fprint(w, response)
}

func handleUpdateUDPListener(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	blbid := vars["blbid"]
	if blbid == "lb-f5d263e5" {
		w.WriteHeader(200)
		return
	}
	w.WriteHeader(400)
}

func handleUpdateTCPListener(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	blbid := vars["blbid"]
	if blbid == "lb-e5b33752" {
		w.WriteHeader(200)
		return
	}
	w.WriteHeader(400)
}

func handleDeleteListeners(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	blbid := vars["blbid"]
	if blbid == "lb-e5b33752" {
		w.WriteHeader(200)
	} else {
		w.WriteHeader(400)
	}

}
