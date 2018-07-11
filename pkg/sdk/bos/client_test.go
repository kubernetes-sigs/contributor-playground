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
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"k8s.io/cloud-provider-baiducloud/pkg/sdk/bce"
	"k8s.io/cloud-provider-baiducloud/pkg/sdk/util"
)

var credentials = bce.NewCredentials(os.Getenv("BAIDU_BCE_AK"), os.Getenv("BAIDU_BCE_SK"))

//var bceConfig = bce.NewConfig(credentials)
var bceConfig = &bce.Config{
	Credentials: credentials,
	Checksum:    true,
	Region:      os.Getenv("BOS_REGION"),
}
var bosConfig = NewConfig(bceConfig)
var bosClient = NewClient(bosConfig)

func TestCheckBucketName(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)

			defer func() {
				if err := recover(); err != nil {
					fmt.Println(err)

					defer func() {
						if err := recover(); err != nil {
							fmt.Println(err)
							t.Error(util.FormatTest("checkBucketName", "panic", "no panic"))
						}
					}()

					checkBucketName("bucket-0")
				}
			}()

			checkBucketName("/bucket-0")
			t.Error(util.FormatTest("checkBucketName", "no panic", "panic"))
		}
	}()

	checkBucketName("")
	t.Error(util.FormatTest("checkBucketName", "no panic", "panic"))
}

func TestCheckObjectKey(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)

			defer func() {
				if err := recover(); err != nil {
					fmt.Println(err)

					defer func() {
						if err := recover(); err != nil {
							fmt.Println(err)
							t.Error(util.FormatTest("checkObjectKey", "panic", "no panic"))
						}
					}()

					checkObjectKey("object-0")
				}
			}()

			checkObjectKey("/object-0")
			t.Error(util.FormatTest("checkObjectKey", "no panic", "panic"))
		}
	}()

	checkObjectKey("")
	t.Error(util.FormatTest("checkObjectKey", "no panic", "panic"))
}

func TestGetURL(t *testing.T) {
	expected := fmt.Sprintf("http://bucket-0.%s.bcebos.com/object-0", bosClient.GetRegion())
	url := bosClient.GetURL("bucket-0", "object-0", nil)

	if url != expected {
		t.Error(util.FormatTest("GetURL", url, expected))
	}

	expected = fmt.Sprintf("http://%s.bcebos.com/object-0", bosClient.GetRegion())
	url = bosClient.GetURL("", "object-0", nil)

	if url != expected {
		t.Error(util.FormatTest("GetURL", url, expected))
	}
}

func TestGetBucketLocation(t *testing.T) {
	bucketNamePrefix := "baidubce-sdk-go-test-for-get-bucket-location-"
	method := "GetBucketLocation"

	bucketName := bucketNamePrefix + strconv.Itoa(int(time.Now().UnixNano()))
	_, err := bosClient.GetBucketLocation(bucketName, nil)

	if err == nil {
		t.Error(util.FormatTest(method, "nil", "error"))
	}

	around(t, method, bucketNamePrefix, "", func(bucketName string) {
		expected := bosClient.GetRegion()
		location, _ := bosClient.GetBucketLocation(bucketName, nil)

		if location.LocationConstraint != expected {
			t.Error(util.FormatTest(method, location.LocationConstraint, expected))
		}
	})
}

func TestListBuckets(t *testing.T) {
	_, err := bosClient.ListBuckets(nil)

	if err != nil {
		t.Error(util.FormatTest("ListBuckets", err.Error(), "nil"))
	}
}

func TestCreateBucket(t *testing.T) {
	bucketNamePrefix := "baidubce-sdk-go-test-for-create-bucket-"
	method := "CreateBucket"

	around(t, method, bucketNamePrefix, "", nil)
}

func TestDoesBucketExist(t *testing.T) {
	bucketNamePrefix := "baidubce-sdk-go-test-for-does-bucket-exist-"
	method := "DoesBucketExist"

	expected := false
	exists, _ := bosClient.DoesBucketExist(bucketNamePrefix, nil)

	if exists != expected {
		t.Error(util.FormatTest(method, strconv.FormatBool(exists), strconv.FormatBool(expected)))
	}

	around(t, method, bucketNamePrefix, "", func(bucketName string) {
		expected := true
		exists, err := bosClient.DoesBucketExist(bucketName, nil)

		if err != nil {
			t.Error(util.FormatTest(method, err.Error(), "nil"))
		} else if exists != expected {
			t.Error(util.FormatTest(method, strconv.FormatBool(exists), strconv.FormatBool(expected)))
		}
	})

}

func TestDeleteBucket(t *testing.T) {
	bucketNamePrefix := "baidubce-sdk-go-test-for-delete-bucket-"
	method := "DeleteBucket"

	around(t, method, bucketNamePrefix, "", nil)
}

func TestSetBucketPrivate(t *testing.T) {
	bucketNamePrefix := "baidubce-sdk-go-test-for-set-bucket-private-"
	method := "SetBucketPrivate"

	around(t, method, bucketNamePrefix, "", func(bucketName string) {
		err := bosClient.SetBucketPrivate(bucketName, nil)
		if err != nil {
			t.Error(util.FormatTest(method, err.Error(), "nil"))
		}
	})
}

func TestSetBucketPublicRead(t *testing.T) {
	bucketNamePrefix := "baidubce-sdk-go-test-for-set-bucket-public-read-"
	method := "SetBucketPublicRead"

	around(t, method, bucketNamePrefix, "", func(bucketName string) {
		err := bosClient.SetBucketPublicRead(bucketName, nil)
		if err != nil {
			t.Error(util.FormatTest(method, err.Error(), "nil"))
		}
	})
}

func TestSetBucketPublicReadWrite(t *testing.T) {
	bucketNamePrefix := "baidubce-sdk-go-test-for-set-bucket-public-rw-"
	method := "SetBucketPublicReadWrite"

	around(t, method, bucketNamePrefix, "", func(bucketName string) {
		err := bosClient.SetBucketPublicReadWrite(bucketName, nil)
		if err != nil {
			t.Error(util.FormatTest(method, err.Error(), "nil"))
		}
	})
}

func TestGetBucketAcl(t *testing.T) {
	bucketNamePrefix := "baidubce-sdk-go-test-for-get-bucket-acl-"
	method := "GetBucketAcl"

	_, err := bosClient.GetBucketAcl(bucketNamePrefix, nil)
	if err == nil {
		t.Error(util.FormatTest(method, "nil", "error"))
	}

	around(t, method, bucketNamePrefix, "", func(bucketName string) {
		_, err := bosClient.GetBucketAcl(bucketName, nil)
		if err != nil {
			t.Error(util.FormatTest(method, err.Error(), "nil"))
		}
	})
}

func TestSetBucketAcl(t *testing.T) {
	bucketNamePrefix := "baidubce-sdk-go-test-for-set-bucket-acl-"
	method := "SetBucketAcl"

	around(t, method, bucketNamePrefix, "", func(bucketName string) {
		bucketAcl := BucketAcl{
			AccessControlList: []Grant{
				Grant{
					Grantee: []BucketGrantee{
						BucketGrantee{Id: "ef5a4b19192f4931adcf0e12f82795e2"},
					},
					Permission: []string{"FULL_CONTROL"},
				},
			},
		}
		if err := bosClient.SetBucketAcl(bucketName, bucketAcl, nil); err != nil {
			t.Error(util.FormatTest(method, err.Error(), "nil"))
		}
	})
}

