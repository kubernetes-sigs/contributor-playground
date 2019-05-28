package cloud_provider

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/golang/glog"

	v1 "k8s.io/api/core/v1"
	"k8s.io/cloud-provider-baiducloud/pkg/cloud-sdk/blb"
	"k8s.io/cloud-provider-baiducloud/pkg/cloud-sdk/eip"
)

func (bc *Baiducloud) ensureEIP(ctx context.Context, clusterName string, service *v1.Service, nodes []*v1.Node, serviceAnnotation *ServiceAnnotation, lb *blb.LoadBalancer) (string, error) {
	if lb == nil {
		return "", fmt.Errorf("[%v %v] EnsureLoadBalancer: ensureEIP need not nil lb", service.Namespace, service.Name)
	}
	if len(service.Spec.LoadBalancerIP) == 0 {
		// not set LoadBalancerIP
		return bc.ensureEIPWithNoSpecificIP(ctx, clusterName, service, nodes, serviceAnnotation, lb)
	} else {
		// set LoadBalancerIP
		return bc.ensureEIPWithSpecificIP(ctx, clusterName, service, nodes, serviceAnnotation, lb)
	}
}

func (bc *Baiducloud) ensureEIPWithNoSpecificIP(ctx context.Context, clusterName string, service *v1.Service, nodes []*v1.Node, serviceAnnotation *ServiceAnnotation, lb *blb.LoadBalancer) (string, error) {
	pubIP := lb.PublicIp
	if len(pubIP) == 0 { // blb not bind eip, mostly case ==>
		glog.V(2).Infof("[%v %v] EnsureLoadBalancer: createEIP!", service.Namespace, service.Name)
		args, err := bc.getEipArgsFromAnnotation(serviceAnnotation)
		if err != nil {
			glog.Errorf("[%v %v] getEipArgsFromAnnotation failed: %v", service.Namespace, service.Name, err)
			return "", err
		}
		if len(args.Name) == 0 {
			args.Name = lb.Name // default EIP name = lb name
		}
		//sometimes there are several times to get EIP
		if !strings.Contains(lb.Desc, "cce_auto_create_eip") {
			lb.Desc = "cce_auto_create_eip" + lb.Desc
			newLbArg := blb.UpdateLoadBalancerArgs{
				LoadBalancerId: lb.BlbId,
				Desc:           lb.Desc,
				Name:           lb.Name,
			}
			err = bc.clientSet.Blb().UpdateLoadBalancer(&newLbArg)
			if err != nil {
				return "", err
			}
		}
		glog.V(3).Infof("lb.Desc: %s", lb.Desc)
		pubIP, err = bc.createEIP(args, lb)
		if err != nil {
			if pubIP != "" {
				args := eip.EipArgs{
					Ip: pubIP,
				}
				bc.clientSet.Eip().DeleteEip(&args)
			}
			return "", err
		}
	} else { // blb already bind eip
		glog.V(3).Infof("[%v %v] EnsureLoadBalancer: blb's eip already exists, start to ensure...", service.Namespace, service.Name)
		argsGet := eip.GetEipsArgs{
			Ip: pubIP,
		}
		eips, err := bc.clientSet.Eip().GetEips(&argsGet)
		if err != nil {
			return "", err
		}
		if len(eips) == 0 {
			err = fmt.Errorf("[%v %v] EnsureLoadBalancer: EIP %s not Exist", service.Namespace, service.Name, pubIP)
			return "", err
		}
		targetEip := eips[0]
		if (len(serviceAnnotation.ElasticIPPaymentTiming) != 0 && serviceAnnotation.ElasticIPPaymentTiming != targetEip.PaymentTiming) ||
			(len(serviceAnnotation.ElasticIPBillingMethod) != 0 && serviceAnnotation.ElasticIPBillingMethod != targetEip.BillingMethod) {
			glog.V(3).Infof("[%v %v] EnsureLoadBalancer: EIP config change, need delete old eip and create new one", service.Namespace, service.Name)
			// TODO
			//pubIP, err = bc.deleteOldAndCreateNewEip(service, serviceAnnotation, pubIP, lb)
			//if err != nil {
			//	return "", err
			//}
			return "", fmt.Errorf("not support change ElasticIP PaymentTiming or ElasticIP BillingMethod, you can delete old and create a new one")
		}
		if serviceAnnotation.ElasticIPBandwidthInMbps != 0 && serviceAnnotation.ElasticIPBandwidthInMbps != targetEip.BandwidthInMbps {
			glog.V(3).Infof("[%v %v] EnsureLoadBalancer: EIP config change, need change ElasticIPBandwidthInMbps", service.Namespace, service.Name)
			// just validate args
			_, err := bc.getEipArgsFromAnnotation(serviceAnnotation)
			if err != nil {
				glog.Errorf("[%v %v] Eip Args error: %v", service.Namespace, service.Name, err)
				return "", err
			}
			err = bc.resizeEip(service, serviceAnnotation, pubIP)
			if err != nil {
				return "", err
			}
		}
	}
	return pubIP, nil
}

