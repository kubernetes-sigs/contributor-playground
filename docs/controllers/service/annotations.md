# Service Annotations

CCE cloud controller manager watches for Services of type `LoadBalancer` and will create corresponding BCE Load Balancers and Elastic IPs matching the Kubernetes service. The Load Balancer and Elastic IP can be configured by applying annotations to the Service resource. The following annotations can be used:

## BLB

### service.beta.kubernetes.io/cce-load-balancer-internal-vpc: "true"
Indicate that the Service only has a BLB and can only be accessed inside the VPC.

### service.beta.kubernetes.io/cce-load-balancer-allocate-vip: "true"
Indicate that the BLB for Service has a VIP.**(Baidu Internal Use)**

## EIP

### service.beta.kubernetes.io/cce-elastic-ip-payment-timing: ""
Set EIP payment. Support value:  
- Prepaid
- Postpaid

### service.beta.kubernetes.io/cce-elastic-ip-billing-method: ""
Set EIP billing method. Support value:  
- ByTraffic
- ByBandwidth

### service.beta.kubernetes.io/cce-elastic-ip-bandwidth-in-mbps: ""
Set EIP bandwidth. Support value:  
- 1~1000 for ByTraffic EIP
- 1~200 for Prepaid and ByBandwidth EIP

### service.beta.kubernetes.io/cce-elastic-ip-reservation-length: ""
Set EIP reservation length in month. Support value:  [1,2,3,4,5,6,7,8,9,12,24,36]