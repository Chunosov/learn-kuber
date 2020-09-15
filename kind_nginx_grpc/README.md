# Start sample gRPC service with nginx ingress controller on kind cluster

https://kubernetes.github.io/ingress-nginx/examples/grpc/
https://github.com/kubernetes/ingress-nginx/tree/master/docs/examples/grpc
https://medium.com/@Alibaba_Cloud/accessing-grpc-services-through-container-service-for-kubernetes-ingress-controller-511f984071e1

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

If it show's an error like:

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
NAME         TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)          AGE
grpcbin      NodePort    10.100.104.58    <none>        9000:30882/TCP   11

kubectl get node -o wide
NAME                 STATUS   ROLES    AGE   VERSION   INTERNAL-IP   EXTERNAL-IP   OS-IMAGE       KERNEL-VERSION       CONTAINER-RUNTIME
kind-control-plane   Ready    master   18m   v1.18.2   172.19.0.8    <none>        Ubuntu 19.10   4.15.0-112-generic   containerd://1.3.3-14-g449e9269

curl http://172.19.0.8:30882
Warning: Binary output can mess up your terminal. Use "--output -" to tell
Warning: curl to output it to your terminal anyway, or consider "--output
Warning: <FILE>" to save to a file.
```

Seems it works and provides some binary output. Check with [grpcurl](../README.md#grpcul):

```bash
grpcurl -v -d '{"greeting": "TEST"}' -plaintext 172.19.0.8:30882 hello.HelloService.SayHello

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

It works as NodePort.

## Check via Ingress

Check if the service is available through ingress controller:

```bash
curl --insecure -H "Host: api.bar" https://localhost:4443
curl: (35) OpenSSL SSL_connect: SSL_ERROR_SYSCALL in connection to localhost:4443

grpcurl -v -d '{"greeting": "TEST"}' -insecure -H "Host: api.bar" -authority api.bar localhost:4443 hello.HelloService.SayHello
Failed to dial target host "localhost:4443": EOF
```

<span style="color:red"><b>It doesn't work.</b></span>
