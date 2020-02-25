package temp_cce

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"icode.baidu.com/baidu/jpaas-caas/bce-sdk-go/bce"
)

func newClient(accessKeyID, secretAccessKey, region, endpoint string) *Client {
	return NewClient(&Config{
		&bce.Config{
			Credentials: bce.NewCredentials(accessKeyID, secretAccessKey),
			Checksum:    true,
			Timeout:     30 * time.Second,
			Region:      region,
			Endpoint:    endpoint,
		},
	})
}

// TODO: just for debug manually
func testListClusterNodes(t *testing.T) {
	ctx := context.Background()

	// Init BCE Client
	ak := "xxxxxxxx"
	sk := "xxxxxxxx"
	region := "sz"
	endpoint := "cce.su.baidubce.com"

	c := newClient(ak, sk, region, endpoint)

	// Test ListClusterNodes
	nodesResq, err := c.ListClusterNodes(ctx, "xxxxxx", nil)
	if err != nil {
		t.Errorf("ListClusterNodes failed: %v", err)
		return
	}

	str, _ := json.Marshal(nodesResq)
	t.Errorf("ListClusterNodes failed: %v", string(str))
}
