package cloud_provider

import (
	"context"
	"testing"

	"k8s.io/apimachinery/pkg/types"
)

func TestGetZone(t *testing.T) {
	ctx := context.Background()

	cloud, _, err := newCluster()
	if err != nil {
		t.Errorf("create cluster error, %v", err)
	}

	zone, err := cloud.GetZone(ctx)
	if err != nil {
		t.Errorf("GetZone error, %v", err)
	}
	if zone.FailureDomain != "unknow" || zone.Region != cloud.Region {
		t.Errorf("GetZone error, want unknow and %s, get %v", cloud.Region, zone)
	}
}

// providerID = ""
// providerID = "test"
// providerID = "test/test"
// providerID = "test//test"
// right input
func TestGetZoneByProviderID(t *testing.T) {
	ctx := context.Background()

	cloud, _, err := newCluster()
	if err != nil {
		t.Errorf("create cluster error, %v", err)
	}

	cases := []string{
		"",
		"test",
		"test/test",
		"test//test",
	}

	for _, c := range cases {
		_, err := cloud.GetZoneByProviderID(ctx, c)
		if err == nil {
			t.Errorf("GetZoneByProviderID err, should be error here!!!")
		}
	}
}

// name = ""
// name = "test"
// right input
func TestGetZoneByNodeName(t *testing.T) {
	ctx := context.Background()

	cloud, resp, err := newCluster()
	if err != nil {
		t.Errorf("create cluster error, %v", err)
	}

	cases := []types.NodeName{
		"",
		"test1",
	}

	for _, c := range cases {
		_, err := cloud.GetZoneByNodeName(ctx, c)
		if err == nil {
			t.Errorf("GetZoneByNodeName should get error")
		}
	}

	cases = []types.NodeName{
		types.NodeName(resp.Nodes[0].Hostname),
		types.NodeName(resp.Nodes[0].IP),
	}

	for _, c := range cases {
		zone, err := cloud.GetZoneByNodeName(ctx, c)
		if err != nil {
			t.Errorf("GetZoneByNodeName error, %v", err)
		}
		if zone.FailureDomain != resp.Nodes[0].AvailableZone ||
			zone.Region != cloud.Region {
			t.Errorf("GetZoneByProviderID err, want %s %s, get %s %s",
				resp.Nodes[0].AvailableZone, cloud.Region, zone.FailureDomain, zone.Region)
		}
	}
}
