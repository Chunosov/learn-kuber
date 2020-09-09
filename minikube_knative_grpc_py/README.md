# Example of gRPC service on Python with minikube and knative

Service is a simple greeter taken from official [gRPC examples](https://github.com/grpc/grpc/tree/master/examples/python/helloworld).

Start minikube cluster and install knative:

```bash
minikube start --driver=virtualbox
../knative_install.sh
```

## Run from sources

```bash
python3.8 -m venv .venv
source .venv/bin/activate
pip install grpcio protobuf
python greeter_server.py

# from another terminal:
python greeter_client.py localhost:50051
Greeter client received: Hello, you!  # <-- it works
Greeter client received: Hello again, you!  # <-- it works
```

## Run from docker image

```bash
docker build -t grpc_greeter:v0 .
docker run --rm -d -p 50051:50051 --name grpc_greeter grpc_greeter:v0
python greeter_client.py localhost:50051
Greeter client received: Hello, you!  # <-- it works
Greeter client received: Hello again, you!  # <-- it works
docker stop grpc_greeter
```

## Run as plain kubernetes service

```bash
eval $(minikube -p minikube docker-env)
docker build -t grpc_greeter:v0 .
kubectl apply -f service_plain.yaml
export GREETER_PORT=$(kubectl get services/grpc-greeter-service -o go-template='{{(index .spec.ports 0).nodePort}}')
python greeter_client.py $(minikube ip):$GREETER_PORT
Greeter client received: Hello, you!  # <-- it works
Greeter client received: Hello again, you!  # <-- it works
kubectl delete -f service_plain.yaml
```

## Run as knative service

```bash
eval $(minikube -p minikube docker-env)
docker build -t dev.local/grpc_greeter:v0 .
# Tags resolution should be disabled for dev.local, see ../README.md
kubectl apply -f service_knative.yaml

kubectl get service
NAME                         TYPE           CLUSTER-IP       EXTERNAL-IP                        PORT(S)                             AGE
grpc-greeter                 ExternalName   <none>           grpc-greeter.default.example.com   <none>                              117s
grpc-greeter-k2l9v           ClusterIP      10.101.146.238   <none>                             80/TCP                              118s
grpc-greeter-k2l9v-private   ClusterIP      10.101.122.58    <none>                             80/TCP,9090/TCP,9091/TCP,8022/TCP   118s
kubernetes                   ClusterIP      10.96.0.1        <none>                             443/TCP                             4d1h

kubectl get service -n kong
NAME                      TYPE           CLUSTER-IP       EXTERNAL-IP   PORT(S)                      AGE
kong-proxy                LoadBalancer   10.98.127.254    <pending>     80:32188/TCP,443:32637/TCP   4d1h
```

## Test with python client

```bash
python greeter_client_knative.py $(minikube ip):32188
```

<span style="color:red"><b>Doesn't work.</b> Python gRPC API is not working through proxy because we can't provide additional headers to gRPC call. Additional "Host" header is required to allow kong-proxy resolving to what service it should requrect a request. There is also something called `:authority` pseudoheader (whatever it means) in HTTP/2 seemingly allowing the same behaviour. Python guys seem not very interesting in fixing this issues as there were two PRs in python gRPC repo providing [authority](https://github.com/grpc/grpc/pull/14077) and [host](https://github.com/grpc/grpc/pull/14361) params, but both of them were closed without merging.</span>
