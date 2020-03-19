package fake

import (
	"context"
	"fmt"

	"icode.baidu.com/baidu/jpaas-caas/bce-sdk-go/bce"
	"icode.baidu.com/baidu/jpaas-caas/bce-sdk-go/blb"
	"icode.baidu.com/baidu/jpaas-caas/bce-sdk-go/util"
)

// FakeClient implement of vpc.Interface
type BlbFakeClient struct {
	LoadBalancerMap  map[string]blb.LoadBalancer
	TCPListenerMap   map[string][]blb.TCPListener
	UDPListenerMap   map[string][]blb.UDPListener
	HTTPListenerMap  map[string][]blb.HTTPListener
	BackendServerMap map[string][]blb.BackendServer
}

// NewFakeClient for VPC fake client
func NewBlbFakeClient() *BlbFakeClient {
	return &BlbFakeClient{
		LoadBalancerMap:  map[string]blb.LoadBalancer{},
		TCPListenerMap:   map[string][]blb.TCPListener{},
		UDPListenerMap:   map[string][]blb.UDPListener{},
		HTTPListenerMap:  map[string][]blb.HTTPListener{},
		BackendServerMap: map[string][]blb.BackendServer{},
	}
}

// LoadBalance fake func
func (f *BlbFakeClient) DescribeLoadBalancers(ctx context.Context, args *blb.DescribeLoadBalancersArgs, option *bce.SignOption) ([]blb.LoadBalancer, error) {
	if args == nil {
		return nil, fmt.Errorf("args is nil")
	}
	loadbalancers := []blb.LoadBalancer{}
	for loadBalancerID, LoadBalancer := range f.LoadBalancerMap {
		if loadBalancerID != "" && args.LoadBalancerId != "" && loadBalancerID == args.LoadBalancerId ||
			LoadBalancer.Name != "" && args.LoadBalancerName != "" && LoadBalancer.Name == args.LoadBalancerName ||
			LoadBalancer.Address != "" && args.Address != "" && LoadBalancer.Address == args.Address {
			loadbalancers = append(loadbalancers, LoadBalancer)
		}
	}
	if len(loadbalancers) == 0 {
		return nil, fmt.Errorf("DescribeLoadBalancers error: can not get LoadBalancer from args: %v", args)
	}
	return loadbalancers, nil
}
func (f *BlbFakeClient) CreateLoadBalancer(ctx context.Context, args *blb.CreateLoadBalancerArgs, option *bce.SignOption) (*blb.CreateLoadBalancerResponse, error) {
	if args == nil {
		return nil, fmt.Errorf("args is nil")
	}
	resp := &blb.CreateLoadBalancerResponse{
		Desc:    args.Desc,
		Name:    args.Name,
		Address: "0.0.0.0",
	}
	loadbalancer := blb.LoadBalancer{
		Name:     args.Name,
		Desc:     args.Desc,
		Status:   "available",
		PublicIp: "0.0.0.0",
	}
	for {
		loadbalancerID := util.GenerateBCEShortID("lb")
		if _, ok := f.LoadBalancerMap[loadbalancerID]; !ok {
			loadbalancer.BlbId = loadbalancerID
			resp.LoadBalancerId = loadbalancerID
			f.LoadBalancerMap[loadbalancerID] = loadbalancer
			break
		}
	}
	if resp.LoadBalancerId == "" {
		return nil, fmt.Errorf("CreateLoadBalancer error: can not create LoadBalancer from args: %v", args)
	}
	return resp, nil
}
func (f *BlbFakeClient) UpdateLoadBalancer(ctx context.Context, args *blb.UpdateLoadBalancerArgs, option *bce.SignOption) error {
	if args == nil {
		return fmt.Errorf("args is nil")
	}
	if _, ok := f.LoadBalancerMap[args.LoadBalancerId]; ok {
		loadblanace := f.LoadBalancerMap[args.LoadBalancerId]
		loadblanace.Desc = args.Desc
		loadblanace.Name = args.Name
		f.LoadBalancerMap[args.LoadBalancerId] = loadblanace
		return nil
	}
	return fmt.Errorf("blbID does not exist")
}
func (f *BlbFakeClient) DeleteLoadBalancer(ctx context.Context, args *blb.DeleteLoadBalancerArgs, option *bce.SignOption) error {
	if args == nil {
		return fmt.Errorf("args is nil")
	}
	if _, ok := f.LoadBalancerMap[args.LoadBalancerId]; ok {
		delete(f.LoadBalancerMap, args.LoadBalancerId)
		return nil
	}
	return fmt.Errorf("LoadBalancerId does not exist")
}

