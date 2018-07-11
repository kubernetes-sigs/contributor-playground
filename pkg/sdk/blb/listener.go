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

package blb

import (
	"bytes"
	"encoding/json"
	"fmt"

	"strconv"

	"k8s.io/cloud-provider-baiducloud/pkg/sdk/bce"
)

type TCPListener struct {
	ListenerPort               int    `json:"listenerPort"`
	BackendPort                int    `json:"backendPort"`
	Scheduler                  string `json:"scheduler"`
	HealthCheckTimeoutInSecond int    `json:"healthCheckTimeoutInSecond"`
	HealthCheckInterval        int    `json:"healthCheckInterval"`
	UnhealthyThreshold         int    `json:"unhealthyThreshold"`
	HealthyThreshold           int    `json:"healthyThreshold"`
}

type UDPListener struct {
	ListenerPort               int    `json:"listenerPort"`
	BackendPort                int    `json:"backendPort"`
	Scheduler                  string `json:"scheduler"`
	HealthCheckTimeoutInSecond int    `json:"healthCheckTimeoutInSecond"`
	HealthCheckInterval        int    `json:"healthCheckInterval"`
	UnhealthyThreshold         int    `json:"unhealthyThreshold"`
	HealthyThreshold           int    `json:"healthyThreshold"`
	HealthCheckString          string `json:"healthCheckString"`
}

type HTTPListener struct {
	ListenerPort               int    `json:"listenerPort"`
	BackendPort                int    `json:"backendPort"`
	Scheduler                  string `json:"scheduler"`
	KeepSession                bool   `json:"keepSession"`
	KeepSessionType            string `json:"keepSessionType"`
	KeepSessionDuration        int    `json:"keepSessionDuration"`
	KeepSessionCookieName      int    `json:"keepSessionCookieName"`
	XForwardFor                bool   `json:"xForwardFor"`
	HealthCheckType            string `json:"healthCheckType"`
	HealthCheckPort            int    `json:"healthCheckPort"`
	HealthCheckURI             string `json:"healthCheckURI"`
	HealthCheckTimeoutInSecond int    `json:"healthCheckTimeoutInSecond"`
	HealthCheckInterval        int    `json:"healthCheckInterval"`
	UnhealthyThreshold         int    `json:"unhealthyThreshold"`
	HealthyThreshold           int    `json:"healthyThreshold"`
	HealthCheckNormalStatus    string `json:"healthCheckNormalStatus"`
	ServerTimeout              int    `json:"serverTimeout"`
	RedirectPort               int    `json:"redirectPort"`
}

type HTTPSListener struct {
	ListenerPort               int      `json:"listenerPort"`
	BackendPort                int      `json:"backendPort"`
	Scheduler                  string   `json:"scheduler"`
	KeepSession                bool     `json:"keepSession"`
	KeepSessionType            string   `json:"keepSessionType"`
	KeepSessionDuration        int      `json:"keepSessionDuration"`
	KeepSessionCookieName      int      `json:"keepSessionCookieName"`
	XForwardFor                bool     `json:"xForwardFor"`
	HealthCheckType            string   `json:"healthCheckType"`
	HealthCheckPort            int      `json:"healthCheckPort"`
	HealthCheckURI             string   `json:"healthCheckURI"`
	HealthCheckTimeoutInSecond int      `json:"healthCheckTimeoutInSecond"`
	HealthCheckInterval        int      `json:"healthCheckInterval"`
	UnhealthyThreshold         int      `json:"unhealthyThreshold"`
	HealthyThreshold           int      `json:"healthyThreshold"`
	HealthCheckNormalStatus    string   `json:"healthCheckNormalStatus"`
	ServerTimeout              int      `json:"serverTimeout"`
	CertIds                    []string `json:"certIds"`
	Ie6Compatible              bool     `json:"ie6Compatible"`
}

type CreateTCPListenerArgs struct {
	LoadBalancerId             string `json:"-"`
	ListenerPort               int    `json:"listenerPort"`
	BackendPort                int    `json:"backendPort"`
	Scheduler                  string `json:"scheduler"`
	HealthCheckTimeoutInSecond int    `json:"healthCheckTimeoutInSecond,omitempty"`
	HealthCheckInterval        int    `json:"healthCheckInterval,omitempty"`
	UnhealthyThreshold         int    `json:"unhealthyThreshold,omitempty"`
	HealthyThreshold           int    `json:"healthyThreshold,omitempty"`
}

