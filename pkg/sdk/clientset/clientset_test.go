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
	"testing"

	"k8s.io/cloud-provider-baiducloud/pkg/sdk/bce"
)

var credentials, _ = bce.NewCredentialsFromFile("../aksk-test.json")

func TestNewFromConfig(t *testing.T) {
	cfg, err := bce.NewConfigFromFile("../aksk-test.json")
	if err != nil {
		t.Error(err)
	} else {
		cs, err := NewFromConfig(cfg)
		if err != nil {
			t.Error(err)
		} else {
			if cs.Bcc().AccessKeyID != credentials.AccessKeyID {
				t.Error("ak error")
			}
			if cs.Blb().AccessKeyID != credentials.AccessKeyID {
				t.Error("ak error")
			}
			if cs.Eip().AccessKeyID != credentials.AccessKeyID {
				t.Error("ak error")
			}
		}

	}

}
