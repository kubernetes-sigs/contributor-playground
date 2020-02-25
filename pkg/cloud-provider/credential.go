package cloud_provider

import (
	"context"
	"io/ioutil"
	"strconv"
	"time"

	"github.com/virtual-kubelet/virtual-kubelet/log"
	"icode.baidu.com/baidu/jpaas-caas/bce-sdk-go/bce"
)

const (
	TokenHeaderKey      = "cce-token"
	ClusterIDHeaderKey  = "cce-cluster"
	RemoteHostHeaderKey = "cce-remote-host"
)

var (
	tokenFilename     = "/var/run/secrets/cce/cce-plugin-token/token"
	expiredAtFilename = "/var/run/secrets/cce/cce-plugin-token/expiredAt"
	token             string
	expiredAt         int64
)

func (bc *Baiducloud) getSignOption(ctx context.Context) *bce.SignOption {
	if token == "" || time.Now().Unix() >= expiredAt {
		tokenBytes, err := ioutil.ReadFile(tokenFilename)
		if err != nil {
			log.G(ctx).WithError(err).Error("read token file failed")
			return nil
		}
		expiredAtBytes, err := ioutil.ReadFile(expiredAtFilename)
		if err != nil {
			log.G(ctx).WithError(err).Error("read expiredAt file failed")
			return nil
		}
		expiredAt, err = strconv.ParseInt(string(expiredAtBytes), 10, 64)
		if err != nil {
			log.G(ctx).WithField("expiredAt", string(expiredAtBytes)).WithError(err).
				Error("fail to parse expiredAt string")
			return nil
		}
		token = string(tokenBytes)
	}
	return &bce.SignOption{
		CustomSignFunc: func(ctx context.Context, req *bce.Request) {
			req.Header.Set(TokenHeaderKey, token)
			req.Header.Set(ClusterIDHeaderKey, bc.CloudConfig.ClusterID)
			req.Header.Set(RemoteHostHeaderKey, req.Host)
			req.Host, _ = getCCEGatewayHostAndPort(bc.CloudConfig.Region)
		},
	}
}
