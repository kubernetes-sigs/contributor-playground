package fake

import (
	"context"
	"fmt"

	"icode.baidu.com/baidu/jpaas-caas/bce-sdk-go/bce"
	"icode.baidu.com/baidu/jpaas-caas/bce-sdk-go/util"
	"icode.baidu.com/baidu/jpaas-caas/bce-sdk-go/vpc"
)

// FakeClient implement of vpc.Interface
type VpcFakeClient struct {
	// vpcID
	VPCMap map[string]*vpc.VPC
	// subnetID
	SubnetMap map[string]*vpc.Subnet
	// routeruleID |
	RouteRuleMap map[string]vpc.RouteRule
	//  RuleTableID | VpcID
	VpcRuleTableMap map[string]string
}

// NewFakeClient for VPC fake client
func NewVpcFakeClient() *VpcFakeClient {
	return &VpcFakeClient{
		VPCMap:          map[string]*vpc.VPC{},
		SubnetMap:       map[string]*vpc.Subnet{},
		RouteRuleMap:    map[string]vpc.RouteRule{},
		VpcRuleTableMap: map[string]string{},
	}
}

// CreateVPC create VPC
func (f *VpcFakeClient) CreateVPC(ctx context.Context, args *vpc.CreateVPCArgs, option *bce.SignOption) (string, error) {
	if args == nil {
		return "", fmt.Errorf("CreateVPC faile: args is nil")
	}
	var routeTableID string
	vpc := &vpc.VPC{
		Name:        args.Name,
		Description: args.Description,
		CIDR:        args.CIDR,
	}
	for {
		vpcID := util.GenerateBCEShortID("vpc")
		if _, ok := f.VPCMap[vpcID]; !ok {
			vpc.VPCID = vpcID
			f.VPCMap[vpcID] = vpc
			routeTableID = f.generateVpcRoutetableMap(vpcID)
			break
		}
	}
	return vpc.VPCID + "/" + routeTableID, nil
}

// help func
func (f *VpcFakeClient) generateVpcRoutetableMap(vpcID string) string {
	for {
		routeTableID := util.GenerateBCEShortID("rt")
		if _, ok := f.VpcRuleTableMap[routeTableID]; !ok {
			f.VpcRuleTableMap[routeTableID] = vpcID
			return routeTableID
		}
	}
}

// ListVPC to list VPC of region
func (f *VpcFakeClient) ListVPC(ctx context.Context, args *vpc.ListVPCArgs, option *bce.SignOption) ([]*vpc.VPC, error) {
	vpcs := []*vpc.VPC{}
	for _, vpc := range f.VPCMap {
		vpcs = append(vpcs, vpc)
	}
	return vpcs, nil
}

// CreateSubnet to Create Subnet under VPC
func (f *VpcFakeClient) CreateSubnet(ctx context.Context, args *vpc.CreateSubnetArgs, option *bce.SignOption) (string, error) {
	if args == nil {
		return "", fmt.Errorf("CreateSubnet faile: args is nil")
	}
	subnet := &vpc.Subnet{
		Name:        args.Name,
		ZoneName:    args.ZoneName,
		CIDR:        args.CIDR,
		VPCID:       args.VPCID,
		SubnetType:  args.SubnetType,
		Description: args.Description,
	}
	for {
		subnetID := util.GenerateBCEShortID("sub")
		if _, ok := f.SubnetMap[subnetID]; !ok {
			subnet.SubnetID = subnetID
			f.SubnetMap[subnetID] = subnet
			break
		}
	}
	return subnet.SubnetID, nil
}

// ListSubnet to List Subnet under VPC
func (f *VpcFakeClient) ListSubnet(ctx context.Context, args *vpc.ListSubnetArgs, option *bce.SignOption) ([]*vpc.Subnet, error) {
	if args == nil {
		return []*vpc.Subnet{}, fmt.Errorf("ListSubnet failed: args is nil")
	}
	subnets := []*vpc.Subnet{}
	for _, subnet := range f.SubnetMap {
		isMatch := true
		if args.VPCID != "" {
			if subnet.VPCID != args.VPCID {
				isMatch = false
			}
		} else if args.SubnetType != "" {
			if subnet.SubnetType != args.SubnetType {
				isMatch = false
			}
		} else if args.ZoneName != "" {
			if subnet.ZoneName != args.ZoneName {
				isMatch = false
			}
		}
		if isMatch {
			subnets = append(subnets, subnet)
		}
	}
	return subnets, nil
}

// DescribeSubnet to Describe Subnet under VPC
func (f *VpcFakeClient) DescribeSubnet(ctx context.Context, subnetID string, option *bce.SignOption) (*vpc.Subnet, error) {
	for _, subnet := range f.SubnetMap {
		if subnet.SubnetID == subnetID {
			return subnet, nil
		}
	}
	return nil, fmt.Errorf("NoSuchObject")
}
func (f *VpcFakeClient) ListRouteTable(ctx context.Context, args *vpc.ListRouteArgs, option *bce.SignOption) ([]vpc.RouteRule, error) {
	if args == nil {
		return nil, fmt.Errorf("args is nil")
	}
	routeTableID := args.RouteTableID
	for k, v := range f.VpcRuleTableMap {
		if v == args.VpcID {
			routeTableID = k
		}
	}
	routerules := []vpc.RouteRule{}
	for _, routerule := range f.RouteRuleMap {
		if routerule.RouteTableID == routeTableID {
			routerules = append(routerules, routerule)
		}
	}
	return routerules, nil
}
func (f *VpcFakeClient) DeleteRoute(ctx context.Context, routeID string, option *bce.SignOption) error {
	if len(routeID) == 0 {
		return fmt.Errorf("routeID is nil")
	}
	for routeruleID, routerule := range f.RouteRuleMap {
		if routerule.RouteRuleID == routeID {
			delete(f.RouteRuleMap, routeruleID)
			return nil
		}
	}
	return fmt.Errorf("DeleteRoute %s not exist", routeID)
}
func (f *VpcFakeClient) CreateRouteRule(ctx context.Context, args *vpc.CreateRouteRuleArgs, option *bce.SignOption) (string, error) {
	if args == nil {
		return "", fmt.Errorf("args is nil")
	}
	routerule := vpc.RouteRule{
		RouteTableID:       args.RouteTableID,
		SourceAddress:      args.SourceAddress,
		DestinationAddress: args.DestinationAddress,
		NexthopID:          args.NexthopID,
		NexthopType:        args.NexthopType,
		Description:        args.Description,
	}
	for {
		routeruleID := util.GenerateBCEShortID("rr")
		if _, ok := f.RouteRuleMap[routeruleID]; !ok {
			routerule.RouteRuleID = routeruleID
			f.RouteRuleMap[routeruleID] = routerule
			break
		}
	}
	return routerule.RouteRuleID, nil
}
