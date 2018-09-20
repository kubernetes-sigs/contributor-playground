# Service Annotations

CCE cloud controller manager watches for Services of type `LoadBalancer` and will create corresponding BCE Load Balancers matching the Kubernetes service. The Load Balancer can be configured by applying annotations to the Service resource. The following annotations can be used:

### service.beta.kubernetes.io/cce-load-balancer-internal-vpc

Indicate that the Service only has a BLB and can only be accessed inside the VPC.