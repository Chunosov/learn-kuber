# Example of gRPC ping app on Go with minikube and knative

This is a reproduction of the official [gRPC Server - Go](https://knative.dev/docs/serving/samples/grpc-ping-go/) example from knative. gRPC support was added to [knative](https://github.com/knative/serving/pull/2539) and to [kong](https://github.com/Kong/kong/pull/4801).

## Start

Start minikube cluster, install knative, connect to the cluster's docker:

```bash
minikube start --driver=virtualbox
../knative_install.sh
eval $(minikube -p minikube docker-env)
```

Build the image:

```bash
docker build -t dev.local/grpc-ping-go .
```

## Run as docker container

```bash
minikube ssh
docker run --rm -d -p 8080:8080 --name grpc-ping-go dev.local/grpc-ping-go
docker run --rm dev.local/grpc-ping-go /client -server_addr="172.17.0.1:8080" -insecure
2020/09/08 09:08:30 Ping got hello - pong  # <-- It works
docker stop grpc-ping-go
```

Here 172.17.0.1 is the address of docker's host machine. Because client is run in another container, it need to connect to a host's port to comminicate with the server container. And this host's port in its turn is bound to a port inside of the server container. Port numbers are the same because we run the server with `-p 8080:8080` and host address is the address of the `docker0` interface (https://stackoverflow.com/questions/31324981/how-to-access-host-port-from-docker-container): `ip addr show docker0`.

## Run as plain kubernetes service

```bash
kubectl apply -f service_plain.yaml
export GREETER_PORT=$(kubectl get services/grpc-ping-go-service -o go-template='{{(index .spec.ports 0).nodePort}}')
docker run --rm dev.local/grpc-ping-go /client -server_addr=172.17.0.1:$GREETER_PORT -insecure
2020/09/08 09:39:55 Ping got hello - pong  # <-- It works
kubectl delete -f service_plain.yaml
```

## Run as knative service

### Tag image

For some reason the `latest` image is failed to pull:

```
  Warning  Failed     2m58s                 kubelet, minikube  Failed to pull image "dev.local/grpc-ping-go": rpc error: code = Unknown desc = Error response from daemon: Get https://dev.local/v2/: dial tcp: lookup dev.local on 10.0.2.3:53: read udp 10.0.2.15:44690->10.0.2.3:53: i/o timeout
```

So we need to add a version tag:

```bash
docker tag dev.local/grpc-ping-go dev.local/grpc-ping-go:v0
```

### Deploy service

Deploy the service and get proxy port (be sure that tag resolution is [disabled](../README.md) for `dev.local`):

```bash
kubectl apply -f service_knative.yaml

kubectl get service
NAME                      TYPE           CLUSTER-IP     EXTERNAL-IP                     PORT(S)                             AGE
grpc-ping                 ExternalName   <none>         grpc-ping.default.example.com   <none>                              69m
grpc-ping-xnbts           ClusterIP      10.109.91.68   <none>                          81/TCP                              69m
grpc-ping-xnbts-private   ClusterIP      10.96.4.148    <none>                          80/TCP,9090/TCP,9091/TCP,8022/TCP   69m
kubernetes                ClusterIP      10.96.0.1      <none>                          443/TCP                             4d23h

kubectl get service -n kong
NAME                      TYPE           CLUSTER-IP       EXTERNAL-IP   PORT(S)                      AGE
kong-proxy                LoadBalancer   10.98.127.254    <pending>     80:32188/TCP,443:32637/TCP   4d23h
```

### Test service

Test steps are described in [gRPC PR](https://github.com/knative/serving/pull/2539) for knative.

#### Use curl

```bash
curl -H "Host: grpc-ping.default.example.com" http://$(minikube ip):32188
net/http: HTTP/1.x transport connection broken: malformed HTTP response "\x00\x00\x06\x04\x00\x00\x00\x00\x00\x00\x05\x00\x00@\x00"

kubectl get pod
NAME                                        READY   STATUS    RESTARTS   AGE
grpc-ping-xnbts-deployment-ddd8f5d4-xhb2n   2/2     Running   0          21s

kubectl logs grpc-ping-xnbts-deployment-ddd8f5d4-xhb2n user-container
Starting server at port 8080  # <-- Server is working
```

curl doesn't understand grpc answer, but we can see that grpc service was started in a pod.

#### Use grpcurl

Download [grpcurl](https://github.com/fullstorydev/grpcurl/releases), move it to `/usr/local/bin`, and make executable.

Try to call the service:

```bash
grpcurl -plaintext -proto ./proto/ping.proto \
  -authority grpc-ping.default.example.com \
  -H "Host: grpc-ping.default.example.com" \
  -format json \
  -d '{"msg": "TEST"}' \
  $(minikube ip):32188 ping.PingService/Ping

Failed to dial target host "192.168.99.107:32188": context deadline exceeded
```

<span style="color:red"><b>Doesn't work.</b> Seems kong refuses to proxy grpc requests at all. No pods are run. Deployment is not kicked - `kubectl describe deployment` shows events "Scaled up replica ... to 1" then "Scaled down replica ... to 0", and these events relate to the revious curl run.</span>

<span style="color:red">Probably some additional kong configuration is required.</span>

#### Use example client

Extract the client app from the container and use is on minikube vm:

```bash
minikube ssh
docker inspect dev.local/grpc-ping-go:v0
...  # Search for something like:
"UpperDir": "/var/lib/docker/overlay2/aa810f01497fc601617cde8aa0d8ec89a4aebfe12440a41b421e6a395c534510/diff",
...

sudo ls /var/lib/docker/overlay2/aa810f01497fc601617cde8aa0d8ec89a4aebfe12440a41b421e6a395c534510/diff
client # <-- Here you are!

# Extract to some common place
sudo cp /var/lib/docker/overlay2/
aa810f01497fc601617cde8aa0d8ec89a4aebfe12440a41b421e6a395c534510/diff/client /tmp

# run outside of container
/tmp/client -server_addr=$PROXY_ADDR -insecure -server_host_override=grpc-ping.default.example.com

/tmp/client -server_addr=$PROXY_ADDR -insecure -server_host_override=grpc-ping.default.svc.cluster.local
```

**TODO**: <span style="color:red"><b>Doesn't work.</b> Tried different address-port combinations that can be obtained from `kubectl describe service kong-proxy -n kong`. Most of them get "Connection refused". Seems one of them got "context deadline exceeded" so it was correct combination. Don't remembed what of them it was, repeat the experiment.</red>
