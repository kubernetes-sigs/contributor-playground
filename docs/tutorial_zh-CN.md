目录
=================
  * [一、使用说明](#一、使用说明)
  * [二、快速开始](#二、快速开始)
  * [三、高级配置](#三、高级配置)
    * [3.1固定EIP](#3.1固定EIP)
    * [3.2自定义EIP配置](#3.2自定义EIP配置)
    * [3.3不分配EIP，即VPC内BLB](#3.3不分配EIP，即VPC内BLB)
    * [3.4UDP-Service](#3.4UDP-Service)
    * [3.5BLB自动分配VIP](#3.5BLB自动分配VIP)

# 一、使用说明
本文档会详细介绍如何在CCE下创建类型是**LoadBalancer**的Service。  
Kubernetes官方教程：[Services](https://kubernetes.io/docs/concepts/services-networking/service/)

# 二、快速开始
当用户创建类型是**LoadBalancer**的Service，默认情况下，CCE会联动的创建BLB，并为此BLB绑定EIP。而当用户删除此Service时，CCE也会联动的删除BLB和EIP。  
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

# 三、高级配置

## 3.1固定EIP
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

## 3.2自定义EIP配置
### EIP支持配置类型
#### 预付费（Prepaid）
| 项目 | 限制 |
| ------ | ------ |
| 公网带宽 | 1-200Mbps，Int |
| 购买时长 | [1,2,3,4,5,6,7,8,9,12,24,36]，时间单位，month |

#### 后付费（Postpaid）
| 计费方式 | 公网带宽 | 费用举例 |
| ------ | ------ | ------ |
| 按使用流量计费（ByTraffic） | 1~1000Mbps，Int | 配置费用：￥0.00032/分钟；流量费用：￥0.76/GB |
| 按使用带宽计费（ByBandwidth） | 1-200Mbps，Int | 配置费用（1Mbps为例）：￥0.00094/分钟 |

### 使用方式
在创建Service时设置相应Annotation如下：
```
// 付费方式，默认：Postpaid；可选：Postpaid、Prepaid
service.beta.kubernetes.io/cce-elastic-ip-payment-timing:"Postpaid"
// 计费方式，默认：ByTraffic；可选：ByTraffic、ByBandwidth
service.beta.kubernetes.io/cce-elastic-ip-billing-method:"ByTraffic"
// 公网带宽，单位为Mbps，默认：1000或者200；对于prepay以及bandwidth类型的EIP，限制为为1~200之间的整数，对于traffic类型的EIP，限制为1~1000之前的整数。
service.beta.kubernetes.io/cce-elastic-ip-bandwidth-in-mbps:"1000"
// 对于预付费，必须设置时长，[1,2,3,4,5,6,7,8,9,12,24,36]，单位月；对于后付费，此设置无效
service.beta.kubernetes.io/cce-elastic-ip-reservation-length:"36"
```

后付费举例：
```
kind: Service
apiVersion: v1
metadata:
  name: nginx-service
  annotations:
    service.beta.kubernetes.io/cce-elastic-ip-payment-timing: "Postpaid"
    service.beta.kubernetes.io/cce-elastic-ip-billing-method: "ByTraffic"
    service.beta.kubernetes.io/cce-elastic-ip-bandwidth-in-mbps: "200"
spec:
  selector:
    app: nginx
  type: LoadBalancer
  ports:
  - name: http
    port: 80
    targetPort: 80
```

预付费举例（请先确定钱够，否则会失败）：
```
kind: Service
apiVersion: v1
metadata:
  name: nginx-service
  annotations:
    service.beta.kubernetes.io/cce-elastic-ip-payment-timing: "Prepaid"
    service.beta.kubernetes.io/cce-elastic-ip-bandwidth-in-mbps: "10"
    service.beta.kubernetes.io/cce-elastic-ip-reservation-length:"1"
spec:
  selector:
    app: nginx
  type: LoadBalancer
  ports:
  - name: http
    port: 80
    targetPort: 80
```

### 说明事项
#### 默认配置
默认为：后付费+按流量+1000M带宽。

#### 用用户已有的EIP，即设置loadBalancerIP的情况
不支持再次配置用户提供的已有EIP，如需修改需要用户自行在console修改。

#### 用户更新Service EIP的配置（即手动编辑annotation）
支持更新的配置有：公网带宽

#### 预付费  
（1）对于预付费，由于EIP API的限制，目前不支持自动续费，需要用户自行到console上续费。
（2）对于预付费，不需要设置计费方式
（3）删除Service时，预付费EIP不会释放，到期后才会释放

## 3.3不分配EIP，即VPC内BLB
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

## 3.4UDP-Service
修改spec.ports.protocol为UDP即可使用UDP Service的功能，如下所示：
```yaml
---
apiVersion: v1
kind: Service
metadata:
  name: udp-server-demo-svc
  labels:
    app: udp-server-demo
spec:
  type: LoadBalancer
  ports:
  - name: udp-server-demo-port
    port: 3005
    targetPort: 3005
    protocol: UDP
  selector:
    app: udp-server-demo
---
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: udp-server-demo
  labels:
    app: udp-server-demo
spec:
  replicas: 1
  selector:
    matchLabels:
      app: udp-server-demo
  template:
    metadata:
      labels:
        app: udp-server-demo
    spec:
      containers:
      - name: udp-server-demo
        image: hub.baidubce.com/jpaas-public/udp-server-demo:latest
        ports:
        - containerPort: 3005
          protocol: UDP
```
（1）部署udp测试服务
```
$ kubectl apply -f udp.yaml
```
（2）UDP Service创建成功
```
$ kubectl get svc
NAME                  TYPE           CLUSTER-IP       EXTERNAL-IP      PORT(S)          AGE
kubernetes            ClusterIP      172.16.0.1       <none>           443/TCP          6h
udp-server-demo-svc   LoadBalancer   172.16.122.139   10.10.10.10      3005:31441/UDP   1m
```
（3）查看服务日志
```
$ kubectl logs -f udp-server-demo-6fdf5d796f-h6595
Received: HealthCheck
Get Health Check, response OK
Received: HealthCheck
Get Health Check, response OK
Received: HealthCheck
Get Health Check, response OK
```
**注：根据百度云负载均衡BLB的要求，对于监听UDP的服务，一定要通过UDP健康检查，BLB才会把流量转发到后端，所以需要用户的后端UDP服务响应健康检查字符串，详见：[UDP健康检查介绍](https://cloud.baidu.com/doc/BLB/GettingStart.html#.E9.85.8D.E7.BD.AEUDP.E7.9B.91.E5.90.AC.E5.99.A8)**

## 3.5BLB自动分配VIP
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
