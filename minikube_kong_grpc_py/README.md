# Run simple Python gRPC service behind kong-proxy on minikube

## Start cluster

Run minikube cluster, [install kong](https://github.com/Kong/kubernetes-ingress-controller/blob/master/docs/deployment/minikube.md), and connect to cluster's docker:

```bash
minikube start --driver=virtualbox
kubectl create -f https://bit.ly/k4k8s
eval $(minikube -p minikube docker-env)

export PROXY_IP=$(minikube service -n kong kong-proxy --url | head -1)
echo $PROXY_IP
http://192.168.99.108:30552

export PROXY_HOST=$(minikube ip)
export PROXY_PORT=$(kubectl get services/kong-proxy -n kong -o go-template='{{(index .spec.ports 0).nodePort}}')
export PROXY_ADDR=$PROXY_HOST:$PROXY_PORT
echo $PROXY_ADDR
```

Build simple python gRPC services image, deploy a service, and make an ingress rule:

```bash
docker build -t grpc_greeter:v0 ../minikube_knative_grpc_py
kubectl apply -f service.yaml
```

Adjust service and ingress to be grpc instead of http:

```bash
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

Download [grpcurl](https://github.com/fullstorydev/grpcurl/releases), move it to `/usr/local/bin`, and make executable.

Test with grpcurl:

```bash
grpcurl -v -plaintext -d '{"name":"TEST"}' $PROXY_ADDR helloworld.Greeter/SayHello
Failed to dial target host "192.168.39.161:30312": context deadline exceeded
```

<span style="color:red"><b>It doesn't work.</b> The same error as in the [previous example](../minikube_kong_grpc/README.md).</span>
