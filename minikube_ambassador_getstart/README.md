# Run get-started example with Ambassador ingress in minikube cluster

This is a reproduction of abmassador's [quick start](https://www.getambassador.io/docs/latest/tutorials/quickstart-demo/) example.

Start minikube cluster, install [Ambassador](https://www.getambassador.io/docs/latest/tutorials/getting-started/), and wait while it gets started:

```bash
minikube start --driver virtualbox

kubectl apply -f https://www.getambassador.io/yaml/aes-crds.yaml
kubectl wait --for condition=established --timeout=90s crd -lproduct=aes
kubectl apply -f https://www.getambassador.io/yaml/aes.yaml
kubectl -n ambassador wait --for condition=available --timeout=90s deploy -lproduct=aes
```

That's what we have:

```bash
kubectl get service -n ambassador
NAME               TYPE           CLUSTER-IP       EXTERNAL-IP   PORT(S)                      AGE
ambassador         LoadBalancer   10.105.13.239    <pending>     80:30942/TCP,443:32299/TCP   5m6s
ambassador-admin   ClusterIP      10.102.233.45    <none>        8877/TCP                     5m6s
ambassador-redis   ClusterIP      10.106.189.201   <none>        6379/TCP                     5m6s
```

```bash
export PROXY_PORT_HTTP=$(kubectl get service ambassador -n ambassador --output jsonpath="{.spec.ports[0].nodePort}")
export PROXY_PORT_HTTPS=$(kubectl get service ambassador -n ambassador --output jsonpath="{.spec.ports[1].nodePort}")
echo $PROXY_PORT_HTTP,$PROXY_PORT_HTTPS
30942,32299
```

Access service behind the ambassador proxy:

```bash
# slash at the end is important!
curl -i http://$(minikube ip):$PROXY_PORT_HTTP/backend/
HTTP/1.1 301 Moved Permanently
location: https://192.168.99.116:30942/backend/
date: Thu, 17 Sep 2020 15:16:20 GMT
server: envoy
content-length: 0

curl -k https://$(minikube ip):$PROXY_PORT_HTTPS/backend/
{
    "server": "ample-blueberry-3i013rg9",
    "quote": "A principal idea is omnipresent, much like candy.",
    "time": "2020-09-17T15:13:04.783293272Z"
}
```

Some Ambassador UI is available at the proxy root address:

```bash
firefox http://$(minikube ip):$PROXY_PORT_HTTP
firefox https://$(minikube ip):$PROXY_PORT_HTTPS
```

Seems working.
