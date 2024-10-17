apiVersion: v1
kind: Pod
metadata:
  name: henry-pod
  labels:
    name: Henry
    location: Melbourne
  annotations:
    description: "This pod is named Henry and located in Melbourne."
    volunteer: "Henry is a volunteer at K8SUG (Kubernetes User Group)"
spec:
  containers:
    - name: nginx-container
      image: nginx:latest
      ports:
        - containerPort: 80
