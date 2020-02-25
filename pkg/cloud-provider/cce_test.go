package cloud_provider

import (
	"icode.baidu.com/baidu/jpaas-caas/cloud-provider-baiducloud/pkg/fake"
)

func NewFakeCloud(clusterID string) *Baiducloud {
	return &Baiducloud{
		CloudConfig: CloudConfig{
			ClusterID: clusterID,
		},
		clientSet: &ClientSet{
			BLBClient: fake.NewBlbFakeClient(),
			VPCClient: fake.NewVpcFakeClient(),
			CCEClient: fake.NewCceFakeClient(),
			EIPClient: fake.NewEipFakeClient(),
		},
	}
}