func TestPutObject(t *testing.T) {
	bucketNamePrefix := "baidubce-sdk-go-test-for-put-object-"
	method := "PutObject"
	objectKey := "put-object-from-string.txt"
	str := "Hello World 你好"

	metadata := &ObjectMetadata{
		CacheControl:       "no-cache",
		ContentDisposition: "attachment",
		ContentType:        "text/plain",
	}

	around(t, method, bucketNamePrefix, objectKey, func(bucketName string) {
		_, err := bosClient.PutObject(bucketName, objectKey, str, metadata, nil)
		if err != nil {
			t.Error(util.FormatTest(method, err.Error(), "nil"))
		}
	})

	byteArray := []byte(str)
	objectKey = "put-object-from-bytes.txt"
	around(t, method, bucketNamePrefix, objectKey, func(bucketName string) {
		_, err := bosClient.PutObject(bucketName, objectKey, byteArray, metadata, nil)
		if err != nil {
			t.Error(util.FormatTest(method, err.Error(), "nil"))
		}
	})

	reader := strings.NewReader(str)
	objectKey = "put-object-from-reader.txt"
	around(t, method, bucketNamePrefix, objectKey, func(bucketName string) {
		_, err := bosClient.PutObject(bucketName, objectKey, reader, metadata, nil)
		if err != nil {
			t.Error(util.FormatTest(method, err.Error(), "nil"))
		}
	})

	defer func() {
		if err := recover(); err == nil {
			t.Error(util.FormatTest(method, "panic", "nil"))
		}
	}()

	around(t, method, bucketNamePrefix, objectKey, func(bucketName string) {
		_, err := bosClient.PutObject(bucketName, objectKey, 1, nil, nil)
		if err != nil {
			t.Error(util.FormatTest(method, err.Error(), "nil"))
		}
	})
}

func TestDeleteObject(t *testing.T) {
	bucketNamePrefix := "baidubce-sdk-go-test-for-delete-object-"
	method := "DeleteObject"
	objectKey := "put-object-from-string.txt"
	str := "Hello World 你好"

	around(t, method, bucketNamePrefix, objectKey, func(bucketName string) {
		_, err := bosClient.PutObject(bucketName, objectKey, str, nil, nil)
		if err != nil {
			t.Error(util.FormatTest(method, err.Error(), "nil"))
		}
	})
}

func TestDeleteMultipleObjects(t *testing.T) {
	bucketNamePrefix := "baidubce-sdk-go-test-for-delete-multiple-objects-"
	method := "DeleteMultipleObjects"
	str := "Hello World 你好"

	objects := []string{
		"examples/delete-multiple-objects/put-object-from-string.txt",
		"examples/delete-multiple-objects/put-object-from-string-2.txt",
		"examples/delete-multiple-objects/put-object-from-string-3.txt",
	}

	around(t, method, bucketNamePrefix, objects, func(bucketName string) {
		for _, objectKey := range objects {
			_, err := bosClient.PutObject(bucketName, objectKey, str, nil, nil)
			if err != nil {
				t.Error(util.FormatTest(method, err.Error(), "nil"))
			}
		}
	})

	around(t, method, bucketNamePrefix, []string{}, nil)

	bucketName := bucketNamePrefix + strconv.Itoa(int(time.Now().UnixNano()))
	if len(bucketName) > 63 {
		bucketName = bucketName[:63]
	}

	err := bosClient.CreateBucket(bucketName, nil)

	if err != nil {
		t.Error(util.FormatTest(method+" at creating bucket", err.Error(), "nil"))
	} else {
		defer func() {
			deleteMultipleObjectsResponse, err := bosClient.DeleteMultipleObjects(bucketName, objects, nil)

			if err != nil {
				t.Error(util.FormatTest(method, err.Error(), "nil"))
			} else if len(deleteMultipleObjectsResponse.Errors) != 1 {
				t.Error(util.FormatTest(method, strconv.Itoa(len(deleteMultipleObjectsResponse.Errors)), strconv.Itoa(1)))
			}

			err = bosClient.DeleteBucket(bucketName, nil)

			if err != nil {
				t.Error(util.FormatTest(method+" at deleting bucket", err.Error(), "nil"))
			}

		}()

		for index, objectKey := range objects {
			if index < len(objects)-1 {
				_, err := bosClient.PutObject(bucketName, objectKey, str, nil, nil)
				if err != nil {
					t.Error(util.FormatTest(method, err.Error(), "nil"))
				}
			}
		}
	}
}

func TestListObjects(t *testing.T) {
	bucketNamePrefix := "baidubce-sdk-go-test-for-list-objects-"
	method := "ListObjects"
	objectKey := "put-object-from-string.txt"
	str := "Hello World 你好"

	around(t, method, bucketNamePrefix, objectKey, func(bucketName string) {
		_, err := bosClient.PutObject(bucketName, objectKey, str, nil, nil)
		if err != nil {
			t.Error(util.FormatTest(method, err.Error(), "nil"))
		} else {
			listObjectsResponse, err := bosClient.ListObjects(bucketName, nil)
			if err != nil {
				t.Error(util.FormatTest(method, err.Error(), "nil"))
			} else if length := len(listObjectsResponse.Contents); length != 1 {
				t.Error(util.FormatTest(method, strconv.Itoa(length), "1"))
			}
		}
	})
}

func TestListObjectsFromRequest(t *testing.T) {
	bucketNamePrefix := "baidubce-sdk-go-test-for-list-objects-from-request-"
	method := "ListObjectsFromRequest"
	str := "Hello World 你好"

	objects := []string{
		"hello.txt",
		"examples/list-objects-from-request/put-object-from-string.txt",
		"examples/list-objects-from-request/put-object-from-string-2.txt",
		"examples2/list-objects-from-request/put-object-from-string-3.txt",
		"examples2/list-objects-from-request/put-object-from-string.txt",
		"examples3/list-objects-from-request/put-object-from-string-2.txt",
		"examples3/list-objects-from-request/put-object-from-string-3.txt",
	}
	objectsLength := len(objects)

	around(t, method, bucketNamePrefix, objects, func(bucketName string) {
		for _, objectKey := range objects {
			_, err := bosClient.PutObject(bucketName, objectKey, str, nil, nil)

			if err != nil {
				t.Error(util.FormatTest(method, err.Error(), "nil"))
				return
			}
		}

		listObjectsRequest := ListObjectsRequest{
			BucketName: bucketName,
		}
		listObjectsResponse, err := bosClient.ListObjectsFromRequest(listObjectsRequest, nil)
		if err != nil {
			t.Error(util.FormatTest(method, err.Error(), "nil"))
		} else if length := len(listObjectsResponse.Contents); length != objectsLength {
			t.Error(util.FormatTest(method, strconv.Itoa(length), strconv.Itoa(objectsLength)))
		}

		listObjectsRequest = ListObjectsRequest{
			BucketName: bucketName,
			Delimiter:  "/",
		}
		listObjectsResponse, err = bosClient.ListObjectsFromRequest(listObjectsRequest, nil)
		if err != nil {
			t.Error(util.FormatTest(method, err.Error(), "nil"))
		} else {
			if length := len(listObjectsResponse.Contents); length != 1 {
				t.Error(util.FormatTest(method, strconv.Itoa(length), "1"))
			}

			if length := len(listObjectsResponse.CommonPrefixes); length != 3 {
				t.Error(util.FormatTest(method, strconv.Itoa(length), "3"))
			}
		}

		listObjectsRequest = ListObjectsRequest{
			BucketName: bucketName,
			Prefix:     "examples2",
		}
		listObjectsResponse, err = bosClient.ListObjectsFromRequest(listObjectsRequest, nil)
		if err != nil {
			t.Error(util.FormatTest(method, err.Error(), "nil"))
		} else if length := len(listObjectsResponse.Contents); length != 2 {
			t.Error(util.FormatTest(method, strconv.Itoa(length), "2"))
		}

		listObjectsRequest = ListObjectsRequest{
			BucketName: bucketName,
			MaxKeys:    2,
		}
		listObjectsResponse, err = bosClient.ListObjectsFromRequest(listObjectsRequest, nil)
		if err != nil {
			t.Error(util.FormatTest(method, err.Error(), "nil"))
		} else if length := len(listObjectsResponse.Contents); length != 2 {
			t.Error(util.FormatTest(method, strconv.Itoa(length), "2"))
		}

		listObjectsRequest = ListObjectsRequest{
			BucketName: bucketName,
			Marker:     "examples2",
		}
		listObjectsResponse, err = bosClient.ListObjectsFromRequest(listObjectsRequest, nil)
		if err != nil {
			t.Error(util.FormatTest(method, err.Error(), "nil"))
		} else if length := len(listObjectsResponse.Contents); length != 5 {
			t.Error(util.FormatTest(method, strconv.Itoa(length), "5"))
		}
	})
}

