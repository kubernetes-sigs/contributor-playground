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

package util

import (
	"testing"
)

func TestGuessMimeType(t *testing.T) {
	expected := "image/png"
	result := GuessMimeType("examples/test.png")

	if expected != result {
		t.Error(FormatTest("GuessMimeType", result, expected))
	}

	expected = "image/png"
	result = GuessMimeType(".png")

	if expected != result {
		t.Error(FormatTest("GuessMimeType", result, expected))
	}

	expected = "application/octet-stream"
	result = GuessMimeType("png")

	if expected != result {
		t.Error(FormatTest("GuessMimeType", result, expected))
	}

	expected = "application/octet-stream"
	result = GuessMimeType("examples/test")

	if expected != result {
		t.Error(FormatTest("GuessMimeType", result, expected))
	}
}
