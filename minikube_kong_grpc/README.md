# Run example gRPC service behind kong-proxy on minikube

This is a reproduction of official [gRPC example](https://github.com/Kong/kubernetes-ingress-controller/blob/main/docs/guides/using-ingress-with-grpc.md) from kong.

## Start cluster

Run minikube cluster and [install kong](https://github.com/Kong/kubernetes-ingress-controller/blob/master/docs/deployment/minikube.md):

```bash
minikube start --driver=virtualbox
kubectl create -f https://bit.ly/k4k8s

kubectl get service kong-proxy -n kong
NAME         TYPE           CLUSTER-IP     EXTERNAL-IP   PORT(S)                      AGE
kong-proxy   LoadBalancer   10.99.52.137   <pending>     80:30493/TCP,443:30363/TCP   21m

export PORT_HTTP=30493
export PORT_HTTPS=30363
```

## Deploy service

Deploy sample gRPC servce:

```bash
kubectl apply -f service_0.yaml

kubectl get service grpcbin
NAME      TYPE       CLUSTER-IP       EXTERNAL-IP   PORT(S)                         AGE
grpcbin   NodePort   10.108.224.113   <none>        9000:31936/TCP,9001:32009/TCP   19s

export SERVICE_PORT_HTTP=31936
export SERVICE_PORT_HTTPS=32009
```

## Try as NodePort

```bash
grpcurl -v -d '{"greeting": "TEST"}' -plaintext $(minikube ip):$SERVICE_PORT_HTTP hello.HelloService.SayHello

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

It works on both ports HTTP and HTTPS.

### Try invalid case

```bash
grpcurl -v -d '{"greeting": "TEST"}' -plaintext $(minikube ip):$SERVICE_PORT_HTTPS hello.HelloService.SayHello
Failed to dial target host "192.168.99.113:32009": context deadline exceeded

grpcurl -v -d '{"greeting": "TEST"}' -insecure $(minikube ip):$SERVICE_PORT_HTTP hello.HelloService.SayHello
Failed to dial target host "192.168.99.113:31936": tls: first record does not look like a TLS handshake
```

The "context deadline exceeded" error happens when we are trying to access HTTPS port with `-plaintext`. When we are trying to access HTTP port with `-insecure`, expected TLS records are not found.

## Try with ingress

Make an ingress rule accessing grpcbin service outside of the cluster. By default, kong assumes all routes and services are HTTP or HTTPs. We must mark them as gRPC via an annotation:

```yaml
konghq.com/protocols: grpc
```

Make sure you have grpcurl [installed](../README.md#grpcurl).

### Very basic ingress rule

```bash
kubectl apply -f ingress_1.yaml

kubectl describe ingress grpcbin
Warning: extensions/v1beta1 Ingress is deprecated in v1.14+, unavailable in v1.22+; use networking.k8s.io/v1 Ingress
Name:             grpcbin
Namespace:        default
Address:
Default backend:  default-http-backend:80 (<error: endpoints "default-http-backend" not found>)
Rules:
  Host        Path  Backends
  ----        ----  --------
  *
              /   grpcbin:insecure   172.17.0.4:9000)
Annotations:  konghq.com/protocols: grpc
Events:       <none>
```

Try to connect:

```bash
curl http://$(minikube ip):$PORT_HTTP
{"message":"no Route matched with those values"}

curl -k https://$(minikube ip):$PORT_HTTPS
{"message":"no Route matched with those values"}

grpcurl -v -d '{"greeting": "TEST"}' -plaintext $(minikube ip):$PORT_HTTP hello.HelloService.SayHello
Failed to dial target host "192.168.99.113:30493": context deadline exceeded

grpcurl -v -d '{"greeting": "TEST"}' -insecure $(minikube ip):$PORT_HTTPS hello.HelloService.SayHello
Error invoking method "hello.HelloService.SayHello": failed to query for service descriptor "hello.HelloService": server does not support the reflection API

grpcurl -proto hello.proto -v -d '{"greeting": "TEST"}' -insecure $(minikube ip):$PORT_HTTPS hello.HelloService.SayHello

Resolved method descriptor:
rpc SayHello ( .hello.HelloRequest ) returns ( .hello.HelloResponse );

Request metadata to send:
(empty)

Response headers received:
(empty)

Response trailers received:
(empty)
Sent 1 request and received 0 responses
ERROR:
  Code: Unimplemented
  Message: Not Found: HTTP status code 404; transport: missing content-type field
```

### Ingress via host

```bash
kubectl delete ingress --all
kubectl apply -f ingress_2.yaml

grpcurl -v -d '{"greeting": "TEST"}' -authority myexample.com -plaintext $(minikube ip):30493 hello.HelloService.SayHello
Failed to dial target host "192.168.99.113:30493": context deadline exceeded

grpcurl -v -d '{"greeting": "TEST"}' -authority myexample.com -insecure $(minikube ip):30363 hello.HelloService.SayHello
Error invoking method "hello.HelloService.SayHello": failed to query for service descriptor "hello.HelloService": server does not support the reflection API

grpcurl -proto hello.proto -v -d '{"greeting": "TEST"}' -authority myexample.com -insecure $(minikube ip):30363 hello.HelloService.SayHello

Resolved method descriptor:
rpc SayHello ( .hello.HelloRequest ) returns ( .hello.HelloResponse );

Request metadata to send:
(empty)

Response headers received:
(empty)

Response trailers received:
(empty)
Sent 1 request and received 0 responses
ERROR:
  Code: Unimplemented
  Message: Not Found: HTTP status code 404; transport: missing content-type field
```

<span style="color:red"><b>It doesn't work.</b> Posted a [question](https://discuss.konghq.com/t/does-grpc-proxy-works-under-minikube/7092) in the official discussion channel.</span>

<details>
  <summary>Direct steps to reproduce (for the question)</summary>

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
</details>