func TestCopyObject(t *testing.T) {
	bucketNamePrefix := "baidubce-sdk-go-test-for-copy-object-"
	method := "CopyObject"
	objectKey := "put-object-from-string.txt"
	str := "Hello World 你好"

	around(t, method, bucketNamePrefix, objectKey, func(bucketName string) {
		_, err := bosClient.PutObject(bucketName, objectKey, str, nil, nil)

		if err != nil {
			t.Error(util.FormatTest(method, err.Error(), "nil"))
		} else {
			destKey := "put-object-from-string-copy.txt"
			_, err := bosClient.CopyObject(bucketName, objectKey, bucketName, destKey, nil)

			if err != nil {
				t.Error(util.FormatTest(method, err.Error(), "nil"))
			} else {
				listObjectsResponse, err := bosClient.ListObjects(bucketName, nil)

				if err != nil {
					t.Error(util.FormatTest(method, err.Error(), "nil"))
				} else if length := len(listObjectsResponse.Contents); length != 2 {
					t.Error(util.FormatTest(method, strconv.Itoa(length), "2"))
				} else {
					err = bosClient.DeleteObject(bucketName, destKey, nil)

					if err != nil {
						t.Error(util.FormatTest(method+" at deleting object", err.Error(), "nil"))
					}
				}
			}
		}
	})
}

func TestCopyObjectFromRequest(t *testing.T) {
	bucketNamePrefix := "baidubce-sdk-go-test-for-copy-object-from-request-"
	method := "CopyObjectFromRequest"
	objectKey := "put-object-from-string.txt"
	str := "Hello World 你好"

	around(t, method, bucketNamePrefix, objectKey, func(bucketName string) {
		_, err := bosClient.PutObject(bucketName, objectKey, str, nil, nil)

		if err != nil {
			t.Error(util.FormatTest(method, err.Error(), "nil"))
		} else {
			destKey := "put-object-from-string-copy.txt"

			copyObjectRequest := CopyObjectRequest{
				SrcBucketName:  bucketName,
				SrcKey:         objectKey,
				DestBucketName: bucketName,
				DestKey:        destKey,
				ObjectMetadata: &ObjectMetadata{
					CacheControl: "no-cache",
					UserMetadata: map[string]string{
						"test-user-metadata": "test user metadata",
						"x-bce-meta-name":    "x-bce-meta-name",
					},
				},
			}

			_, err := bosClient.CopyObjectFromRequest(copyObjectRequest, nil)

			if err != nil {
				t.Error(util.FormatTest(method, err.Error(), "nil"))
			} else {
				listObjectsResponse, err := bosClient.ListObjects(bucketName, nil)

				if err != nil {
					t.Error(util.FormatTest(method, err.Error(), "nil"))
				} else if length := len(listObjectsResponse.Contents); length != 2 {
					t.Error(util.FormatTest(method, strconv.Itoa(length), "2"))
				} else {
					err = bosClient.DeleteObject(bucketName, destKey, nil)

					if err != nil {
						t.Error(util.FormatTest(method+" at deleting object", err.Error(), "nil"))
					}
				}
			}
		}
	})
}

func TestGetObject(t *testing.T) {
	bucketNamePrefix := "baidubce-sdk-go-test-for-get-object-"
	method := "GetObject"
	objectKey := "put-object-from-string.txt"
	str := "Hello World 你好"

	around(t, method, bucketNamePrefix, objectKey, func(bucketName string) {
		_, err := bosClient.PutObject(bucketName, objectKey, str, nil, nil)

		if err != nil {
			t.Error(util.FormatTest(method, err.Error(), "nil"))
		} else {
			object, err := bosClient.GetObject(bucketName, objectKey, nil)

			if err != nil {
				t.Error(util.FormatTest(method, err.Error(), "nil"))
			} else if object.ObjectMetadata.ETag == "" {
				t.Error(util.FormatTest(method, "etag is empty", "non empty etag"))
			} else {
				byteArray, err := ioutil.ReadAll(object.ObjectContent)

				if err != nil {
					t.Error(util.FormatTest(method, err.Error(), "nil"))
				} else if len(byteArray) == 0 {
					t.Error(util.FormatTest(method, "body is empty", "non empty body"))
				}
			}
		}
	})
}

func TestGetObjectFromRequest(t *testing.T) {
	bucketNamePrefix := "baidubce-sdk-go-test-for-get-object-from-request-"
	method := "GetObjectFromRequest"
	objectKey := "put-object-from-string.txt"
	str := "Hello World 你好"

	around(t, method, bucketNamePrefix, objectKey, func(bucketName string) {
		_, err := bosClient.PutObject(bucketName, objectKey, str, nil, nil)

		if err != nil {
			t.Error(util.FormatTest(method, err.Error(), "nil"))
		} else {
			getObjectRequest := GetObjectRequest{
				BucketName: bucketName,
				ObjectKey:  objectKey,
			}
			getObjectRequest.SetRange(0, 1000)
			object, err := bosClient.GetObjectFromRequest(getObjectRequest, nil)

			if err != nil {
				t.Error(util.FormatTest(method, err.Error(), "nil"))
			} else if object.ObjectMetadata.ETag == "" {
				t.Error(util.FormatTest(method, "etag is empty", "non empty etag"))
			} else {
				byteArray, err := ioutil.ReadAll(object.ObjectContent)

				if err != nil {
					t.Error(util.FormatTest(method, err.Error(), "nil"))
				} else if len(byteArray) == 0 {
					t.Error(util.FormatTest(method, "body is empty", "non empty body"))
				}
			}
		}
	})
}

