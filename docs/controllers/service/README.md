# Loadbalancers

CCE cloud controller manager runs service controller, which is responsible for watching services of type `LoadBalancer` and creating BCE loadbalancers to satisfy its requirements. Here are some examples of how it's used.

## HTTP loadbalancer
Here's an example on how to create a simple http loadbalancer backed by nginx pods:
```
$ kubectl apply -f http-nginx.yml
service "nginx-service" created
deployment "nginx-deployment" created
```
Get more info about Service:
```
kubectl get svc nginx-service
NAME            TYPE           CLUSTER-IP       EXTERNAL-IP      PORT(S)        AGE
nginx-service   LoadBalancer   172.18.127.136   180.76.177.103   80:30601/TCP   1m
```
So, the EIP is `180.76.177.103`

## HTTP loadbalancer with loadBalancerIP
If loadBalancerIP is specified, the load-balancer will be created with the user-specified loadBalancerIP:
```
$ kubectl apply -f http-with-loadBalancerIP.yml
service "nginx-service" created
deployment "nginx-deployment" created
```
Get more info about Service:
```
kubectl get svc nginx-service
NAME            TYPE           CLUSTER-IP       EXTERNAL-IP      PORT(S)        AGE
nginx-service   LoadBalancer   172.18.127.136   178.11.24.19     80:30601/TCP   1m
```
As you can see, the EXTERNAL-IP is the user-specified loadBalancerIP.

## HTTP loadbalancer support internal VPC BLB
In a mixed environment it is sometimes necessary to route traffic from services inside the same VPC.
This can be achieved by adding the annotation to the service:
```
$ kubectl apply -f http-support-internal-vpc-blb.yml
service "nginx-service" created
deployment "nginx-deployment" created
```
Get more info about Service:
```
kubectl get svc nginx-service
NAME            TYPE           CLUSTER-IP       EXTERNAL-IP      PORT(S)        AGE
nginx-service   LoadBalancer   172.18.127.136   192.168.9.222    80:30601/TCP   1m
```
As you can see, the EXTERNAL-IP `192.168.9.222` can only be accessed inside the VPC.
