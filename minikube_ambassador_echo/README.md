# Run simple echo service with Ambassador ingress in minikube cluster

Start minikube cluster, install [Ambassador](https://www.getambassador.io/docs/latest/tutorials/getting-started/), and wait while it gets started:

```bash
minikube start --driver virtualbox

kubectl apply -f https://www.getambassador.io/yaml/aes-crds.yaml
kubectl wait --for condition=established --timeout=90s crd -lproduct=aes
kubectl apply -f https://www.getambassador.io/yaml/aes.yaml
kubectl -n ambassador wait --for condition=available --timeout=90s deploy -lproduct=aes

export PROXY_PORT_HTTP=$(kubectl get service ambassador -n ambassador --output jsonpath="{.spec.ports[0].nodePort}")
export PROXY_PORT_HTTPS=$(kubectl get service ambassador -n ambassador --output jsonpath="{.spec.ports[1].nodePort}")
echo $PROXY_PORT_HTTP,$PROXY_PORT_HTTPS
30942,32299
```

Deploy simple echo services and make sure they work as NodePort:

```bash
kubectl apply -f service.yaml

export APPLE_PORT=$(kubectl get service apple-service --output jsonpath="{.spec.ports[0].nodePort}")
export BANANA_PORT=$(kubectl get service banana-service --output jsonpath="{.spec.ports[0].nodePort}")

curl http://$(minikube ip):$APPLE_PORT
apple

curl http://$(minikube ip):$BANANA_PORT
banana
```

Deploy mapping and check if services are available through the proxy:

```bash
kubectl apply -f mapping.yaml

curl -k https://$(minikube ip):$PROXY_PORT_HTTPS/apple/
apple

curl -k https://$(minikube ip):$PROXY_PORT_HTTPS/banana/
banana
```

By default services are available via HTTPS though [docs say](https://www.getambassador.io/docs/latest/topics/using/intro-mappings/) the protocol should be HTTP if scheme is not specified. Even when we explicitly specify `http://` scheme in mapping, services still will be available only on `https://`. When we call them on HTTP port, it says:

```bash
curl -i http://$(minikube ip):$PROXY_PORT_HTTP/apple/
HTTP/1.1 301 Moved Permanently
location: https://192.168.99.116:30942/apple/
date: Thu, 17 Sep 2020 15:51:50 GMT
server: envoy
content-length: 0
```

**Question:** how to expose a service on HTTP port?

Seems working, at least partially.