func TestGetObjectToFile(t *testing.T) {
	bucketNamePrefix := "baidubce-sdk-go-test-for-get-object-to-file-"
	method := "GetObjectToFile"
	objectKey := "put-object-from-string.txt"
	str := "Hello World 你好"

	around(t, method, bucketNamePrefix, objectKey, func(bucketName string) {
		_, err := bosClient.PutObject(bucketName, objectKey, str, nil, nil)

		if err != nil {
			t.Error(util.FormatTest(method, err.Error(), "nil"))
		} else {
			getObjectRequest := &GetObjectRequest{
				BucketName: bucketName,
				ObjectKey:  objectKey,
			}
			getObjectRequest.SetRange(0, 1000)

			file, err := os.OpenFile(objectKey, os.O_WRONLY|os.O_CREATE, 0666)

			if err != nil {
				t.Error(util.FormatTest(method, err.Error(), "nil"))
			} else {
				objectMetadata, err := bosClient.GetObjectToFile(getObjectRequest, file, nil)

				if err != nil {
					t.Error(util.FormatTest(method, err.Error(), "nil"))
				} else if objectMetadata.ETag == "" {
					t.Error(util.FormatTest(method, "etag is empty", "non empty etag"))
				} else if !util.CheckFileExists(objectKey) {
					t.Error(util.FormatTest(method, "file is not saved to local", "file saved to local"))
				} else {
					err := os.Remove(objectKey)

					if err != nil {
						t.Error(util.FormatTest(method, err.Error(), "nil"))
					}
				}
			}
		}
	})
}

func TestGetObjectMetadata(t *testing.T) {
	bucketNamePrefix := "baidubce-sdk-go-test-for-get-object-metadata-"
	method := "GetObjectMetadata"
	objectKey := "put-object-from-string.txt"
	str := "Hello World 你好"

	around(t, method, bucketNamePrefix, objectKey, func(bucketName string) {
		_, err := bosClient.PutObject(bucketName, objectKey, str, nil, nil)

		if err != nil {
			t.Error(util.FormatTest(method, err.Error(), "nil"))
		} else {
			objectMetadata, err := bosClient.GetObjectMetadata(bucketName, objectKey, nil)

			if err != nil {
				t.Error(util.FormatTest(method, err.Error(), "nil"))
			} else if objectMetadata.ETag == "" {
				t.Error(util.FormatTest(method, "etag is empty", "non empty etag"))
			}
		}
	})
}

func TestGeneratePresignedUrl(t *testing.T) {
	bucketNamePrefix := "baidubce-sdk-go-test-for-generate-presigned-url-"
	method := "GeneratePresignedUrl"
	objectKey := "put-object-from-string.txt"
	str := "Hello World 你好"

	around(t, method, bucketNamePrefix, objectKey, func(bucketName string) {
		_, err := bosClient.PutObject(bucketName, objectKey, str, nil, nil)

		if err != nil {
			t.Error(util.FormatTest(method, err.Error(), "nil"))
		} else {
			url, err := bosClient.GeneratePresignedUrl(bucketName, objectKey, nil)

			if err != nil {
				t.Error(util.FormatTest(method, err.Error(), "nil"))
			} else {
				req, err := http.NewRequest("GET", url, nil)

				if err != nil {
					t.Error(util.FormatTest(method, err.Error(), "nil"))
				} else {
					httpClient := http.Client{}
					res, err := httpClient.Do(req)

					if err != nil {
						t.Error(util.FormatTest(method, err.Error(), "nil"))
					} else if res.StatusCode != 200 {
						t.Error(util.FormatTest(method, fmt.Sprintf("status code: %v", res.StatusCode), "status code: 200"))
					}
				}
			}
		}
	})
}

func TestAppendObject(t *testing.T) {
	bucketNamePrefix := "baidubce-sdk-go-test-for-append-object-"
	method := "AppendObject"
	objectKey := "append-object-from-string.txt"
	str := "Hello World 你好"
	contentLength := len(str)
	offset := 0

	around(t, method, bucketNamePrefix, objectKey, func(bucketName string) {
		appendObjectResponse, err := bosClient.AppendObject(bucketName, objectKey, offset, str, nil, nil)
		if err != nil {
			t.Error(util.FormatTest(method, err.Error(), "nil"))
		} else if appendObjectResponse.GetETag() == "" || appendObjectResponse.GetNextAppendOffset() == "" {
			t.Error(util.FormatTest(method, "etag and next append offset are not exists", "etag and next append offset"))
		} else {
			length, err := strconv.Atoi(appendObjectResponse.GetNextAppendOffset())

			if err != nil {
				t.Error(util.FormatTest(method, err.Error(), "nil"))
			} else {
				md5 := appendObjectResponse.GetMD5()
				if md5 == "" {
					t.Error(util.FormatTest(method, md5, "md5 value of object"))
				} else {
					offset = length
					appendObjectResponse, err := bosClient.AppendObject(bucketName, objectKey, offset, []byte(str), nil, nil)

					if err != nil {
						t.Error(util.FormatTest(method, err.Error(), "nil"))
					} else if appendObjectResponse.GetETag() == "" || appendObjectResponse.GetNextAppendOffset() == "" {
						t.Error(util.FormatTest(method, "etag and next append offset are not exists", "etag and next append offset"))
					} else {
						offset = length * 2
						appendObjectResponse, err := bosClient.AppendObject(
							bucketName,
							objectKey,
							offset,
							strings.NewReader(str),
							&ObjectMetadata{CacheControl: "no-cache"},
							nil,
						)
						if err != nil {
							t.Error(util.FormatTest(method, err.Error(), "nil"))
						} else if appendObjectResponse.GetETag() == "" || appendObjectResponse.GetNextAppendOffset() == "" {
							t.Error(util.FormatTest(method, "etag and next append offset are not exists", "etag and next append offset"))
						} else {
							offset, err := strconv.Atoi(appendObjectResponse.GetNextAppendOffset())

							if err != nil {
								t.Error(util.FormatTest(method, err.Error(), "nil"))
							} else if offset != contentLength*3 {
								t.Error(util.FormatTest(method, strconv.Itoa(offset), strconv.Itoa(contentLength*3)))
							} else {
								file, err := os.OpenFile(objectKey, os.O_WRONLY|os.O_CREATE, 0666)
								if err != nil {
									t.Error(util.FormatTest(method, err.Error(), "nil"))
								} else {
									defer func() {
										file.Close()
										os.Remove(file.Name())
									}()

									_, err = bosClient.AppendObject(bucketName, objectKey, offset, file, nil, nil)
									if err == nil {
										t.Error(util.FormatTest(method, "nil", "error"))
									} else {
										defer func() {
											if err := recover(); err == nil {
												t.Error(util.FormatTest(method, "nil", "error"))
											}
										}()

										_, err = bosClient.AppendObject(bucketName, objectKey, offset, 12, nil, nil)
									}
								}
							}
						}
					}
				}
			}
		}
	})
}

func TestMultipartUploadFromFile(t *testing.T) {
	bucketNamePrefix := "baidubce-sdk-go-test-for-multipart-upload-from-file-"
	method := "MultipartUploadFromFile"
	objectKey := "test-multipart-upload"

	around(t, method, bucketNamePrefix, objectKey, func(bucketName string) {
		file, err := util.TempFileWithSize(1024 * 1024 * 6)

		defer func() {
			if file != nil {
				file.Close()
				os.Remove(file.Name())
			}
		}()

		if err != nil {
			t.Error(util.FormatTest(method, err.Error(), "nil"))
			return
		}

		var partSize int64 = 1024 * 1024 * 2

		completeMultipartUploadResponse, err := bosClient.MultipartUploadFromFile(bucketName,
			objectKey, file.Name(), partSize)

		if err != nil {
			t.Error(util.FormatTest(method, err.Error(), "nil"))
		} else if completeMultipartUploadResponse.ETag == "" {
			t.Error(util.FormatTest(method, "etag is not exists", "etag"))
		}
	})
}

