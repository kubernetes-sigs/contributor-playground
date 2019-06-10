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
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/golang/glog"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/kubernetes/pkg/cloudprovider"

	"k8s.io/cloud-provider-baiducloud/pkg/cloud-sdk/cce"
)

// Instances returns an instances interface. Also returns true if the interface is supported, false otherwise.
func (bc *Baiducloud) Instances() (cloudprovider.Instances, bool) {
	return bc, true
}

// NodeAddresses returns the addresses of the specified instance.
// TODO(roberthbailey): This currently is only used in such a way that it
// returns the address of the calling instance. We should do a rename to
// make this clearer.
func (bc *Baiducloud) NodeAddresses(ctx context.Context, name types.NodeName) ([]v1.NodeAddress, error) {
	nameStr := string(name)
	nodeIP := net.ParseIP(nameStr)
	if nodeIP == nil {
		return nil, fmt.Errorf("Node name: %s should be an IP address\n", nameStr)
	}
	return []v1.NodeAddress{
		{Type: v1.NodeInternalIP, Address: nameStr},
		{Type: v1.NodeHostName, Address: nameStr},
	}, nil
}

// NodeAddressesByProviderID returns the addresses of the specified instance.
// The instance is specified using the providerID of the node. The
// ProviderID is a unique identifier of the node. This will not be called
// from the node whose nodeaddresses are being queried. i.e. local metadata
// services cannot be used in this method to obtain nodeaddresses
func (bc *Baiducloud) NodeAddressesByProviderID(ctx context.Context, providerID string) ([]v1.NodeAddress, error) {
	splitted := strings.Split(providerID, "//")
	if len(splitted) != 2 {
		return nil, fmt.Errorf("parse ProviderID failed: %v", providerID)
	}
	instanceId := splitted[1]
	var addresses []v1.NodeAddress
	instances, err := bc.clientSet.Cce().ListInstances(bc.ClusterID)
	if err != nil {
		return nil, err
	}
	if len(instances) == 0 {
		return addresses, nil
	}
	for _, instance := range instances {
		if instance.InstanceId == instanceId {
			addresses = append(addresses, v1.NodeAddress{Type: v1.NodeHostName, Address: instance.InternalIP})
			addresses = append(addresses, v1.NodeAddress{Type: v1.NodeInternalIP, Address: instance.InternalIP})
			return addresses, nil
		}
	}
	return nil, fmt.Errorf("NodeAddressesByProviderID faill, not found target providerID: %v", providerID)
}

// InstanceID returns the cloud provider ID of the node with the specified NodeName.
// Note that if the instance does not exist or is no longer running, we must return ("", cloudprovider.InstanceNotFound)
func (bc *Baiducloud) InstanceID(ctx context.Context, name types.NodeName) (string, error) {
	instance, err := bc.getInstanceByNodeName(name)
	if err != nil {
		return "", err
	}
	return instance.InstanceId, nil
}

// InstanceType returns the type of the specified instance.
func (bc *Baiducloud) InstanceType(ctx context.Context, name types.NodeName) (string, error) {
	ins, err := bc.clientSet.Cce().DescribeCluster(bc.ClusterID)
	if err != nil {
		return "", err
	}
	if ins.NodeConfig.GpuCount > 0 {
		return string("GPU"), nil

	}
	return string("BCC"), nil

}

// InstanceTypeByProviderID returns the type of the specified instance.
func (bc *Baiducloud) InstanceTypeByProviderID(ctx context.Context, providerID string) (string, error) {
	ins, err := bc.clientSet.Cce().DescribeCluster(bc.ClusterID)
	if err != nil {
		return "", err
	}
	if ins.NodeConfig.GpuCount > 0 {
		return string("GPU"), nil

	}
	return string("BCC"), nil

}

// AddSSHKeyToAllInstances adds an SSH public key as a legal identity for all instances
// expected format for the key is standard ssh-keygen format: <protocol> <blob>
func (bc *Baiducloud) AddSSHKeyToAllInstances(ctx context.Context, user string, keyData []byte) error {
	return fmt.Errorf("AddSSHKeyToAllInstances not supported")
}

// CurrentNodeName returns the name of the node we are currently running on
// On most clouds (e.g. GCE) this is the hostname, so we provide the hostname
func (bc *Baiducloud) CurrentNodeName(ctx context.Context, hostname string) (types.NodeName, error) {
	// excepting hostname is an ip address
	nodeIP := net.ParseIP(hostname)
	if nodeIP != nil {
		bc.NodeIP = hostname
	}
	return types.NodeName(hostname), nil
}

// InstanceExistsByProviderID returns true if the instance with the given provider id still exists and is running.
// If false is returned with no error, the instance will be immediately deleted by the cloud controller manager.
func (bc *Baiducloud) InstanceExistsByProviderID(ctx context.Context, providerID string) (bool, error) {
	splitted := strings.Split(providerID, "//")
	if len(splitted) != 2 {
		return false, fmt.Errorf("parse ProviderID failed: %v", providerID)
	}
	instance, err := bc.getInstanceByID(string(splitted[1]))
	if err != nil {
		return false, err
	}
	if instance.Status != cce.InstanceStatusRunning {
		return false, fmt.Errorf("target instance %v not running", instance)
	}
	return true, nil
}

// InstanceShutdownByProviderID returns true if the instance is shutdown in cloudprovider
func (bc *Baiducloud) InstanceShutdownByProviderID(ctx context.Context, providerID string) (bool, error) {
	// TODO
	glog.V(2).Infoln("InstanceShutdownByProviderID unimplemented, return false temp")
	return false, nil
}

func (bc *Baiducloud) getInstanceByNodeName(name types.NodeName) (vm cce.CceInstance, err error) {
	nameStr := string(name)
	nodeIP := net.ParseIP(nameStr)
	if nodeIP == nil {
		return vm, fmt.Errorf("Node name: %s should be an IP address\n", nameStr)
	}
	ins, err := bc.clientSet.Cce().ListInstances(bc.ClusterID)
	if err != nil {
		return vm, err
	}
	for _, i := range ins {
		if i.InternalIP == nameStr {
			return i, nil
		}
	}
	return vm, cloudprovider.InstanceNotFound
}

// Returns the instance with the specified ID
func (bc *Baiducloud) getInstanceByID(instanceID string) (*cce.CceInstance, error) {
	ins, err := bc.clientSet.Cce().ListInstances(bc.ClusterID)
	if err != nil {
		return nil, err
	}
	if len(ins) == 0 {
		return nil, cloudprovider.InstanceNotFound
	}
	for _, i := range ins {
		if i.InstanceId == instanceID {
			return &i, nil
		}
	}

	return nil, cloudprovider.InstanceNotFound
}
