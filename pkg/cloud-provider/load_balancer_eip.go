package cloud_provider

import (
	"context"
	"fmt"
	"strings"
	"time"

	"k8s.io/klog"

	"icode.baidu.com/baidu/jpaas-caas/bce-sdk-go/blb"
	"icode.baidu.com/baidu/jpaas-caas/bce-sdk-go/eip"
	v1 "k8s.io/api/core/v1"
)

func (bc *Baiducloud) ensureEIP(ctx context.Context, clusterName string, service *v1.Service) (string, error) {
	startTime := time.Now()
	serviceKey := fmt.Sprintf("%s/%s", service.Namespace, service.Name)
	defer func() {
		klog.V(4).Infof(Message(ctx, fmt.Sprintf("Finished ensureEIP for service %q (%v)", serviceKey, time.Since(startTime))))
	}()
	lb, _, err := bc.getServiceAssociatedBLB(ctx, clusterName, service)
	if err != nil {
		return "", err
	}
	if lb == nil {
		return "", fmt.Errorf("ensureEIP for service %s failed: lb is nil", serviceKey)
	}

	if len(service.Spec.LoadBalancerIP) == 0 {
		// not set LoadBalancerIP
		return bc.ensureEIPWithNoSpecificIP(ctx, service, lb)
	} else {
		// set LoadBalancerIP
		return bc.ensureEIPWithSpecificIP(ctx, service, lb)
	}
}

func (bc *Baiducloud) ensureEIPWithNoSpecificIP(ctx context.Context, service *v1.Service, lb *blb.LoadBalancer) (string, error) {
	serviceAnnotation, err := ExtractServiceAnnotation(service)
	if err != nil {
		return "", err
	}
	pubIP := lb.PublicIp
	if len(pubIP) == 0 { // blb not bind eip, mostly case ==>
		klog.V(2).Infof("[%v %v] EnsureLoadBalancer: createEIP!", service.Namespace, service.Name)
		args, err := bc.getEipArgsFromAnnotation(serviceAnnotation)
		if err != nil {
			klog.Errorf("[%v %v] getEipArgsFromAnnotation failed: %v", service.Namespace, service.Name, err)
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
			err = bc.clientSet.BLBClient.UpdateLoadBalancer(ctx, &newLbArg, bc.getSignOption(ctx))
			if err != nil {
				return "", err
			}
		}
		klog.V(3).Infof("lb.Desc: %s", lb.Desc)

		pubIP, err = bc.getServiceAssociatedEip(ctx, service)
		if err != nil {
			return "", err
		}
		if len(pubIP) == 0 {
			pubIP, err = bc.createEIP(ctx, args)
			if err != nil {
				return "", err
			}
		}

		_, err = bc.bindEip(ctx, lb, pubIP, service)
		if err != nil {
			return "", err
		}
	} else { // blb already bind eip
		klog.V(3).Infof("[%v %v] EnsureLoadBalancer: blb's eip already exists, start to ensure...", service.Namespace, service.Name)
		eips, err := bc.getEipByIP(ctx, pubIP)
		if err != nil {
			return "", err
		}
		if eips == nil || len(eips) == 0 {
			err = fmt.Errorf("[%v %v] EnsureLoadBalancer: EIP %s not Exist", service.Namespace, service.Name, pubIP)
			return "", err
		}
		targetEip := eips[0]
		if (len(serviceAnnotation.ElasticIPPaymentTiming) != 0 && serviceAnnotation.ElasticIPPaymentTiming != targetEip.PaymentTiming) ||
			(len(serviceAnnotation.ElasticIPBillingMethod) != 0 && serviceAnnotation.ElasticIPBillingMethod != targetEip.BillingMethod) {
			klog.V(3).Infof("[%v %v] EnsureLoadBalancer: EIP config change, need delete old eip and create new one", service.Namespace, service.Name)
			// TODO
			//pubIP, err = bc.deleteOldAndCreateNewEip(service, serviceAnnotation, pubIP, lb)
			//if err != nil {
			//	return "", err
			//}
			return "", fmt.Errorf("not support change ElasticIP PaymentTiming or ElasticIP BillingMethod, you can delete old and create a new one")
		}
		if serviceAnnotation.ElasticIPBandwidthInMbps != 0 && serviceAnnotation.ElasticIPBandwidthInMbps != targetEip.BandwidthInMbps {
			klog.V(3).Infof("[%v %v] EnsureLoadBalancer: EIP config change, need change ElasticIPBandwidthInMbps", service.Namespace, service.Name)
			// just validate args
			_, err := bc.getEipArgsFromAnnotation(serviceAnnotation)
			if err != nil {
				klog.Errorf("[%v %v] Eip Args error: %v", service.Namespace, service.Name, err)
				return "", err
			}
			err = bc.resizeEip(ctx, serviceAnnotation, pubIP)
			if err != nil {
				return "", err
			}
		}
	}
	return pubIP, nil
}