func TestAbortMultipartUpload(t *testing.T) {
	bucketNamePrefix := "baidubce-sdk-go-test-for-abort-multipart-upload-"
	method := "AbortMultipartUpload"
	objectKey := "test-multipart-upload"

	around(t, method, bucketNamePrefix, "", func(bucketName string) {
		initiateMultipartUploadRequest := InitiateMultipartUploadRequest{
			BucketName:     bucketName,
			ObjectKey:      objectKey,
			ObjectMetadata: &ObjectMetadata{},
		}

		initiateMultipartUploadResponse, err := bosClient.InitiateMultipartUpload(initiateMultipartUploadRequest, nil)

		if err != nil {
			t.Error(util.FormatTest(method, err.Error(), "nil"))
			return
		}

		uploadId := initiateMultipartUploadResponse.UploadId

		abortMultipartUploadRequest := AbortMultipartUploadRequest{
			BucketName: bucketName,
			ObjectKey:  objectKey,
			UploadId:   uploadId,
		}

		err = bosClient.AbortMultipartUpload(abortMultipartUploadRequest, nil)

		if err != nil {
			t.Error(util.FormatTest(method, err.Error(), "nil"))
			return
		}
	})
}

func TestListParts(t *testing.T) {
	bucketNamePrefix := "baidubce-sdk-go-test-for-list-parts-"
	objectKey := "test-list-parts"
	method := "ListParts"

	around(t, method, bucketNamePrefix, "", func(bucketName string) {
		initiateMultipartUploadRequest := InitiateMultipartUploadRequest{
			BucketName: bucketName,
			ObjectKey:  objectKey,
		}

		initiateMultipartUploadResponse, err := bosClient.InitiateMultipartUpload(initiateMultipartUploadRequest, nil)

		if err != nil {
			t.Error(util.FormatTest(method, err.Error(), "nil"))
			return
		}

		defer func() {
			if initiateMultipartUploadResponse != nil {
				abortMultipartUploadRequest := AbortMultipartUploadRequest{
					BucketName: bucketName,
					ObjectKey:  objectKey,
					UploadId:   initiateMultipartUploadResponse.UploadId,
				}

				err := bosClient.AbortMultipartUpload(abortMultipartUploadRequest, nil)

				if err != nil {
					t.Error(util.FormatTest(method, err.Error(), "nil"))
				}
			}
		}()

		files := make([]*os.File, 0)
		file, err := util.TempFileWithSize(1024 * 1024 * 6)
		files = append(files, file)

		if err != nil {
			t.Error(util.FormatTest(method, err.Error(), "nil"))
			return
		}

		defer func() {
			for _, f := range files {
				f.Close()
				os.Remove(f.Name())
			}
		}()

		fileInfo, err := file.Stat()

		if err != nil {
			t.Error(util.FormatTest(method, err.Error(), "nil"))
			return
		}

		var partSize int64 = 1024 * 1024 * 5
		var totalSize int64 = fileInfo.Size()
		var partCount int = int(math.Ceil(float64(totalSize) / float64(partSize)))

		var waitGroup sync.WaitGroup
		parts := make([]PartSummary, 0, partCount)

		for i := 0; i < partCount; i++ {
			var skipBytes int64 = partSize * int64(i)
			var size int64 = int64(math.Min(float64(totalSize-skipBytes), float64(partSize)))

			tempFile, err := util.TempFile(nil, "", "")
			files = append(files, tempFile)

			if err != nil {
				t.Error(util.FormatTest(method, err.Error(), "nil"))
				return
			}

			limitReader := io.LimitReader(file, size)
			_, err = io.Copy(tempFile, limitReader)

			if err != nil {
				t.Error(util.FormatTest(method, err.Error(), "nil"))
				return
			}

			partNumber := i + 1

			uploadPartRequest := UploadPartRequest{
				BucketName: bucketName,
				ObjectKey:  objectKey,
				UploadId:   initiateMultipartUploadResponse.UploadId,
				PartSize:   size,
				PartNumber: partNumber,
				PartData:   tempFile,
			}

			waitGroup.Add(1)

			parts = append(parts, PartSummary{PartNumber: partNumber})

			go func(partNumber int) {
				defer waitGroup.Done()

				uploadPartResponse, err := bosClient.UploadPart(uploadPartRequest, nil)

				if err != nil {
					t.Error(util.FormatTest(method, err.Error(), "nil"))
					return
				}

				parts[partNumber-1].ETag = uploadPartResponse.GetETag()
			}(partNumber)
		}

		waitGroup.Wait()

		listPartsResponse, err := bosClient.ListParts(bucketName, objectKey, initiateMultipartUploadResponse.UploadId, nil)

		if err != nil {
			t.Error(util.FormatTest(method, err.Error(), "nil"))
			return
		}

		partCount = len(listPartsResponse.Parts)

		if partCount != 2 {
			t.Error(util.FormatTest(method, fmt.Sprintf("part count is %d", partCount), "part count should be 2"))
		}

		listPartsResponse, err = bosClient.ListPartsFromRequest(
			ListPartsRequest{
				BucketName:       bucketName,
				ObjectKey:        objectKey,
				UploadId:         initiateMultipartUploadResponse.UploadId,
				MaxParts:         100,
				PartNumberMarker: "1",
			},
			nil,
		)

		if err != nil {
			t.Error(util.FormatTest(method, err.Error(), "nil"))
			return
		}

		partCount = len(listPartsResponse.Parts)

		if partCount != 1 {
			t.Error(util.FormatTest(method, fmt.Sprintf("part count is %d", partCount), "part count should be 1"))
		}

		func() {
			defer func() {
				if err := recover(); err == nil {
					t.Error(util.FormatTest(method, "nil", "error"))
				}
			}()

			uploadPartRequest := UploadPartRequest{
				BucketName: bucketName,
				ObjectKey:  objectKey,
				UploadId:   initiateMultipartUploadResponse.UploadId,
				PartSize:   1024*1024*1024*5 + 1,
				PartNumber: partCount + 1,
				PartData:   nil,
			}

			bosClient.UploadPart(uploadPartRequest, nil)
		}()

		func() {
			defer func() {
				if err := recover(); err == nil {
					t.Error(util.FormatTest(method, "nil", "error"))
				}
			}()

			uploadPartRequest := UploadPartRequest{
				BucketName: bucketName,
				ObjectKey:  objectKey,
				UploadId:   initiateMultipartUploadResponse.UploadId,
				PartSize:   1024 * 1024 * 1024 * 5,
				PartNumber: MAX_PART_NUMBER + 1,
				PartData:   nil,
			}

			bosClient.UploadPart(uploadPartRequest, nil)
		}()
	})
}

func TestListMultipartUploads(t *testing.T) {
	bucketNamePrefix := "baidubce-sdk-go-test-for-list-multipart-uploads-"
	objectKey := "test-multipart-upload"
	method := "ListMultipartUploads"

	around(t, method, bucketNamePrefix, "", func(bucketName string) {
		initiateMultipartUploadRequest := InitiateMultipartUploadRequest{
			BucketName: bucketName,
			ObjectKey:  objectKey,
		}

		initiateMultipartUploadResponse, err := bosClient.InitiateMultipartUpload(initiateMultipartUploadRequest, nil)

		if err != nil {
			t.Error(util.FormatTest(method, err.Error(), "nil"))
			return
		}

		defer func() {
			if initiateMultipartUploadResponse != nil {
				abortMultipartUploadRequest := AbortMultipartUploadRequest{
					BucketName: bucketName,
					ObjectKey:  objectKey,
					UploadId:   initiateMultipartUploadResponse.UploadId,
				}

				err = bosClient.AbortMultipartUpload(abortMultipartUploadRequest, nil)

				if err != nil {
					t.Error(util.FormatTest(method, err.Error(), "nil"))
				}
			}
		}()

		listMultipartUploadsResponse, err := bosClient.ListMultipartUploads(bucketName, nil)

		if err != nil {
			t.Error(util.FormatTest(method, err.Error(), "nil"))
			return
		}

		partCount := len(listMultipartUploadsResponse.Uploads)

		if partCount != 1 {
			t.Error(util.FormatTest(method, fmt.Sprintf("part count is %d", partCount), "part count should be 1"))
		}

		commonPrefixes := listMultipartUploadsResponse.GetCommonPrefixes()

		if len(commonPrefixes) != 0 {
			t.Error(util.FormatTest(method, fmt.Sprintf("length of common prefixes is %d", len(commonPrefixes)), "length of common prefixes should be 0"))
		}
	})
}

