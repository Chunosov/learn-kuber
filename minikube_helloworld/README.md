# Run helloworld service in minikube cluster via CLI

Start demo cluster:

```bash
minikube start --driver=kvm2
```

Make a deployment. This command runs container in a pod, but it is unavailable from outside of cluster.

```bash
kubectl create deployment hello-minikube --image=k8s.gcr.io/echoserver:1.10
deployment.apps/hello-minikube created

kubectl get deployment
NAME             READY   UP-TO-DATE   AVAILABLE   AGE
hello-minikube   1/1     1            1           61s

kubectl get pods
NAME                              READY   STATUS    RESTARTS   AGE
hello-minikube-64b64df8c9-v5t7r   1/1     Running   0          2m35s

docker ps | grep echo
603c7183f61b        k8s.gcr.io/echoserver   "/usr/local/bin/run.…"   3 minutes ago       Up 3 minutes                            k8s_echoserver_hello-minikube-64b64df8c9-v5t7r_default_6e0075ca-ae88-4526-a7f7-231a65485ea1_0

```

To make running app available from outside of the cluster we need to make a service:

```bash
kubectl expose deployment hello-minikube --type=NodePort --port=8080
service/hello-minikube exposed

kubectl get services
NAME             TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)          AGE
hello-minikube   NodePort    10.101.132.171   <none>        8080:31186/TCP   45s
```

Service type NodePort means that a port from a node will be exposed at the cluster address. Here 8080 is port on the node, 31186 is respective port on the cluster address. Cluster address can be obtained with `minikube ip`.

```bash
export NODE_PORT=$(kubectl get services/hello-minikube -o go-template='{{(index .spec.ports 0).nodePort}}')
echo NODE_PORT=$NODE_PORT
```

Or service address can be obtained with:

```bash
minikube service hello-minikube --url
http://192.168.39.196:31186
```

Now we can call the service:

```bash
curl $(minikube ip):$NODE_PORT


Hostname: hello-minikube-64b64df8c9-v5t7r

Pod Information:
	-no pod information available-

Server values:
	server_version=nginx: 1.13.3 - lua: 10008

Request Information:
	client_address=172.17.0.1
	method=GET
	real path=/
	query=
	request_version=1.1
	request_scheme=http
	request_uri=http://192.168.39.196:8080/

Request Headers:
	accept=*/*
	host=192.168.39.196:31186
	user-agent=curl/7.58.0

Request Body:
	-no body in request
```

Stop the cluster and delete it:

```bash
minikube stop
✋  Stopping node "minikube"  ...
🛑  1 nodes stopped.

minikube delete
🔥  Deleting "minikube" in kvm2 ...
💀  Removed all traces of the "minikube" cluster.

```
