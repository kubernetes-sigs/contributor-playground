package temp_cce

import (
	"context"
	"encoding/json"
	"fmt"

	"icode.baidu.com/baidu/jpaas-caas/bce-sdk-go/bce"
)

const EnableClusterRBAC FeatureType = "EnableClusterRBAC"

type FeatureType string

type CheckWhiteListResponse struct {
	IsExist bool `json:"isExist"`
}

func (c *Client) CheckWhiteList(ctx context.Context, featureType FeatureType, signOpt *bce.SignOption) (bool, error) {
	if len(featureType) == 0 {
		return false, fmt.Errorf("featureType cannot be empty")
	}
	params := map[string]string{
		"featureType": string(featureType),
	}
	req, err := bce.NewRequest("GET", c.GetURL("/v1/cluster/check_white_list", params), nil)
	if err != nil {
		return false, err
	}

	resp, err := c.SendRequest(ctx, req, signOpt)
	if err != nil {
		return false, err
	}

	bodyContent, err := resp.GetBodyContent()
	if err != nil {
		return false, err
	}

	checkWhiteListResponse := new(CheckWhiteListResponse)
	err = json.Unmarshal(bodyContent, checkWhiteListResponse)
	if err != nil {
		return false, err
	}

	return checkWhiteListResponse.IsExist, nil
}