func TestListMultipartUploadsFromRequest(t *testing.T) {
	bucketNamePrefix := "baidubce-sdk-go-test-for-list-multipart-uploads-from-request-"
	objectKey := "test-multipart-upload"
	method := "ListMultipartUploadsFromRequest"

	around(t, method, bucketNamePrefix, "", func(bucketName string) {
		initiateMultipartUploadRequest := InitiateMultipartUploadRequest{
			BucketName: bucketName,
			ObjectKey:  objectKey,
		}

		initiateMultipartUploadResponse, err := bosClient.InitiateMultipartUpload(initiateMultipartUploadRequest, nil)

		if err != nil {
			t.Error(util.FormatTest(method, err.Error(), "nil"))
			return
		}

		defer func() {
			if initiateMultipartUploadResponse != nil {
				abortMultipartUploadRequest := AbortMultipartUploadRequest{
					BucketName: bucketName,
					ObjectKey:  objectKey,
					UploadId:   initiateMultipartUploadResponse.UploadId,
				}

				err = bosClient.AbortMultipartUpload(abortMultipartUploadRequest, nil)

				if err != nil {
					t.Error(util.FormatTest(method, err.Error(), "nil"))
				}
			}
		}()

		listMultipartUploadsRequest := ListMultipartUploadsRequest{
			BucketName: bucketName,
			Delimiter:  "/",
			Prefix:     objectKey,
			MaxUploads: 100,
			KeyMarker:  "",
		}
		listMultipartUploadsResponse, err := bosClient.ListMultipartUploadsFromRequest(listMultipartUploadsRequest, nil)

		if err != nil {
			t.Error(util.FormatTest(method, err.Error(), "nil"))
			return
		}

		partCount := len(listMultipartUploadsResponse.Uploads)

		if partCount != 1 {
			t.Error(util.FormatTest(method, fmt.Sprintf("part count is %d", partCount), "part count should be 1"))
		}

		listMultipartUploadsRequest = ListMultipartUploadsRequest{
			BucketName: bucketName,
			Delimiter:  "/",
			Prefix:     objectKey,
			MaxUploads: 100,
			KeyMarker:  objectKey,
		}
		listMultipartUploadsResponse, err = bosClient.ListMultipartUploadsFromRequest(listMultipartUploadsRequest, nil)

		if err != nil {
			t.Error(util.FormatTest(method, err.Error(), "nil"))
			return
		}

		partCount = len(listMultipartUploadsResponse.Uploads)

		if partCount != 0 {
			t.Error(util.FormatTest(method, fmt.Sprintf("part count is %d", partCount), "part count should be 0"))
		}
	})
}

func TestSetBucketCors(t *testing.T) {
	bucketNamePrefix := "baidubce-sdk-go-test-for-set-bucket-cors-"
	method := "SetBucketCors"

	around(t, method, bucketNamePrefix, "", func(bucketName string) {
		bucketCors := BucketCors{
			CorsConfiguration: []BucketCorsItem{
				BucketCorsItem{
					AllowedOrigins:       []string{"http://*", "https://*"},
					AllowedMethods:       []string{"GET", "HEAD", "POST", "PUT"},
					AllowedHeaders:       []string{"*"},
					AllowedExposeHeaders: []string{"ETag", "x-bce-request-id", "Content-Type"},
					MaxAgeSeconds:        3600,
				},
				BucketCorsItem{
					AllowedOrigins:       []string{"http://www.example.com", "www.example2.com"},
					AllowedMethods:       []string{"GET", "HEAD", "DELETE"},
					AllowedHeaders:       []string{"Authorization", "x-bce-test", "x-bce-test2"},
					AllowedExposeHeaders: []string{"user-custom-expose-header"},
					MaxAgeSeconds:        2000,
				},
			},
		}

		if err := bosClient.SetBucketCors(bucketName, bucketCors, nil); err != nil {
			t.Error(util.FormatTest(method, err.Error(), "nil"))
		}
	})
}

func TestGetBucketCors(t *testing.T) {
	bucketNamePrefix := "baidubce-sdk-go-test-for-get-bucket-cors-"
	method := "GetBucketCors"

	around(t, method, bucketNamePrefix, "", func(bucketName string) {
		bucketCors := BucketCors{
			CorsConfiguration: []BucketCorsItem{
				BucketCorsItem{
					AllowedOrigins:       []string{"http://*", "https://*"},
					AllowedMethods:       []string{"GET", "HEAD", "POST", "PUT"},
					AllowedHeaders:       []string{"*"},
					AllowedExposeHeaders: []string{"ETag", "x-bce-request-id", "Content-Type"},
					MaxAgeSeconds:        3600,
				},
				BucketCorsItem{
					AllowedOrigins:       []string{"http://www.example.com", "www.example2.com"},
					AllowedMethods:       []string{"GET", "HEAD", "DELETE"},
					AllowedHeaders:       []string{"Authorization", "x-bce-test", "x-bce-test2"},
					AllowedExposeHeaders: []string{"user-custom-expose-header"},
					MaxAgeSeconds:        2000,
				},
			},
		}

		if err := bosClient.SetBucketCors(bucketName, bucketCors, nil); err != nil {
			t.Error(util.FormatTest(method+":SetBucketCors", err.Error(), "nil"))
		}

		_, err := bosClient.GetBucketCors(bucketName, nil)

		if err != nil {
			t.Error(util.FormatTest(method, err.Error(), "nil"))
		}
	})
}

func TestDeleteBucketCors(t *testing.T) {
	bucketNamePrefix := "baidubce-sdk-go-test-for-delete-bucket-cors-"
	method := "DeleteBucketCors"

	around(t, method, bucketNamePrefix, "", func(bucketName string) {
		err := bosClient.DeleteBucketCors(bucketName, nil)
		if err != nil {
			t.Error(util.FormatTest(method, err.Error(), "nil"))
		}
	})
}