// Listenr fake func
func (f *BlbFakeClient) CreateTCPListener(ctx context.Context, args *blb.CreateTCPListenerArgs, option *bce.SignOption) (err error) {
	if args == nil {
		return fmt.Errorf("args is nil")
	}
	tcp := blb.TCPListener{
		ListenerPort:               args.ListenerPort,
		BackendPort:                args.BackendPort,
		Scheduler:                  args.Scheduler,
		HealthCheckTimeoutInSecond: args.HealthCheckTimeoutInSecond,
		HealthCheckInterval:        args.HealthCheckInterval,
		UnhealthyThreshold:         args.UnhealthyThreshold,
		HealthyThreshold:           args.HealthyThreshold,
	}
	// check LoadBalancerId if nil
	argsLb := &blb.DescribeLoadBalancersArgs{
		LoadBalancerId: args.LoadBalancerId,
	}
	lbs, err := f.DescribeLoadBalancers(ctx, argsLb, nil)
	if err != nil || len(lbs) == 0 {
		return fmt.Errorf("can not get lb according to args’ BlbID err: %v", err)
	}
	f.TCPListenerMap[args.LoadBalancerId] = append(f.TCPListenerMap[args.LoadBalancerId], tcp)
	return nil
}
func (f *BlbFakeClient) CreateUDPListener(ctx context.Context, args *blb.CreateUDPListenerArgs, option *bce.SignOption) (err error) {
	if args == nil {
		return fmt.Errorf("args is nil")
	}
	udp := blb.UDPListener{
		ListenerPort:               args.ListenerPort,
		BackendPort:                args.BackendPort,
		Scheduler:                  args.Scheduler,
		HealthCheckTimeoutInSecond: args.HealthCheckTimeoutInSecond,
		HealthCheckInterval:        args.HealthCheckInterval,
		UnhealthyThreshold:         args.UnhealthyThreshold,
		HealthyThreshold:           args.HealthyThreshold,
		HealthCheckString:          args.HealthCheckString,
	}
	// check LoadBalancerId if nil
	// check LoadBalancerId if nil
	argsLb := &blb.DescribeLoadBalancersArgs{
		LoadBalancerId: args.LoadBalancerId,
	}
	lbs, err := f.DescribeLoadBalancers(ctx, argsLb, nil)
	if err != nil || len(lbs) == 0 {
		return fmt.Errorf("can not get lb according to args’ BlbID err: %v", err)
	}
	f.UDPListenerMap[args.LoadBalancerId] = append(f.UDPListenerMap[args.LoadBalancerId], udp)
	return nil
}
func (f *BlbFakeClient) CreateHTTPListener(ctx context.Context, args *blb.CreateHTTPListenerArgs, option *bce.SignOption) (err error) {
	if args == nil {
		return fmt.Errorf("args is nil")
	}
	http := blb.HTTPListener{
		ListenerPort:               args.ListenerPort,
		BackendPort:                args.BackendPort,
		Scheduler:                  args.Scheduler,
		KeepSession:                args.KeepSession,
		KeepSessionType:            args.KeepSessionType,
		KeepSessionDuration:        args.KeepSessionDuration,
		KeepSessionCookieName:      args.KeepSessionCookieName,
		XForwardFor:                args.XForwardFor,
		HealthCheckType:            args.HealthCheckType,
		HealthCheckURI:             args.HealthCheckURI,
		HealthCheckTimeoutInSecond: args.HealthCheckTimeoutInSecond,
		UnhealthyThreshold:         args.UnhealthyThreshold,
		HealthyThreshold:           args.HealthyThreshold,
		HealthCheckNormalStatus:    args.HealthCheckNormalStatus,
		ServerTimeout:              args.ServerTimeout,
		RedirectPort:               args.RedirectPort,
	}
	f.HTTPListenerMap[args.LoadBalancerId] = append(f.HTTPListenerMap[args.LoadBalancerId], http)
	return nil
}
func (f *BlbFakeClient) DescribeTCPListener(ctx context.Context, args *blb.DescribeTCPListenerArgs, option *bce.SignOption) ([]blb.TCPListener, error) {
	if args == nil {
		return nil, fmt.Errorf("args is nil")
	}
	tcpListeners := []blb.TCPListener{}
	if _, ok := f.TCPListenerMap[args.LoadBalancerId]; ok {
		for _, t := range f.TCPListenerMap[args.LoadBalancerId] {
			tcpListeners = append(tcpListeners, t)
		}
		return tcpListeners, nil
	}
	return nil, fmt.Errorf("DescribeTCPListener failed, can not get tcpListeners from args %v", args)
}
func (f *BlbFakeClient) DescribeUDPListener(ctx context.Context, args *blb.DescribeUDPListenerArgs, option *bce.SignOption) ([]blb.UDPListener, error) {
	if args == nil {
		return nil, fmt.Errorf("DescribeUDPListeners need args")
	}
	if args.LoadBalancerId == "" {
		return nil, fmt.Errorf("DescribeUDPListeners args need loadbalancerId")
	}
	udpListenerList, found := f.UDPListenerMap[args.LoadBalancerId]
	if !found {
		return nil, fmt.Errorf("Sepcified BLB %s not found", args.LoadBalancerId)
	}
	result := make([]blb.UDPListener, 0)
	for _, u := range udpListenerList {
		if args.ListenerPort != 0 && u.ListenerPort != args.ListenerPort {
			continue
		}
		result = append(result, u)
	}
	return result, nil
}
func (f *BlbFakeClient) UpdateTCPListener(ctx context.Context, args *blb.UpdateTCPListenerArgs, option *bce.SignOption) error {
	if args == nil || args.LoadBalancerId == "" || args.ListenerPort == 0 {
		return fmt.Errorf("UpdateTCPListener need args")
	}
	rawTcpList, found := f.TCPListenerMap[args.LoadBalancerId]
	if !found {
		return fmt.Errorf("Specified BLB %s not found", args.LoadBalancerId)
	}
	tcpList := make([]blb.TCPListener, 0)
	for _, t := range rawTcpList {
		if t.ListenerPort != args.ListenerPort {
			tcpList = append(tcpList, t)
		}
	}
	newTcpListner := blb.TCPListener{
		ListenerPort:               args.ListenerPort,
		HealthCheckInterval:        args.HealthCheckInterval,
		HealthCheckTimeoutInSecond: args.HealthCheckTimeoutInSecond,
		HealthyThreshold:           args.HealthyThreshold,
		UnhealthyThreshold:         args.UnhealthyThreshold,
		BackendPort:                args.BackendPort,
		Scheduler:                  args.Scheduler,
	}
	tcpList = append(tcpList, newTcpListner)
	return nil
}
func (f *BlbFakeClient) UpdateUDPListener(ctx context.Context, args *blb.UpdateUDPListenerArgs, option *bce.SignOption) error {
	err := validateUpdateUDPListenerArgs(args)
	if err != nil {
		return err
	}
	rawUdpList, found := f.UDPListenerMap[args.LoadBalancerId]
	if !found {
		return fmt.Errorf("Specified BLB %s not found", args.LoadBalancerId)
	}
	udpList := make([]blb.UDPListener, 0)
	for _, u := range rawUdpList {
		if u.ListenerPort != args.ListenerPort {
			udpList = append(udpList, u)
		}
	}
	newUdpListener := blb.UDPListener{
		ListenerPort:               args.ListenerPort,
		HealthCheckInterval:        args.HealthCheckInterval,
		HealthCheckString:          args.HealthCheckString,
		HealthCheckTimeoutInSecond: args.HealthCheckTimeoutInSecond,
		HealthyThreshold:           args.HealthyThreshold,
		UnhealthyThreshold:         args.UnhealthyThreshold,
		BackendPort:                args.BackendPort,
		Scheduler:                  args.Scheduler,
	}
	udpList = append(udpList, newUdpListener)
	return nil
}
func (f *BlbFakeClient) DeleteListeners(ctx context.Context, args *blb.DeleteListenersArgs, option *bce.SignOption) error {
	err := validateDeleteListenersArgs(args)
	if err != nil {
		return err
	}
	// listener port to remove
	listenerToRemove := make(map[int]int, len(args.PortList))
	for _, p := range args.PortList {
		listenerToRemove[p] = p
	}
	// tcp
	rawTcpList, found := f.TCPListenerMap[args.LoadBalancerId]
	if found {
		tcpList := make([]blb.TCPListener, 0)
		for _, t := range rawTcpList {
			if _, in := listenerToRemove[t.ListenerPort]; !in {
				tcpList = append(tcpList, t)
			}
		}
		f.TCPListenerMap[args.LoadBalancerId] = tcpList
	}
	// udp
	rawUdpList, found := f.UDPListenerMap[args.LoadBalancerId]
	if found {
		udpList := make([]blb.UDPListener, 0)
		for _, u := range rawUdpList {
			if _, in := listenerToRemove[u.ListenerPort]; !in {
				udpList = append(udpList, u)
			}
		}
		f.UDPListenerMap[args.LoadBalancerId] = udpList
	}
	// http
	rawHttpList, found := f.HTTPListenerMap[args.LoadBalancerId]
	if found {
		httpList := make([]blb.HTTPListener, 0)
		for _, h := range rawHttpList {
			if _, in := listenerToRemove[h.ListenerPort]; !in {
				httpList = append(httpList, h)
			}
		}
		f.HTTPListenerMap[args.LoadBalancerId] = httpList
	}
	// TODO: https
	return nil
}

