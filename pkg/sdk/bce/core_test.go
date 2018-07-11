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
	"errors"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	"fmt"

	"k8s.io/cloud-provider-baiducloud/pkg/sdk/util"
)

func getConfig() *Config {
	return &Config{
		Credentials: NewCredentials(os.Getenv("BAIDU_BCE_AK"), os.Getenv("BAIDU_BCE_SK")),
		Checksum:    true,
		Region:      os.Getenv("BOS_REGION"),
	}
}

var bceClient = NewClient(getConfig())

func TestNewConfig(t *testing.T) {
	credentials := Credentials{
		AccessKeyID:     "0b0f67dfb88244b289b72b142befad0c",
		SecretAccessKey: "bad522c2126a4618a8125f4b6cf6356f",
	}
	config := NewConfig(&credentials)

	if config == nil {
		t.Error(util.FormatTest("NewConfig", "nil", "not nil"))
	}
}

func TestGetRegion(t *testing.T) {
	expected := os.Getenv("BOS_REGION")

	if expected == "" {
		expected = "bj"
	}

	config := &Config{Region: os.Getenv("BOS_REGION")}
	region := config.GetRegion()

	if region != expected {
		t.Error(util.FormatTest("GetRegion", region, expected))
	}

	expected = "bj"
	config = &Config{}
	region = config.GetRegion()

	if region != expected {
		t.Error(util.FormatTest("GetRegion", region, expected))
	}
}

func TestGetUserAgent(t *testing.T) {
	expected := DefaultUserAgent
	config := &Config{}
	userAgent := config.GetUserAgent()

	if userAgent != expected {
		t.Error(util.FormatTest("GetUserAgent", userAgent, expected))
	}

	expected = "baidubce-sdk-go"
	config = &Config{UserAgent: "baidubce-sdk-go"}
	userAgent = config.GetUserAgent()

	if userAgent != expected {
		t.Error(util.FormatTest("GetUserAgent", userAgent, expected))
	}
}

func TestGetMaxErrorRetry(t *testing.T) {
	expected := 3
	retryPolicy := NewDefaultRetryPolicy(expected, 20*time.Second)

	if retryPolicy.MaxErrorRetry != expected {
		t.Error(util.FormatTest("GetMaxErrorRetry", strconv.Itoa(retryPolicy.MaxErrorRetry), strconv.Itoa(expected)))
	}
}

func TestGetMaxDelay(t *testing.T) {
	expected := 20 * time.Second
	retryPolicy := NewDefaultRetryPolicy(3, expected)

	if retryPolicy.MaxDelay != expected {
		t.Error(util.FormatTest("GetMaxDelay", retryPolicy.MaxDelay.String(), expected.String()))
	}
}

func TestGetDelayBeforeNextRetry(t *testing.T) {
	maxErrorRetry := 3
	maxDelay := 20 * time.Second
	retryPolicy := NewDefaultRetryPolicy(maxErrorRetry, maxDelay)
	delay := retryPolicy.GetDelayBeforeNextRetry(errors.New("Unknown Error"), 5)

	if delay != -1 {
		t.Error(util.FormatTest("GetDelayBeforeNextRetry", delay.String(), strconv.Itoa(-1)))
	}

	delay = retryPolicy.GetDelayBeforeNextRetry(errors.New("Unknown Error"), 1)

	if delay != -1 {
		t.Error(util.FormatTest("GetDelayBeforeNextRetry", delay.String(), strconv.Itoa(-1)))
	}
	delay = retryPolicy.GetDelayBeforeNextRetry(&Error{StatusCode: http.StatusInternalServerError}, 1)
	expected := (1 << 1) * 300 * time.Millisecond

	if delay != expected {
		t.Error(util.FormatTest("GetDelayBeforeNextRetry", delay.String(), expected.String()))
	}

	delay = retryPolicy.GetDelayBeforeNextRetry(&Error{StatusCode: http.StatusServiceUnavailable}, 2)
	expected = (1 << 2) * 300 * time.Millisecond

	if delay != expected {
		t.Error(util.FormatTest("GetDelayBeforeNextRetry", delay.String(), expected.String()))
	}

	maxDelay = 1 * time.Second
	retryPolicy = NewDefaultRetryPolicy(maxErrorRetry, maxDelay)
	delay = retryPolicy.GetDelayBeforeNextRetry(&Error{StatusCode: http.StatusServiceUnavailable}, 2)
	expected = retryPolicy.GetMaxDelay()

	if delay != expected {
		t.Error(util.FormatTest("GetDelayBeforeNextRetry", delay.String(), expected.String()))
	}
}