func (bc *Baiducloud) ensureEIPWithSpecificIP(ctx context.Context, service *v1.Service, lb *blb.LoadBalancer) (string, error) {
	pubIP := lb.PublicIp
	loadBalancerIP := service.Spec.LoadBalancerIP
	klog.V(3).Infof("[%v %v] EnsureLoadBalancer: Try to bind Custom LoadBalancerIP %s to BLB %s.", service.Namespace, service.Name, loadBalancerIP, lb.BlbId)
	if len(pubIP) == 0 { // blb not bind target eip
		// check eip status & bind blb
		lb, err := bc.bindEip(ctx, lb, loadBalancerIP, service)
		if err != nil {
			return "", err
		}
		lb.PublicIp = loadBalancerIP
		pubIP = loadBalancerIP
		klog.V(3).Infof("[%v %v] EnsureLoadBalancer: Bind EIP to BLB success.", service.Namespace, service.Name)
	} else { // blb already bind eip
		if pubIP == loadBalancerIP { // blb bind correct LoadBalancerIP
			klog.V(3).Infof("[%v %v] EnsureLoadBalancer: BLB %s already bind EIP %s.", service.Namespace, service.Name, lb.BlbId, pubIP)
		} else { // blb not bind correct LoadBalancerIP, need update
			klog.V(3).Infof("[%v %v] EnsureLoadBalancer: BLB %s already bind EIP %s, but need updating to %s.", service.Namespace, service.Name, lb.BlbId, pubIP, loadBalancerIP)
			err := bc.unbindEip(ctx, lb, pubIP)
			if err != nil {
				return "", err
			}

			lb, err = bc.bindEip(ctx, lb, loadBalancerIP, service)
			if err != nil {
				return "", err
			}
			lb.PublicIp = loadBalancerIP
			pubIP = loadBalancerIP
			klog.V(3).Infof("[%v %v] EnsureLoadBalancer: Bind EIP to BLB success.", service.Namespace, service.Name)
		}
	}
	return pubIP, nil
}

func (bc *Baiducloud) unbindEip(ctx context.Context, lb *blb.LoadBalancer, ip string) error {
	eips, err := bc.getEipByIP(ctx, ip)
	if err != nil {
		return err
	}
	if eips == nil || len(eips) == 0 {
		klog.Warningf("EIP %s not found", ip)
		return nil
	}
	err = bc.clientSet.EIPClient.UnbindEIP(ctx, ip, bc.getSignOption(ctx))
	if err != nil {
		klog.V(3).Infof("Unbind Eip error : %s", err.Error())
		return err
	}
	return nil
}

