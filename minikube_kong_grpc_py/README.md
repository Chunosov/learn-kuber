# Run simple Python gRPC service behind kong-proxy on minikube

Run minikube cluster and [install kong](https://github.com/Kong/kubernetes-ingress-controller/blob/master/docs/deployment/minikube.md):

```bash
minikube start --driver=virtualbox
kubectl create -f https://bit.ly/k4k8s

export PROXY_HOST=$(minikube ip)
export PROXY_PORT=$(kubectl get services/kong-proxy -n kong -o go-template='{{(index .spec.ports 0).nodePort}}')
export PROXY_ADDR=$PROXY_HOST:$PROXY_PORT
echo $PROXY_ADDR
```

Connect to cluster's docker and build simple python gRPC services image:

```bash
eval $(minikube -p minikube docker-env)
docker build -t grpc_greeter:v0 ../minikube_knative_grpc_py
```

Deploy a service, make an ingress rule, adjust service and ingress to be grpc instead of http:

```bash
kubectl apply -f service.yaml
kubectl patch ingress demo-grpc-py -p '{"metadata":{"annotations":{"konghq.com/protocols":"grpc"}}}'
kubectl patch service grpc-greeter -p '{"metadata":{"annotations":{"konghq.com/protocols":"grpc"}}}'
```

Check if all applied correctly:

```bash
kubectl describe ingress demo-grpc-py
Name:             demo-grpc-py
Namespace:        default
Address:
Default backend:  default-http-backend:80 (<error: endpoints "default-http-backend" not found>)
Rules:
  Host        Path  Backends
  ----        ----  --------
  *
              /   grpc-greeter:50051 (172.17.0.5:50051)
Annotations:  konghq.com/protocols: grpc
Events:       <none>

kubectl describe service grpc-greeter
Name:              grpc-greeter
Namespace:         default
Labels:            app=grpc-greeter
Annotations:       konghq.com/protocols: grpc
Selector:          app=grpc-greeter
Type:              ClusterIP
IP:                10.98.223.216
Port:              grpc  50051/TCP
TargetPort:        50051/TCP
Endpoints:         172.17.0.5:50051
Session Affinity:  None
Events:            <none>
```

Test with curl:

```bash
curl $PROXY_ADDR
{"message":"Non-gRPC request matched gRPC route"}
```

So we have a route, the route is of gRPC type.

Test with grpcurl (make sure you have it [installed](../README.md#grpcurl)):

```bash
grpcurl -v -plaintext -d '{"name":"TEST"}' $PROXY_ADDR helloworld.Greeter/SayHello
Failed to dial target host "192.168.39.161:30312": context deadline exceeded

grpcurl -v -insecure -d '{"name":"TEST"}' $PROXY_ADDR helloworld.Greeter/SayHello
Failed to dial target host "192.168.39.161:30312": tls: first record does not look like a TLS handshake
```

<span style="color:red"><b>It doesn't work.</b> The same error as in the [previous example](../minikube_kong_grpc/README.md).</span>