func TestShouldRetry(t *testing.T) {
	maxErrorRetry := 3
	maxDelay := 20 * time.Second
	retryPolicy := NewDefaultRetryPolicy(maxErrorRetry, maxDelay)
	expected := false
	unknownError := errors.New("Unknown Error")
	shouldRetry := retryPolicy.shouldRetry(unknownError, maxErrorRetry+1)

	if shouldRetry != expected {
		t.Error(util.FormatTest("shouldRetry", strconv.FormatBool(shouldRetry), strconv.FormatBool(expected)))
	}

	shouldRetry = retryPolicy.shouldRetry(unknownError, maxErrorRetry-1)

	if shouldRetry != expected {
		t.Error(util.FormatTest("shouldRetry", strconv.FormatBool(shouldRetry), strconv.FormatBool(expected)))
	}

	bceError := &Error{}
	shouldRetry = retryPolicy.shouldRetry(bceError, maxErrorRetry-1)

	if shouldRetry != expected {
		t.Error(util.FormatTest("shouldRetry", strconv.FormatBool(shouldRetry), strconv.FormatBool(expected)))
	}

	expected = true
	bceError.StatusCode = http.StatusInternalServerError
	shouldRetry = retryPolicy.shouldRetry(bceError, maxErrorRetry-1)

	if shouldRetry != expected {
		t.Error(util.FormatTest("shouldRetry", strconv.FormatBool(shouldRetry), strconv.FormatBool(expected)))
	}

	bceError.StatusCode = http.StatusServiceUnavailable
	shouldRetry = retryPolicy.shouldRetry(bceError, maxErrorRetry-1)

	if shouldRetry != expected {
		t.Error(util.FormatTest("shouldRetry", strconv.FormatBool(shouldRetry), strconv.FormatBool(expected)))
	}
}

func TestCheckSignOption(t *testing.T) {
	signOption := CheckSignOption(nil)

	if signOption == nil {
		t.Error(util.FormatTest("CheckSignOption", "nil", "not nil"))
	}

	originSignOption := NewSignOption(
		"2015-04-27T08:23:49Z",
		ExpirationPeriodInSeconds,
		getHeaders(),
		nil,
	)
	signOption = CheckSignOption(originSignOption)

	if signOption != originSignOption {
		t.Error(util.FormatTest("CheckSignOption", "new signOption", "originSignOption"))
	}
}

func TestAddHeadersToSign(t *testing.T) {
	signOption := NewSignOption(
		"2015-04-27T08:23:49Z",
		ExpirationPeriodInSeconds,
		nil,
		nil,
	)
	signOption.AddHeadersToSign("content-type", "content-length")

	if len(signOption.HeadersToSign) != 2 {
		t.Error(util.FormatTest("AddHeadersToSign", strconv.Itoa(len(signOption.HeadersToSign)), strconv.Itoa(2)))
	}

	if signOption.HeadersToSign[0] != "content-type" {
		t.Error(util.FormatTest("AddHeadersToSign", signOption.HeadersToSign[0], "content-type"))
	}

	signOption = NewSignOption(
		"2015-04-27T08:23:49Z",
		ExpirationPeriodInSeconds,
		nil,
		[]string{"content-type"},
	)
	signOption.AddHeadersToSign("content-type", "content-length")

	if len(signOption.HeadersToSign) != 2 {
		t.Error(util.FormatTest("AddHeadersToSign", strconv.Itoa(len(signOption.HeadersToSign)), strconv.Itoa(2)))
	}
}