func (bc *Baiducloud) bindEip(ctx context.Context, lb *blb.LoadBalancer, ip string, service *v1.Service) (*blb.LoadBalancer, error) {
	for i := 0; i < 10; i++ {
		eips, err := bc.getEipByIP(ctx, ip)
		if err != nil {
			return nil, err
		}
		if eips == nil || len(eips) == 0 {
			err = fmt.Errorf("[%v %v] EnsureLoadBalancer: EIP %s not Exist", service.Namespace, service.Name, ip)
			return nil, err
		}
		if eips[0].Status == eip.EIPAvailable {
			break
		}
		klog.Infof(Message(ctx, fmt.Sprintf("eip %s status is %s, not available, wait...", ip, eips[0].Status)))
		time.Sleep(3 * time.Second)
	}

	// bind blb
	argsBind := &eip.BindEIPArgs{
		IP:           ip,
		InstanceID:   lb.BlbId,
		InstanceType: eip.BLB,
	}
	klog.V(3).Infof("[%v %v] Bind EIP: %v", service.Namespace, service.Name, argsBind)
	klog.V(3).Infof("[%v %v] Bind BLB: %v", service.Namespace, service.Name, lb)
	err := bc.clientSet.EIPClient.BindEIP(ctx, ip, argsBind, bc.getSignOption(ctx))
	if err != nil {
		klog.V(3).Infof("BindEip error: %v", err)
		return nil, err
	}
	return lb, nil
}

func (bc *Baiducloud) refreshBlb(ctx context.Context, lb *blb.LoadBalancer) (*blb.LoadBalancer, error) {
	newlb, exist, err := bc.getBLBByID(ctx, lb.BlbId)
	if err != nil {
		klog.V(3).Infof("getBLBByName error: %s", lb.BlbId)
		return nil, err
	}
	if !exist {
		klog.V(3).Infof("getBLBByName not exist: %s", lb.BlbId)
		return nil, fmt.Errorf("BLB not exists:%s", lb.BlbId)
	}
	lb = newlb
	klog.V(3).Infof("BLB status is : %s", lb.Status)
	return lb, nil
}

func (bc *Baiducloud) createEIP(ctx context.Context, args *eip.CreateEIPArgs) (string, error) {
	// according to eip api doc, limit to 65
	if len(args.Name) > 65 {
		args.Name = args.Name[:65]
	}
	klog.Infof(Message(ctx, fmt.Sprintf("CreateEip:  %v", args)))
	ip, err := bc.clientSet.EIPClient.CreateEIP(ctx, args, bc.getSignOption(ctx))
	if err != nil {
		return "", err
	}
	klog.Infof(Message(ctx, fmt.Sprintf("CreatedEIP is %s", ip)))
	return ip, nil
}

func (bc *Baiducloud) deleteEIP(ctx context.Context, ip string) error {
	argsGet := eip.GetEIPsArgs{
		EIP: ip,
	}
	eips, err := bc.clientSet.EIPClient.GetEIPs(ctx, &argsGet, bc.getSignOption(ctx))
	if err != nil {
		return err
	}
	if len(eips) == 0 {
		klog.Infof(Message(ctx, fmt.Sprintf("specified eip %s not found, skip delete", ip)))
		return nil
	}

	if eips[0].Status != "Available" {
		time.Sleep(5 * time.Second)
	}

	err = bc.clientSet.EIPClient.DeleteEIP(ctx, ip, bc.getSignOption(ctx))
	if err != nil {
		return err
	}
	return nil
}