func (bc *Baiducloud) ensureEIPWithSpecificIP(ctx context.Context, clusterName string, service *v1.Service, nodes []*v1.Node, serviceAnnotation *ServiceAnnotation, lb *blb.LoadBalancer) (string, error) {
	pubIP := lb.PublicIp
	loadBalancerIP := service.Spec.LoadBalancerIP
	glog.V(3).Infof("[%v %v] EnsureLoadBalancer: Try to bind Custom LoadBalancerIP %s to BLB %s.", service.Namespace, service.Name, loadBalancerIP, lb.BlbId)
	if len(pubIP) == 0 { // blb not bind target eip
		// check eip status
		argsGet := eip.GetEipsArgs{
			Ip: loadBalancerIP,
		}
		eips, err := bc.clientSet.Eip().GetEips(&argsGet)
		if err != nil {
			return "", err
		}
		if len(eips) == 0 {
			err = fmt.Errorf("[%v %v] EnsureLoadBalancer: EIP %s not Exist", service.Namespace, service.Name, loadBalancerIP)
			return "", err
		} else {
			eipStatus := eips[0].Status
			for index := 0; (index < 10) && (eipStatus != "available"); index++ {
				glog.V(3).Infof("[%v %v] Eip: %s is not available, retry:  %d", service.Namespace, service.Name, loadBalancerIP, index)
				time.Sleep(10 * time.Second)
				eips, err := bc.clientSet.Eip().GetEips(&argsGet)
				if err != nil {
					return "", err
				}
				eipStatus = eips[0].Status
			}
			glog.V(3).Infof("Eip final status is: %s", eipStatus)
			if eipStatus != "available" {
				return "", fmt.Errorf("[%v %v] EnsureLoadBalancer: target eip not available", service.Namespace, service.Name)
			}
		}

		// bind
		lb.Status = "unknown" // add here to do loop
		for index := 0; (index < 10) && (lb.Status != "available"); index++ {
			glog.V(3).Infof("[%v %v] BLB: %s is not available, retry:  %d", service.Namespace, service.Name, lb.BlbId, index)
			time.Sleep(10 * time.Second)
			newlb, exist, err := bc.getBCELoadBalancerById(lb.BlbId)
			if err != nil {
				glog.V(3).Infof("getBCELoadBalancer error: %s", lb.BlbId)
				return "", err
			}
			if !exist {
				glog.V(3).Infof("getBCELoadBalancer not exist: %s", lb.BlbId)
				return "", fmt.Errorf("BLB not exists:%s", lb.BlbId)
			}
			lb = newlb
			glog.V(3).Infof("[%v %v] BLB status is : %s", service.Namespace, service.Name, lb.Status)
		}
		argsBind := &eip.BindEipArgs{
			Ip:           loadBalancerIP,
			InstanceId:   lb.BlbId,
			InstanceType: eip.BLB,
		}
		glog.V(3).Infof("[%v %v] Bind EIP: %v", service.Namespace, service.Name, argsBind)
		glog.V(3).Infof("[%v %v] Bind BLB: %v", service.Namespace, service.Name, lb)
		err = bc.clientSet.Eip().BindEip(argsBind)
		if err != nil {
			glog.V(3).Infof("BindEip error: %v", err)
			return "", err
		}
		lb.PublicIp = loadBalancerIP
		pubIP = loadBalancerIP
		glog.V(3).Infof("[%v %v] EnsureLoadBalancer: Bind EIP to BLB success.", service.Namespace, service.Name)
	} else { // blb already bind eip
		if pubIP == loadBalancerIP { // blb bind correct LoadBalancerIP
			glog.V(3).Infof("[%v %v] EnsureLoadBalancer: BLB %s already bind EIP %s.", service.Namespace, service.Name, lb.BlbId, pubIP)
		} else { // blb not bind correct LoadBalancerIP, need update
			glog.V(3).Infof("[%v %v] EnsureLoadBalancer: BLB %s already bind EIP %s, but need updating to %s.", service.Namespace, service.Name, lb.BlbId, pubIP, loadBalancerIP)
			// check eip status
			argsGet := eip.GetEipsArgs{
				Ip: pubIP,
			}
			eips, err := bc.clientSet.Eip().GetEips(&argsGet)
			if err != nil {
				return "", err
			}
			if len(eips) > 0 {
				unbindArgs := eip.EipArgs{
					Ip: pubIP,
				}
				// just unbind, not delete
				err := bc.clientSet.Eip().UnbindEip(&unbindArgs)
				if err != nil {
					glog.V(3).Infof("Unbind Eip error : %s", err.Error())
					return "", err
				}
			}
			// bind
			lb.Status = "unknown" // add here to do loop
			for index := 0; (index < 10) && (lb.Status != "available"); index++ {
				glog.V(3).Infof("[%v %v] BLB: %s is not available, retry:  %d", service.Namespace, service.Name, lb.BlbId, index)
				time.Sleep(10 * time.Second)
				newlb, exist, err := bc.getBCELoadBalancerById(lb.BlbId)
				if err != nil {
					glog.V(3).Infof("getBCELoadBalancer error: %s", lb.BlbId)
					return "", err
				}
				if !exist {
					glog.V(3).Infof("getBCELoadBalancer not exist: %s", lb.BlbId)
					return "", fmt.Errorf("BLB not exists:%s", lb.BlbId)
				}
				lb = newlb
				glog.V(3).Infof("[%v %v] BLB status is : %s", service.Namespace, service.Name, lb.Status)
			}
			argsBind := &eip.BindEipArgs{
				Ip:           loadBalancerIP,
				InstanceId:   lb.BlbId,
				InstanceType: eip.BLB,
			}
			glog.V(3).Infof("[%v %v] Bind EIP: %v", service.Namespace, service.Name, argsBind)
			glog.V(3).Infof("[%v %v] Bind BLB: %v", service.Namespace, service.Name, lb)
			err = bc.clientSet.Eip().BindEip(argsBind)
			if err != nil {
				glog.V(3).Infof("BindEip error: %v", err)
				return "", err
			}
			lb.PublicIp = loadBalancerIP
			pubIP = loadBalancerIP
			glog.V(3).Infof("[%v %v] EnsureLoadBalancer: Bind EIP to BLB success.", service.Namespace, service.Name)
		}
	}
	return pubIP, nil
}

