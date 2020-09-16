# Run gRPC service with kong ingress controller in kind cluster

Start a cluster:

```bash
kind create cluster --config cluster.yaml

kubectl get node -owide
NAME                 STATUS   ROLES    AGE   VERSION   INTERNAL-IP   EXTERNAL-IP   OS-IMAGE       KERNEL-VERSION       CONTAINER-RUNTIME
kind-control-plane   Ready    master   18m   v1.18.2   172.19.0.2    <none>        Ubuntu 19.10   4.15.0-112-generic   containerd://1.3.3-14-g449e9269

export NODE_ADDR=172.19.0.2
```

Install kong ingress controller:

```bash
kubectl apply --filename https://raw.githubusercontent.com/Kong/kubernetes-ingress-controller/0.10.x/deploy/single/all-in-one-dbless.yaml

kubectl get service -n kong
NAME                      TYPE           CLUSTER-IP       EXTERNAL-IP   PORT(S)                      AGE
kong-proxy                LoadBalancer   10.100.50.116    <pending>     80:32581/TCP,443:31682/TCP   9s
kong-validation-webhook   ClusterIP      10.107.249.147   <none>        443/TCP                      9s

export PROXY_PORT_HTTP=32581
export PROXY_PORT_HTTPS=31682
```

Deploy s service:

```bash
kubectl apply -f service.yaml
```

## Try if all work via NodePort

Make shure you have grpcurl [installed](../README.md#grpcurl).

```bash
kubectl get service
kubectl get service
NAME         TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)                         AGE
grpcbin      NodePort    10.104.203.141   <none>        9000:32388/TCP,9001:32216/TCP   3m39s


export SERVICE_PORT_HTTP=32388
export SERVICE_PORT_HTTPS=32216
```

```bash
grpcurl -v -d '{"greeting": "TEST"}' -plaintext $NODE_ADDR:$SERVICE_PORT_HTTP hello.HelloService.SayHello

Resolved method descriptor:
rpc SayHello ( .hello.HelloRequest ) returns ( .hello.HelloResponse );

Request metadata to send:
(empty)

Response headers received:
content-type: application/grpc

Response contents:
{
  "reply": "hello TEST"
}

Response trailers received:
(empty)
Sent 1 request and received 1 response
```

```bash
grpcurl -v -d '{"greeting": "TEST"}' -insecure $NODE_ADDR:$SERVICE_PORT_HTTPS hello.HelloService.SayHello

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

## Try via ingress

From [here we know](../kind_kong_ingress/README.md) that kong exposes services at kong-proxy NodePort, not on cluster address as nginx does.

```bash
kubectl apply -f ingress_0.yaml

grpcurl -v -d '{"greeting": "TEST"}' -plaintext $NODE_ADDR:$PROXY_PORT_HTTP hello.HelloService.SayHello
Failed to dial target host "172.19.0.2:32581": context deadline exceeded

grpcurl -v -d '{"greeting": "TEST"}' -insecure $NODE_ADDR:$PROXY_PORT_HTTPS hello.HelloService.SayHello
# it hangs forever
```

Do this, because it is in the official (not working) [grpc example](https://github.com/Kong/kubernetes-ingress-controller/blob/main/docs/guides/using-ingress-with-grpc.md) from kong, but it changes nothing:

```bash
kubectl patch svc grpcbin -p '{"metadata":{"annotations":{"konghq.com/protocols":"grpc"}}}'
```

<span style="color:red"><b>It doesn't work.</b> But it doesn't work differently than on [minikube cluster](../minikube_kong_grpc/README.md).</span>
