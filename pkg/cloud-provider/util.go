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
	"fmt"
	"net/url"
	"strings"
)

// kubernetesInstanceID represents the id for an instance in the kubernetes API;
// the following form
//  * baidubce:///<zone>/<InstanceId>
//  * baidubce:////<InstanceId>
//  * <InstanceId>
type kubernetesInstanceID string

// bccInstanceID represents the ID of the instance in the API, e.g. i-12345678
// The "traditional" format is "i-12345678"
// A new longer format is also being introduced: "i-12345678abcdef01"
// We should not assume anything about the length or format, though it seems
// reasonable to assume that instances will continue to start with "i-".
type bccInstanceID string

// mapToBCCInstanceID extracts the bccInstanceID from the kubernetesInstanceID
func (name kubernetesInstanceID) mapToBCCInstanceID() (bccInstanceID, error) {
	s := string(name)

	if !strings.HasPrefix(s, "baidubce://") {
		// Assume a bare aws volume id (vol-1234...)
		// Build a URL with an empty host (AZ)
		s = "baidubce://" + "/" + "/" + s
	}
	url, err := url.Parse(s)
	if err != nil {
		return "", fmt.Errorf("Invalid instance name (%s): %v", name, err)
	}
	if url.Scheme != "baidubce" {
		return "", fmt.Errorf("Invalid scheme for BCE instance (%s)", name)
	}

	bccID := ""
	tokens := strings.Split(strings.Trim(url.Path, "/"), "/")
	if len(tokens) == 1 {
		// instanceId
		bccID = tokens[0]
	} else if len(tokens) == 2 {
		// az/instanceId
		bccID = tokens[1]
	}

	// We sanity check the resulting volume; the two known formats are
	// i-12345678 and i-12345678abcdef01
	// TODO: Regex match?
	if bccID == "" || strings.Contains(bccID, "/") || !strings.HasPrefix(bccID, "i-") {
		return "", fmt.Errorf("Invalid format for BCE instance (%s)", name)
	}

	return bccInstanceID(bccID), nil
}
