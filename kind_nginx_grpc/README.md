# Start sample gRPC service with nginx ingress controller on kind cluster

https://kubernetes.github.io/ingress-nginx/examples/grpc/
https://github.com/kubernetes/ingress-nginx/tree/master/docs/examples/grpc
https://medium.com/@Alibaba_Cloud/accessing-grpc-services-through-container-service-for-kubernetes-ingress-controller-511f984071e1
https://github.com/kubernetes/ingress-nginx/issues/5886

## Start

Start a cluster, and install nginx ingess controller, and wait while it is started:

```bash
kind create cluster --config cluster.yaml

kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/master/deploy/static/provider/kind/deploy.yaml

kubectl wait --namespace ingress-nginx \
  --for=condition=ready pod \
  --selector=app.kubernetes.io/component=controller \
  --timeout=120s
```

## Prepare TLS

Generate TLS certificate, and make a secret, and deploy a service:

```bash
openssl req -x509 -sha256 -nodes -days 365 -newkey rsa:2048 -keyout tls.key -out tls.crt -subj "/CN=api.bar/O=api.bar"
```

If it shows an error like:

```
Can't load /home/nikolay/.rnd into RNG
140664075555264:error:2406F079:random number generator:RAND_load_file:Cannot open file:../crypto/rand/randfile.c:88:Filename=/home/nikolay/.rnd
```