func (bc *Baiducloud) createEIP(args *eip.CreateEipArgs, lb *blb.LoadBalancer) (string, error) {
	// according to eip api doc, limit to 65
	if len(args.Name) > 65 {
		args.Name = args.Name[:65]
	}
	glog.V(3).Infof("CreateEip:  %v", args)
	ip, err := bc.clientSet.Eip().CreateEip(args)
	if err != nil {
		return "", err
	}
	argsGet := eip.GetEipsArgs{
		Ip: ip,
	}
	eips, err := bc.clientSet.Eip().GetEips(&argsGet)
	if err != nil {
		return "", err
	}
	if len(eips) > 0 {
		eipStatus := eips[0].Status
		for index := 0; (index < 10) && (eipStatus != "available"); index++ {
			glog.V(3).Infof("Eip: %s is not available, retry:  %d", ip, index)
			time.Sleep(10 * time.Second)
			eips, err := bc.clientSet.Eip().GetEips(&argsGet)
			if err != nil {
				return "", err
			}
			if len(eips) == 0 {
				return "", fmt.Errorf("createEIP failed: CreateEip success but query failed")
			}
			eipStatus = eips[0].Status
		}
		glog.V(3).Infof("Eip status is: %s", eipStatus)
	}
	lb.Status = "unknown" // add here to do loop
	for index := 0; (index < 10) && (lb.Status != "available"); index++ {
		glog.V(3).Infof("BLB: %s is not available, retry:  %d", lb.BlbId, index)
		time.Sleep(10 * time.Second)
		newlb, exist, err := bc.getBCELoadBalancerById(lb.BlbId)
		if err != nil {
			glog.V(3).Infof("getBCELoadBalancer error: %s", lb.BlbId)
			return "", err
		}
		if !exist {
			glog.V(3).Infof("getBCELoadBalancer not exist: %s", lb.BlbId)
			return "", fmt.Errorf("BLB not exists:%s", lb.BlbId)
		}
		lb = newlb
		glog.V(3).Infof("BLB status is : %s", lb.Status)
	}
	argsBind := &eip.BindEipArgs{
		Ip:           ip,
		InstanceId:   lb.BlbId,
		InstanceType: eip.BLB,
	}
	glog.V(3).Infof("BindEip:  %v", argsBind)
	glog.V(3).Infof("Bind BLB: %v", lb)
	err = bc.clientSet.Eip().BindEip(argsBind)
	if err != nil {
		glog.V(3).Infof("BindEip error: %v", err)
		return ip, err
	}
	lb.PublicIp = ip
	glog.V(3).Infof("createEIP: lb.PublicIp is %s", lb.PublicIp)
	return ip, nil
}

