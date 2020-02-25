package cloud_provider

import (
	"context"
	"testing"

	"icode.baidu.com/baidu/jpaas-caas/bce-sdk-go/bce"
	"icode.baidu.com/baidu/jpaas-caas/bce-sdk-go/blb"
	"icode.baidu.com/baidu/jpaas-caas/bce-sdk-go/eip"
)

func TestGetEipArgsFromAnnotation(t *testing.T) {
	cloud, _, _, err := beforeTestBlb()
	if err != nil {
		t.Errorf("beforeTestBlb err , %v", err)
	}
	sa := &ServiceAnnotation{
		ElasticIPName:              "",
		ElasticIPBillingMethod:     "",
		ElasticIPBandwidthInMbps:   1,
		ElasticIPReservationLength: 1,
	}
	_, err = cloud.getEipArgsFromAnnotation(sa)
	if err != nil {
		t.Errorf("getEipArgsFromAnnotation err, err: %s", err)
	}
	sa = &ServiceAnnotation{
		ElasticIPName:              "da",
		ElasticIPBillingMethod:     "test",
		ElasticIPBandwidthInMbps:   1000000,
		ElasticIPReservationLength: 10000,
	}
	_, err = cloud.getEipArgsFromAnnotation(sa)
	if err == nil {
		t.Errorf("getEipArgsFromAnnotation err, should be not support target ElasticIPBillingMethod")
	}
}
func TestCreateEIP(t *testing.T) {
	cloud, _, _, err := beforeTestBlb()
	if err != nil {
		t.Errorf("beforeTestBlb err , %v", err)
	}
	ctx := context.Background()
	sa := &ServiceAnnotation{
		ElasticIPName:              "da",
		ElasticIPBillingMethod:     "",
		ElasticIPBandwidthInMbps:   100,
		ElasticIPReservationLength: 100,
	}
	args, err := cloud.getEipArgsFromAnnotation(sa)
	if err != nil {
		t.Errorf("getEipArgsFromAnnotation err, err: %s", err)
	}
	ip, err := cloud.createEIP(ctx, args)
	if err != nil {
		t.Errorf("createEIP err, err: %s", err)
	}
	if len(ip) == 0 {
		t.Errorf("createEIP err, ip is nil")
	}
}
func TestDeleteEIP(t *testing.T) {
	cloud, _, _, err := beforeTestBlb()
	if err != nil {
		t.Errorf("beforeTestBlb err , %v", err)
	}
	ctx := context.Background()
	// invalid case
	eips := []string{
		"",
		"test",
	}
	for _, eip := range eips {
		err = cloud.deleteEIP(ctx, eip)
		if err != nil {
			t.Errorf("deleteEIP err, err: %v", err)
		}
	}
	// right case
	sa := &ServiceAnnotation{
		ElasticIPName:              "da",
		ElasticIPBillingMethod:     "",
		ElasticIPBandwidthInMbps:   100,
		ElasticIPReservationLength: 100,
	}
	args, err := cloud.getEipArgsFromAnnotation(sa)
	if err != nil {
		t.Errorf("getEipArgsFromAnnotation err, err: %s", err)
	}
	ip, err := cloud.clientSet.EIPClient.CreateEIP(ctx, args, nil)
	if err != nil {
		t.Errorf("createEIP err, err: %s", err)
	}
	if len(ip) == 0 {
		t.Errorf("createEIP err, ip is nil")
	}
	err = cloud.deleteEIP(ctx, ip)
	if err != nil {
		t.Errorf("deleteEIP err, err: %s", err)
	}
	argsGet := eip.GetEIPsArgs{
		EIP: ip,
	}
	ips, err := cloud.clientSet.EIPClient.GetEIPs(ctx, &argsGet, nil)
	if err != nil {
		t.Errorf("GetEIPs err, err: %s", err)
	}
	if len(ips) != 0 {
		t.Errorf("deleteEIP err, eips left: %v", ips)
	}
}
func TestResizeEip(t *testing.T) {
	cloud, _, _, err := beforeTestBlb()
	if err != nil {
		t.Errorf("beforeTestBlb err , %v", err)
	}
	ctx := context.Background()
	sa := &ServiceAnnotation{
		ElasticIPBandwidthInMbps: 11,
	}
	ip := ""
	err = cloud.resizeEip(ctx, sa, ip)
	if err == nil {
		t.Errorf("resizeEip err, should be an error")
	}
	// right cases
	sa = &ServiceAnnotation{
		ElasticIPName:              "da",
		ElasticIPBillingMethod:     "",
		ElasticIPBandwidthInMbps:   100,
		ElasticIPReservationLength: 100,
	}
	args, err := cloud.getEipArgsFromAnnotation(sa)
	if err != nil {
		t.Errorf("getEipArgsFromAnnotation err, err: %s", err)
	}
	ip, err = cloud.clientSet.EIPClient.CreateEIP(ctx, args, nil)
	if err != nil {
		t.Errorf("createEIP err, err: %s", err)
	}
	if len(ip) == 0 {
		t.Errorf("createEIP err, ip is nil")
	}
	sa.ElasticIPBandwidthInMbps = 1000
	err = cloud.resizeEip(ctx, sa, ip)
	if err != nil {
		t.Errorf("resizeEip err, err: %s", err)
	}
}