func TestAddHeader(t *testing.T) {
	signOption := NewSignOption(
		"2015-04-27T08:23:49Z",
		ExpirationPeriodInSeconds,
		nil,
		nil,
	)
	signOption.AddHeader("content-type", "text/plain")
	signOption.AddHeader("content-length", "1024")

	if len(signOption.Headers) != 2 {
		t.Error(util.FormatTest("AddHeader", strconv.Itoa(len(signOption.Headers)), strconv.Itoa(2)))
	}

	if signOption.Headers["content-type"] != "text/plain" {
		t.Error(util.FormatTest("AddHeader", signOption.Headers["content-type"], "text/plain"))
	}

	signOption = NewSignOption(
		"2015-04-27T08:23:49Z",
		ExpirationPeriodInSeconds,
		map[string]string{"content-type": "text/plain"},
		nil,
	)
	originHeaders := signOption.Headers
	signOption.AddHeader("content-type", "text/plain")
	signOption.AddHeader("content-length", "1024")

	if len(signOption.Headers) != 2 {
		t.Error(util.FormatTest("AddHeader", strconv.Itoa(len(signOption.Headers)), strconv.Itoa(2)))
	}

	if len(originHeaders) != len(signOption.Headers) {
		t.Error(util.FormatTest("AddHeader", strconv.Itoa(len(signOption.Headers)), strconv.Itoa(len(originHeaders))))
	}
}

func TestAddHeaders(t *testing.T) {
	signOption := NewSignOption(
		"2015-04-27T08:23:49Z",
		ExpirationPeriodInSeconds,
		nil,
		nil,
	)

	var headers map[string]string
	signOption.AddHeaders(headers)

	if len(signOption.Headers) != len(headers) {
		t.Error(util.FormatTest("AddHeaders", strconv.Itoa(len(signOption.Headers)), strconv.Itoa(len(headers))))
	}

	headers = map[string]string{
		"content-type":   "text/plain",
		"content-length": "1024",
	}
	signOption.AddHeaders(headers)

	if len(signOption.Headers) != len(headers) {
		t.Error(util.FormatTest("AddHeaders", strconv.Itoa(len(signOption.Headers)), strconv.Itoa(len(headers))))
	}

	signOption = NewSignOption(
		"2015-04-27T08:23:49Z",
		ExpirationPeriodInSeconds,
		map[string]string{"content-type": "text/plain"},
		nil,
	)
	signOption.AddHeaders(headers)

	if len(signOption.Headers) != len(headers) {
		t.Error(util.FormatTest("AddHeaders", strconv.Itoa(len(signOption.Headers)), strconv.Itoa(len(headers))))
	}

	signOption = NewSignOption(
		"2015-04-27T08:23:49Z",
		ExpirationPeriodInSeconds,
		map[string]string{"content-encoding": "gzip"},
		nil,
	)
	signOption.AddHeaders(headers)

	if len(signOption.Headers) != len(headers)+1 {
		t.Error(util.FormatTest("AddHeaders", strconv.Itoa(len(signOption.Headers)), strconv.Itoa(len(headers)+1)))
	}
}

func TestInit(t *testing.T) {
	signOption := &SignOption{}

	signOption.init()

	if signOption.initialized != true {
		t.Error(util.FormatTest("init", strconv.FormatBool(signOption.initialized), strconv.FormatBool(true)))
	}

	timestamp := signOption.Timestamp

	if timestamp == "" {
		t.Error(util.FormatTest("init", timestamp, "timestamp string"))
	}

	if signOption.ExpirationPeriodInSeconds != ExpirationPeriodInSeconds {
		t.Error(util.FormatTest("init", strconv.Itoa(signOption.ExpirationPeriodInSeconds), strconv.Itoa(ExpirationPeriodInSeconds)))
	}

	signOption.init()

	if signOption.Timestamp != timestamp {
		t.Error(util.FormatTest("init", signOption.Timestamp, timestamp))
	}

	signOption = &SignOption{
		Headers:       map[string]string{"content-type": "text/plain"},
		HeadersToSign: []string{"content-type", "date"},
	}
	signOption.init()

	if !util.Contains(signOption.HeadersToSign, "host", true) {
		t.Error(util.FormatTest("init", "", "host should be specified"))
	}

	if !util.MapContains(signOption.Headers, func(key, value string) bool {
		return key == "date"
	}) {
		t.Error(util.FormatTest("init", "", "date should be specified"))
	}

	if util.MapContains(signOption.Headers, func(key, value string) bool {
		return key == "x-bce-date"
	}) {
		t.Error(util.FormatTest("init", "x-bce-date", "x-bce-date should be specified"))
	}

	signOption = &SignOption{
		Headers:       map[string]string{"content-type": "text/plain", "date": "2015-11-16T07:33:15Z"},
		HeadersToSign: []string{"content-type", "date"},
	}
	expectedDate := util.TimeStringToRFC1123(signOption.Headers["date"])
	signOption.init()

	if signOption.Headers["date"] != expectedDate {
		t.Error(util.FormatTest("init", signOption.Headers["date"], expectedDate))
	}

	signOption = &SignOption{
		Headers:       map[string]string{"content-type": "text/plain"},
		HeadersToSign: []string{"content-type"},
	}
	signOption.init()

	if util.MapContains(signOption.Headers, func(key, value string) bool {
		return key == "date"
	}) {
		t.Error(util.FormatTest("init", "date", "date should not be specified"))
	}

	if !util.MapContains(signOption.Headers, func(key, value string) bool {
		return key == "x-bce-date"
	}) {
		t.Error(util.FormatTest("init", "", "x-bce-date should not be specified"))
	}
}

