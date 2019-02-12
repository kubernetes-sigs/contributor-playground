目录
=================
  * [使用说明](#使用说明)
  * [快速开始](#快速开始)
  * [高级配置](#高级配置)
    * [固定EIP](#固定EIP)
    * [不分配EIP，即VPC内BLB](#不分配EIP，即VPC内BLB)
    * [BLB自动分配VIP](#BLB自动分配VIP)

# 使用说明
本文档会详细介绍如何在CCE下创建类型是**LoadBalancer**的Service。  
Kubernetes官方教程：[Services](https://kubernetes.io/docs/concepts/services-networking/service/)

# 快速开始
当用户创建类型是**LoadBalancer**的Service，默认情况下，CCE会联动的创建BLB，并为此BLB绑定EIP。  
以创建一个简单的Nginx为例：
```yaml
---
kind: Service
apiVersion: v1
metadata:
  name: nginx-service
spec:
  selector:
    app: nginx
  type: LoadBalancer
  ports:
  - name: nginx-port
    port: 80
    targetPort: 80
    protocol: TCP
---
apiVersion: apps/v1beta1
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  replicas: 1
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx
        ports:
        - containerPort: 80
```
（1）创建
```bash
$ kubectl create -f nginx.yaml
```
（2）查询EIP  
  IP 8.8.8.8 即为此Nginx的EIP。 
```bash
$ kubectl get svc
NAME            CLUSTER-IP     EXTERNAL-IP      PORT(S)        AGE
nginx-service   1.1.1.1        8.8.8.8          80:30274/TCP   5m
```
（3）查询BLB
```bash
$ kubectl get svc nginx-service -o jsonpath={.metadata.annotations}
map[service.beta.kubernetes.io/cce-load-balancer-id:lb-xxxxxx]
```
  lb-xxxxxx即为此Service的BLB的id。
 
（4）访问测试
```bash
$ curl -i http://8.8.8.8
```

# 高级配置
## 固定EIP
当用户删除Service并重新创建的时候，EIP会变，这样就需要去更改依赖于此IP的其他所有服务，所以CCE提供一种方式来固定此EIP。  
固定EIP的方案：  
（1）用户预先百度云上购买一个EIP实例  
（2）在创建Service时，设置loadBalancerIP为此EIP  
（3）创建Service，此时EXTERNAL-IP即为此EIP  
（4）删除Service，CCE只会解绑此EIP而不会释放此EIP，用户下次还可以继续使用  
示例如下：  
```yaml
---
kind: Service
apiVersion: v1
metadata:
  name: nginx-service-eip-with-load-balancer-ip
spec:
  selector:
    app: nginx-eip-with-load-balancer-ip
  type: LoadBalancer
  loadBalancerIP: 8.8.8.8
  ports:
  - name: nginx-port
    port: 80
    targetPort: 80
    protocol: TCP
---
apiVersion: apps/v1beta1
kind: Deployment
metadata:
  name: nginx-deployment-eip-with-load-balancer-ip
spec:
  replicas: 1
  template:
    metadata:
      labels:
        app: nginx-eip-with-load-balancer-ip
    spec:
      containers:
      - name: nginx
        image: nginx
        ports:
        - containerPort: 80
```
这样查到的EXTERNAL-IP即为此EIP：  
```
kubectl get svc nginx-service
NAME                                    TYPE           CLUSTER-IP       EXTERNAL-IP      PORT(S)        AGE
nginx-service-eip-with-loadBalancerIP   LoadBalancer   1.1.1.1          8.8.8.8          80:30601/TCP   1m
```

## 不分配EIP，即VPC内BLB
用户使用时：  
（1）设置Service.Spec.Type=LoadBalancer  
（2）为Service添加annotations，即service.beta.kubernetes.io/cce-load-balancer-internal-vpc: "true"  
示例如下：  
```yaml
---
kind: Service
apiVersion: v1
metadata:
  name: nginx-service-blb-internal-vpc
  annotations:
    service.beta.kubernetes.io/cce-load-balancer-internal-vpc: "true"
spec:
  selector:
    app: nginx-blb-internal-vpc
  type: LoadBalancer
  ports:
  - name: nginx-port
    port: 80
    targetPort: 80
    protocol: TCP
---
apiVersion: apps/v1beta1
kind: Deployment
metadata:
  name: nginx-deployment-blb-internal-vpc
spec:
  replicas: 1
  template:
    metadata:
      labels:
        app: nginx-blb-internal-vpc
    spec:
      containers:
      - name: nginx
        image: nginx
        ports:
        - containerPort: 80
```
这样查到的EXTERNAL-IP只能在VPC内访问：  
```
kubectl get svc nginx-service
NAME                             TYPE           CLUSTER-IP       EXTERNAL-IP      PORT(S)        AGE
nginx-service-blb-internal-vpc   LoadBalancer   1.1.1.1          2.2.2.2          80:30601/TCP   1m
```
**注：此内网BLB只能在一个VPC内的集群间正常使用；在使用同一个集群内的内网BLB时，会存在问题，建议在同一个集群内直接使用Service的ClusterIP**

## BLB自动分配VIP
通过给Service添加Annotation，此时创建的BLB会自动分配VIP。  
示例如下： 
```yaml
---
kind: Service
apiVersion: v1
metadata:
  name: nginx-service-blb-allocate-vip
  annotations:
    service.beta.kubernetes.io/cce-load-balancer-allocate-vip: "true"
spec:
  selector:
    app: nginx
  type: LoadBalancer
  ports:
  - name: nginx-port
    port: 80
    targetPort: 80
    protocol: TCP
---
apiVersion: apps/v1beta1
kind: Deployment
metadata:
  name: nginx-deployment-blb-allocate-vip
spec:
  replicas: 1
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx
        ports:
        - containerPort: 80
```
**注：此VIP只能在百度内网使用，查询VIP请参考BLB的API**