type CreateUDPListenerArgs struct {
	LoadBalancerId             string `json:"-"`
	ListenerPort               int    `json:"listenerPort"`
	BackendPort                int    `json:"backendPort"`
	Scheduler                  string `json:"scheduler"`
	HealthCheckTimeoutInSecond int    `json:"healthCheckTimeoutInSecond,omitempty"`
	HealthCheckInterval        int    `json:"healthCheckInterval,omitempty"`
	UnhealthyThreshold         int    `json:"unhealthyThreshold,omitempty"`
	HealthyThreshold           int    `json:"healthyThreshold,omitempty"`
	HealthCheckString          string `json:"healthCheckString"`
}

type CreateHTTPListenerArgs struct {
	LoadBalancerId             string `json:"-"`
	ListenerPort               int    `json:"listenerPort"`
	BackendPort                int    `json:"backendPort"`
	Scheduler                  string `json:"scheduler"`
	KeepSession                bool   `json:"keepSession,omitempty"`
	KeepSessionType            string `json:"keepSessionType,omitempty"`
	KeepSessionDuration        int    `json:"keepSessionDuration,omitempty"`
	KeepSessionCookieName      int    `json:"keepSessionCookieName,omitempty"`
	XForwardFor                bool   `json:"xForwardFor,omitempty"`
	HealthCheckType            string `json:"healthCheckType,omitempty"`
	HealthCheckPort            int    `json:"healthCheckPort,omitempty"`
	HealthCheckURI             string `json:"healthCheckURI,omitempty"`
	HealthCheckTimeoutInSecond int    `json:"healthCheckTimeoutInSecond,omitempty"`
	HealthCheckInterval        int    `json:"healthCheckInterval,omitempty"`
	UnhealthyThreshold         int    `json:"unhealthyThreshold,omitempty"`
	HealthyThreshold           int    `json:"healthyThreshold,omitempty"`
	HealthCheckNormalStatus    string `json:"healthCheckNormalStatus,omitempty"`
	ServerTimeout              int    `json:"serverTimeout,omitempty"`
	RedirectPort               int    `json:"redirectPort,omitempty"`
}

// CreateLoadBalancerHTTPListener create HTTP listener on loadbalancer
// You can read doc at https://cloud.baidu.com/doc/BLB/API.html#.C0.F3.F3.ED.5C.D8.4D.66.19.FF.DA.7A.0F.75.05.7C
func (c *Client) CreateTCPListener(args *CreateTCPListenerArgs) (err error) {
	params := map[string]string{
		"clientToken": c.GenerateClientToken(),
	}
	if args == nil {
		return fmt.Errorf("CreateTCPListener need args")
	}
	postContent, err := json.Marshal(args)
	if err != nil {
		return err
	}
	req, err := bce.NewRequest("POST", c.GetURL("v1/blb"+"/"+args.LoadBalancerId+"/TCPlistener", params), bytes.NewBuffer(postContent))
	if err != nil {
		return err
	}
	_, err = c.SendRequest(req, nil)
	if err != nil {
		return err
	}
	return nil
}

// CreateLoadBalancerUDPListener create UDP listener on loadbalancer
//
// You can read doc at https://cloud.baidu.com/doc/BLB/API.html#.D7.A3.9B.E1.45.BD.9E.FA.B0.2F.60.12.B3.39.E8.9D
func (c *Client) CreateUDPListener(args *CreateUDPListenerArgs) (err error) {
	params := map[string]string{
		"clientToken": c.GenerateClientToken(),
	}
	if args == nil {
		return fmt.Errorf("CreateUDPListener need args")
	}
	postContent, err := json.Marshal(args)
	if err != nil {
		return err
	}
	req, err := bce.NewRequest("POST", c.GetURL("v1/blb"+"/"+args.LoadBalancerId+"/UDPlistener", params), bytes.NewBuffer(postContent))
	if err != nil {
		return err
	}
	_, err = c.SendRequest(req, nil)
	if err != nil {
		return err
	}
	return nil
}