func TestSignedHeadersToString(t *testing.T) {
	signOption := NewSignOption(
		"2015-04-27T08:23:49Z",
		ExpirationPeriodInSeconds,
		nil,
		[]string{"date", "content-type", "server", "content-length"},
	)
	result := signOption.signedHeadersToString()
	expected := "content-length;content-type;date;server"

	if result != expected {
		t.Error(util.FormatTest("signedHeadersToString", result, expected))
	}

	signOption = NewSignOption(
		"2015-04-27T08:23:49Z",
		ExpirationPeriodInSeconds,
		map[string]string{"content-type": "text/plain", "server": "nginx"},
		[]string{"date", "content-type", "server", "content-length"},
	)
	result = signOption.signedHeadersToString()
	expected = "content-length;content-type;date;server"

	if result != expected {
		t.Error(util.FormatTest("signedHeadersToString", result, expected))
	}

	signOption = NewSignOption(
		"2015-04-27T08:23:49Z",
		ExpirationPeriodInSeconds,
		map[string]string{"content-type": "text/plain", "server": "nginx", "content-length": "1024"},
		nil,
	)
	result = signOption.signedHeadersToString()
	expected = "content-length;content-type"

	if result != expected {
		t.Error(util.FormatTest("signedHeadersToString", result, expected))
	}
}

func TestGetSigningKey(t *testing.T) {
	const expected = "d9f35aaba8a5f3efa654851917114b6f22cd831116fd7d8431e08af22dcff24c"
	credentials := Credentials{
		AccessKeyID:     "0b0f67dfb88244b289b72b142befad0c",
		SecretAccessKey: "bad522c2126a4618a8125f4b6cf6356f",
	}
	signOption := NewSignOption(
		"2015-04-27T08:23:49Z",
		ExpirationPeriodInSeconds,
		getHeaders(),
		nil,
	)
	signingKey := getSigningKey(credentials, signOption)

	if signingKey != expected {
		t.Error(util.FormatTest("getSigningKey", signingKey, expected))
	}
}

func TestSign(t *testing.T) {
	expected := "a19e6386e990691aca1114a20357c83713f1cb4be3d74942bb4ed37469ecdacf"
	credentials := Credentials{
		AccessKeyID:     "0b0f67dfb88244b289b72b142befad0c",
		SecretAccessKey: "bad522c2126a4618a8125f4b6cf6356f",
	}
	req := getRequest()
	signOption := NewSignOption(
		"2015-04-27T08:23:49Z",
		ExpirationPeriodInSeconds,
		getHeaders(),
		nil,
	)
	signature := sign(credentials, *req, signOption)

	if signature != expected {
		t.Error(util.FormatTest("sign", signature, expected))
	}
}

func TestGenerateAuthorization(t *testing.T) {
	expected := "bce-auth-v1/0b0f67dfb88244b289b72b142befad0c/2015-04-27T08:23:49Z/1800/content-length;content-md5;" +
		"content-type;host;x-bce-date/a19e6386e990691aca1114a20357c83713f1cb4be3d74942bb4ed37469ecdacf"
	credentials := Credentials{
		AccessKeyID:     "0b0f67dfb88244b289b72b142befad0c",
		SecretAccessKey: "bad522c2126a4618a8125f4b6cf6356f",
	}
	req := getRequest()
	signOption := NewSignOption(
		"2015-04-27T08:23:49Z",
		ExpirationPeriodInSeconds,
		getHeaders(),
		nil,
	)
	authorization := GenerateAuthorization(credentials, *req, signOption)

	if authorization != expected {
		t.Error(util.FormatTest("GenerateAuthorization", authorization, expected))
	}

	req = getRequest()
	GenerateAuthorization(credentials, *req, nil)
}

