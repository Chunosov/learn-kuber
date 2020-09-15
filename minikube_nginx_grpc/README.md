# Run gRPC service with nginx ingress controller

Run sample minikube cluster and enable `ingress` addon (minikube uses [ingress-nginx](https://kubernetes.github.io/ingress-nginx/deploy/#minikube) by default):

```bash
minikube start --driver=virtualbox
minikube addons enable ingress
```

Deploy a service:

```bash
kubectl apply -f service.yaml

kubectl get service
NAME         TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)                         AGE
grpcbin      NodePort    10.108.93.188   <none>        9000:32150/TCP,9001:30547/TCP   43s
kubernetes   ClusterIP   10.96.0.1       <none>        443/TCP                         102m
```

## Try as NodePort

Try with curl

```bash
curl http://$(minikube ip):32150
Warning: Binary output can mess up your terminal. Use "--output -" to tell
Warning: curl to output it to your terminal anyway, or consider "--output
Warning: <FILE>" to save to a file.

curl -k https://$(minikube ip):30547
invalid gRPC request method
```

So it can connect to both HTTP and HTTPS ports but obviously can't understand the answer.

Try with [grpcurl](../README.md#grpcurl):

```bash
grpcurl -plaintext $(minikube ip):32150 list
addsvc.Add
grpc.gateway.examples.examplepb.ABitOfEverythingService
grpc.reflection.v1alpha.ServerReflection
grpcbin.GRPCBin
hello.HelloService

grpcurl -insecure $(minikube ip):30547 list
addsvc.Add
grpc.gateway.examples.examplepb.ABitOfEverythingService
grpc.reflection.v1alpha.ServerReflection
grpcbin.GRPCBin
hello.HelloService
```

Call method on HTTP port:

```bash
grpcurl -v -d '{"greeting": "TEST"}' -plaintext $(minikube ip):32150 hello.HelloService.SayHello

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

Call method on HTTPS port:

```bash
grpcurl -v -d '{"greeting": "TEST"}' -insecure $(minikube ip):30547 hello.HelloService.SayHello

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

## Try via ingress

### Very basic ingress rule

The only required thing for proxying gRPC is to add one proper annotation to a common HTTP ingress rule:

```yaml
nginx.ingress.kubernetes.io/backend-protocol: "GRPC"
```

```bash
kubectl apply -f ingress_1.yaml

grpcurl -v -d '{"greeting": "TEST"}' -plaintext $(minikube ip):80 hello.HelloService.SayHello
Failed to dial target host "192.168.99.112:80": context deadline exceeded

grpcurl -v -d '{"greeting": "TEST"}' -insecure $(minikube ip):443 hello.HelloService.SayHello

Resolved method descriptor:
rpc SayHello ( .hello.HelloRequest ) returns ( .hello.HelloResponse );

Request metadata to send:
(empty)

Response headers received:
content-type: application/grpc
date: Tue, 15 Sep 2020 15:01:43 GMT
server: nginx/1.19.1

Response contents:
{
  "reply": "hello TEST"
}

Response trailers received:
(empty)
Sent 1 request and received 1 response
```

It seems gRPC proxy is only working on HTTPS port but not on HTTP. But in the target service we use the `insecure` port, because TLS is terminated at the ingress level and inside of cluster gGRPC traffic travells unencripted.

### Forward encripted traffic inside the cluster

Then we should use GRPCS in annotation instead of GRPC:

```yaml
nginx.ingress.kubernetes.io/backend-protocol: "GRPCS"
```

And then we can use the `secure` port of the service.

```bash
kubectl delete ingress --all
kubectl apply -f ingress_2.yaml

grpcurl -v -d '{"greeting": "TEST"}' -plaintext $(minikube ip):80 hello.HelloService.SayHello
Failed to dial target host "192.168.99.112:80": context deadline exceeded
nikolay@leo:~/Projects/learn-kuber/minikube_nginx_grpc$ grpcurl -v -d '{"greeting": "TEST"}' -insecure $(minikube ip):443 hello.HelloService.SayHello

Resolved method descriptor:
rpc SayHello ( .hello.HelloRequest ) returns ( .hello.HelloResponse );

Request metadata to send:
(empty)

Response headers received:
content-type: application/grpc
date: Tue, 15 Sep 2020 15:11:36 GMT
server: nginx/1.19.1
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

The result the same as the previous.

### Ingress rule with host

Now we have to use the `-authority` parameter of grpcurl:

```bash
kubectl delete ingress --all
kubectl apply -f ingress_3.yaml

grpcurl -v -d '{"greeting": "TEST"}' -authority myexample.com -plaintext $(minikube ip):80 hello.HelloService.SayHello
Failed to dial target host "192.168.99.112:80": context deadline exceeded

grpcurl -v -d '{"greeting": "TEST"}' -authority myexample.com -insecure $(minikube ip):443 hello.HelloService.SayHello

Resolved method descriptor:
rpc SayHello ( .hello.HelloRequest ) returns ( .hello.HelloResponse );

Request metadata to send:
(empty)

Response headers received:
content-type: application/grpc
date: Tue, 15 Sep 2020 15:23:41 GMT
server: nginx/1.19.1

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
kubectl apply -f ingress_4.yaml

grpcurl -v -d '{"greeting": "TEST"}' -authority myexample.com -plaintext $(minikube ip):80 hello.HelloService.SayHello
Failed to dial target host "192.168.99.112:80": context deadline exceeded

grpcurl -v -d '{"greeting": "TEST"}' -authority myexample.com -insecure $(minikube ip):443 hello.HelloService.SayHello

Resolved method descriptor:
rpc SayHello ( .hello.HelloRequest ) returns ( .hello.HelloResponse );

Request metadata to send:
(empty)

Response headers received:
content-type: application/grpc
date: Tue, 15 Sep 2020 15:39:52 GMT
server: nginx/1.19.1

Response contents:
{
  "reply": "hello TEST"
}

Response trailers received:
(empty)
Sent 1 request and received 1 response
```

Output is the same. I'm not sure if the certificate is used or not, at least it doesn't break anything.

## Conclusion

nginx ingress processes gRPC requests only on secured port (443). When we call insecured port (80) is always issues the "context deadline exceeded" error that doesn't tell a much about what happens.
