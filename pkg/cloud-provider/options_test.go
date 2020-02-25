package cloud_provider

import (
	"strconv"
	"testing"

	api "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//测试案例1：不对annotation赋值，测试默认值是否解析正确
//测试案例2：赋值正确，测试能否正确解析结果
//测试案例3：赋值异常，测试能否可以识别错误

func buildService() *api.Service {
	return &api.Service{
		ObjectMeta: meta_v1.ObjectMeta{
			Name:      "foo",
			Namespace: api.NamespaceDefault,
		},
		Spec: api.ServiceSpec{},
	}
}

func TestExtractServiceAnnotationBLB(t *testing.T) {
	svc := buildService()

	result, err := ExtractServiceAnnotation(svc)
	if err != nil {
		t.Errorf("failed to extract service annotation: %v", err)
	}
	if result.CceAutoAddEip != "" {
		t.Errorf("extract service CceAutoAddEip annotation wrong")
	}
	if result.CceAutoAddLoadBalancerID != "" {
		t.Errorf("extract service CceAutoAddLoadBalancerID annotation wrong")
	}
	if result.LoadBalancerAllocateVip != "" {
		t.Errorf("extract service LoadBalancerAllocateVip annotation wrong")
	}
	if result.LoadBalancerExistID != "" {
		t.Errorf("extract service LoadBalancerExistID annotation wrong")
	}
	if strconv.Itoa(result.LoadBalancerHealthCheckInterval) != "0" {
		t.Errorf("extract service LoadBalancerHealthCheckInterval annotation wrong")
	}
	if result.LoadBalancerHealthCheckString != "" {
		t.Errorf("extract service LoadBalancerHealthCheckString annotation wrong")
	}
	if strconv.Itoa(result.LoadBalancerHealthCheckTimeoutInSecond) != "0" {
		t.Errorf("extract service LoadBalancerHealthCheckTimeoutInSecond annotation wrong")
	}
	if strconv.Itoa(result.LoadBalancerHealthyThreshold) != "0" {
		t.Errorf("extract service LoadBalancerHealthyThreshold annotation wrong")
	}
	if result.LoadBalancerInternalVpc != "" {
		t.Errorf("extract service LoadBalancerInternalVpc annotation wrong")
	}
	if strconv.Itoa(result.LoadBalancerRsMaxNum) != "0" {
		t.Errorf("extract service LoadBalancerRsMaxNum annotation wrong")
	}
	if result.LoadBalancerScheduler != "" {
		t.Errorf("extract service LoadBalancerScheduler annotation wrong")
	}
	if strconv.Itoa(result.LoadBalancerUnhealthyThreshold) != "0" {
		t.Errorf("extract service LoadBalancerUnhealthyThreshold annotation wrong")
	}
	if result.LoadBalancerSubnetID != "" {
		t.Errorf("extract service LoadBalancerSubnetID annotation wrong")
	}

	data := map[string]string{}
	data[ServiceAnnotationCceAutoAddLoadBalancerID] = "dsada11"
	data[ServiceAnnotationCceAutoAddEip] = "10.12.1.1"
	data[ServiceAnnotationLoadBalancerExistID] = "dsada11"
	data[ServiceAnnotationLoadBalancerInternalVpc] = "10.12.1.1"
	data[ServiceAnnotationLoadBalancerAllocateVip] = "10.12.1.1"
	data[ServiceAnnotationLoadBalancerSubnetID] = "10.12.1.1"
	data[ServiceAnnotationLoadBalancerRsMaxNum] = "11"
	data[ServiceAnnotationLoadBalancerScheduler] = "dd"
	data[ServiceAnnotationLoadBalancerHealthCheckTimeoutInSecond] = "11"
	data[ServiceAnnotationLoadBalancerHealthCheckInterval] = "11"
	data[ServiceAnnotationLoadBalancerUnhealthyThreshold] = "11"
	data[ServiceAnnotationLoadBalancerHealthyThreshold] = "11"
	data[ServiceAnnotationLoadBalancerHealthCheckString] = "dsada11"

	svc.SetAnnotations(data)
	result, err = ExtractServiceAnnotation(svc)
	if err != nil {
		t.Errorf("failed to extract service annotation: %v", err)
	}
	if result.CceAutoAddEip != "10.12.1.1" {
		t.Errorf("extract service CceAutoAddEip annotation wrong")
	}
	if result.CceAutoAddLoadBalancerID != "dsada11" {
		t.Errorf("extract service CceAutoAddLoadBalancerID annotation wrong")
	}
	if result.LoadBalancerAllocateVip != "10.12.1.1" {
		t.Errorf("extract service LoadBalancerAllocateVip annotation wrong")
	}
	if result.LoadBalancerExistID != "dsada11" {
		t.Errorf("extract service LoadBalancerExistID annotation wrong")
	}
	if strconv.Itoa(result.LoadBalancerHealthCheckInterval) != "11" {
		t.Errorf("extract service LoadBalancerHealthCheckInterval annotation wrong")
	}
	if result.LoadBalancerHealthCheckString != "dsada11" {
		t.Errorf("extract service LoadBalancerHealthCheckString annotation wrong")
	}
	if strconv.Itoa(result.LoadBalancerHealthCheckTimeoutInSecond) != "11" {
		t.Errorf("extract service LoadBalancerHealthCheckTimeoutInSecond annotation wrong")
	}
	if strconv.Itoa(result.LoadBalancerHealthyThreshold) != "11" {
		t.Errorf("extract service LoadBalancerHealthyThreshold annotation wrong")
	}
	if result.LoadBalancerInternalVpc != "10.12.1.1" {
		t.Errorf("extract service LoadBalancerInternalVpc annotation wrong")
	}
	if strconv.Itoa(result.LoadBalancerRsMaxNum) != "11" {
		t.Errorf("extract service LoadBalancerRsMaxNum annotation wrong")
	}
	if result.LoadBalancerScheduler != "dd" {
		t.Errorf("extract service LoadBalancerScheduler annotation wrong")
	}
	if strconv.Itoa(result.LoadBalancerUnhealthyThreshold) != "11" {
		t.Errorf("extract service LoadBalancerUnhealthyThreshold annotation wrong")
	}
	if result.LoadBalancerSubnetID != "10.12.1.1" {
		t.Errorf("extract service LoadBalancerSubnetID annotation wrong")
	}

	data2 := map[string]string{}

	data2[ServiceAnnotationLoadBalancerHealthCheckTimeoutInSecond] = "10.12.1.1"
	svc.SetAnnotations(data2)
	result, err = ExtractServiceAnnotation(svc)
	if err == nil {
		t.Errorf("extract service LoadBalancerHealthCheckTimeoutInSecond annotation wrong, should exist wrong")
	}

	data3 := map[string]string{}
	data3[ServiceAnnotationLoadBalancerHealthCheckInterval] = "dsada11"
	svc.SetAnnotations(data3)
	result, err = ExtractServiceAnnotation(svc)
	if err == nil {
		t.Errorf("extract service LoadBalancerHealthCheckInterval annotation wrong, should exist wrong")
	}

	data4 := map[string]string{}
	data4[ServiceAnnotationLoadBalancerUnhealthyThreshold] = "10.12.1.1"
	svc.SetAnnotations(data4)
	result, err = ExtractServiceAnnotation(svc)
	if err == nil {
		t.Errorf("extract service LoadBalancerUnhealthyThreshold annotation wrong, should exist wrong")
	}

	data5 := map[string]string{}
	data5[ServiceAnnotationLoadBalancerHealthyThreshold] = "10.12.1.1"
	svc.SetAnnotations(data5)
	result, err = ExtractServiceAnnotation(svc)
	if err == nil {
		t.Errorf("extract service LoadBalancerHealthyThreshold annotation wrong, should exist wrong")
	}

	data6 := map[string]string{}
	data6[ServiceAnnotationLoadBalancerRsMaxNum] = "1000"
	svc.SetAnnotations(data6)
	result, err = ExtractServiceAnnotation(svc)
	if err == nil {
		t.Errorf("extract service LoadBalancerRsMaxNum annotation wrong, should exist wrong")
	}

}

func TestExtractServiceAnnotationEIP(t *testing.T) {
	svc := buildService()

	result, err := ExtractServiceAnnotation(svc)
	if err != nil {
		t.Errorf("failed to extract service annotation: %v", err)
	}
	if strconv.Itoa(result.ElasticIPBandwidthInMbps) != "0" {
		t.Errorf("extract service ElasticIPBandwidthInMbps annotation wrong")
	}
	if result.ElasticIPBillingMethod != "" {
		t.Errorf("extract service ElasticIPBillingMethod annotation wrong")
	}
	if result.ElasticIPName != "" {
		t.Errorf("extract service ElasticIPName annotation wrong")
	}
	if result.ElasticIPPaymentTiming != "" {
		t.Errorf("extract service ElasticIPPaymentTiming annotation wrong")
	}
	if strconv.Itoa(result.ElasticIPReservationLength) != "0" {
		t.Errorf("extract service ElasticIPReservationLength annotation wrong")
	}
	data := map[string]string{}
	data[ServiceAnnotationElasticIPName] = "dasdawwdc11"
	data[ServiceAnnotationElasticIPPaymentTiming] = "wa111"
	data[ServiceAnnotationElasticIPBillingMethod] = "dd11"
	data[ServiceAnnotationElasticIPBandwidthInMbps] = "11"
	data[ServiceAnnotationElasticIPReservationLength] = "22"

	svc.SetAnnotations(data)
	result, err = ExtractServiceAnnotation(svc)
	if err != nil {
		t.Errorf("failed to extract service annotation: %v", err)
	}
	if strconv.Itoa(result.ElasticIPBandwidthInMbps) != "11" {
		t.Errorf("extract service ElasticIPBandwidthInMbps annotation wrong")
	}
	if result.ElasticIPBillingMethod != "dd11" {
		t.Errorf("extract service ElasticIPBillingMethod annotation wrong")
	}
	if result.ElasticIPName != "dasdawwdc11" {
		t.Errorf("extract service ElasticIPName annotation wrong")
	}
	if result.ElasticIPPaymentTiming != "wa111" {
		t.Errorf("extract service ElasticIPPaymentTiming annotation wrong")
	}
	if strconv.Itoa(result.ElasticIPReservationLength) != "22" {
		t.Errorf("extract service ElasticIPReservationLength annotation wrong")
	}

	data2 := map[string]string{}
	data2[ServiceAnnotationElasticIPBandwidthInMbps] = "1dd"
	svc.SetAnnotations(data2)
	result, err = ExtractServiceAnnotation(svc)
	if err == nil {
		t.Errorf("extract service ElasticIPBandwidthInMbps annotation wrong, should exist wrong")
	}

	data3 := map[string]string{}
	data3[ServiceAnnotationElasticIPReservationLength] = "22ddd"
	svc.SetAnnotations(data3)
	result, err = ExtractServiceAnnotation(svc)
	if err == nil {
		t.Errorf("extract service ElasticIPReservationLength annotation wrong, should exist wrong")
	}
}

func TestExtractNodeAnnotation(t *testing.T) {
	case1 := &api.Node{
		ObjectMeta: meta_v1.ObjectMeta{
			Name:      "foo",
			Namespace: api.NamespaceDefault,
		},
	}

	result, err := ExtractNodeAnnotation(case1)
	if err != nil {
		t.Errorf("extract node annotation wrong, %s", err)
	}
	if result.VpcID != "" {
		t.Errorf("extract node VpcId annotation wrong")
	}
	if result.VpcRouteTableID != "" {
		t.Errorf("extract node VpcRouteTableId annotation wrong")
	}
	if result.VpcRouteRuleID != "" {
		t.Errorf("extract node VpcRouteRuleId annotation wrong")
	}
	if result.CCMVersion != "" {
		t.Errorf("extract node CCMVersion annotation wrong")
	}
	if strconv.FormatBool(result.AdvertiseRoute) != "true" {
		t.Errorf("extract node AdvertiseRoute annotation wrong")
	}

	data := map[string]string{}

	data[NodeAnnotationVpcID] = "test11"
	data[NodeAnnotationVpcRouteTableID] = "test11"
	data[NodeAnnotationVpcRouteRuleID] = "test11"
	data[NodeAnnotationCCMVersion] = "test11"
	data[NodeAnnotationAdvertiseRoute] = "false"

	case1.SetAnnotations(data)

	result, err = ExtractNodeAnnotation(case1)
	if err != nil {
		t.Errorf("extract node annotation wrong, %s", err)
	}
	if result.VpcID != "test11" {
		t.Errorf("extract node VpcId annotation wrong")
	}
	if result.VpcRouteTableID != "test11" {
		t.Errorf("extract node VpcRouteTableId annotation wrong")
	}
	if result.VpcRouteRuleID != "test11" {
		t.Errorf("extract node VpcRouteRuleId annotation wrong")
	}
	if result.CCMVersion != "test11" {
		t.Errorf("extract node CCMVersion annotation wrong")
	}
	if strconv.FormatBool(result.AdvertiseRoute) != "false" {
		t.Errorf("extract node AdvertiseRoute annotation wrong")
	}

	data2 := map[string]string{}
	data2[NodeAnnotationAdvertiseRoute] = "ddd11"
	case1.SetAnnotations(data2)

	result, err = ExtractNodeAnnotation(case1)
	if err == nil {
		t.Errorf("extract node NodeAnnotationAdvertiseRoute annotation wrong, should exist wrong")
	}
}
