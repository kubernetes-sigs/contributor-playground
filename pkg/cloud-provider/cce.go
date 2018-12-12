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

package cloud_provider

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"time"

	"github.com/golang/glog"
	"k8s.io/kubernetes/pkg/cloudprovider"
	"k8s.io/kubernetes/pkg/controller"

	"k8s.io/cloud-provider-baiducloud/pkg/sdk/bce"
	"k8s.io/cloud-provider-baiducloud/pkg/sdk/clientset"
)

// ProviderName is the name of this cloud provider.
const ProviderName = "cce"

// CceUserAgent is prefix of http header UserAgent
const CceUserAgent = "cce-k8s:"

type BCECloud struct {
	CloudConfig
	clientSet clientset.Interface
}

type CloudConfig struct {
	ClusterID       string `json:"ClusterId"`
	ClusterName     string `json:"ClusterName"`
	AccessKeyID     string `json:"AccessKeyID"`
	SecretAccessKey string `json:"SecretAccessKey"`
	Region          string `json:"Region"`
	VpcID           string `json:"VpcId"`
	SubnetID        string `json:"SubnetId"`
	MasterID        string `json:"MasterId"`
	Endpoint        string `json:"Endpoint"`
	NodeIP          string `json:"NodeIP"`
	Debug           bool   `json:"Debug"`
}

func init() {
	cloudprovider.RegisterCloudProvider(ProviderName, func(configReader io.Reader) (cloudprovider.Interface, error) {
		return newCloud(configReader)
	})
}

// NewCloud returns a Cloud with initialized clients
func newCloud(configReader io.Reader) (cloudprovider.Interface, error) {
	var cloud BCECloud
	var cloudConfig CloudConfig
	configContents, err := ioutil.ReadAll(configReader)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(configContents, &cloudConfig)
	if err != nil {
		return nil, err
	}
	glog.V(4).Infof("Init CCE cloud with cloudConfig: %v\n", cloudConfig)
	if cloudConfig.MasterID == "" {
		return nil, fmt.Errorf("Cloud config mast have a Master ID\n")
	}
	if cloudConfig.ClusterID == "" {
		return nil, fmt.Errorf("Cloud config mast have a ClusterID\n")
	}
	if cloudConfig.Endpoint == "" {
		return nil, fmt.Errorf("Cloud config mast have a Endpoint\n")
	}
	cred := bce.NewCredentials(cloudConfig.AccessKeyID, cloudConfig.SecretAccessKey)
	bceConfig := bce.NewConfig(cred)
	bceConfig.Region = cloudConfig.Region
	// timeout need to set
	bceConfig.Timeout = 10 * time.Second
	// fix endpoint
	fixEndpoint := cloudConfig.Endpoint + "/internal-api"
	bceConfig.Endpoint = fixEndpoint
	// http request from cce's kubernetes has an useragent header
	// example: useragent: cce-k8s:c-adfdf
	bceConfig.UserAgent = CceUserAgent + cloudConfig.ClusterID
	cloud.CloudConfig = cloudConfig
	cloud.clientSet, err = clientset.NewFromConfig(bceConfig)
	if err != nil {
		return nil, err
	}
	cloud.clientSet.Blb().SetDebug(true)
	cloud.clientSet.Eip().SetDebug(true)
	cloud.clientSet.Bcc().SetDebug(true)
	cloud.clientSet.Cce().SetDebug(true)
	cloud.clientSet.Vpc().SetDebug(true)
	return &cloud, nil
}

// LoadBalancer returns a balancer interface. Also returns true if the interface is supported, false otherwise.
func (bc *BCECloud) LoadBalancer() (cloudprovider.LoadBalancer, bool) {
	return bc, true
}

// Instances returns an instances interface. Also returns true if the interface is supported, false otherwise.
func (bc *BCECloud) Instances() (cloudprovider.Instances, bool) {
	return bc, true
}

// Zones returns a zones interface. Also returns true if the interface is supported, false otherwise.
func (bc *BCECloud) Zones() (cloudprovider.Zones, bool) {
	return bc, true
}

// Clusters returns a clusters interface.  Also returns true if the interface is supported, false otherwise.
func (bc *BCECloud) Clusters() (cloudprovider.Clusters, bool) {
	return nil, false
}

// Routes returns a routes interface along with whether the interface is supported.
func (bc *BCECloud) Routes() (cloudprovider.Routes, bool) {
	return bc, true
}

// ProviderName returns the cloud provider ID.
func (bc *BCECloud) ProviderName() string {
	return ProviderName
}

// HasClusterID returns true if the cluster has a clusterID
func (bc *BCECloud) HasClusterID() bool {
	return true
}

// Initialize passes a Kubernetes clientBuilder interface to the cloud provider
func (bc *BCECloud) Initialize(clientBuilder controller.ControllerClientBuilder) {}
