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
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/kubernetes/pkg/cloudprovider"
)

// GetZone returns the Zone containing the current failure zone and locality region that the program is running in
func (bc *BCECloud) GetZone() (cloudprovider.Zone, error) {
	zone := cloudprovider.Zone{
		FailureDomain: "unknow",
		Region:        bc.Region,
	}
	if bc.NodeIP != "" {
		ins, err := bc.getInstanceByNodeName(types.NodeName(bc.NodeIP))
		// ins, err := bc.getVirtualMachine(types.NodeName(bc.NodeIP))
		if err != nil {
			return zone, err
		}
		zone.FailureDomain = ins.ZoneName
	}
	return zone, nil
}

// GetZoneByProviderID implements Zones.GetZoneByProviderID
// This is particularly useful in external cloud providers where the kubelet
// does not initialize node data.
func (bc *BCECloud) GetZoneByProviderID(providerID string) (cloudprovider.Zone, error) {
	instanceID, err := kubernetesInstanceID(providerID).mapToBCCInstanceID()
	if err != nil {
		return cloudprovider.Zone{}, err
	}
	instance, err := bc.getInstanceByID(string(instanceID))
	if err != nil {
		return cloudprovider.Zone{}, err
	}

	return cloudprovider.Zone{
		FailureDomain: instance.ZoneName,
		Region:        bc.Region,
	}, nil
}

// GetZoneByNodeName implements Zones.GetZoneByNodeName
// This is particularly useful in external cloud providers where the kubelet
// does not initialize node data.
func (bc *BCECloud) GetZoneByNodeName(nodeName types.NodeName) (cloudprovider.Zone, error) {
	instance, err := bc.getInstanceByNodeName(nodeName)
	if err != nil {
		return cloudprovider.Zone{}, err
	}
	zone := cloudprovider.Zone{
		FailureDomain: instance.ZoneName,
		Region:        bc.Region,
	}

	return zone, nil
}
