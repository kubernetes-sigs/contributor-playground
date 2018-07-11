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
	"errors"
	"fmt"
	"net"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/kubernetes/pkg/cloudprovider"

	"k8s.io/cloud-provider-baiducloud/pkg/sdk/cce"
)

// NodeAddresses returns the addresses of the specified instance.
func (bc *BCECloud) NodeAddresses(name types.NodeName) ([]v1.NodeAddress, error) {
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

func (bc *BCECloud) getIPForMachine(name types.NodeName) (string, error) {
	ins, err := bc.clientSet.Cce().ListInstances(bc.ClusterID)
	if err != nil {
		return "", err
	}
	insName := string(name)
	for _, i := range ins {
		if i.InstanceName == insName {
			return i.InternalIP, nil
		}
	}
	return "", cloudprovider.InstanceNotFound
}

func (bc *BCECloud) getVpcID() (string, error) {
	if bc.VpcID == "" {
		ins, err := bc.clientSet.Cce().ListInstances(bc.ClusterID)
		if err != nil {
			return "", err
		}
		if len(ins) > 0 {
			bc.VpcID = ins[0].VpcId
			bc.SubnetID = ins[0].SubnetId
		} else {
			return "", fmt.Errorf("Get vpcid error\n")
		}
	}
	return bc.VpcID, nil
}

// getVirtualMachine get instance info by OPENAPI
func (bc *BCECloud) getVirtualMachine(name types.NodeName) (vm cce.CceInstance, err error) {
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

func (bc *BCECloud) getInstanceByNodeName(name types.NodeName) (vm cce.CceInstance, err error) {
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

// ExternalID returns the cloud provider ID of the specified instance (deprecated).
func (bc *BCECloud) ExternalID(name types.NodeName) (string, error) {
	return bc.InstanceID(name)
}

// InstanceID returns the cloud provider ID of the specified instance.
// Note that if the instance does not exist or is no longer running, we must return ("", cloudprovider.InstanceNotFound)
func (bc *BCECloud) InstanceID(name types.NodeName) (string, error) {
	machine, err := bc.getInstanceByNodeName(name)
	if err != nil {
		return "", err
	}
	return machine.InstanceId, nil
}

// InstanceType returns the type of the specified instance.
// Note that if the instance does not exist or is no longer running, we must return ("", cloudprovider.InstanceNotFound)
// (Implementer Note): This is used by kubelet. Kubelet will label the node. Real log from kubelet:
//       Adding node label from cloud provider: beta.kubernetes.io/instance-type=[value]
func (bc *BCECloud) InstanceType(name types.NodeName) (string, error) {
	return string("BCC"), nil
}

// AddSSHKeyToAllInstances adds an SSH public key as a legal identity for all instances
// expected format for the key is standard ssh-keygen format: <protocol> <blob>
func (bc *BCECloud) AddSSHKeyToAllInstances(user string, keyData []byte) error {
	return fmt.Errorf("not supported")
}

// CurrentNodeName returns the name of the node we are currently running on
// On most clouds (e.g. GCE) this is the hostname, so we provide the hostname
func (bc *BCECloud) CurrentNodeName(hostname string) (types.NodeName, error) {
	// excepting hostname is an ip address
	nodeIP := net.ParseIP(hostname)
	if nodeIP != nil {
		bc.NodeIP = hostname
	}
	return types.NodeName(hostname), nil
}

// NodeAddressesByProviderID returns the node addresses of an instances with the specified unique providerID
// e.g. BCE providerID: baidubce://i-8TokkCDO
func (bc *BCECloud) NodeAddressesByProviderID(providerID string) ([]v1.NodeAddress, error) {
	return nil, errors.New("unimplemented")
}

// InstanceTypeByProviderID returns the cloudprovider instance type of the node with the specified unique providerID
// This method will not be called from the node that is requesting this ID. i.e. metadata service
// and other local methods cannot be used here
func (bc *BCECloud) InstanceTypeByProviderID(providerID string) (string, error) {
	return string("BCC"), nil
}

// InstanceExistsByProviderID returns true if the instance with the given provider id still exists and is running.
// If false is returned with no error, the instance will be immediately deleted by the cloud controller manager.
func (bc *BCECloud) InstanceExistsByProviderID(providerID string) (bool, error) {
	return false, errors.New("unimplemented")
}

// Returns the instance with the specified ID
func (bc *BCECloud) getInstanceByID(instanceID string) (*cce.CceInstance, error) {
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