func TestOptionsObject(t *testing.T) {
	bucketNamePrefix := "baidubce-sdk-go-test-for-options-object-"
	objectKey := "put-object-from-string.txt"
	str := "Hello World 你好"
	method := "OptionsObject"

	around(t, method, bucketNamePrefix, objectKey, func(bucketName string) {
		_, err := bosClient.PutObject(bucketName, objectKey, str, nil, nil)

		if err != nil {
			t.Error(util.FormatTest(method+":PutObject", err.Error(), "nil"))
		}

		bucketCors := BucketCors{
			CorsConfiguration: []BucketCorsItem{
				BucketCorsItem{
					AllowedOrigins:       []string{"http://*", "https://*"},
					AllowedMethods:       []string{"GET", "HEAD", "POST", "PUT"},
					AllowedHeaders:       []string{"*"},
					AllowedExposeHeaders: []string{"ETag", "x-bce-request-id", "Content-Type"},
					MaxAgeSeconds:        3600,
				},
				BucketCorsItem{
					AllowedOrigins:       []string{"http://www.example.com", "www.example2.com"},
					AllowedMethods:       []string{"GET", "HEAD", "DELETE"},
					AllowedHeaders:       []string{"Authorization", "x-bce-test", "x-bce-test2"},
					AllowedExposeHeaders: []string{"user-custom-expose-header"},
					MaxAgeSeconds:        2000,
				},
			},
		}

		if err := bosClient.SetBucketCors(bucketName, bucketCors, nil); err != nil {
			t.Error(util.FormatTest(method+":SetBucketCors", err.Error(), "nil"))
		}
		_, err = bosClient.OptionsObject(bucketName, objectKey, "http://www.example.com", "GET", "x-bce-test")

		if err != nil {
			t.Error(util.FormatTest(method, err.Error(), "nil"))
		}
	})
}

func TestSetBucketLogging(t *testing.T) {
	bucketNamePrefix := "baidubce-sdk-go-test-for-set-bucket-logging-"
	targetBucket := "baidubce-sdk-go-test-for-set-bucket-logging-logs" + strconv.Itoa(int(time.Now().UnixNano()))

	if len(targetBucket) > 63 {
		targetBucket = targetBucket[:63]
	}

	targetPrefix := "baidubce-sdk-go"
	method := "SetBucketLogging"

	around(t, method, bucketNamePrefix, "", func(bucketName string) {
		defer bosClient.DeleteBucket(targetBucket, nil)

		err := bosClient.CreateBucket(targetBucket, nil)

		if err != nil {
			t.Error(util.FormatTest(method+":CreateBucket", err.Error(), "nil"))
		} else {
			err = bosClient.SetBucketLogging(bucketName, targetBucket, targetPrefix, nil)

			if err != nil {
				t.Error(util.FormatTest(method, err.Error(), "nil"))
			}
		}
	})
}

func TestGetBucketLogging(t *testing.T) {
	bucketNamePrefix := "baidubce-sdk-go-test-for-get-bucket-logging-"
	targetBucket := "baidubce-sdk-go-test-for-get-bucket-logging-logs" + strconv.Itoa(int(time.Now().UnixNano()))

	if len(targetBucket) > 63 {
		targetBucket = targetBucket[:63]
	}

	targetPrefix := "baidubce-sdk-go"
	method := "GetBucketLogging"

	around(t, method, bucketNamePrefix, "", func(bucketName string) {
		defer bosClient.DeleteBucket(targetBucket, nil)

		err := bosClient.CreateBucket(targetBucket, nil)

		if err != nil {
			t.Error(util.FormatTest(method+":CreateBucket", err.Error(), "nil"))
		} else {
			err = bosClient.SetBucketLogging(bucketName, targetBucket, targetPrefix, nil)

			if err != nil {
				t.Error(util.FormatTest(method+":SetBucketLogging", err.Error(), "nil"))
			} else {
				bucketLogging, err := bosClient.GetBucketLogging(bucketName, nil)

				if err != nil {
					t.Error(util.FormatTest(method, err.Error(), "nil"))
				} else if bucketLogging.Status != "enabled" {
					t.Error(util.FormatTest(method, bucketLogging.Status, "enabled"))
				}
			}
		}
	})
}

func TestDeleteBucketLogging(t *testing.T) {
	bucketNamePrefix := "baidubce-sdk-go-test-for-delete-bucket-logging-"
	targetBucket := "baidubce-sdk-go-test-for-delete-bucket-logging-logs" + strconv.Itoa(int(time.Now().UnixNano()))

	if len(targetBucket) > 63 {
		targetBucket = targetBucket[:63]
	}

	targetPrefix := "baidubce-sdk-go"
	method := "DeleteBucketLogging"

	around(t, method, bucketNamePrefix, "", func(bucketName string) {
		defer bosClient.DeleteBucket(targetBucket, nil)

		err := bosClient.CreateBucket(targetBucket, nil)

		if err != nil {
			t.Error(util.FormatTest(method+":CreateBucket", err.Error(), "nil"))
		} else {
			err := bosClient.SetBucketLogging(bucketName, targetBucket, targetPrefix, nil)

			if err != nil {
				t.Error(util.FormatTest(method+":SetBucketLogging", err.Error(), "nil"))
			} else {
				err := bosClient.DeleteBucketLogging(bucketName, nil)

				if err != nil {
					t.Error(util.FormatTest(method, err.Error(), "nil"))
				}
			}
		}
	})
}

func TestSetBucketLifecycle(t *testing.T) {
	bucketNamePrefix := "baidubce-sdk-go-test-for-set-bucket-lifecycle-"
	method := "SetBucketLifecycle"

	around(t, method, bucketNamePrefix, "", func(bucketName string) {
		bucketLifecycle := BucketLifecycle{
			Rule: []BucketLifecycleItem{
				BucketLifecycleItem{
					Id:       "1",
					Status:   "disabled",
					Resource: []string{bucketName + "/*"},
					Condition: BucketLifecycleItemCondition{
						Time: BucketLifecycleItemConditionTime{DateGreaterThan: "2017-01-30T00:00:00Z"},
					},
					Action: BucketLifecycleItemAction{Name: "DeleteObject"},
				},
				BucketLifecycleItem{
					Id:       "2",
					Status:   "enabled",
					Resource: []string{bucketName + "/test/*"},
					Condition: BucketLifecycleItemCondition{
						Time: BucketLifecycleItemConditionTime{DateGreaterThan: "$(lastModified)+P7D"},
					},
					Action: BucketLifecycleItemAction{Name: "Transition", StorageClass: "STANDARD_IA"},
				},
				BucketLifecycleItem{
					Id:       "3",
					Status:   "enabled",
					Resource: []string{bucketName + "/multi/*"},
					Condition: BucketLifecycleItemCondition{
						Time: BucketLifecycleItemConditionTime{DateGreaterThan: "$(lastModified)+P7D"},
					},
					Action: BucketLifecycleItemAction{Name: "AbortMultipartUpload"},
				},
			},
		}

		if err := bosClient.SetBucketLifecycle(bucketName, bucketLifecycle, nil); err != nil {
			t.Error(util.FormatTest(method, err.Error(), "nil"))
		}
	})
}

