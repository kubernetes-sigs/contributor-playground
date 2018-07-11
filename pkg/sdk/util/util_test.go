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
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestGetURL(t *testing.T) {
	expected := "http://bos.cn-n1.baidubce.com/v1/example/测试"
	url := GetURL("", "bos.cn-n1.baidubce.com", "v1/example/测试", nil)

	if url != expected {
		t.Error(FormatTest("GetURL", url, expected))
	}

	url = GetURL("", "bos.cn-n1.baidubce.com", "/v1/example/测试", nil)

	if url != expected {
		t.Error(FormatTest("GetURL", url, expected))
	}

	expected = expected + "?id=123&name=abc"
	params := map[string]string{"id": "123", "name": "abc"}
	url = GetURL("", "bos.cn-n1.baidubce.com", "v1/example/测试", params)

	if url != expected {
		t.Error(FormatTest("GetURL", url, expected))
	}

	expected = "https://bos.cn-n1.baidubce.com/v1/example/测试" + "?id=123&name=abc"
	url = GetURL("https", "bos.cn-n1.baidubce.com", "v1/example/测试", params)

	if url != expected {
		t.Error(FormatTest("GetURL", url, expected))
	}
}

func TestGetURIPath(t *testing.T) {
	const URI string = "http://bos.cn-n1.baidubce.com/v1/example/测试"
	expected := "/v1/example/测试"
	path := GetURIPath(URI)

	if path != expected {
		t.Error(FormatTest("GetURIPath", path, expected))
	}
}

func TestURIEncodeExceptSlash(t *testing.T) {
	const URI string = "http://bos.cn-n1.baidubce.com/v1/example/测试"
	expected := "/v1/example/%E6%B5%8B%E8%AF%95"
	path := GetURIPath(URI)
	path = URIEncodeExceptSlash(path)

	if path != expected {
		t.Error(FormatTest("URIEncodeExceptSlash", path, expected))
	}
}

func TestHmacSha256Hex(t *testing.T) {
	expected := "6e9ef29b75fffc5b7abae527d58fdadb2fe42e7219011976917343065f58ed4a"
	encrypted := HmacSha256Hex("key", "message")

	if encrypted != expected {
		t.Error(FormatTest("HmacSha256Hex", encrypted, expected))
	}
}

func TestGetMD5(t *testing.T) {
	expected := "de22e061b93b832dd8af907ca9002fd7"
	result := GetMD5("baidubce-sdk-go", false)

	if result != expected {
		t.Error(FormatTest("GetMD5", result, expected))
	}

	result = GetMD5([]byte("baidubce-sdk-go"), false)

	if result != expected {
		t.Error(FormatTest("GetMD5", result, expected))
	}

	result = GetMD5(strings.NewReader("baidubce-sdk-go"), false)

	if result != expected {
		t.Error(FormatTest("GetMD5", result, expected))
	}

	f, err := TempFile([]byte("baidubce-sdk-go"), "", "")

	if err != nil {
		t.Error(FormatTest("GetMD5", err.Error(), "nil"))
	} else {
		defer func() {
			f.Close()
			os.Remove(f.Name())
		}()

		result = GetMD5(f, false)

		if result != expected {
			t.Error(FormatTest("GetMD5", result, expected))
		}
	}

	result = GetMD5(bufio.NewReader(strings.NewReader("baidubce-sdk-go")), false)

	if result != expected {
		t.Error(FormatTest("GetMD5", result, expected))
	}

	result = GetMD5("baidubce-sdk-go", true)
	expected = "3iLgYbk7gy3Yr5B8qQAv1w=="

	if result != expected {
		t.Error(FormatTest("GetMD5", result, expected))
	}

	defer func() {
		if err := recover(); err == nil {
			t.Error(FormatTest("GetMD5", "nil", "error"))
		}
	}()

	result = GetMD5(1, false)
}

