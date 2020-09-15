# Run simple echo service in kind cluster

Start a cluster and deploy sample service:

```bash
kind create cluster
kubectl apply -f service.yaml
```

The service has the NodePort type so we need a node address and service port to connect to it:

```bash
kubectl get services
NAME         TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)          AGE
echo         NodePort    10.109.142.143   <none>        8080:31854/TCP   48s
kubernetes   ClusterIP   10.96.0.1        <none>        443/TCP          2m47s

kubectl get node -o wide
NAME                 STATUS   ROLES    AGE    VERSION   INTERNAL-IP   EXTERNAL-IP   OS-IMAGE       KERNEL-VERSION       CONTAINER-RUNTIME
kind-control-plane   Ready    master   6m3s   v1.18.2   172.19.0.2    <none>        Ubuntu 19.10   4.15.0-112-generic   containerd://1.3.3-14-g449e9269
```

Call the service:

```bash
curl 172.19.0.2:31854

Hostname: echo-65c5c6c87f-qq7m8

Pod Information:
	-no pod information available-

Server values:
	server_version=nginx: 1.13.3 - lua: 10008

Request Information:
	client_address=10.244.0.1
	method=GET
	real path=/
	query=
	request_version=1.1
	request_scheme=http
	request_uri=http://172.19.0.2:8080/

Request Headers:
	accept=*/*
	host=172.19.0.2:31854
	user-agent=curl/7.58.0

Request Body:
	-no body in request-
```

It works.
