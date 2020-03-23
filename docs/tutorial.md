# Tutorial

This example will show you how to use the CCM to create a service of `type:LoadBalancer`.

### Load Balancer example

When you create a service with `type: LoadBalancer` an CCE load balancer will
be created.

The example below will create an NGINX deployment and expose it via a load
balancer. Note that the service **type** is set to **LoadBalancer**.

```yaml
---
apiVersion: apps/v1beta1
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  replicas: 2
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
  - name: http
    port: 80
    targetPort: 80
```

Create it

```bash
$ kubectl create -f examples/nginx-demo-svc.yaml
```

Watch the service and await a public IP address. This will be the load balancer
IP which you can use to connect to your service.

```bash
$ kubectl get svc --watch
NAME            CLUSTER-IP     EXTERNAL-IP      PORT(S)        AGE
nginx-service   1.1.1.1        8.8.8.8          80:30274/TCP   5m
```

You can now access your service via the provisioned load balancer

```bash
$ curl -i http://8.8.8.8
```