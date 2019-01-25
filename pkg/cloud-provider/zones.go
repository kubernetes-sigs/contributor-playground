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
	"strings"

	"k8s.io/apimachinery/pkg/types"
	"k8s.io/kubernetes/pkg/cloudprovider"
)

// Zones returns a zones interface. Also returns true if the interface is supported, false otherwise.
func (bc *Baiducloud) Zones() (cloudprovider.Zones, bool) {
	return bc, true
}

// GetZone returns the Zone containing the current failure zone and locality region that the program is running in
// In most cases, this method is called from the kubelet querying a local metadata service to acquire its zone.
// For the case of external cloud providers, use GetZoneByProviderID or GetZoneByNodeName since GetZone
// can no longer be called from the kubelets.
func (bc *Baiducloud) GetZone(ctx context.Context) (cloudprovider.Zone, error) {
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

// GetZoneByProviderID returns the Zone containing the current zone and locality region of the node specified by providerId
// This method is particularly used in the context of external cloud providers where node initialization must be down
// outside the kubelets.
func (bc *Baiducloud) GetZoneByProviderID(ctx context.Context, providerID string) (cloudprovider.Zone, error) {
	splitted := strings.Split(providerID, "//")
	if len(splitted) != 2 {
		return cloudprovider.Zone{}, fmt.Errorf("parse ProviderID failed: %v", providerID)
	}
	instance, err := bc.getInstanceByID(string(splitted[1]))
	if err != nil {
		return cloudprovider.Zone{}, err
	}

	return cloudprovider.Zone{
		FailureDomain: instance.ZoneName,
		Region:        bc.Region,
	}, nil
}

// GetZoneByNodeName returns the Zone containing the current zone and locality region of the node specified by node name
// This method is particularly used in the context of external cloud providers where node initialization must be down
// outside the kubelets.
func (bc *Baiducloud) GetZoneByNodeName(ctx context.Context, nodeName types.NodeName) (cloudprovider.Zone, error) {
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