func TestGetSha256(t *testing.T) {
	expected := "b39aa8e24bcfc4b20c77f7ab36021e5c23cce79df034279ca9991e0472368b89"
	result := GetSha256("baidubce-sdk-go")

	if result != expected {
		t.Error(FormatTest("GetSha256", result, expected))
	}

	result = GetSha256([]byte("baidubce-sdk-go"))

	if result != expected {
		t.Error(FormatTest("GetSha256", result, expected))
	}

	result = GetSha256(strings.NewReader("baidubce-sdk-go"))

	if result != expected {
		t.Error(FormatTest("GetSha256", result, expected))
	}

	f, err := TempFile([]byte("baidubce-sdk-go"), "", "")

	if err != nil {
		t.Error(FormatTest("GetSha256", err.Error(), "nil"))
	} else {
		defer func() {
			f.Close()
			os.Remove(f.Name())
		}()

		result = GetSha256(f)

		if result != expected {
			t.Error(FormatTest("GetSha256", result, expected))
		}
	}

	result = GetSha256(bufio.NewReader(strings.NewReader("baidubce-sdk-go")))

	if result != expected {
		t.Error(FormatTest("GetSha256", result, expected))
	}

	defer func() {
		if err := recover(); err == nil {
			t.Error(FormatTest("GetSha256", "nil", "error"))
		}
	}()

	result = GetSha256(1)
}

func TestBase64Encode(t *testing.T) {
	expected := "YmFpZHViY2Utc2RrLWdv"
	result := Base64Encode([]byte("baidubce-sdk-go"))

	if result != expected {
		t.Error(FormatTest("Base64Encode", result, expected))
	}
}

func TestContains(t *testing.T) {
	expected := true
	arr := []string{"abc", "XYz"}
	result := Contains(arr, "abc", true)

	if result != expected {
		t.Error(FormatTest("Contains", strconv.FormatBool(result), strconv.FormatBool(expected)))
	}

	expected = false
	result = Contains(arr, "Abc", false)

	if result != expected {
		t.Error(FormatTest("Contains", strconv.FormatBool(result), strconv.FormatBool(expected)))
	}

	expected = true
	result = Contains(arr, "xyz", true)

	if result != expected {
		t.Error(FormatTest("Contains", strconv.FormatBool(result), strconv.FormatBool(expected)))
	}

	result = Contains(arr, "Xyz", true)

	if result != expected {
		t.Error(FormatTest("Contains", strconv.FormatBool(result), strconv.FormatBool(expected)))
	}

	result = Contains(arr, "xYZ", true)

	if result != expected {
		t.Error(FormatTest("Contains", strconv.FormatBool(result), strconv.FormatBool(expected)))
	}

	result = Contains(arr, "XYZ", true)

	if result != expected {
		t.Error(FormatTest("Contains", strconv.FormatBool(result), strconv.FormatBool(expected)))
	}
}

func TestMapContains(t *testing.T) {
	expected := true
	m := map[string]string{"id": "123", "name": "Matt"}
	result := MapContains(m, func(key, value string) bool {
		return key == "id"
	})

	if result != expected {
		t.Error(FormatTest("MapContains", strconv.FormatBool(result), strconv.FormatBool(expected)))
	}

	result = MapContains(m, func(key, value string) bool {
		return value == "123"
	})

	if result != expected {
		t.Error(FormatTest("MapContains", strconv.FormatBool(result), strconv.FormatBool(expected)))
	}

	expected = false
	result = MapContains(m, func(key, value string) bool {
		return value == "matt"
	})

	if result != expected {
		t.Error(FormatTest("MapContains", strconv.FormatBool(result), strconv.FormatBool(expected)))
	}
}

func TestGetMapKey(t *testing.T) {
	expected := "id"
	m := map[string]string{"id": "123", "Name": "Matt"}
	result := GetMapKey(m, "id", true)

	if result != expected {
		t.Error(FormatTest("GetMapKey", result, expected))
	}

	result = GetMapKey(m, "id", false)

	if result != expected {
		t.Error(FormatTest("GetMapKey", result, expected))
	}

	result = GetMapKey(m, "Id", true)

	if result != expected {
		t.Error(FormatTest("GetMapKey", result, expected))
	}

	expected = ""
	result = GetMapKey(m, "Id", false)

	if result != expected {
		t.Error(FormatTest("GetMapKey", result, expected))
	}

	expected = ""
	result = GetMapKey(m, "age", true)

	if result != expected {
		t.Error(FormatTest("GetMapKey", result, expected))
	}

	expected = ""
	result = GetMapKey(m, "age", false)

	if result != expected {
		t.Error(FormatTest("GetMapKey", result, expected))
	}
}