// CreateHTTPListener create HTTP listener on loadbalancer
//
// You can read doc at https://cloud.baidu.com/doc/BLB/API.html#.D7.A3.9B.E1.45.BD.9E.FA.B0.2F.60.12.B3.39.E8.9D
func (c *Client) CreateHTTPListener(args *CreateHTTPListenerArgs) (err error) {
	params := map[string]string{
		"clientToken": c.GenerateClientToken(),
	}
	if args == nil {
		return fmt.Errorf("CreateHTTPListener need args")
	}
	postContent, err := json.Marshal(args)
	if err != nil {
		return err
	}
	req, err := bce.NewRequest("POST", c.GetURL("v1/blb"+"/"+args.LoadBalancerId+"/HTTPlistener", params), bytes.NewBuffer(postContent))
	if err != nil {
		return err
	}
	_, err = c.SendRequest(req, nil)
	if err != nil {
		return err
	}
	return nil
}

type DescribeTCPListenerArgs struct {
	LoadBalancerId string
	ListenerPort   int
}
type DescribeTCPListenerResponse struct {
	Marker          string        `json:"marker"`
	IsTruncated     bool          `json:"isTruncated"`
	NextMarker      string        `json:"nextMarker"`
	MaxKeys         int           `json:"maxKeys"`
	TCPListenerList []TCPListener `json:"listenerList"`
}

// DescribeTCPListener Describe TCPListener
// TODO: args need to validate
func (c *Client) DescribeTCPListener(args *DescribeTCPListenerArgs) ([]TCPListener, error) {
	if args == nil {
		return nil, fmt.Errorf("DescribeTCPListener need args")
	}
	var params map[string]string
	if args.ListenerPort != 0 {
		params = map[string]string{
			"listenerPort": strconv.Itoa(args.ListenerPort),
		}
	}
	req, err := bce.NewRequest("GET", c.GetURL("v1/blb"+"/"+args.LoadBalancerId+"/TCPlistener", params), nil)

	if err != nil {
		return nil, err
	}
	resp, err := c.SendRequest(req, nil)

	if err != nil {
		return nil, err
	}
	bodyContent, err := resp.GetBodyContent()

	if err != nil {
		return nil, err
	}
	var blbsResp *DescribeTCPListenerResponse
	err = json.Unmarshal(bodyContent, &blbsResp)

	if err != nil {
		return nil, err
	}
	return blbsResp.TCPListenerList, nil
}

type DescribeUDPListenerArgs struct {
	LoadBalancerId string
	ListenerPort   int
}
type DescribeUDPListenerResponse struct {
	Marker          string        `json:"marker"`
	IsTruncated     bool          `json:"isTruncated"`
	NextMarker      string        `json:"nextMarker"`
	MaxKeys         int           `json:"maxKeys"`
	UDPListenerList []UDPListener `json:"listenerList"`
}

// DescribeUDPListeners Describe UDPListeners
// TODO: args need to validate
func (c *Client) DescribeUDPListener(args *DescribeUDPListenerArgs) ([]UDPListener, error) {
	if args == nil {
		return nil, fmt.Errorf("DescribeUDPListeners need args")
	}
	var params map[string]string
	if args.ListenerPort != 0 {
		params = map[string]string{
			"listenerPort": strconv.Itoa(args.ListenerPort),
		}
	}
	req, err := bce.NewRequest("GET", c.GetURL("v1/blb"+"/"+args.LoadBalancerId+"/UDPlistener", params), nil)

	if err != nil {
		return nil, err
	}
	resp, err := c.SendRequest(req, nil)

	if err != nil {
		return nil, err
	}
	bodyContent, err := resp.GetBodyContent()

	if err != nil {
		return nil, err
	}
	var blbsResp *DescribeUDPListenerResponse
	err = json.Unmarshal(bodyContent, &blbsResp)

	if err != nil {
		return nil, err
	}
	return blbsResp.UDPListenerList, nil
}

type UpdateTCPListenerArgs struct {
	LoadBalancerId             string `json:"-"`
	ListenerPort               int    `json:"-"`
	BackendPort                int    `json:"backendPort,omitempty"`
	Scheduler                  string `json:"scheduler,omitempty"`
	HealthCheckTimeoutInSecond int    `json:"healthCheckTimeoutInSecond,omitempty"`
	HealthCheckInterval        int    `json:"healthCheckInterval,omitempty"`
	UnhealthyThreshold         int    `json:"unhealthyThreshold,omitempty"`
	HealthyThreshold           int    `json:"healthyThreshold,omitempty"`
}