func TestDeleteEipFinally(t *testing.T) {
	cloud, _, lbResp, err := beforeTestBlb()
	if err != nil {
		t.Errorf("beforeTestBlb err , %v", err)
	}
	ctx := context.Background()
	// nil cases
	lb := &blb.LoadBalancer{
		BlbId: lbResp.LoadBalancerId,
		Name:  "test",
		Desc:  "cce_auto_create_eip/test",
	}
	svc := buildService()
	err = cloud.ensureEipDeleted(ctx, svc, lb)
	if err != nil {
		t.Errorf("ensureEipDeleted err, err: %s", err)
	}
	// right cases
	// create eip
	sa := &ServiceAnnotation{
		ElasticIPName:              "da",
		ElasticIPBillingMethod:     "",
		ElasticIPBandwidthInMbps:   100,
		ElasticIPReservationLength: 100,
	}
	args, err := cloud.getEipArgsFromAnnotation(sa)
	if err != nil {
		t.Errorf("getEipArgsFromAnnotation err, err: %s", err)
	}
	ip, err := cloud.createEIP(ctx, args)
	if err != nil {
		t.Errorf("createEIP err, err: %s", err)
	}
	if len(ip) == 0 {
		t.Errorf("createEIP err, ip is nil")
	}
	lb.PublicIp = ip
	err = cloud.ensureEipDeleted(ctx, svc, lb)
	if err != nil {
		t.Errorf("ensureEipDeleted err, err: %s", err)
	}
	argsGet := eip.GetEIPsArgs{
		EIP: ip,
	}
	eips, err := cloud.clientSet.EIPClient.GetEIPs(ctx, &argsGet, nil)
	if err != nil {
		t.Errorf("GetEIPs err, err: %s", err)
	}
	if len(eips) != 0 {
		t.Errorf("ensureEipDeleted err, eip: %v is still existing", eips[0])
	}
}
func TestRefreshBlb(t *testing.T) {
	cloud, _, lbResp, err := beforeTestBlb()
	if err != nil {
		t.Errorf("beforeTestBlb err , %v", err)
	}
	ctx := context.Background()
	lb := &blb.LoadBalancer{
		// BlbId: lbResp.LoadBalancerId,
	}
	_, err = cloud.refreshBlb(ctx, lb)
	if err == nil {
		t.Errorf("refreshBlb err, should be an error")
	}
	lb = &blb.LoadBalancer{
		BlbId: lbResp.LoadBalancerId,
	}
	newlb, err := cloud.refreshBlb(ctx, lb)
	if err != nil {
		t.Errorf("refreshBlb err, err: %s", err)
	}
	if newlb.Status != "available" {
		t.Errorf("refreshBlb err, get newLb: %v", newlb)
	}
}