func (bc *Baiducloud) deleteEIP(ip string) error {
	argsGet := eip.GetEipsArgs{
		Ip: ip,
	}
	eips, err := bc.clientSet.Eip().GetEips(&argsGet)
	if err != nil {
		return err
	}
	if len(eips) > 0 {
		eipStatus := eips[0].Status
		for index := 0; (index < 10) && (eipStatus != "available"); index++ {
			glog.V(3).Infof("Eip: %s is not available, retry:  %d", ip, index)
			time.Sleep(10 * time.Second)
			eips, err := bc.clientSet.Eip().GetEips(&argsGet)
			if err != nil {
				return err
			}
			eipStatus = eips[0].Status
		}
	}
	args := eip.EipArgs{
		Ip: ip,
	}
	err = bc.clientSet.Eip().DeleteEip(&args)
	if err != nil {
		return err
	}
	return nil
}

func (bc *Baiducloud) getEipArgsFromAnnotation(serviceAnnotation *ServiceAnnotation) (*eip.CreateEipArgs, error) {
	var args *eip.CreateEipArgs

	paymentTiming := serviceAnnotation.ElasticIPPaymentTiming
	if len(paymentTiming) == 0 {
		paymentTiming = eip.PAYMENTTIMING_POSTPAID // default Postpaid
	}
	billingMethod := serviceAnnotation.ElasticIPBillingMethod
	if len(billingMethod) == 0 {
		billingMethod = eip.BILLINGMETHOD_BYTRAFFIC // default ByTraffic
	}
	bandwidthInMbps := serviceAnnotation.ElasticIPBandwidthInMbps
	reservationLength := serviceAnnotation.ElasticIPReservationLength
	switch paymentTiming {
	case eip.PAYMENTTIMING_PREPAID:
		if len(serviceAnnotation.ElasticIPBillingMethod) != 0 {
			return nil, fmt.Errorf("when using Prepaid EIP, do not need to set ElasticIPBillingMethod")
		}
		if bandwidthInMbps == 0 { // not set bandwidthInMbps
			bandwidthInMbps = 200
		} else {
			if bandwidthInMbps < 1 || bandwidthInMbps > 200 {
				return nil, fmt.Errorf("prepaid EIP bandwidthInMbps should in [1, 200]")
			}
		}
		reservationLengthAllowed := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36}
		rightReservationLength := false
		for _, length := range reservationLengthAllowed {
			if reservationLength == length {
				rightReservationLength = true
			}
		}
		if !rightReservationLength {
			return nil, fmt.Errorf("prepaid EIP reservationLength should in [1,2,3,4,5,6,7,8,9,12,24,36]")
		}
		args = &eip.CreateEipArgs{
			Name:            serviceAnnotation.ElasticIPName,
			BandwidthInMbps: bandwidthInMbps,
			Billing: &eip.Billing{
				PaymentTiming: paymentTiming,
				Reservation: &eip.Reservation{
					ReservationLength:   reservationLength,
					ReservationTimeUnit: "Month",
				},
			},
		}
	case eip.PAYMENTTIMING_POSTPAID:
		switch billingMethod {
		case eip.BILLINGMETHOD_BYTRAFFIC:
			if bandwidthInMbps == 0 { // not set bandwidthInMbps
				bandwidthInMbps = 1000
			} else {
				if bandwidthInMbps < 1 || bandwidthInMbps > 1000 {
					return nil, fmt.Errorf("postpaid ByTraffic EIP bandwidthInMbps should in [1, 1000]")
				}
			}
		case eip.BILLINGMETHOD_BYBANDWIDTH:
			if bandwidthInMbps == 0 { // not set bandwidthInMbps
				bandwidthInMbps = 200
			} else {
				if bandwidthInMbps < 1 || bandwidthInMbps > 200 {
					return nil, fmt.Errorf("postpaid ByBandwidth EIP bandwidthInMbps should in [1, 200]")
				}
			}
		default:
			return nil, fmt.Errorf("not support target ElasticIPBillingMethod: %v", billingMethod)
		}
		args = &eip.CreateEipArgs{
			Name:            serviceAnnotation.ElasticIPName,
			BandwidthInMbps: bandwidthInMbps,
			Billing: &eip.Billing{
				PaymentTiming: paymentTiming,
				BillingMethod: billingMethod,
			},
		}
	default:
		return nil, fmt.Errorf("not support target ElasticIPPaymentTiming: %v", paymentTiming)
	}

	return args, nil
}