// UpdateTCPListener update a TCPListener
// TODO: args need to validate
func (c *Client) UpdateTCPListener(args *UpdateTCPListenerArgs) error {

	if args == nil || args.LoadBalancerId == "" || args.ListenerPort == 0 {
		return fmt.Errorf("UpdateTCPListener need args")
	}
	var params map[string]string
	if args.ListenerPort != 0 {
		params = map[string]string{
			"listenerPort": strconv.Itoa(args.ListenerPort),
		}
	}
	postContent, err := json.Marshal(args)
	if err != nil {
		return err
	}
	req, err := bce.NewRequest("PUT", c.GetURL("v1/blb"+"/"+args.LoadBalancerId+"/TCPlistener", params), bytes.NewBuffer(postContent))
	if err != nil {
		return err
	}
	_, err = c.SendRequest(req, nil)
	if err != nil {
		return err
	}
	return nil
}

type UpdateUDPListenerArgs struct {
	LoadBalancerId             string `json:"-"`
	ListenerPort               int    `json:"listenerPort"`
	BackendPort                int    `json:"backendPort"`
	Scheduler                  string `json:"scheduler"`
	HealthCheckTimeoutInSecond int    `json:"healthCheckTimeoutInSecond,omitempty"`
	HealthCheckInterval        int    `json:"healthCheckInterval,omitempty"`
	UnhealthyThreshold         int    `json:"unhealthyThreshold,omitempty"`
	HealthyThreshold           int    `json:"healthyThreshold,omitempty"`
	HealthCheckString          string `json:"healthCheckString"`
}

func (args *UpdateUDPListenerArgs) validate() error {
	if args.LoadBalancerId == "" {
		return fmt.Errorf("UpdateUDPListener need LoadBalancerId")
	}
	if args.ListenerPort == 0 {
		return fmt.Errorf("UpdateUDPListener need ListenerPort")
	}
	if args.BackendPort == 0 {
		return fmt.Errorf("UpdateUDPListener need BackendPort")
	}
	if args.Scheduler == "" {
		return fmt.Errorf("UpdateUDPListener need Scheduler")
	}
	if args.HealthCheckString == "" {
		return fmt.Errorf("UpdateUDPListener need HealthCheckString")
	}
	return nil
}

// UpdateUDPListener update a UDPListener
func (c *Client) UpdateUDPListener(args *UpdateUDPListenerArgs) error {
	err := args.validate()
	if err != nil {
		return err
	}
	params := map[string]string{
		"listenerPort": strconv.Itoa(args.ListenerPort),
	}
	postContent, err := json.Marshal(args)
	if err != nil {
		return err
	}
	req, err := bce.NewRequest("PUT", c.GetURL("v1/blb"+"/"+args.LoadBalancerId+"/UDPlistener", params), bytes.NewBuffer(postContent))
	if err != nil {
		return err
	}
	_, err = c.SendRequest(req, nil)
	if err != nil {
		return err
	}
	return nil
}

type DeleteListenersArgs struct {
	LoadBalancerId string `json:"-"`
	// action         string `json:"-"`
	PortList []int `json:"portList"`
}

func (args *DeleteListenersArgs) validate() error {
	if args.LoadBalancerId == "" {
		return fmt.Errorf("DeleteListenersArgs need LoadBalancerId")
	}
	if args.PortList == nil {
		return fmt.Errorf("DeleteListenersArgs need PortList")
	}
	return nil
}

// UpdateUDPListener update a UDPListener
func (c *Client) DeleteListeners(args *DeleteListenersArgs) error {
	err := args.validate()
	if err != nil {
		return err
	}
	params := map[string]string{
		"batchdelete": "",
		"clientToken": c.GenerateClientToken(),
	}
	postContent, err := json.Marshal(args)
	if err != nil {
		return err
	}
	// url := "http://" + Endpoint[c.GetRegion()] + "/v1/blb" + "/" + args.LoadBalancerId + "/listener?" + "batchdelete=&" + "clientToken=" + c.GenerateClientToken()
	req, err := bce.NewRequest("PUT", c.GetURL("v1/blb"+"/"+args.LoadBalancerId+"/listener", params), bytes.NewBuffer(postContent))
	if err != nil {
		return err
	}
	_, err = c.SendRequest(req, nil)
	if err != nil {
		return err
	}
	return nil
}
