# Run echo service and nginx ingress controller in kind cluster

This is an [official example](https://kind.sigs.k8s.io/docs/user/ingress/#ingress-nginx) of usage nginx ingress in kind cluster.

Start a cluster. `ingress-ready=true` node label is required for the controller to be installed, so we use a yaml config instead of makind default cluster:

```bash
kind create cluster --config cluster.yaml
```

Install nginx ingress controller and wait while it is started:

```bash
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/master/deploy/static/provider/kind/deploy.yaml

kubectl wait --namespace ingress-nginx \
  --for=condition=ready pod \
  --selector=app.kubernetes.io/component=controller \
  --timeout=90s
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
curl localhost:8080/foo
foo

curl localhost:8080/bar
bar
```
