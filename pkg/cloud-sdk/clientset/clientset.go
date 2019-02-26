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

package clientset

import (
	"fmt"

	"k8s.io/cloud-provider-baiducloud/pkg/cloud-sdk/bcc"
	"k8s.io/cloud-provider-baiducloud/pkg/cloud-sdk/bce"
	"k8s.io/cloud-provider-baiducloud/pkg/cloud-sdk/blb"
	"k8s.io/cloud-provider-baiducloud/pkg/cloud-sdk/cce"
	"k8s.io/cloud-provider-baiducloud/pkg/cloud-sdk/eip"
	"k8s.io/cloud-provider-baiducloud/pkg/cloud-sdk/vpc"
)

// Interface contains all methods of clients
type Interface interface {
	Bcc() *bcc.Client
	Blb() *blb.Client
	Eip() *eip.Client
	Cce() *cce.Client
	Vpc() *vpc.Client
}

// Clientset contains the clients for groups.
type Clientset struct {
	BccClient *bcc.Client
	BlbClient *blb.Client
	EipClient *eip.Client
	CceClient *cce.Client
	VpcClient *vpc.Client
}

// Bcc retrieves the BccClient
func (c *Clientset) Bcc() *bcc.Client {
	if c == nil {
		return nil
	}
	return c.BccClient
}

// Blb retrieves the BccClient
func (c *Clientset) Blb() *blb.Client {
	if c == nil {
		return nil
	}
	return c.BlbClient
}

// Eip retrieves the BccClient
func (c *Clientset) Eip() *eip.Client {
	if c == nil {
		return nil
	}
	return c.EipClient
}

// Cce retrieves the CceClient
func (c *Clientset) Cce() *cce.Client {
	if c == nil {
		return nil
	}
	return c.CceClient
}

// Vpc retrieves the VpcClient
func (c *Clientset) Vpc() *vpc.Client {
	if c == nil {
		return nil
	}
	return c.VpcClient
}

// NewFromConfig create a new Clientset for the given config.
func NewFromConfig(cfg *bce.Config) (*Clientset, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}
	var cs Clientset
	var cceCfg = *cfg
	bccConfig := bcc.NewConfig(cfg)
	blbConfig := blb.NewConfig(cfg)
	eipConfig := eip.NewConfig(cfg)
	// cce endpoint is different
	cceConfig := cce.NewConfig(&cceCfg)
	vpcConfig := vpc.NewConfig(cfg)
	cs.BccClient = bcc.NewClient(bccConfig)
	cs.BlbClient = blb.NewBLBClient(blbConfig)
	cs.EipClient = eip.NewEIPClient(eipConfig)
	cs.CceClient = cce.NewClient(cceConfig)
	cs.VpcClient = vpc.NewVPCClient(vpcConfig)
	return &cs, nil
}