func (bc *Baiducloud) getEipArgsFromAnnotation(serviceAnnotation *ServiceAnnotation) (*eip.CreateEIPArgs, error) {
	var args *eip.CreateEIPArgs

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
		args = &eip.CreateEIPArgs{
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
		args = &eip.CreateEIPArgs{
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

func (bc *Baiducloud) resizeEip(ctx context.Context, serviceAnnotation *ServiceAnnotation, ip string) error {
	return bc.clientSet.EIPClient.ResizeEIP(ctx, ip, &eip.ResizeEIPArgs{
		BandwidthInMbps: serviceAnnotation.ElasticIPBandwidthInMbps,
		IP:              ip,
	}, bc.getSignOption(ctx))
}

func (bc *Baiducloud) ensureEipDeleted(ctx context.Context, service *v1.Service, lb *blb.LoadBalancer) error {
	startTime := time.Now()
	serviceKey := fmt.Sprintf("%s/%s", service.Namespace, service.Name)
	defer func() {
		klog.V(4).Infof(Message(ctx, fmt.Sprintf("Finished ensureEipDeleted for service %q (%v)", serviceKey, time.Since(startTime))))
	}()
	if len(service.Spec.LoadBalancerIP) != 0 {
		msg := fmt.Sprintf("service %s has fixed EIP %s, not delete it", serviceKey, service.Spec.LoadBalancerIP)
		klog.Info(Message(ctx, msg))
		if lb != nil {
			msg := fmt.Sprintf("service %s has fixed EIP %s, unbind it", serviceKey, service.Spec.LoadBalancerIP)
			klog.Info(Message(ctx, msg))
			err := bc.unbindEip(ctx, lb, service.Spec.LoadBalancerIP)
			if err != nil {
				return err
			}
		}
		return nil
	}
	// get eip
	var targetEip string
	if len(service.Status.LoadBalancer.Ingress) != 0 { // P0: use service EXTERNAL_IP
		targetEip = service.Status.LoadBalancer.Ingress[0].IP
		klog.Infof(Message(ctx, fmt.Sprintf("selected ip %s for service %s from status to delete", targetEip, serviceKey)))
	}
	if len(targetEip) == 0 && lb != nil { // P1: use BLB public ip
		targetEip = lb.PublicIp
		klog.Infof(Message(ctx, fmt.Sprintf("selected ip %s for service %s from lb.PublicIp to delete", targetEip, serviceKey)))
	}
	if len(targetEip) == 0 {
		ip, ok := service.Annotations[ServiceAnnotationCceAutoAddEip]
		if ok {
			targetEip = ip
			klog.Infof(Message(ctx, fmt.Sprintf("selected ip %s for service %s from annotation to delete", ip, serviceKey)))
		}
	}
	if len(targetEip) == 0 { // get none EIP
		klog.V(3).Infof(Message(ctx, fmt.Sprintf("Eip for service %s not exist, skipping ensureEipDeleted", serviceKey)))
		return nil
	}

	if lb != nil {
		err := bc.unbindEip(ctx, lb, targetEip)
		if err != nil {
			return err
		}
	}
	// delete eip
	err := bc.deleteEIP(ctx, targetEip)
	if err != nil {
		return err
	}
	return nil
}

func (bc *Baiducloud) getEipByIP(ctx context.Context, ip string) ([]*eip.EIP, error) {
	argsGet := eip.GetEIPsArgs{
		EIP: ip,
	}
	eips, err := bc.clientSet.EIPClient.GetEIPs(ctx, &argsGet, bc.getSignOption(ctx))
	if err != nil {
		return nil, err
	}
	return eips, nil
}

func (bc *Baiducloud) getEipsByName(ctx context.Context, name string) ([]*eip.EIP, error) {
	result := make([]*eip.EIP, 0)
	argsGet := eip.GetEIPsArgs{}
	eips, err := bc.clientSet.EIPClient.GetEIPs(ctx, &argsGet, bc.getSignOption(ctx))
	if err != nil {
		return nil, err
	}
	for _, ip := range eips {
		if ip.Name == name {
			result = append(result, ip)
		}
	}
	return result, nil
}

func (bc *Baiducloud) getServiceAssociatedEip(ctx context.Context, service *v1.Service) (string, error) {
	result, err := ExtractServiceAnnotation(service)
	if err != nil {
		return "", err
	}
	annotationEip := result.CceAutoAddEip
	if len(annotationEip) > 0 {
		klog.Infof(Message(ctx, fmt.Sprintf("getServiceAssociatedEip from annotation: %s", annotationEip)))
		return annotationEip, nil
	}

	sameNameEips, err := bc.getEipsByName(ctx, getBlbName(bc.ClusterID, service))
	if err != nil {
		return "", err
	}
	if len(sameNameEips) > 1 {
		return "", fmt.Errorf("has multi eips created for service %s/%s: %d", service.Namespace, service.Name, len(sameNameEips))
	}
	if len(sameNameEips) == 0 {
		return "", nil
	}
	return sameNameEips[0].EIP, nil
}