func TestGetMapValue(t *testing.T) {
	expected := "123"
	m := map[string]string{"id": "123", "Name": "Matt"}
	result := GetMapValue(m, "id", true)

	if result != expected {
		t.Error(FormatTest("GetMapValue", result, expected))
	}

	result = GetMapValue(m, "id", false)

	if result != expected {
		t.Error(FormatTest("GetMapValue", result, expected))
	}

	result = GetMapValue(m, "Id", true)

	if result != expected {
		t.Error(FormatTest("GetMapValue", result, expected))
	}

	expected = ""
	result = GetMapValue(m, "Id", false)

	if result != expected {
		t.Error(FormatTest("GetMapValue", result, expected))
	}
}

func TestTimeToUTCString(t *testing.T) {
	expected := "2015-11-16T15:33:15Z"
	datetime, _ := time.Parse(time.RFC1123, "Mon, 16 Nov 2015 15:33:15 CST")
	_, offset := datetime.Zone()
	offset = offset / 3600
	datetime = datetime.Add(time.Duration(offset) * time.Hour)
	utc := TimeToUTCString(datetime)

	if utc != expected {
		t.Error(FormatTest("TimeToUTCString", utc, expected))
	}
}

func TestTimeStringToRFC1123(t *testing.T) {
	expected := "Mon, 16 Nov 2015 07:33:15 UTC"
	result := TimeStringToRFC1123("2015-11-16T07:33:15Z")

	if result != expected {
		t.Error(FormatTest("TimeStringToRFC1123", result, expected))
	}

	defer func() {
		if err := recover(); err == nil {
			t.Error(FormatTest("TimeStringToRFC1123", "nil", "error"))
		}
	}()

	TimeStringToRFC1123("Invalid")
}

func TestHostToURL(t *testing.T) {
	expected := "http://bj.bcebos.com"
	host := "bj.bcebos.com"
	url := HostToURL(host, "http")

	if url != expected {
		t.Error(FormatTest("HostToURL", url, expected))
	}

	host = "http://bj.bcebos.com"
	url = HostToURL(host, "http")

	if url != expected {
		t.Error(FormatTest("HostToURL", url, expected))
	}

	host = "bj.bcebos.com"
	url = HostToURL(host, "")

	if url != expected {
		t.Error(FormatTest("HostToURL", url, expected))
	}

	host = "http://bj.bcebos.com"
	url = HostToURL(host, "")

	if url != expected {
		t.Error(FormatTest("HostToURL", url, expected))
	}
}

func TestToCanonicalQueryString(t *testing.T) {
	expected := ""
	encodedQueryString := ToCanonicalQueryString(nil)

	if encodedQueryString != expected {
		t.Error(FormatTest("ToCanonicalQueryString", encodedQueryString, expected))
	}

	expected = "text10=test&text1=%E6%B5%8B%E8%AF%95&text="
	params := map[string]string{
		"text":   "",
		"text1":  "测试",
		"text10": "test",
	}
	encodedQueryString = ToCanonicalQueryString(params)

	if encodedQueryString != expected {
		t.Error(FormatTest("ToCanonicalQueryString", encodedQueryString, expected))
	}
}

func TestToCanonicalHeaderString(t *testing.T) {
	expected := strings.Join([]string{
		"content-length:8",
		"content-md5:0a52730597fb4ffa01fc117d9e71e3a9",
		"content-type:text%2Fplain",
		"host:bj.bcebos.com",
		"x-bce-date:2015-04-27T08%3A23%3A49Z",
	}, "\n")

	header := map[string]string{
		"Host":           "bj.bcebos.com",
		"Content-Type":   "text/plain",
		"Content-Length": "8",
		"Content-Md5":    "0a52730597fb4ffa01fc117d9e71e3a9",
		"x-bce-date":     "2015-04-27T08:23:49Z",
	}

	canonicalHeader := ToCanonicalHeaderString(header)

	if canonicalHeader != expected {
		t.Error(FormatTest("ToCanonicalHeaderString", canonicalHeader, expected))
	}
}

