# Run simple Go gRPC service behind kong-proxy on minikube

This example uses the same [code](https://github.com/evankanderson/sia) as knative guys have used to test their [gRPC PR](https://github.com/knative/serving/pull/2539#issuecomment-459148556).

## Start

Run minikube cluster and [install kong](https://github.com/Kong/kubernetes-ingress-controller/blob/master/docs/deployment/minikube.md):

```bash
minikube start --driver=virtualbox
kubectl create -f https://bit.ly/k4k8s

export PROXY_HOST=$(minikube ip)
export PROXY_PORT=$(kubectl get services/kong-proxy -n kong -o go-template='{{(index .spec.ports 0).nodePort}}')
export PROXY_ADDR=$PROXY_HOST:$PROXY_PORT
echo $PROXY_ADDR
```

Connect to cluster's docker and build an image:

```bash
eval $(minikube -p minikube docker-env)
docker build -t grpc_sia:v0 .
```

## Test image in cluster

Test that image is ok by accessing container from inside of the cluster:

```bash
docker run --rm -p 5000:8080 --name grpc_sia grpc_sia:v0

# in another terminal
minikube ssh

curl -d "test" localhost:5000
You POST for "/"

Got: test
```

Server is up and can answer something.

## Deploy service

Our servce is of NodePort type to make is testable from outside of the cluster:

```bash
kubectl apply -f service.yaml

kubectl get service grpc-sia
NAME       TYPE       CLUSTER-IP      EXTERNAL-IP   PORT(S)          AGE
grpc-sia   NodePort   10.108.21.143   <none>        8080:30625/TCP   12s
```

## Test as plain kubernetes service

Make shure yout have grpcurl [installed](../README.md#grpcurl). Call the service from outside of the cluster:

```bash
curl -d "test" $(minikube ip):30625
You POST for "/"

Got: test

grpcurl -plaintext -proto doer/doer.proto -format json -d '{"thing":"SOMETHING"}' $(minikube ip):30625 doer.Doer/DoIt
{
  "words": "Did: SOMETHING"
}
```

It works.

## Test through kong-proxy

kong proxies http request because the route is http (by default):

```bash
curl -d "TEST" $PROXY_ADDR
You POST for "/"

Got: TEST
```

But it can't proxy grpc requests:

```bash
grpcurl -plaintext -proto doer/doer.proto -format json -d '{"thing":"SOMETHING"}' $PROXY_ADDR doer.Doer/DoIt
Failed to dial target host "192.168.39.161:30312": context deadline exceeded
```

Tell kong that our ingress rule and service is grpc:

```bash
kubectl patch ingress demo-grpc-go -p '{"metadata":{"annotations":{"konghq.com/protocols":"grpc"}}}'
kubectl patch service grpc-sia -p '{"metadata":{"annotations":{"konghq.com/protocols":"grpc"}}}'

kubectl describe ingress demo-grpc-go
Name:             demo-grpc-go
Namespace:        default
Address:
Default backend:  default-http-backend:80 (<error: endpoints "default-http-backend" not found>)
Rules:
  Host        Path  Backends
  ----        ----  --------
  *
              /   grpc-sia:8080 (172.17.0.6:8080)
Annotations:  konghq.com/protocols: grpc
Events:       <none>

kubectl describe service grpc-sia
Name:                     grpc-sia
Namespace:                default
Labels:                   app=grpc-sia
Annotations:              konghq.com/protocols: grpc
Selector:                 app=grpc-sia
Type:                     NodePort
IP:                       10.108.21.143
Port:                     grpc  8080/TCP
TargetPort:               8080/TCP
NodePort:                 grpc  30625/TCP
Endpoints:                172.17.0.6:8080
Session Affinity:         None
External Traffic Policy:  Cluster
Events:                   <none>
```

The route is grpc now, kong refuses to proxy http request:

```bash
curl -d "TEST" $PROXY_ADDR
{"message":"Non-gRPC request matched gRPC route"}
```

<span style="color:red"><b>gRPC still doesn't work.</b> The same error as in the [previous example](../minikube_kong_grpc/README.md):</span>

```bash
grpcurl -plaintext -proto doer/doer.proto -format json -d '{"thing":"SOMETHING"}' $PROXY_ADDR doer.Doer/DoIt
Failed to dial target host "192.168.39.161:30312": context deadline exceeded
```

### Try to use KongIngess with service

```bash
echo "
apiVersion: configuration.konghq.com/v1
kind: KongIngress
metadata:
  name: grpc-upstream
proxy:
  protocol: grpc
" | kubectl apply -f -

kubectl patch service grpc-sia -p '{"metadata":{"annotations":{"configuration.konghq.com":"grpc-upstream"}}}'
kubectl patch ingress demo-grpc-go -p '{"metadata":{"annotations":{"configuration.konghq.com":"grpc-upstream"}}}'
```