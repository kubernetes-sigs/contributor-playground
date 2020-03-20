package fake

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"icode.baidu.com/baidu/jpaas-caas/bce-sdk-go/bce"
	"icode.baidu.com/baidu/jpaas-caas/bce-sdk-go/eip"
)

// FakeClient for unit test
type EipFakeClient struct {
	EIPMap map[string]*eip.EIP
}

// NewFakeClient for EIP fake client
func NewEipFakeClient() *EipFakeClient {
	return &EipFakeClient{
		EIPMap: map[string]*eip.EIP{},
	}
}

// CreateEIP create EIP
func (f *EipFakeClient) CreateEIP(ctx context.Context, args *eip.CreateEIPArgs, option *bce.SignOption) (string, error) {
	if args == nil {
		return "", fmt.Errorf("CreateEIP faile: args is nil")
	}
	eip := &eip.EIP{
		Name:            args.Name,
		Status:          eip.EIPAvailable,
		BandwidthInMbps: args.BandwidthInMbps,
	}
	for {
		ip := generateRandomEIP()
		if _, ok := f.EIPMap[ip]; !ok {
			eip.EIP = ip
			f.EIPMap[ip] = eip
			break
		}
	}
	return eip.EIP, nil
}

// BindEIP bind eip with instance
func (f *EipFakeClient) BindEIP(ctx context.Context, ip string, args *eip.BindEIPArgs, option *bce.SignOption) error {
	e, ok := f.EIPMap[ip]
	if !ok {
		return fmt.Errorf("EIP %s not exist", ip)
	}
	if e.InstanceType != "" || e.InstanceID != "" {
		return fmt.Errorf("EIP %s has already been binded", ip)
	}
	e.Status = eip.EIPBinded
	e.InstanceType = args.InstanceType
	e.InstanceID = args.InstanceID
	return nil
}

// UnbindEIP unbind EIP with instance
// If eip.status == EIPAvailable, return nil
func (f *EipFakeClient) UnbindEIP(ctx context.Context, ip string, option *bce.SignOption) error {
	e, ok := f.EIPMap[ip]
	if !ok {
		return fmt.Errorf("EIP %s not exist", ip)
	}
	e.Status = eip.EIPAvailable
	e.InstanceType = ""
	e.InstanceID = ""
	return nil
}

// DeleteEIP delete pointed EIP
func (f *EipFakeClient) DeleteEIP(ctx context.Context, eip string, option *bce.SignOption) error {
	if _, ok := f.EIPMap[eip]; ok {
		delete(f.EIPMap, eip)
		return nil
	}
	return fmt.Errorf("DeleteEIP %s not exist", eip)
}

// ResizeEIP resize EIP bindwidth
func (f *EipFakeClient) ResizeEIP(ctx context.Context, eip string, args *eip.ResizeEIPArgs, option *bce.SignOption) error {
	if args == nil {
		return fmt.Errorf("ResizeEIP failed: args is nil")
	}
	if args.BandwidthInMbps < 1 || args.BandwidthInMbps > 1000 {
		return fmt.Errorf("ResizeEIP failed: %d out of range", args.BandwidthInMbps)
	}
	e, ok := f.EIPMap[eip]
	if !ok {
		return fmt.Errorf("EIP %s not exist", eip)
	}
	e.BandwidthInMbps = args.BandwidthInMbps
	return nil
}

// GetEIPs to get eips by condition
func (f *EipFakeClient) GetEIPs(ctx context.Context, args *eip.GetEIPsArgs, option *bce.SignOption) ([]*eip.EIP, error) {
	result := []*eip.EIP{}
	// Return all EIPs
	if args == nil || args.EIP == "" {
		for _, eip := range f.EIPMap {
			result = append(result, eip)
		}
		return result, nil
	}
	// Only process EIP
	if args != nil && args.EIP != "" {
		for _, eip := range f.EIPMap {
			if eip.EIP == args.EIP {
				result = append(result, eip)
				return result, nil
			}
		}
	}
	return []*eip.EIP{}, nil
}
func generateRandomEIP() string {
	rand.Seed(time.Now().Unix())
	ip := fmt.Sprintf("100.%d.%d.%d", rand.Intn(255), rand.Intn(255), rand.Intn(255))
	return ip
}
