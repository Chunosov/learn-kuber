# Run echo service and nginx ingress controller in kind cluster

This is an [official example](https://kind.sigs.k8s.io/docs/user/ingress/#ingress-nginx) of usage nginx ingress in kind cluster.

Start a cluster. `ingress-ready=true` node label is required for the controller to be installed, so we use a yaml config instead of makind default cluster:

```bash
kind create cluster --config cluster.yaml
```

<details>
  <summary>Error "port is already allocated"</summary>

If you see an error similar to this:

```
Creating cluster "kind" ...
 ✓ Ensuring node image (kindest/node:v1.18.2) 🖼
 ✗ Preparing nodes 📦
ERROR: failed to create cluster: docker run error: command "docker run --hostname kind-control-plane --name kind-control-plane --label io.x-k8s.kind.role=control-plane --privileged --security-opt seccomp=unconfined --security-opt apparmor=unconfined --tmpfs /tmp --tmpfs /run --volume /var --volume /lib/modules:/lib/modules:ro --detach --tty --label io.x-k8s.kind.cluster=kind --net kind --restart=on-failure:1 --publish=0.0.0.0:80:80/TCP --publish=127.0.0.1:32865:6443/TCP kindest/node:v1.18.2@sha256:7b27a6d0f2517ff88ba444025beae41491b016bc6af573ba467b70c5e8e0d85f" failed with error: exit status 125
Command Output: 502bfa850c7786918aadbbe962f3cc645ebdd7373a37549811850c4c112f4393
docker: Error response from daemon: driver failed programming external connectivity on endpoint kind-control-plane (9e581c838fd744b7ca264b40c9f10de6351b00ec3d1f6c0ac86b2976db977869): Bind for 0.0.0.0:80 failed: port is already allocated.
```

Then instead of searching for what application is holding the port, search for another forgotten kind cluster and delete it:

```bash
kind get clusters
some-cluster-you-have-forget-about

kind delete cluster --name some-cluster-you-have-forget-about
```
---
</details>

&nbsp;

Install nginx ingress controller and wait while it is started:

```bash
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/master/deploy/static/provider/kind/deploy.yaml

kubectl wait --namespace ingress-nginx \
  --for=condition=ready pod \
  --selector=app.kubernetes.io/component=controller \
  --timeout=120s
```

That's what we have now:

```bash
kubectl get service -n ingress-nginx
NAME                                 TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)                      AGE
ingress-nginx-controller             NodePort    10.105.193.214   <none>        80:31268/TCP,443:30886/TCP   3m42s
ingress-nginx-controller-admission   ClusterIP   10.106.238.180   <none>        443/TCP                      3m43s
```

Deploy two example services:

```bash
kubectl apply -f service.yaml
```

Check if they work:

```bash
curl http://localhost/foo
foo

curl http://localhost/bar
bar
```

The same services available at HTTPS port:

```bash
curl -k https://localhost/foo
foo

curl -k https://localhost/bar
bar
```