func TestNewHttpClient(t *testing.T) {
	credentials := Credentials{
		AccessKeyID:     "0b0f67dfb88244b289b72b142befad0c",
		SecretAccessKey: "bad522c2126a4618a8125f4b6cf6356f",
	}
	config := &Config{
		Credentials:    &credentials,
		ProxyHost:      "guoyao.me",
		ProxyPort:      8000,
		MaxConnections: 10,
	}

	httpClient := newHttpClient(config)

	if httpClient == nil {
		t.Error(util.FormatTest("newHttpClient", "nil", "valid http client"))
	}
}

func TestSetDebug(t *testing.T) {
	config := getConfig()
	client := NewClient(config)
	expected := false

	if client.debug != expected {
		t.Error(util.FormatTest("SetDebug", strconv.FormatBool(client.debug), strconv.FormatBool(expected)))
	}

	expected = true
	client.SetDebug(true)

	if client.debug != expected {
		t.Error(util.FormatTest("SetDebug", strconv.FormatBool(client.debug), strconv.FormatBool(expected)))
	}
}

func TestGetURL(t *testing.T) {
	client := NewClient(&Config{})
	host := "guoyao.me"
	uriPath := "articals"
	params := map[string]string{"pageNo": "2", "pageSize": "10"}
	url := client.GetURL(host, uriPath, params)
	expected := "http://guoyao.me/articals?pageNo=2&pageSize=10"

	if url != expected {
		t.Error(util.FormatTest("GetURL", url, expected))
	}

	uriPath = "/articals"
	url = client.GetURL(host, uriPath, params)

	if url != expected {
		t.Error(util.FormatTest("GetURL", url, expected))
	}

	client = NewClient(&Config{APIVersion: "v1"})
	url = client.GetURL(host, uriPath, params)
	expected = "http://guoyao.me/v1/articals?pageNo=2&pageSize=10"

	if url != expected {
		t.Error(util.FormatTest("GetURL", url, expected))
	}
}

func TestGetSessionToken(t *testing.T) {
	method := "GetSessionToken"
	config := getConfig()

	req := SessionTokenRequest{
		DurationSeconds: 600,
		Id:              "ef5a4b19192f4931adcf0e12f82795e2",
		AccessControlList: []AccessControlListItem{
			AccessControlListItem{
				Service:    "bce:bos",
				Region:     config.GetRegion(),
				Effect:     "Allow",
				Resource:   []string{"baidubce-sdk-go/*"},
				Permission: []string{"READ"},
			},
		},
	}

	sessionTokenResponse, err := bceClient.GetSessionToken(req, nil)

	if err != nil {
		t.Error(util.FormatTest(method, err.Error(), "nil"))
	} else if sessionTokenResponse.SessionToken == "" {
		t.Error(util.FormatTest(method, "sessionToken is empty", "non empty sessionToken"))
	}
}

