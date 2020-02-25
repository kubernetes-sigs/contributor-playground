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

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	cloudprovider "k8s.io/cloud-provider"
	"k8s.io/klog"

	cce "icode.baidu.com/baidu/jpaas-caas/cloud-provider-baiducloud/pkg/temp-cce"
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
	var IP string

	// TODO if hostname is x.x.x.x ?
	nameStr := string(name)
	nodeIP := net.ParseIP(nameStr)
	if nodeIP != nil {
		IP = nameStr
	} else {
		instance, err := bc.getInstanceByNodeName(ctx, name)
		if err != nil {
			return nil, err
		}
		IP = instance.IP
	}

	return []v1.NodeAddress{
		{Type: v1.NodeInternalIP, Address: IP},
		{Type: v1.NodeHostName, Address: IP},
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
	instanceID := splitted[1]
	var addresses []v1.NodeAddress
	instanceResponse, err := bc.clientSet.CCEClient.ListClusterNodes(ctx, bc.ClusterID, bc.getSignOption(ctx))
	if err != nil {
		return nil, err
	}
	instances := instanceResponse.Nodes
	if len(instances) == 0 {
		return addresses, nil
	}
	for _, instance := range instances {
		if instance.InstanceID == instanceID {
			addresses = append(addresses, v1.NodeAddress{Type: v1.NodeHostName, Address: instance.IP})
			addresses = append(addresses, v1.NodeAddress{Type: v1.NodeInternalIP, Address: instance.IP})
			return addresses, nil
		}
	}
	return nil, fmt.Errorf("NodeAddressesByProviderID faill, not found target providerID: %v", providerID)
}

// InstanceID returns the cloud provider ID of the node with the specified NodeName.
// Note that if the instance does not exist or is no longer running, we must return ("", cloudprovider.InstanceNotFound)
func (bc *Baiducloud) InstanceID(ctx context.Context, name types.NodeName) (string, error) {
	instance, err := bc.getInstanceByNodeName(ctx, name)
	if err != nil {
		return "", err
	}
	return instance.InstanceID, nil
}

// InstanceType returns the type of the specified instance.
func (bc *Baiducloud) InstanceType(ctx context.Context, name types.NodeName) (string, error) {
	ins, err := bc.getInstanceByNodeName(ctx, name)
	if err != nil {
		return "", err
	}
	if ins.InstanceType == "9" {
		return string("GPU"), nil

	}
	return string("BCC"), nil

}

// InstanceTypeByProviderID returns the type of the specified instance.
func (bc *Baiducloud) InstanceTypeByProviderID(ctx context.Context, providerID string) (string, error) {
	ins, err := bc.getInstanceByProviderID(ctx, providerID)
	if err != nil {
		return "", err
	}
	if ins.InstanceType == "9" {
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
	if len(hostname) != 0 {
		bc.NodeName = hostname
	}
	return types.NodeName(hostname), nil
}

// InstanceExistsByProviderID returns true if the instance with the given provider id still exists and is running.
// If false is returned with no error, the instance will be immediately deleted by the cloud controller manager.
func (bc *Baiducloud) InstanceExistsByProviderID(ctx context.Context, providerID string) (bool, error) {
	instance, err := bc.getInstanceByProviderID(ctx, providerID)
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
	klog.V(2).Infoln("InstanceShutdownByProviderID unimplemented, return false temp")
	return false, nil
}

func (bc *Baiducloud) getInstanceByNodeName(ctx context.Context, name types.NodeName) (vm *cce.Node, err error) {
	nameStr := string(name)
	if len(nameStr) == 0 {
		return vm, fmt.Errorf("Node name: %s is nil\n ", nameStr)
	}
	instanceResponse, err := bc.clientSet.CCEClient.ListClusterNodes(ctx, bc.ClusterID, bc.getSignOption(ctx))
	if err != nil {
		return vm, err
	}
	ins := instanceResponse.Nodes
	for _, i := range ins {
		// nodeName can be a ip or a hostname
		if i.Hostname == nameStr || i.IP == nameStr {
			return i, nil
		}
	}
	return vm, cloudprovider.InstanceNotFound
}

// Returns the instance with the providerID
func (bc *Baiducloud) getInstanceByProviderID(ctx context.Context, providerID string) (*cce.Node, error) {
	// when node.spec.providerID is not set, providerID is only instanceID, not start with cce://
	if !strings.HasPrefix(providerID, bc.ProviderName()+"://") {
		providerID = bc.ProviderName() + "://" + providerID
	}
	splitted := strings.Split(providerID, "//")
	if len(splitted) != 2 {
		return nil, fmt.Errorf("parse ProviderID failed: %v", providerID)
	}
	instanceID := splitted[1]
	instanceResponse, err := bc.clientSet.CCEClient.ListClusterNodes(ctx, bc.ClusterID, bc.getSignOption(ctx))
	if err != nil {
		return nil, err
	}
	ins := instanceResponse.Nodes
	if len(ins) == 0 {
		return nil, cloudprovider.InstanceNotFound
	}
	for _, i := range ins {
		if i.InstanceID == instanceID {
			return i, nil
		}
	}

	return nil, cloudprovider.InstanceNotFound
}