// case is nil
// right case
// todo case1 eip is not available
// todo case2 eip is bound
// todo case3 blb is not available
// the timeout is too long..
func TestBindEip(t *testing.T) {
	cloud, _, lbResp, err := beforeTestBlb()
	if err != nil {
		t.Errorf("beforeTestBlb err , %v", err)
	}
	ctx := context.Background()
	lb := &blb.LoadBalancer{
		BlbId: lbResp.LoadBalancerId,
	}
	loadBalancerIP := ""
	svc := buildService()
	_, err = cloud.bindEip(ctx, lb, loadBalancerIP, svc)
	if err == nil {
		t.Errorf("bindEip err, should be an error")
	}
	// create eip what is available
	args := &eip.CreateEIPArgs{
		Name:            "test",
		BandwidthInMbps: 100,
		Billing: &eip.Billing{
			PaymentTiming: "",
			BillingMethod: "",
		},
	}
	ip, err := cloud.clientSet.EIPClient.CreateEIP(ctx, args, &bce.SignOption{
		CustomSignFunc: CCEServiceSign,
	})
	if err != nil {
		t.Errorf("CreateEIP err, err； %s", err)
	}
	_, err = cloud.bindEip(ctx, lb, ip, svc)
	if err != nil {
		t.Errorf("bindEip err, err； %s", err)
	}
}

// case is nil
// right case
func TestEnsureEIPWithNoSpecificIP(t *testing.T) {
	cloud, _, blbResp, err := beforeTestBlb()
	if err != nil {
		t.Errorf("beforeTestBlb err , %v", err)
	}
	ctx := context.Background()
	// cases is nil
	svc := buildService()
	// clusterName svc nodes are not important
	serviceAnnotation := &ServiceAnnotation{}
	lb := &blb.LoadBalancer{}
	_, err = cloud.ensureEIPWithNoSpecificIP(ctx, svc, lb)
	if err == nil {
		t.Errorf("ensureEIPWithNoSpecificIP err, should be an error")
	}
	// right case1
	serviceAnnotation = &ServiceAnnotation{
		ElasticIPName:              "da",
		ElasticIPBillingMethod:     "",
		ElasticIPBandwidthInMbps:   100,
		ElasticIPReservationLength: 100,
	}
	lb = &blb.LoadBalancer{
		BlbId: blbResp.LoadBalancerId,
	}
	ip, err := cloud.ensureEIPWithNoSpecificIP(ctx, svc, lb)
	if err != nil {
		t.Errorf("ensureEIPWithNoSpecificIP err, err: %s", err)
	}
	if len(ip) == 0 {
		t.Errorf("ensureEIPWithNoSpecificIP err, ip is nil")
	}
	// right case2
	// create eip and bind blb
	lb2 := &blb.LoadBalancer{
		BlbId:  blbResp.LoadBalancerId,
		Status: "available",
	}
	args, err := cloud.getEipArgsFromAnnotation(serviceAnnotation)
	if err != nil {
		t.Errorf("getEipArgsFromAnnotation err, err: %s", err)
	}
	eip, err := cloud.createEIP(ctx, args)
	if err != nil {
		t.Errorf("createEIP err, err: %s", err)
	}
	if len(eip) == 0 {
		t.Errorf("createEIP err, ip is nil")
	}
	ip2, err := cloud.ensureEIPWithNoSpecificIP(ctx, svc, lb2)
	if err != nil {
		t.Errorf("ensureEIPWithNoSpecificIP err, err: %s", err)
	}
	if len(ip2) == 0 {
		t.Errorf("ensureEIPWithNoSpecificIP err, ip is nil")
	}
}