func TestURLEncode(t *testing.T) {
	expected := "test-%E6%B5%8B%E8%AF%95"
	result := URLEncode("test-测试")

	if result != expected {
		t.Error(FormatTest("URLEncode", result, expected))
	}
}

func TestSliceToLower(t *testing.T) {
	expected := "name age"
	arr := []string{"Name", "Age"}
	SliceToLower(arr)

	result := fmt.Sprintf("%s %s", arr[0], arr[1])

	if result != expected {
		t.Error(FormatTest("SliceToLower", result, expected))
	}
}

func TestMapKeyToLower(t *testing.T) {
	expected := "name gender"
	m := map[string]string{"Name": "guoyao", "Gender": "male"}
	MapKeyToLower(m)

	result := ""

	if _, ok := m["name"]; ok {
		result += "name"
	}

	if _, ok := m["gender"]; ok {
		result += " gender"
	}

	if result != expected {
		t.Error(FormatTest("MapKeyToLower", result, expected))
	}
}

func TestToMap(t *testing.T) {
	expected := "guoyao:10"

	str := "{\"Name\": \"guoyao\", \"Age\": \"10\", \"Gender\": \"male\"}"
	m, err := ToMap(str, "Name", "Age")

	if err != nil {
		t.Error(FormatTest("ToMap", err.Error(), "nil"))
	} else {
		result := fmt.Sprintf("%s:%v", m["Name"], m["Age"])

		if result != expected {
			t.Error(FormatTest("ToMap", result, expected))
		}
	}

	byteArray := []byte(str)
	m, err = ToMap(byteArray, "Name", "Age")

	if err != nil {
		t.Error(FormatTest("ToMap", err.Error(), "nil"))
	} else {
		result := fmt.Sprintf("%s:%v", m["Name"], m["Age"])

		if result != expected {
			t.Error(FormatTest("ToMap", result, expected))
		}
	}

	p := struct {
		Name   string
		Age    int
		Gender string
	}{"guoyao", 10, "male"}

	m, err = ToMap(p, "Name", "Age")

	if err != nil {
		t.Error(FormatTest("ToMap", err.Error(), "nil"))
	} else {
		result := fmt.Sprintf("%s:%v", m["Name"], m["Age"])

		if result != expected {
			t.Error(FormatTest("ToMap", result, expected))
		}
	}

	m, err = ToMap(1)

	if err == nil {
		t.Error(FormatTest("ToMap", "nil", "error"))
	}
}

func TestToJson(t *testing.T) {
	p := struct {
		Name   string `json:"name"`
		Age    int    `json:"age"`
		Gender string `json:"gender"`
	}{"guoyao", 10, "male"}

	byteArray, err := ToJson(p)

	if err != nil {
		t.Error(FormatTest("ToMap", err.Error(), "nil"))
	} else {
		expected := "{\"name\":\"guoyao\",\"age\":10,\"gender\":\"male\"}"
		result := string(byteArray)

		if result != expected {
			t.Error(FormatTest("ToMap", result, expected))
		}
	}

	byteArray, err = ToJson(p, "name", "age")

	if err != nil {
		t.Error(FormatTest("ToMap", err.Error(), "nil"))
	} else {
		expected := "{\"age\":10,\"name\":\"guoyao\"}"
		result := string(byteArray)

		if result != expected {
			t.Error(FormatTest("ToMap", result, expected))
		}
	}

	byteArray, err = ToJson(1, "name")

	if err == nil {
		t.Error(FormatTest("ToJson", "nil", "error"))
	}
}

func TestCheckFileExists(t *testing.T) {
	expected := true
	result := CheckFileExists("util_test.go")

	if result != expected {
		t.Error(FormatTest("CheckFileExists", strconv.FormatBool(result), strconv.FormatBool(expected)))
	}

	expected = false
	result = CheckFileExists("util_test_2.go")

	if result != expected {
		t.Error(FormatTest("CheckFileExists", strconv.FormatBool(result), strconv.FormatBool(expected)))
	}
}