then edit the config `/etc/ssl/openssl.cnf` and comment line `RANDFILE		= $ENV::HOME/.rnd` ([from here](https://github.com/openssl/openssl/issues/7754)).

Add the certificate as a secret to the cluster:

```bash
kubectl create secret tls tls-secret --key tls.key --cert tls.crt
```

## Deploy

```bash
kubectl apply -f service.yaml
```

## Check as NodePort

Check if the service is available as NodePort:

```bash
kubectl get service
NAME         TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)                         AGE
grpcbin      NodePort    10.111.126.18   <none>        9000:31002/TCP,9001:31149/TCP   11s
kubernetes   ClusterIP   10.96.0.1       <none>        443/TCP                         8m7s

kubectl get node -o wide
NAME                 STATUS   ROLES    AGE     VERSION   INTERNAL-IP   EXTERNAL-IP   OS-IMAGE       KERNEL-VERSION       CONTAINER-RUNTIME
kind-control-plane   Ready    master   8m34s   v1.18.2   172.19.0.2    <none>        Ubuntu 19.10   4.15.0-112-generic   containerd://1.3.3-14-g449e9269

curl http://172.19.0.2:31002
Warning: Binary output can mess up your terminal. Use "--output -" to tell
Warning: curl to output it to your terminal anyway, or consider "--output
Warning: <FILE>" to save to a file.

curl -k https://172.19.0.2:31149
invalid gRPC request method
```

Seems it works but curl obviously doesn't understand the grpc answer.

Check it with [grpcurl](../README.md#grpcul):

```bash
grpcurl -v -d '{"greeting": "TEST"}' -plaintext 172.19.0.2:31002 hello.HelloService.SayHello

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
grpcurl -v -d '{"greeting": "TEST"}' -insecure 172.19.0.2:31149 hello.HelloService.SayHello

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

It works as NodePort.

## Check via Ingress

Check if the service is available through ingress controller.

### Very basic ingress rule

The only required thing for proxying gRPC is to add one proper annotation to a common HTTP ingress rule:

```yaml
nginx.ingress.kubernetes.io/backend-protocol: "GRPC"
```

```bash
kubectl apply -f ingress_0.yaml

grpcurl -v -d '{"greeting": "TEST"}' -plaintext localhost:80 hello.HelloService.SayHello
Failed to dial target host "localhost:80": context deadline exceeded

grpcurl -v -d '{"greeting": "TEST"}' -insecure localhost:443 hello.HelloService.SayHello

Resolved method descriptor:
rpc SayHello ( .hello.HelloRequest ) returns ( .hello.HelloResponse );

Request metadata to send:
(empty)

Response headers received:
content-type: application/grpc
date: Tue, 15 Sep 2020 17:31:25 GMT
server: nginx/1.19.2
strict-transport-security: max-age=15724800; includeSubDomains

Response contents:
{
  "reply": "hello TEST"
}

Response trailers received:
(empty)
Sent 1 request and received 1 response
```

It seems gRPC proxy is only working on HTTPS port but not on HTTP. But in the target service we use the `insecure` port, because TLS is terminated at the ingress level and inside of cluster gGRPC traffic travells unencripted.

### Ingress rule with host

Now we have to use the `-authority` parameter of grpcurl:

```bash
kubectl delete ingress --all
kubectl apply -f ingress_1.yaml

grpcurl -v -d '{"greeting": "TEST"}' -plaintext -authority myexample.com localhost:80 hello.HelloService.SayHello
Failed to dial target host "localhost:80": context deadline exceeded

grpcurl -v -d '{"greeting": "TEST"}' -insecure -authority myexample.com localhost:443 hello.HelloService.SayHello

Resolved method descriptor:
rpc SayHello ( .hello.HelloRequest ) returns ( .hello.HelloResponse );

Request metadata to send:
(empty)

Response headers received:
content-type: application/grpc
date: Tue, 15 Sep 2020 17:36:38 GMT
server: nginx/1.19.2
strict-transport-security: max-age=15724800; includeSubDomains

Response contents:
{
  "reply": "hello TEST"
}

Response trailers received:
(empty)
Sent 1 request and received 1 response
```

### Use custom certificate

It's used in [this example](https://github.com/kubernetes/ingress-nginx/tree/master/docs/examples/grpc) from nginx (or [this link](https://kubernetes.github.io/ingress-nginx/examples/grpc/)), so just try it for the full picture.

Delete all existed ingresses:

```bash
kubectl delete ingress --all
```

Then generate a certificate and make a secret, it should be done before making new ingresses:

```bash
openssl req -x509 -sha256 -nodes -days 365 -newkey rsa:2048 -keyout tls.key -out tls.crt -subj "/CN=myexample.com/O=myexample.com"
```

<details>
  <summary>Error "Cannot open file .rnd"</summary>

If it shows an error like:

```
Can't load /home/nikolay/.rnd into RNG
140664075555264:error:2406F079:random number generator:RAND_load_file:Cannot open file:../crypto/rand/randfile.c:88:Filename=/home/nikolay/.rnd
```

then edit the config `/etc/ssl/openssl.cnf` and comment line `RANDFILE		= $ENV::HOME/.rnd` ([solution from here](https://github.com/openssl/openssl/issues/7754)).

---
</details>

&nbsp;

Add the certificate as a secret to the cluster:

```bash
kubectl create secret tls tls-secret --key tls.key --cert tls.crt
```

Then we can make an ingress using this certificate:

```bash
kubectl apply -f ingress_2.yaml

grpcurl -v -d '{"greeting": "TEST"}' -plaintext -authority myexample.com localhost:80 hello.HelloService.SayHello
Failed to dial target host "localhost:80": context deadline exceeded

grpcurl -v -d '{"greeting": "TEST"}' -insecure -authority myexample.com localhost:443 hello.HelloService.SayHello

Resolved method descriptor:
rpc SayHello ( .hello.HelloRequest ) returns ( .hello.HelloResponse );

Request metadata to send:
(empty)

Response headers received:
content-type: application/grpc
date: Tue, 15 Sep 2020 17:42:21 GMT
server: nginx/1.19.2
strict-transport-security: max-age=15724800; includeSubDomains

Response contents:
{
  "reply": "hello TEST"
}

Response trailers received:
(empty)
Sent 1 request and received 1 response
```

Output is the same. I'm not sure if the certificate is used or not, at least it doesn't break anything.

### Use magic domain name

Delete existed ingresses and secrets:

```bash
kubectl delete ingress --all
kubectl delete secret tls-secret
```

Then recreate and test ingress (IP-address 127.0.0.1 must be used, not "localhost"):

```bash
kubectl apply -f ingress_3.yaml

grpcurl -v -d '{"greeting": "TEST"}' -plaintext 127.0.0.1.xip.io:80 hello.HelloService.SayHello
Failed to dial target host "127.0.0.1.xip.io:80": context deadline exceeded
nikolay@leo:~/Projects/learn-kuber/kind_nginx_grpc$ grpcurl -v -d '{"greeting": "TEST"}' -insecure 127.0.0.1.xip.io:443 hello.HelloService.SayHello

Resolved method descriptor:
rpc SayHello ( .hello.HelloRequest ) returns ( .hello.HelloResponse );

Request metadata to send:
(empty)

Response headers received:
content-type: application/grpc
date: Tue, 15 Sep 2020 17:51:56 GMT
server: nginx/1.19.2
strict-transport-security: max-age=15724800; includeSubDomains

Response contents:
{
  "reply": "hello TEST"
}

Response trailers received:
(empty)
Sent 1 request and received 1 response
```
