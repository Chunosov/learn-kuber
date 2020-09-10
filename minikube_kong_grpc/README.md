# Run example gRPC service behind kong-proxy on minikube

This is a reproduction of official [gRPC example](https://github.com/Kong/kubernetes-ingress-controller/blob/main/docs/guides/using-ingress-with-grpc.md) from kong.

## Start cluster

Run minikube cluster and [install kong](https://github.com/Kong/kubernetes-ingress-controller/blob/master/docs/deployment/minikube.md):

```bash
minikube start --driver=virtualbox
kubectl create -f https://bit.ly/k4k8s

export PROXY_IP=$(minikube service -n kong kong-proxy --url | head -1)
echo $PROXY_IP
http://192.168.99.108:30552
```

## Deploy service

Deploy sample gRPC servce:

```bash
kubectl apply -f service.yaml
# or
kubectl apply -f https://bit.ly/grpcbin-service

kubectl get service grpcbin
NAME      TYPE        CLUSTER-IP     EXTERNAL-IP   PORT(S)    AGE
grpcbin   ClusterIP   10.97.115.99   <none>        9001/TCP   12s
```

## Access service

Make an ingress rule accessing grpcbin service outside of the cluster:

```bash
kubectl apply -f ingress.yaml

kubectl describe ingress demo-grpc
Name:             demo-grpc
Namespace:        default
Address:
Default backend:  default-http-backend:80 (<error: endpoints "default-http-backend" not found>)
Rules:
  Host        Path  Backends
  ----        ----  --------
  *
              /        grpcbin:9001 (172.17.0.5:9001)
Annotations:  Events:  <none>
```

### Use curl

Now the service can be accessed outside of the cluster (though it doesn't work as expected):

```bash
curl $PROXY_IP
Client sent an HTTP request to an HTTPS server.

# the same inside of cluster
minikube ssh
curl 172.17.0.5:9001
Client sent an HTTP request to an HTTPS server.
```

### Use grpcurl

Make sure you have grpcurl [installed](../README.md#grpcurl).

```bash
grpcurl -v -d '{"greeting": "TEST"}' -plaintext $(minikube ip):30552 hello.HelloService.SayHello
Failed to dial target host "192.168.99.108:30552": context deadline exceeded
```

By default, kong assumes all routes and services are HTTP or HTTPs. We must mark them as gRPC:

```bash
kubectl patch ingress demo-grpc -p '{"metadata":{"annotations":{"konghq.com/protocols":"grpc"}}}'
kubectl patch service grpcbin -p '{"metadata":{"annotations":{"konghq.com/protocols":"grpc"}}}'
```

<span style="color:red"><b>It doesn't work.</b> It says "context deadline exceeded". Posted a [question](https://discuss.konghq.com/t/does-grpc-proxy-works-under-minikube/7092) in the official discussion channel.</span>

---

Direct steps to reproduce (for the question):

```bash
minikube start --driver=kvm2
kubectl create -f https://bit.ly/k4k8s
export PROXY_IP=$(minikube ip)
export PROXY_PORT=$(kubectl get services/kong-proxy -n kong -o go-template='{{(index .spec.ports 0).nodePort}}')
echo $PROXY_IP:$PROXY_PORT
kubectl apply -f https://bit.ly/grpcbin-service
echo "apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: demo
spec:
  rules:
  - http:
      paths:
      - path: /
        backend:
          serviceName: grpcbin
          servicePort: 9001" | kubectl apply -f -
kubectl patch ingress demo -p '{"metadata":{"annotations":{"konghq.com/protocols":"grpc,grpcs"}}}'
kubectl patch svc grpcbin -p '{"metadata":{"annotations":{"konghq.com/protocols":"grpcs"}}}'
grpcurl -v -d '{"greeting": "Kong Hello world!"}' -insecure $PROXY_IP:$PROXY_PORT hello.HelloService.SayHello

# another try without grpcs:
kubectl patch ingress demo -p '{"metadata":{"annotations":{"konghq.com/protocols":"grpc"}}}'
kubectl patch svc grpcbin -p '{"metadata":{"annotations":{"konghq.com/protocols":"grpc"}}}'
grpcurl -v -d '{"greeting": "Kong Hello world!"}' -insecure $PROXY_IP:$PROXY_PORT hello.HelloService.SayHello
```