// backendserver fake func
func (f *BlbFakeClient) AddBackendServers(ctx context.Context, args *blb.AddBackendServersArgs, option *bce.SignOption) error {
	if err := validateAddBackendServersArgs(args); err != nil {
		return err
	}
	_, found := f.LoadBalancerMap[args.LoadBalancerId]
	if !found {
		return fmt.Errorf("Specified BLB %s not found", args.LoadBalancerId)
	}
	backendList := make([]blb.BackendServer, 0)
	for _, rs := range args.BackendServerList {
		backendList = append(backendList, blb.BackendServer{
			InstanceId: rs.InstanceId,
			Weight:     rs.Weight,
		})
	}
	f.BackendServerMap[args.LoadBalancerId] = append(f.BackendServerMap[args.LoadBalancerId], backendList...)
	return nil
}
func (f *BlbFakeClient) DescribeBackendServers(ctx context.Context, args *blb.DescribeBackendServersArgs, option *bce.SignOption) ([]blb.BackendServer, error) {
	err := validateDescribeBackendServersArgs(args)
	if err != nil {
		return nil, err
	}
	result := make([]blb.BackendServer, 0)
	rsList, _ := f.BackendServerMap[args.LoadBalancerId]
	for _, rs := range rsList {
		result = append(result, blb.BackendServer{
			InstanceId: rs.InstanceId,
			Weight:     rs.Weight,
		})
	}
	return result, nil
}
func (f *BlbFakeClient) UpdateBackendServers(ctx context.Context, args *blb.UpdateBackendServersArgs, option *bce.SignOption) error {
	err := validateUpdateBackendServersArgs(args)
	if err != nil {
		return err
	}
	rawBackendList, found := f.BackendServerMap[args.LoadBalancerId]
	if !found {
		return fmt.Errorf("Specified BLB %s not found", args.LoadBalancerId)
	}
	backendList := make([]blb.BackendServer, 0)
	backendsToUpdate := make(map[string]string)
	for _, rs := range args.BackendServerList {
		backendsToUpdate[rs.InstanceId] = rs.InstanceId
		backendList = append(backendList, blb.BackendServer{
			InstanceId: rs.InstanceId,
			Weight:     rs.Weight,
		})
	}
	for _, b := range rawBackendList {
		if _, found := backendsToUpdate[b.InstanceId]; !found {
			backendList = append(backendList, b)
		}
	}
	f.BackendServerMap[args.LoadBalancerId] = backendList
	return nil
}
func (f *BlbFakeClient) RemoveBackendServers(ctx context.Context, args *blb.RemoveBackendServersArgs, option *bce.SignOption) error {
	err := validateRemoveBackendServersArgs(args)
	if err != nil {
		return err
	}
	rsList, found := f.BackendServerMap[args.LoadBalancerId]
	if !found {
		return fmt.Errorf("BLB %s not found", args.LoadBalancerId)
	}
	rsToRemove := make(map[string]string, len(args.BackendServerList))
	for _, instanceID := range args.BackendServerList {
		rsToRemove[instanceID] = instanceID
	}
	leftRs := make([]blb.BackendServer, 0)
	for _, rs := range rsList {
		if _, found := rsToRemove[rs.InstanceId]; !found {
			leftRs = append(leftRs, rs)
		}
	}
	f.BackendServerMap[args.LoadBalancerId] = leftRs
	return nil
}
func validateUpdateUDPListenerArgs(args *blb.UpdateUDPListenerArgs) error {
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
func validateDeleteListenersArgs(args *blb.DeleteListenersArgs) error {
	if args.LoadBalancerId == "" {
		return fmt.Errorf("DeleteListenersArgs need LoadBalancerId")
	}
	if args.PortList == nil {
		return fmt.Errorf("DeleteListenersArgs need PortList")
	}
	return nil
}
func validateAddBackendServersArgs(args *blb.AddBackendServersArgs) error {
	if args == nil {
		return fmt.Errorf("AddBackendServersArgs need args")
	}
	if args.LoadBalancerId == "" {
		return fmt.Errorf("AddBackendServersArgs need LoadBalancerId")
	}
	if args.BackendServerList == nil {
		return fmt.Errorf("UpdateUDPListener need BackendServerList")
	}
	return nil
}
func validateDescribeBackendServersArgs(args *blb.DescribeBackendServersArgs) error {
	if args == nil {
		return fmt.Errorf("DescribeBackendServersArgs need args")
	}
	if args.LoadBalancerId == "" {
		return fmt.Errorf("DescribeBackendServersArgs need LoadBalancerId")
	}
	return nil
}
func validateUpdateBackendServersArgs(args *blb.UpdateBackendServersArgs) error {
	if args == nil {
		return fmt.Errorf("UpdateBackendServersArgs need args")
	}
	if args.LoadBalancerId == "" {
		return fmt.Errorf("UpdateBackendServersArgs need LoadBalancerId")
	}
	if len(args.BackendServerList) == 0 {
		return fmt.Errorf("UpdateBackendServersArgs need BackendServerList")
	}
	return nil
}
func validateRemoveBackendServersArgs(args *blb.RemoveBackendServersArgs) error {
	if args == nil {
		return fmt.Errorf("UpdateBackendServersArgs need args")
	}
	if args.LoadBalancerId == "" {
		return fmt.Errorf("UpdateBackendServersArgs need LoadBalancerId")
	}
	if len(args.BackendServerList) == 0 {
		return fmt.Errorf("UpdateBackendServersArgs need BackendServerList")
	}
	return nil
}
