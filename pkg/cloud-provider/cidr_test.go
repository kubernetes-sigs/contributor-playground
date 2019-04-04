package cloud_provider

import (
	"net"
	"testing"
)

func TestVerifyNoOverlap(t *testing.T) {
	_, cidrBlock, err := net.ParseCIDR("0.0.0.0/0")
	if err != nil {
		t.Error(err)
	}
	_, cceCidr, err := net.ParseCIDR("192.168.1.0/24")
	if err != nil {
		t.Error(err)
	}
	_, customRightCidr, err := net.ParseCIDR("172.16.0.0/29")
	if err != nil {
		t.Error(err)
	}
	_, customWrongCidr, err := net.ParseCIDR("192.168.1.0/26")
	if err != nil {
		t.Error(err)
	}
	err = VerifyNoOverlap([]*net.IPNet{cceCidr, customRightCidr,}, cidrBlock)
	if err != nil {
		t.Error(err)
	}
	err = VerifyNoOverlap([]*net.IPNet{cceCidr, customWrongCidr,}, cidrBlock)
	if err != nil {
		t.Log(err)
	}
}
