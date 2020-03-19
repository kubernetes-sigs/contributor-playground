package temp_cce

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"icode.baidu.com/baidu/jpaas-caas/bce-sdk-go/bce"
)

// GenerateSignatureArgs define signature args
type GenerateSignatureArgs struct {
	ClusterID      string      `json:"clusterid,omitempty"`
	Method         string      `json:"method"`
	URI            string      `json:"uri"`
	Headers        http.Header `json:"headers"`
	Queries        url.Values  `json:"queries"`
	KeyOnlyQueries []string    `json:"keyOnlyQueries"`
}

// Signature define signature
type Signature struct {
	Authorization string `json:"authorization"`
	SecurityToken string `json:"securityToken"`
}

// GenerateSignature to generate signature
func (c *Client) GenerateSignature(ctx context.Context, args *GenerateSignatureArgs) (*Signature, error) {
	params := map[string]string{
		"clientToken": c.GenerateClientToken(),
	}
	// not post ClusterId field in post body
	clusterID := args.ClusterID

	args.ClusterID = ""
	postContent, err := json.Marshal(args)
	if err != nil {
		return nil, err
	}
	req, err := bce.NewRequest("POST", c.GetURL("internal-api/v1/signature", params), bytes.NewBuffer(postContent))
	if err != nil {
		return nil, err
	}

	resp, err := c.SendRequest(ctx, req, &bce.SignOption{
		CustomUserAgent: fmt.Sprintf("cce-k8s:%s", clusterID),
	})
	if err != nil {
		return nil, err
	}

	bodyContent, err := resp.GetBodyContent()
	if err != nil {
		return nil, err
	}

	sig := new(Signature)
	err = json.Unmarshal(bodyContent, sig)
	if err != nil {
		return nil, err
	}

	return sig, nil
}