func (bc *Baiducloud) deleteOldAndCreateNewEip(service *v1.Service, serviceAnnotation *ServiceAnnotation, oldEip string, lb *blb.LoadBalancer) (string, error) {
	err := bc.deleteEIP(oldEip)
	if err != nil {
		return "", err
	}
	glog.V(2).Infof("[%v %v] EnsureLoadBalancer: createEIP!", service.Namespace, service.Name)
	args, err := bc.getEipArgsFromAnnotation(serviceAnnotation)
	if err != nil {
		glog.Errorf("[%v %v] getEipArgsFromAnnotation failed: %v", service.Namespace, service.Name, err)
		return "", err
	}
	if len(args.Name) == 0 {
		args.Name = lb.Name // default EIP name = lb name
	}
	pubIP, err := bc.createEIP(args, lb)
	if err != nil {
		if pubIP != "" {
			args := eip.EipArgs{
				Ip: pubIP,
			}
			bc.clientSet.Eip().DeleteEip(&args)
		}
		return "", err
	}
	return pubIP, nil
}

func (bc *Baiducloud) resizeEip(service *v1.Service, serviceAnnotation *ServiceAnnotation, targetEip string) error {
	return bc.clientSet.Eip().ResizeEip(&eip.ResizeEipArgs{
		BandwidthInMbps: serviceAnnotation.ElasticIPBandwidthInMbps,
		Ip:              targetEip,
	})
}