// cases is nil
// right case1 blb already bound the eip
// right case2 blb does not bind the eip
// right case3 blb already bound eip what is not expected
func TestEnsureEIPWithSpecificIP(t *testing.T) {
	cloud, _, lbResp, err := beforeTestBlb()
	if err != nil {
		t.Errorf("beforeTestBlb err , %v", err)
	}
	ctx := context.Background()
	// cases is nil

	// clusterName nodes are not important
	svc := buildService()
	svc.Spec.LoadBalancerIP = "test"

	lb := &blb.LoadBalancer{
		BlbId: lbResp.LoadBalancerId,
	}
	_, err = cloud.ensureEIPWithSpecificIP(ctx, svc, lb)
	if err == nil {
		t.Errorf("ensureEIPWithSpecificIP err, should be an error")
	}
	// right case1 blb already bound the eip
	// create eip bind blb
	sa := &ServiceAnnotation{
		ElasticIPName:              "da",
		ElasticIPBillingMethod:     "",
		ElasticIPBandwidthInMbps:   100,
		ElasticIPReservationLength: 100,
	}
	args, err := cloud.getEipArgsFromAnnotation(sa)
	if err != nil {
		t.Errorf("getEipArgsFromAnnotation err, err: %s", err)
	}
	ip1, err := cloud.createEIP(ctx, args)
	if err != nil {
		t.Errorf("createEIP err, err: %s", err)
	}
	if len(ip1) == 0 {
		t.Errorf("createEIP err, ip is nil")
	}
	svc.Spec.LoadBalancerIP = ip1
	lb.PublicIp = ip1
	_, err = cloud.ensureEIPWithSpecificIP(ctx, svc, lb)
	if err != nil {
		t.Errorf("ensureEIPWithSpecificIP err, err: %s", err)
	}
	// right case2 blb does not bind the eip
	cloud, _, lbResp, err = beforeTestBlb()
	if err != nil {
		t.Errorf("beforeTestBlb err , %v", err)
	}
	lb = &blb.LoadBalancer{
		BlbId: lbResp.LoadBalancerId,
	}
	eipArgs := &eip.CreateEIPArgs{
		Name:            "test",
		BandwidthInMbps: 100,
		Billing: &eip.Billing{
			PaymentTiming: "",
			BillingMethod: "",
		},
	}
	ip2, err := cloud.clientSet.EIPClient.CreateEIP(ctx, eipArgs, nil)
	if err != nil {
		t.Errorf("CreateEIP err, err； %s", err)
	}
	svc.Spec.LoadBalancerIP = ip2
	lb.PublicIp = ""
	_, err = cloud.ensureEIPWithSpecificIP(ctx, svc, lb)
	if err != nil {
		t.Errorf("ensureEIPWithSpecificIP err, err: %s", err)
	}
	// right case3 blb already bound eip what is not expected
	cloud, _, lbResp, err = beforeTestBlb()
	if err != nil {
		t.Errorf("beforeTestBlb err , %v", err)
	}
	lb = &blb.LoadBalancer{
		BlbId: lbResp.LoadBalancerId,
	}
	sa = &ServiceAnnotation{
		ElasticIPName:              "test1",
		ElasticIPBillingMethod:     "",
		ElasticIPBandwidthInMbps:   100,
		ElasticIPReservationLength: 100,
	}
	args, err = cloud.getEipArgsFromAnnotation(sa)
	if err != nil {
		t.Errorf("getEipArgsFromAnnotation err, err: %s", err)
	}
	ip3, err := cloud.createEIP(ctx, args)
	if err != nil {
		t.Errorf("createEIP err, err: %s", err)
	}
	if len(ip3) == 0 {
		t.Errorf("createEIP err, ip is nil")
	}
	eipArgs = &eip.CreateEIPArgs{
		Name:            "test2",
		BandwidthInMbps: 100,
		Billing: &eip.Billing{
			PaymentTiming: "",
			BillingMethod: "",
		},
	}
	ip4, err := cloud.clientSet.EIPClient.CreateEIP(ctx, eipArgs, nil)
	if err != nil {
		t.Errorf("CreateEIP err, err； %s", err)
	}
	svc.Spec.LoadBalancerIP = ip4
	lb.PublicIp = ip3
	_, err = cloud.ensureEIPWithSpecificIP(ctx, svc, lb)
	if err != nil {
		t.Errorf("ensureEIPWithSpecificIP err, err: %s", err)
	}
}
