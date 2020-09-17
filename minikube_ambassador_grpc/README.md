# Run gRPC service with Ambassador ingress in minikube cluster

Start minikube cluster, install [Ambassador](https://www.getambassador.io/docs/latest/tutorials/getting-started/), and wait while it gets started:

```bash
minikube start --driver virtualbox

kubectl apply -f https://www.getambassador.io/yaml/aes-crds.yaml
kubectl wait --for condition=established --timeout=90s crd -lproduct=aes
kubectl apply -f https://www.getambassador.io/yaml/aes.yaml
kubectl -n ambassador wait --for condition=available --timeout=90s deploy -lproduct=aes

export PROXY_PORT_HTTP=$(kubectl get service ambassador -n ambassador --output jsonpath="{.spec.ports[0].nodePort}")
export PROXY_PORT_HTTPS=$(kubectl get service ambassador -n ambassador --output jsonpath="{.spec.ports[1].nodePort}")
echo $PROXY_PORT_HTTP,$PROXY_PORT_HTTPS
30942,32299
```

Deploy simple echo services:

```bash
kubectl apply -f service.yaml

export SERVICE_PORT_HTTP=$(kubectl get service grpcbin --output jsonpath="{.spec.ports[0].nodePort}")
export SERVICE_PORT_HTTPS=$(kubectl get service grpcbin --output jsonpath="{.spec.ports[1].nodePort}")
echo $SERVICE_PORT_HTTP,$SERVICE_PORT_HTTPS
31577,31632
```

## Test via NodePort

Make sure they work as NodePort (check that you have grpcurl [installed](../README.md#grpcurl)):

```bash
grpcurl -v -d '{"greeting": "TEST"}' -plaintext $(minikube ip):$SERVICE_PORT_HTTP hello.HelloService.SayHello

Resolved method descriptor:
rpc SayHello ( .hello.HelloRequest ) returns ( .hello.HelloResponse );

Request metadata to send:
(empty)

Response headers received:
content-type: application/grpc
date: Thu, 17 Sep 2020 08:09:37 GMT
server: istio-envoy
x-envoy-decorator-operation: grpcbin.default.svc.cluster.local:9000/*
x-envoy-upstream-service-time: 1

Response contents:
{
  "reply": "hello TEST"
}

Response trailers received:
(empty)
Sent 1 request and received 1 response
```

```bash
grpcurl -v -d '{"greeting": "TEST"}' -insecure $(minikube ip):$SERVICE_PORT_HTTPS hello.HelloService.SayHello

Resolved method descriptor:
rpc SayHello ( .hello.HelloRequest ) returns ( .hello.HelloResponse );

Request metadata to send:
(empty)

Response headers received:
content-type: application/grpc
trailer: Grpc-Status
trailer: Grpc-Message
trailer: Grpc-Status-Details-Bin

Response contents:
{
  "reply": "hello TEST"
}

Response trailers received:
(empty)
Sent 1 request and received 1 response
```

It works on both HTTP and HTTPS ports.

## Test via Ambassador maping

Ambassador mappings are based on URL prefixes; for gRPC, the URL prefix is the full-service name, including the package path (package.service) see [mapping file](./mapping.yaml):

```bash
kubectl apply -f mapping.yaml
```

```bash
grpcurl -v -d '{"greeting": "TEST"}' -plaintext $(minikube ip):$PROXY_PORT_HTTP hello.HelloService.SayHello
Error invoking method "hello.HelloService.SayHello": failed to query for service descriptor "hello.HelloService": rpc error: code = Unknown desc =

grpcurl -v -d '{"greeting": "TEST"}' -insecure $(minikube ip):$PROXY_PORT_HTTPS hello.HelloService.SayHello
Error invoking method "hello.HelloService.SayHello": failed to query for service descriptor "hello.HelloService": server does not support the reflection API
```

The same under `minikube tunnel`:

```bash
grpcurl -v -d '{"greeting": "TEST"}' -plaintext 10.105.13.239:80 hello.HelloService.SayHello
Error invoking method "hello.HelloService.SayHello": failed to query for service descriptor "hello.HelloService": rpc error: code = Unknown desc =

grpcurl -v -d '{"greeting": "TEST"}' -insecure 10.105.13.239:443 hello.HelloService.SayHello
Error invoking method "hello.HelloService.SayHello": failed to query for service descriptor "hello.HelloService": server does not support the reflection API
```

<span style="color:red"><b>It doesn't work.</b></span>
