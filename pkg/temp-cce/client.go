package temp_cce

import (
	"icode.baidu.com/baidu/jpaas-caas/bce-sdk-go/bce"
)

// Endpoint contains all endpoints of Baidu Cloud CCE.
var Endpoint = map[string]string{
	"bj":  "cce.bj.baidubce.com",
	"gz":  "cce.gz.baidubce.com",
	"su":  "cce.su.baidubce.com",
	"hkg": "cce.hkg.baidubce.com",
	"fwh": "cce.fwh.baidubce.com",
	"bd":  "cce.bd.baidubce.com",
}

// Config contains all options for bos.Client.
type Config struct {
	*bce.Config
}

// NewConfig config of CCE Client
func NewConfig(config *bce.Config) *Config {
	return &Config{config}
}

// Client is the bos client implemention for Baidu Cloud BOS API.
type Client struct {
	*bce.Client
}

// NewClient client of CCE
func NewClient(config *Config) *Client {
	bceClient := bce.NewClient(config.Config)
	return &Client{bceClient}
}

// GetURL generates the full URL of http request for Baidu Cloud BOS API.
func (c *Client) GetURL(objectKey string, params map[string]string) string {
	host := c.Endpoint

	if host == "" {
		host = Endpoint[c.GetRegion()]
	}

	uriPath := objectKey

	return c.Client.GetURL(host, uriPath, params)
}
