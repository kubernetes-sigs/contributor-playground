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

package bos

import (
	"net/http"
	"strconv"
	"testing"

	"k8s.io/cloud-provider-baiducloud/pkg/sdk/bce"
	"k8s.io/cloud-provider-baiducloud/pkg/sdk/util"
)

func TestNewObjectMetadataFromHeader(t *testing.T) {
	header := http.Header{
		"cache-control":       []string{"no-cache", "no-store"},
		"Content-Disposition": []string{"inline"},
		"Content-Length":      []string{"1024"},
		"Content-Range":       []string{"bytes=0-1024"},
		"Content-Type":        []string{"text/plain"},
		"Expires":             []string{"Tue, 03 Jan 2017 05:30:19 GMT"},
		"Etag":                []string{"abc123"},
		"x-bce-meta-name":     []string{"hello"},
	}
	metadata := NewObjectMetadataFromHeader(header)

	if metadata.CacheControl != "no-cache" {
		t.Error(util.FormatTest("NewObjectMetadataFromHeader", metadata.CacheControl, "no-cache"))
	}

	if metadata.ContentDisposition != "inline" {
		t.Error(util.FormatTest("NewObjectMetadataFromHeader", metadata.ContentDisposition, "inline"))
	}

	if metadata.ContentLength != 1024 {
		t.Error(util.FormatTest("NewObjectMetadataFromHeader", strconv.FormatInt(metadata.ContentLength, 10), strconv.Itoa(1024)))
	}

	if metadata.ContentRange != "bytes=0-1024" {
		t.Error(util.FormatTest("NewObjectMetadataFromHeader", metadata.ContentRange, "bytes=0-1024"))
	}

	if metadata.ContentType != "text/plain" {
		t.Error(util.FormatTest("NewObjectMetadataFromHeader", metadata.ContentType, "text/plain"))
	}

	if metadata.Expires != "Tue, 03 Jan 2017 05:30:19 GMT" {
		t.Error(util.FormatTest("NewObjectMetadataFromHeader", metadata.Expires, "Tue, 03 Jan 2017 05:30:19 GMT"))
	}

	if metadata.ETag != "abc123" {
		t.Error(util.FormatTest("NewObjectMetadataFromHeader", metadata.ETag, "abc123"))
	}

	if metadata.UserMetadata["x-bce-meta-name"] != "hello" {
		t.Error(util.FormatTest("NewObjectMetadataFromHeader", metadata.UserMetadata["x-bce-meta-name"], "hello"))
	}
}

func TestAddUserMetadata(t *testing.T) {
	metadata := &ObjectMetadata{}
	metadata.AddUserMetadata("x-bce-meta-name", "hello")

	if metadata.UserMetadata["x-bce-meta-name"] != "hello" {
		t.Error(util.FormatTest("AddUserMetadata", metadata.UserMetadata["x-bce-meta-name"], "hello"))
	}
}

func TestMergeToSignOption(t *testing.T) {
	option := &bce.SignOption{}
	header := http.Header{
		"cache-control":       []string{"no-cache", "no-store"},
		"Content-Disposition": []string{"inline"},
		"Content-Length":      []string{"1024"},
		"Content-Range":       []string{"bytes=0-1024"},
		"Content-Type":        []string{"text/plain"},
		"Expires":             []string{"Tue, 03 Jan 2017 05:30:19 GMT"},
		"Etag":                []string{"abc123"},
		"x-bce-meta-name":     []string{"hello"},
	}
	metadata := NewObjectMetadataFromHeader(header)
	metadata.ContentMD5 = "md5-value"
	metadata.ContentSha256 = "sha256-value"
	metadata.AddUserMetadata("server", "nginx")
	metadata.mergeToSignOption(option)

	if option.Headers["Cache-Control"] != "no-cache" {
		t.Error(util.FormatTest("mergeToSignOption", option.Headers["Cache-Control"], "no-cache"))
	}

	if option.Headers["Content-Disposition"] != "inline" {
		t.Error(util.FormatTest("mergeToSignOption", option.Headers["Content-Disposition"], "inline"))
	}

	if option.Headers["x-bce-meta-name"] != "hello" {
		t.Error(util.FormatTest("mergeToSignOption", option.Headers["x-bce-meta-name"], "hello"))
	}

	if option.Headers["x-bce-meta-server"] != "nginx" {
		t.Error(util.FormatTest("mergeToSignOption", option.Headers["x-bce-meta-server"], "nginx"))
	}

	if option.Headers["Content-MD5"] != "md5-value" {
		t.Error(util.FormatTest("mergeToSignOption", option.Headers["Content-MD5"], "md5-value"))
	}

	if option.Headers["x-bce-content-sha256"] != "sha256-value" {
		t.Error(util.FormatTest("mergeToSignOption", option.Headers["x-bce-content-sha256"], "sha256-value"))
	}
}

func TestGetETag(t *testing.T) {
	header := http.Header{"Etag": []string{"\"12345\""}}
	response := NewPutObjectResponse(header)
	etag := response.GetETag()

	if etag != "12345" {
		t.Error(util.FormatTest("GetETag", etag, "12345"))
	}
}

