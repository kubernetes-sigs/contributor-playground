module icode.baidu.com/baidu/jpaas-caas/cloud-provider-baiducloud

go 1.13

replace k8s.io/api => k8s.io/api v0.0.0-20191016110408-35e52d86657a

replace k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.0.0-20191016113550-5357c4baaf65

replace k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20191004115801-a2eda9f80ab8

replace k8s.io/apiserver => k8s.io/apiserver v0.0.0-20191016112112-5190913f932d

replace k8s.io/cli-runtime => k8s.io/cli-runtime v0.0.0-20191016114015-74ad18325ed5

replace k8s.io/client-go => k8s.io/client-go v0.0.0-20191016111102-bec269661e48

replace k8s.io/cloud-provider => k8s.io/cloud-provider v0.0.0-20191016115326-20453efc2458

replace k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.0.0-20191016115129-c07a134afb42

replace k8s.io/code-generator => k8s.io/code-generator v0.0.0-20191004115455-8e001e5d1894

replace k8s.io/component-base => k8s.io/component-base v0.0.0-20191016111319-039242c015a9

replace k8s.io/cri-api => k8s.io/cri-api v0.0.0-20190828162817-608eb1dad4ac

replace k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.0.0-20191016115521-756ffa5af0bd

replace k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.0.0-20191016112429-9587704a8ad4

replace k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.0.0-20191016114939-2b2b218dc1df

replace k8s.io/kube-proxy => k8s.io/kube-proxy v0.0.0-20191016114407-2e83b6f20229

replace k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.0.0-20191016114748-65049c67a58b

replace k8s.io/kubectl => k8s.io/kubectl v0.0.0-20191016120415-2ed914427d51

replace k8s.io/kubelet => k8s.io/kubelet v0.0.0-20191016114556-7841ed97f1b2

replace k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.0.0-20191016115753-cf0698c3a16b

replace k8s.io/metrics => k8s.io/metrics v0.0.0-20191016113814-3b1a734dba6e

replace k8s.io/node-api => k8s.io/node-api v0.0.0-20191016115955-b0b11a2622b0

replace k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.0.0-20191016112829-06bb3c9d77c9

replace k8s.io/sample-cli-plugin => k8s.io/sample-cli-plugin v0.0.0-20191016114214-d25a4244b17f

replace k8s.io/sample-controller => k8s.io/sample-controller v0.0.0-20191016113152-0c2dd40eec0c

require (
	github.com/astaxie/beego v1.12.0 // indirect
	github.com/gorilla/mux v1.7.0
	github.com/satori/go.uuid v1.2.0
	github.com/spf13/cobra v0.0.5
	github.com/virtual-kubelet/virtual-kubelet v1.2.1
	icode.baidu.com/baidu/jpaas-caas/bce-sdk-go v0.0.0-20191028131239-bdaf2e2d0a5d
	k8s.io/api v0.0.0
	k8s.io/apimachinery v0.0.0
	k8s.io/apiserver v0.0.0
	k8s.io/client-go v10.0.0+incompatible
	k8s.io/cloud-provider v0.0.0
	k8s.io/component-base v0.0.0
	k8s.io/klog v0.4.0
	k8s.io/kube-controller-manager v0.0.0
	k8s.io/kubernetes v1.16.2
	k8s.io/utils v0.0.0-20190801114015-581e00157fb1
)
