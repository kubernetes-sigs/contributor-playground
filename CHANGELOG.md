# v1.0.8
**New Feature:**
- Support UDP LB Service [#41](https://github.com/baidu/cloud-provider-baiducloud/issues/41)

**Bug Fix:**
- Fix "创建了service 制定了LoadBalancer，把service删掉之后eip还在" [#52](https://github.com/baidu/cloud-provider-baiducloud/issues/52)
- Fix "LB Service创建失败：BLB和EIP对name长度有65的限制" [#58](https://github.com/baidu/cloud-provider-baiducloud/issues/58)
- Fix "反复Edit LB Service导致的创建失败" [#57](https://github.com/baidu/cloud-provider-baiducloud/issues/57)

# v1.0.7
**New Feature:**
- None

**Bug Fix:**
- Many bug fix for BLB && EIP leak.

# v1.0.4
**New Features:**
- add detail information about VPC on node annotation [#29](https://github.com/baidu/cloud-provider-baiducloud/issues/29)
- add Chinese README: https://github.com/baidu/cloud-provider-baiducloud/blob/master/docs/tutorial_zh-CN.md

**Bug Fix:**
- fix blb leak issue [#38](https://github.com/baidu/cloud-provider-baiducloud/issues/38)

# v1.0.3
**New Features:**
- support blb allocate vip [#27](https://github.com/baidu/cloud-provider-baiducloud/issues/27)

**Bug Fix:**
- handle error:create eip success but query failed [#26](https://github.com/baidu/cloud-provider-baiducloud/issues/26)

# v1.0.2
**New Features:**
- Support k8s v1.11.5
- support change blb name in console
- Support internal VPC BLB.

**Bug Fix:**
- Remove unused ScrubDNS interface from cloudprovider
- fix create BLB failed issue
- Fix router issue when node is recover from failure