func TestTempFileWithSize(t *testing.T) {
	var size int64 = 1024
	f, err := TempFileWithSize(size)

	defer func() {
		f.Close()
		os.Remove(f.Name())
	}()

	if err != nil {
		t.Error(FormatTest("TempFileWithSize", err.Error(), "nil"))
	} else {
		stat, err := f.Stat()

		if err != nil {
			t.Error(FormatTest("TempFileWithSize", err.Error(), "nil"))
		} else if stat.Size() != size {
			t.Error(FormatTest("TempFileWithSize", strconv.FormatInt(stat.Size(), 10), strconv.FormatInt(size, 10)))
		}
	}
}

func TestTempFile(t *testing.T) {
	content := "hello"
	f, err := TempFile([]byte(content), "", "")

	defer func() {
		f.Close()
		os.Remove(f.Name())
	}()

	if err != nil {
		t.Error(FormatTest("TempFile", err.Error(), "nil"))
	} else {
		byteArray, err := ioutil.ReadAll(f)

		if err != nil {
			t.Error(FormatTest("TempFile", err.Error(), "nil"))
		} else if string(byteArray) != content {
			t.Error(FormatTest("TempFile", string(byteArray), content))
		}
	}

	pwd, err := os.Getwd()

	if err != nil {
		t.Error(FormatTest("TempFile", err.Error(), "nil"))
	} else {
		content = "world"
		dir := path.Join(pwd, "guoyao")
		f2, err := TempFile([]byte(content), dir, "temp")

		defer func() {
			f2.Close()
			os.RemoveAll(dir)
		}()

		if err != nil {
			t.Error(FormatTest("TempFile", err.Error(), "nil"))
		} else {
			byteArray, err := ioutil.ReadAll(f2)

			if err != nil {
				t.Error(FormatTest("TempFile", err.Error(), "nil"))
			} else if string(byteArray) != content {
				t.Error(FormatTest("TempFile", string(byteArray), content))
			}
		}
	}
}

func TestHomeDir(t *testing.T) {
	home, err := HomeDir()

	if err != nil {
		t.Error(FormatTest("HomeDir", err.Error(), "nil"))
	} else if home == "" {
		t.Error(FormatTest("HomeDir", home, "non empty path"))
	}
}

func TestDirUnix(t *testing.T) {
	home, err := dirUnix()

	if err != nil {
		t.Error(FormatTest("dirUnix", err.Error(), "nil"))
	} else if home == "" {
		t.Error(FormatTest("dirUnix", home, "non empty path"))
	} else {
		os.Setenv("HOME", "")
		home2, err := dirUnix()
		os.Setenv("HOME", home)

		if err != nil {
			t.Error(FormatTest("dirUnix", err.Error(), "nil"))
		} else if home2 == "" {
			t.Error(FormatTest("dirUnix", home2, "non empty path"))
		}
	}
}

func TestDirWindows(t *testing.T) {
	home, err := dirWindows()

	if err != nil {
		t.Error(FormatTest("dirWindows", err.Error(), "nil"))
	} else if home == "" {
		t.Error(FormatTest("dirWindows", home, "non empty path"))
	} else {
		os.Setenv("HOME", "")
		os.Setenv("USERPROFILE", "C:\\Users\\guoyao")
		home2, err := dirWindows()
		os.Setenv("HOME", home)
		os.Setenv("USERPROFILE", "")

		if err != nil {
			t.Error(FormatTest("dirWindows", err.Error(), "nil"))
		} else if home2 == "" {
			t.Error(FormatTest("dirWindows", home2, "non empty path"))
		}
	}
}

func TestDebug(t *testing.T) {
	Debug("title", "message")
}

func TestFormatTest(t *testing.T) {
	expected := "funcName failed. Got a, expected b"
	str := FormatTest("funcName", "a", "b")

	if str != expected {
		t.Error(FormatTest("FormatTest", str, expected))
	}
}

func TestCreateRandomString(t *testing.T) {
	for i := 0; i < 10; i++ {
		s := CreateRandomString()
		t.Logf("Generated Random String: %s", s)
	}
}