func TestGetCommonPrefixes(t *testing.T) {
	response := &ListObjectsResponse{
		CommonPrefixes: []map[string]string{
			map[string]string{"prefix": "prefix-0"},
			map[string]string{"prefix": "prefix-1"},
			map[string]string{"prefix": "prefix-2"},
		},
	}

	prefixes := response.GetCommonPrefixes()

	if len(prefixes) != len(response.CommonPrefixes) {
		t.Error(util.FormatTest("GetCommonPrefixes", strconv.Itoa(len(prefixes)), strconv.Itoa(len(response.CommonPrefixes))))
	}

	if prefixes[1] != response.CommonPrefixes[1]["prefix"] {
		t.Error(util.FormatTest("GetCommonPrefixes", prefixes[1], response.CommonPrefixes[1]["prefix"]))
	}
}

func TestMergeToSignOptionForCopyObjectRequest(t *testing.T) {
	option := &bce.SignOption{}
	request := CopyObjectRequest{
		SrcBucketName:  "source-bucket",
		SrcKey:         "source-bucket-key",
		DestBucketName: "dest-bucket",
		DestKey:        "dest-bucket-key",
		ObjectMetadata: &ObjectMetadata{
			CacheControl: "no-cache",
			ContentType:  "text/plain",
			UserMetadata: map[string]string{"request-id": "12345"},
		},
		SourceMatch:           "xxx",
		SourceNoneMatch:       "xxx",
		SourceModifiedSince:   "xxx",
		SourceUnmodifiedSince: "xxx",
	}

	request.mergeToSignOption(option)

	if len(option.Headers) != 8 {
		t.Error(util.FormatTest("CopyObjectRequest: MergeToSignOption", strconv.Itoa(len(option.Headers)), strconv.Itoa(8)))
	}

	if option.Headers["x-bce-metadata-directive"] != "replace" {
		t.Error(util.FormatTest("CopyObjectRequest: MergeToSignOption", option.Headers["x-bce-metadata-directive"], "replace"))
	}
}

func TestMergeToSignOptionForGetObjectRequest(t *testing.T) {
	lengthRange := "0-1024"
	expected := "bytes=" + lengthRange

	option := &bce.SignOption{}
	request := &GetObjectRequest{
		BucketName: "test-bucket-name",
		ObjectKey:  "test-object-key",
		Range:      lengthRange,
	}

	request.MergeToSignOption(option)

	if option.Headers["Range"] != expected {
		t.Error(util.FormatTest("GetObjectRequest: MergeToSignOption", option.Headers["Range"], expected))
	}
}

func TestErrorForDeleteMultipleObjectsError(t *testing.T) {
	deleteMultipleObjectsError := DeleteMultipleObjectsError{
		Key:     "error-key",
		Code:    "error-code",
		Message: "error-message",
	}

	result := deleteMultipleObjectsError.Error()
	expected := "error-message"

	if result != expected {
		t.Error(util.FormatTest("DeleteMultipleObjectsError: Error", result, expected))
	}

	deleteMultipleObjectsError = DeleteMultipleObjectsError{
		Key:  "error-key",
		Code: "error-code",
	}

	result = deleteMultipleObjectsError.Error()
	expected = "error-code"

	if result != expected {
		t.Error(util.FormatTest("DeleteMultipleObjectsError: Error", result, expected))
	}
}

func TestSort(t *testing.T) {
	request := CompleteMultipartUploadRequest{
		BucketName: "test-bucket",
		ObjectKey:  "test-object-key",
		UploadId:   "test-upload-id",
		Parts: []PartSummary{
			PartSummary{PartNumber: 2},
			PartSummary{PartNumber: 3},
			PartSummary{PartNumber: 1},
		},
	}

	result := request.Parts[0].PartNumber
	expected := 2

	if result != expected {
		t.Error(util.FormatTest("CompleteMultipartUploadRequest: Sort", strconv.Itoa(result), strconv.Itoa(expected)))
	}

	request.sort()

	result = request.Parts[0].PartNumber
	expected = 1

	if result != expected {
		t.Error(util.FormatTest("CompleteMultipartUploadRequest: Sort", strconv.Itoa(result), strconv.Itoa(expected)))
	}
}

func TestIsUserDefinedMetadata(t *testing.T) {
	expected := true
	result := IsUserDefinedMetadata("x-bce-meta-name")

	if result != expected {
		t.Error(util.FormatTest("IsUserDefinedMetadata", strconv.FormatBool(result), strconv.FormatBool(expected)))
	}

	expected = false
	result = IsUserDefinedMetadata("content-type")

	if result != expected {
		t.Error(util.FormatTest("IsUserDefinedMetadata", strconv.FormatBool(result), strconv.FormatBool(expected)))
	}
}

func TestToUserDefinedMetadata(t *testing.T) {
	expected := "x-bce-meta-name"
	result := ToUserDefinedMetadata("x-bce-meta-name")

	if result != expected {
		t.Error(util.FormatTest("ToUserDefinedMetadata", result, expected))
	}

	result = ToUserDefinedMetadata("name")

	if result != expected {
		t.Error(util.FormatTest("ToUserDefinedMetadata", result, expected))
	}

	expected = "x-bce-meta-content-type"
	result = ToUserDefinedMetadata("content-type")

	if result != expected {
		t.Error(util.FormatTest("ToUserDefinedMetadata", result, expected))
	}
}
