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

package bce

import (
	"encoding/json"
	"net/http"
	"strconv"
	"testing"

	"k8s.io/cloud-provider-baiducloud/pkg/sdk/util"
)

func TestError(t *testing.T) {
	bceError := &Error{
		StatusCode: 500,
		Code:       "StatusInternalServerError",
		Message:    "failed",
		RequestID:  "123",
	}

	result := bceError.Error()
	expected := "Error Message: \"failed\", Error Code: \"StatusInternalServerError\", Status Code: 500, Request Id: \"123\""

	if result != expected {
		t.Error(util.FormatTest("Error", result, expected))
	}
}

func TestBuildError(t *testing.T) {
	resp := &Response{BodyContent: []byte{}}
	err := buildError(resp)

	if _, ok := err.(*Error); ok {
		t.Error(util.FormatTest("buildError", "bceError", "error"))
	}

	bceError := &Error{
		StatusCode: 500,
		Code:       strconv.Itoa(http.StatusInternalServerError),
		Message:    "failed",
		RequestID:  "123",
	}
	byteArray, err := json.Marshal(bceError)

	if err != nil {
		t.Error(util.FormatTest("buildError", err.Error(), "nil"))
	}

	httpResponse := &http.Response{StatusCode: http.StatusInternalServerError}
	resp = &Response{BodyContent: byteArray, Response: httpResponse}
	err = buildError(resp)

	if _, ok := err.(*Error); !ok {
		t.Error(util.FormatTest("buildError", "error", "bceError"))
	}

	resp = &Response{BodyContent: []byte("Unknown Error"), Response: httpResponse}
	err = buildError(resp)

	if _, ok := err.(*Error); ok {
		t.Error(util.FormatTest("buildError", "bceError", "error"))
	}
}