func TestSendRequest(t *testing.T) {
	config := getConfig()
	client := NewClient(config)
	client.SetDebug(true)
	url := "http://www.baidu.com"
	request, _ := NewRequest("GET", url, nil)
	resp, err := client.SendRequest(request, &SignOption{
		Credentials: &Credentials{
			AccessKeyID:     "0b0f67dfb88244b289b72b142befad0c",
			SecretAccessKey: "bad522c2126a4618a8125f4b6cf6356f",
		},
	})

	if err != nil {
		t.Error(util.FormatTest("SendRequest", err.Error(), "nil"))
	}

	url = "http://guoyao.me.no-such-host"
	request, _ = NewRequest("GET", url, nil)
	resp, err = client.SendRequest(request, nil)

	if _, ok := err.(*Error); ok {
		t.Error(util.FormatTest("SendRequest", "bceError", "error"))
	}

	url = "http://guoyao.me/no-exist-path"
	request, _ = NewRequest("GET", url, nil)
	resp, err = client.SendRequest(request, nil)

	if err == nil && resp.StatusCode != http.StatusNotFound {
		t.Error(util.FormatTest("SendRequest", strconv.Itoa(resp.StatusCode), strconv.Itoa(http.StatusNotFound)))
	}

	duration := client.RetryPolicy.GetDelayBeforeNextRetry(err, 1)

	if duration != -1 {
		t.Error(util.FormatTest("SendRequest", duration.String(), strconv.Itoa(-1)))
	}

	resp, err = client.SendRequest(getRequest(), nil)

	if err == nil {
		t.Error(util.FormatTest("SendRequest", "nil", "error"))
	}

	duration = client.RetryPolicy.GetDelayBeforeNextRetry(err, 1)

	if duration != -1 {
		t.Error(util.FormatTest("SendRequest", duration.String(), strconv.Itoa(-1)))
	}

	if bceError, ok := err.(*Error); ok {
		bceError.StatusCode = http.StatusInternalServerError
		retriesAttempted := 1
		duration = client.RetryPolicy.GetDelayBeforeNextRetry(err, retriesAttempted)
		expected := (1 << uint(retriesAttempted)) * 300 * time.Millisecond

		if duration != expected {
			t.Error(util.FormatTest("SendRequest", duration.String(), expected.String()))
		}

		retriesAttempted = 2
		duration = client.RetryPolicy.GetDelayBeforeNextRetry(err, retriesAttempted)
		expected = (1 << uint(retriesAttempted)) * 300 * time.Millisecond

		if duration != expected {
			t.Error(util.FormatTest("SendRequest", duration.String(), expected.String()))
		}

		retriesAttempted = 3
		duration = client.RetryPolicy.GetDelayBeforeNextRetry(err, retriesAttempted)
		expected = (1 << uint(retriesAttempted)) * 300 * time.Millisecond

		if duration != expected {
			t.Error(util.FormatTest("SendRequest", duration.String(), expected.String()))
		}

		retriesAttempted = client.RetryPolicy.GetMaxErrorRetry() + 1
		duration = client.RetryPolicy.GetDelayBeforeNextRetry(err, retriesAttempted)
		expected = -1

		if duration != expected {
			t.Error(util.FormatTest("SendRequest", duration.String(), expected.String()))
		}
	}

	client.SetDebug(false)
}

func getRequest() *Request {
	params := map[string]string{
		"partNumber": "9",
		"uploadId":   "VXBsb2FkIElpZS5tMnRzIHVwbG9hZA",
	}

	url := util.GetURL("http", "bj.bcebos.com", "/v1/test/myfolder/readme.txt", params)

	request, _ := NewRequest("PUT", url, nil)

	return request
}

func getHeaders() map[string]string {
	headers := map[string]string{
		"Host":           "bj.bcebos.com",
		"Date":           "Mon, 27 Apr 2015 16:23:49 +0800",
		"Content-Type":   "text/plain",
		"Content-Length": "8",
		"Content-Md5":    "0a52730597fb4ffa01fc117d9e71e3a9",
		"x-bce-date":     "2015-04-27T08:23:49Z",
	}

	return headers
}

func TestBCCSendRequest(t *testing.T) {
	config := getConfig()
	client := NewClient(config)
	client.SetDebug(true)
	url := "https://bcc.bj.baidubce.com/v2/instance"
	request, _ := NewRequest("GET", url, nil)
	resp, err := client.SendRequest(request, &SignOption{
		Credentials: &Credentials{
			AccessKeyID:     "0b565874be53440ba579d8a70cac8f98",
			SecretAccessKey: "d8ad18e630c64bf4b4661ce9b49164c8",
		},
	})
	bodyContent, err := resp.GetBodyContent()

	if err != nil {
		t.Error("resp error")
	}

	fmt.Printf("resp is : %s \n", bodyContent)
	if err != nil {
		t.Error(util.FormatTest("SendRequest", err.Error(), "nil"))
	}

	client.SetDebug(false)
}

func TestNewCredentialsFromFile(t *testing.T) {
	c, err := NewCredentialsFromFile("aksk.json")
	if err != nil {
		t.Error("newcredentialsfromfile error")
	} else {
		if c.AccessKeyID != "aaa" {
			t.Error("error ak")
		}
	}
}

func TestNewConfigFromFile(t *testing.T) {
	c, err := NewConfigFromFile("aksk.json")
	if err != nil {
		t.Error(err)
	} else {
		fmt.Printf("%+v", c)
	}
}