func TestGetBucketLifecycle(t *testing.T) {
	bucketNamePrefix := "baidubce-sdk-go-test-for-get-bucket-lifecycle-"
	method := "GetBucketLifecycle"

	around(t, method, bucketNamePrefix, "", func(bucketName string) {
		bucketLifecycle := BucketLifecycle{
			Rule: []BucketLifecycleItem{
				BucketLifecycleItem{
					Id:       "1",
					Status:   "disabled",
					Resource: []string{bucketName + "/*"},
					Condition: BucketLifecycleItemCondition{
						Time: BucketLifecycleItemConditionTime{DateGreaterThan: "2017-01-30T00:00:00Z"},
					},
					Action: BucketLifecycleItemAction{Name: "DeleteObject"},
				},
				BucketLifecycleItem{
					Id:       "2",
					Status:   "enabled",
					Resource: []string{bucketName + "/test/*"},
					Condition: BucketLifecycleItemCondition{
						Time: BucketLifecycleItemConditionTime{DateGreaterThan: "$(lastModified)+P7D"},
					},
					Action: BucketLifecycleItemAction{Name: "Transition", StorageClass: "STANDARD_IA"},
				},
				BucketLifecycleItem{
					Id:       "3",
					Status:   "enabled",
					Resource: []string{bucketName + "/multi/*"},
					Condition: BucketLifecycleItemCondition{
						Time: BucketLifecycleItemConditionTime{DateGreaterThan: "$(lastModified)+P7D"},
					},
					Action: BucketLifecycleItemAction{Name: "AbortMultipartUpload"},
				},
			},
		}

		if err := bosClient.SetBucketLifecycle(bucketName, bucketLifecycle, nil); err != nil {
			t.Error(util.FormatTest(method, err.Error(), "nil"))
		} else {
			gotBucketLifecycle, err := bosClient.GetBucketLifecycle(bucketName, nil)

			if err != nil {
				t.Error(util.FormatTest(method, err.Error(), "nil"))
			} else {
				byteArrayOfBucketLifecycle, err := util.ToJson(bucketLifecycle, "rule")
				if err != nil {
					t.Error(util.FormatTest(method, err.Error(), "nil"))
				} else {
					byteArrayOfGotBucketLifecycle, err := util.ToJson(gotBucketLifecycle, "rule")
					if err != nil {
						t.Error(util.FormatTest(method, err.Error(), "nil"))
					} else if string(byteArrayOfBucketLifecycle) != string(byteArrayOfGotBucketLifecycle) {
						t.Error(util.FormatTest(method, string(byteArrayOfBucketLifecycle), string(byteArrayOfGotBucketLifecycle)))
					}
				}
			}
		}
	})
}

func TestDeleteBucketLifecycle(t *testing.T) {
	bucketNamePrefix := "baidubce-sdk-go-test-for-delete-bucket-lifecycle-"
	method := "DeleteBucketLifecycle"

	around(t, method, bucketNamePrefix, "", func(bucketName string) {
		bucketLifecycle := BucketLifecycle{
			Rule: []BucketLifecycleItem{
				BucketLifecycleItem{
					Id:       "1",
					Status:   "disabled",
					Resource: []string{bucketName + "/*"},
					Condition: BucketLifecycleItemCondition{
						Time: BucketLifecycleItemConditionTime{DateGreaterThan: "2017-01-30T00:00:00Z"},
					},
					Action: BucketLifecycleItemAction{Name: "DeleteObject"},
				},
				BucketLifecycleItem{
					Id:       "2",
					Status:   "enabled",
					Resource: []string{bucketName + "/test/*"},
					Condition: BucketLifecycleItemCondition{
						Time: BucketLifecycleItemConditionTime{DateGreaterThan: "$(lastModified)+P7D"},
					},
					Action: BucketLifecycleItemAction{Name: "Transition", StorageClass: "STANDARD_IA"},
				},
				BucketLifecycleItem{
					Id:       "3",
					Status:   "enabled",
					Resource: []string{bucketName + "/multi/*"},
					Condition: BucketLifecycleItemCondition{
						Time: BucketLifecycleItemConditionTime{DateGreaterThan: "$(lastModified)+P7D"},
					},
					Action: BucketLifecycleItemAction{Name: "AbortMultipartUpload"},
				},
			},
		}

		if err := bosClient.SetBucketLifecycle(bucketName, bucketLifecycle, nil); err != nil {
			t.Error(util.FormatTest(method, err.Error(), "nil"))
		} else {
			gotBucketLifecycle, err := bosClient.GetBucketLifecycle(bucketName, nil)

			if err != nil {
				t.Error(util.FormatTest(method, err.Error(), "nil"))
			} else if len(gotBucketLifecycle.Rule) != 3 {
				t.Error(util.FormatTest(method, strconv.Itoa(len(gotBucketLifecycle.Rule)), "3"))
			} else {
				err = bosClient.DeleteBucketLifecycle(bucketName, nil)
				if err != nil {
					t.Error(util.FormatTest(method, err.Error(), "nil"))
				} else {
					_, err := bosClient.GetBucketLifecycle(bucketName, nil)
					if err == nil {
						t.Error(util.FormatTest(method, "nil", "error"))
					} else if bceError, ok := err.(*bce.Error); ok && bceError.StatusCode != 404 {
						t.Error(util.FormatTest(method, err.Error(), "bce.Error with status code 404"))
					}
				}
			}
		}
	})
}

func TestPubObjectBySTS(t *testing.T) {
	bucketNamePrefix := "baidubce-sdk-go-test-for-put-object-by-sts-"
	method := "PutObject"
	objectKey := "put-object-from-string.txt"
	str := "Hello World 你好"

	around(t, method, bucketNamePrefix, objectKey, func(bucketName string) {
		req := bce.SessionTokenRequest{
			DurationSeconds: 600,
			Id:              "ef5a4b19192f4931adcf0e12f82795e2",
			AccessControlList: []bce.AccessControlListItem{
				bce.AccessControlListItem{
					Service:    "bce:bos",
					Region:     bceConfig.GetRegion(),
					Effect:     "Allow",
					Resource:   []string{bucketName + "/*"},
					Permission: []string{"READ", "WRITE"},
				},
			},
		}

		sessionTokenResponse, err := bosClient.GetSessionToken(req, nil)

		if err != nil {
			t.Error(util.FormatTest(method+":GetSessionToken", err.Error(), "nil"))
		} else {
			option := &bce.SignOption{
				Credentials: bce.NewCredentials(sessionTokenResponse.AccessKeyId, sessionTokenResponse.SecretAccessKey),
				Headers:     map[string]string{"x-bce-security-token": sessionTokenResponse.SessionToken},
			}

			_, err := bosClient.PutObject(bucketName, objectKey, str, nil, option)

			if err != nil {
				t.Error(util.FormatTest(method, err.Error(), "nil"))
			}
		}
	})
}

func around(t *testing.T, method, bucketNamePrefix string, objectKey interface{}, f func(string)) {
	bucketName := bucketNamePrefix + strconv.Itoa(int(time.Now().UnixNano()))
	if len(bucketName) > 63 {
		bucketName = bucketName[:63]
	}

	err := bosClient.CreateBucket(bucketName, nil)

	if err != nil {
		t.Error(util.FormatTest(method+" at creating bucket", err.Error(), "nil"))
	} else {
		defer func() {
			if key, ok := objectKey.(string); ok {
				if key != "" {
					err = bosClient.DeleteObject(bucketName, key, nil)

					if bceError, ok := err.(*bce.Error); ok && bceError.StatusCode != 404 {
						t.Error(util.FormatTest(method+" at deleting object", err.Error(), "nil"))
					}
				}
			} else if keys, ok := objectKey.([]string); ok {
				deleteMultipleObjectsResponse, err := bosClient.DeleteMultipleObjects(bucketName, keys, nil)

				if err != nil {
					t.Error(util.FormatTest(method, err.Error(), "nil"))
				} else if deleteMultipleObjectsResponse != nil {
					str := ""

					for _, deleteMultipleObjectsError := range deleteMultipleObjectsResponse.Errors {
						if deleteMultipleObjectsError.Code != "NoSuchKey" {
							str += deleteMultipleObjectsError.Error()
						}
					}

					if str != "" {
						t.Error(util.FormatTest(method, str, "empty string"))
					}
				}
			} else {
				t.Error(util.FormatTest(method, "objectKey is not valid", "objectKey should be string or []string"))
			}

			err = bosClient.DeleteBucket(bucketName, nil)

			if err != nil {
				t.Error(util.FormatTest(method+" at deleting bucket", err.Error(), "nil"))
			}
		}()

		if f != nil {
			f(bucketName)
		}
	}
}
