# Run example echo server behind kong-proxy on minikube

This is an extended reproduction of official [getting-started](https://github.com/Kong/kubernetes-ingress-controller/blob/master/docs/guides/getting-started.md) kong example.

## Start cluster

Run minikube cluster and [install kong](https://github.com/Kong/kubernetes-ingress-controller/blob/master/docs/deployment/minikube.md):

```bash
minikube start --driver=virtualbox
kubectl create -f https://bit.ly/k4k8s

export PROXY_IP=$(minikube service -n kong kong-proxy --url | head -1)
echo $PROXY_IP
http://192.168.99.108:30552
```

Currently there are no routes:

```bash
curl -i $PROXY_IP
HTTP/1.1 404 Not Found
Date: Wed, 09 Sep 2020 08:44:39 GMT
Content-Type: application/json; charset=utf-8
Connection: keep-alive
Content-Length: 48
X-Kong-Response-Latency: 0
Server: kong/2.1.3

{"message":"no Route matched with those values"}
```

## Deploy service

Deploy sample echo servce:

```bash
kubectl apply -f service.yaml
# or
kubectl apply -f https://bit.ly/echo-service

kubectl get service
NAME         TYPE        CLUSTER-IP     EXTERNAL-IP   PORT(S)           AGE
echo         ClusterIP   10.96.183.53   <none>        8080/TCP,80/TCP   8m13s
kubernetes   ClusterIP   10.96.0.1      <none>        443/TCP           19m
```

## Check in cluster

Service has the ClusterIP type and is available from inside of cluster:

```bash
kubectl describe service echo
Name:              echo
Namespace:         default
Labels:            app=echo
Annotations:       Selector:  app=echo
Type:              ClusterIP
IP:                10.96.183.53
Port:              high  8080/TCP
TargetPort:        8080/TCP
Endpoints:         172.17.0.4:8080
Port:              low  80/TCP
TargetPort:        8080/TCP
Endpoints:         172.17.0.4:8080
Session Affinity:  None
Events:            <none>

minikube ssh
curl -i 172.17.0.4:8080
HTTP/1.1 200 OK
Date: Wed, 09 Sep 2020 08:59:07 GMT
Content-Type: text/plain
Transfer-Encoding: chunked
Connection: keep-alive
Server: echoserver

Hostname: echo-78b867555-rstj6

Pod Information:
	node name:	minikube
	pod name:	echo-78b867555-rstj6
	pod namespace:	default
	pod IP:	172.17.0.4

Server values:
	server_version=nginx: 1.12.2 - lua: 10010

Request Information:
	client_address=172.17.0.1
	method=GET
	real path=/
	query=
	request_version=1.1
	request_scheme=http
	request_uri=http://172.17.0.4:8080/

Request Headers:
	accept=*/*
	host=172.17.0.4:8080
	user-agent=curl/7.66.0

Request Body:
	-no body in request-
```

## Check outside of cluster

### Resolve route by path

Make an ingress rule to make it available from outside of the cluster:

```bash
kubectl apply -f ingress_path.yaml

kubectl describe ingress
Name:             demo
Namespace:        default
Address:
Default backend:  default-http-backend:80 (<error: endpoints "default-http-backend" not found>)
Rules:
  Host        Path  Backends
  ----        ----  --------
  *
              /foo     echo:80 (172.17.0.4:8080)
Annotations:  Events:  <none>
```

Service can be accessed from host machine:

```bash
curl -i $PROXY_IP/foo

HTTP/1.1 200 OK
Content-Type: text/plain; charset=UTF-8
Transfer-Encoding: chunked
Connection: keep-alive
Date: Wed, 09 Sep 2020 09:10:24 GMT
Server: echoserver
X-Kong-Upstream-Latency: 4
X-Kong-Proxy-Latency: 2
Via: kong/2.1.3

Hostname: echo-78b867555-rstj6

Pod Information:
	node name:	minikube
	pod name:	echo-78b867555-rstj6
	pod namespace:	default
	pod IP:	172.17.0.4

Server values:
	server_version=nginx: 1.12.2 - lua: 10010

Request Information:
	client_address=172.17.0.2
	method=GET
	real path=/foo
	query=
	request_version=1.1
	request_scheme=http
	request_uri=http://192.168.99.108:8080/foo

Request Headers:
	accept=*/*
	connection=keep-alive
	host=192.168.99.108:30552
	user-agent=curl/7.58.0
	x-forwarded-for=172.17.0.1
	x-forwarded-host=192.168.99.108
	x-forwarded-port=8000
	x-forwarded-proto=http
	x-real-ip=172.17.0.1

Request Body:
	-no body in request-
```

### Resolve route by host

Test another ingress rule:

```bash
# service is not available before
curl -H "Host: api.foo.bar" $PROXY_IP
{"message":"no Route matched with those values"}

# make a new rule
kubectl apply -f ingress_host.yaml

# now it can be accessed with Host header in request
curl -H "Host: api.foo.bar" $PROXY_IP
Hostname: echo-78b867555-rstj6

Pod Information:
	node name:	minikube
	pod name:	echo-78b867555-rstj6
	pod namespace:	default
	pod IP:	172.17.0.4

Server values:
	server_version=nginx: 1.12.2 - lua: 10010

Request Information:
	client_address=172.17.0.2
	method=GET
	real path=/
	query=
	request_version=1.1
	request_scheme=http
	request_uri=http://api.foo.bar:8080/

Request Headers:
	accept=*/*
	connection=keep-alive
	host=api.foo.bar
	user-agent=curl/7.58.0
	x-forwarded-for=172.17.0.1
	x-forwarded-host=api.foo.bar
	x-forwarded-port=8000
	x-forwarded-proto=http
	x-real-ip=172.17.0.1

Request Body:
	-no body in request-
```

### Resolve route by host and path

```bash
# service is not available before
curl -H "Host: api.bar" $PROXY_IP/bar
{"message":"no Route matched with those values"}

# make a new rule
kubectl apply -f ingress_host_path.yaml

# now it can be accessed with Host header and path in request
curl -H "Host: api.bar" $PROXY_IP/bar
Hostname: echo-78b867555-rstj6

Pod Information:
	node name:	minikube
	pod name:	echo-78b867555-rstj6
	pod namespace:	default
	pod IP:	172.17.0.4

Server values:
	server_version=nginx: 1.12.2 - lua: 10010

Request Information:
	client_address=172.17.0.2
	method=GET
	real path=/bar
	query=
	request_version=1.1
	request_scheme=http
	request_uri=http://api.bar:8080/bar

Request Headers:
	accept=*/*
	connection=keep-alive
	host=api.bar
	user-agent=curl/7.58.0
	x-forwarded-for=172.17.0.1
	x-forwarded-host=api.bar
	x-forwarded-port=8000
	x-forwarded-proto=http
	x-real-ip=172.17.0.1

Request Body:
	-no body in request-
```

There are all routes we have:

```bash
kubectl get ingress
NAME    CLASS    HOSTS         ADDRESS   PORTS   AGE
demo    <none>   *                       80      21m
demo1   <none>   api.foo.bar             80      6m38s
demo2   <none>   api.bar                 80      70s
```
