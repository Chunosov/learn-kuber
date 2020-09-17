# Run simple gRPC service with Istio in minikube cluster

Start a cluster, [install istio](https://istio.io/latest/docs/setup/getting-started/#install), and get ingress HTTP port number:

```bash
minikube start --driver=virtualbox
istioctl install --set profile=demo
kubectl label namespace default istio-injection=enabled
```

Deploy simple echo services:

```bash
kubectl apply -f service.yaml

export SERVICE_PORT_HTTP=$(kubectl get service grpcbin --output jsonpath="{.spec.ports[0].nodePort}")
export SERVICE_PORT_HTTPS=$(kubectl get service grpcbin --output jsonpath="{.spec.ports[1].nodePort}")
echo $SERVICE_PORT_HTTP,$SERVICE_PORT_HTTPS
32414,31328
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

## Test via Istio gateway

### Use 'intended' grpc port

It's recommended [here](https://stackoverflow.com/questions/44760416/istio-ingress-with-grpc-and-http) to add another port to `istio-ingressgateway` service. Use either `kubectl edit` command:

```bash
kubectl edit svc -n istio-system istio-ingressgateway
```

or an alternative way (for some reason the edit command failed for me):

```bash
kubectl get svc -n istio-system istio-ingressgateway -o yaml > istio.yaml
```

then add another port definition into `spec.ports:` section

```yaml
  - name: grpc
    nodePort: 30000
    port: 50051
    protocol: TCP
    targetPort: 50051
```

and apply changes:

```bash
kubectl apply -f istio.yaml

kubectl get svc -n istio-system istio-ingressgateway
NAME                   TYPE       CLUSTER-IP      EXTERNAL-IP   PORT(S)                                                                                      AGE
istio-ingressgateway   NodePort   10.102.252.82   <none>        15021:31560/TCP,80:32061/TCP,443:31189/TCP,31400:31864/TCP,15443:31850/TCP,50051:30000/TCP   80m
```

**Note:** I'm not sure how it can help because we already have port definition for TCP protocol, why not to use that port instead of introducing another one but differently named.

Make a gateway and check connection. On minikube, LoadBalancer services can be accessed in the same way as NodePort ones:

```bash
kubectl apply -f gateway_1.yaml

grpcurl -proto hello.proto -v -d '{"greeting": "TEST"}' -plaintext $(minikube ip):30000 hello.HelloService.SayHello

Resolved method descriptor:
rpc SayHello ( .hello.HelloRequest ) returns ( .hello.HelloResponse );

Request metadata to send:
(empty)

Response headers received:
(empty)

Response trailers received:
content-type: application/grpc
date: Thu, 17 Sep 2020 12:09:41 GMT
server: istio-envoy
Sent 1 request and received 0 responses
ERROR:
  Code: Unimplemented
  Message:

grpcurl -proto hello.proto -v -d '{"greeting": "TEST"}' -insecure $(minikube ip):30000 hello.HelloService.SayHello
Failed to dial target host "192.168.99.115:30000": tls: first record does not look like a TLS handshake
```

Try at cluster IP under `minikube tunnel` (run tunnel in another terminal):

```bash
 grpcurl -proto hello.proto -v -d '{"greeting": "TEST"}' -plaintext 10.102.252.82:50051 hello.HelloService.SayHello

Resolved method descriptor:
rpc SayHello ( .hello.HelloRequest ) returns ( .hello.HelloResponse );

Request metadata to send:
(empty)

Response headers received:
(empty)

Response trailers received:
content-type: application/grpc
date: Thu, 17 Sep 2020 12:08:18 GMT
server: istio-envoy
Sent 1 request and received 0 responses
ERROR:
  Code: Unimplemented
  Message:

grpcurl -proto hello.proto -v -d '{"greeting": "TEST"}' -insecure 10.102.252.82:50051 hello.HelloService.SayHello
Failed to dial target host "10.102.252.82:50051": tls: first record does not look like a TLS handshake
```

### Use existed TCP port

```bash
kubectl delete -f gateway_1.yaml
kubectl apply -f gateway_2.yaml

# under minikube tunnel
grpcurl -proto hello.proto -v -d '{"greeting": "TEST"}' -plaintext 10.102.252.82:31400 hello.HelloService.SayHello

Resolved method descriptor:
rpc SayHello ( .hello.HelloRequest ) returns ( .hello.HelloResponse );

Request metadata to send:
(empty)

Response headers received:
(empty)

Response trailers received:
content-type: application/grpc
date: Thu, 17 Sep 2020 12:12:57 GMT
server: istio-envoy
Sent 1 request and received 0 responses
ERROR:
  Code: Unimplemented
  Message:
```

### Use HTTP port

```bash
kubectl delete -f gateway_2.yaml
kubectl apply -f gateway_3.yaml

# under minikube tunnel
grpcurl -proto hello.proto -v -d '{"greeting": "TEST"}' -plaintext 10.102.252.82:80 hello.HelloService.SayHello

Resolved method descriptor:
rpc SayHello ( .hello.HelloRequest ) returns ( .hello.HelloResponse );

Request metadata to send:
(empty)

Response headers received:
(empty)

Response trailers received:
content-type: application/grpc
date: Thu, 17 Sep 2020 12:19:36 GMT
server: istio-envoy
Sent 1 request and received 0 responses
ERROR:
  Code: Unimplemented
  Message:
```

<span style="color:red"><b>It doesn't work.</b> It gives the same result as [in this question](https://discuss.istio.io/t/error-calling-grpc-from-client-outside-cluster/2837/2), no answer there. [This bug](https://github.com/istio/istio/issues/5295) asks about not working grpc, and the only aswer "It works" (Thanks for nothing, guys). [This question](https://stackoverflow.com/questions/44760416/istio-ingress-with-grpc-and-http) about the same, no useful answer. [Another question](https://discuss.istio.io/t/error-calling-grpc-from-client-outside-cluster/2837/2) about the same, the only strange answer about configuring a "caller side". [Another issue](https://github.com/istio/istio/issues/7909), closed with no reliable answer.</span>
