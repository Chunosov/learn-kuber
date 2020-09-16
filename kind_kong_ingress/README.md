# Run simple service with kong ingress controller in kind cluster

Start a cluster:

```bash
kind create cluster --config cluster.yaml

kubectl get node -owide
NAME                 STATUS   ROLES    AGE   VERSION   INTERNAL-IP   EXTERNAL-IP   OS-IMAGE       KERNEL-VERSION       CONTAINER-RUNTIME
kind-control-plane   Ready    master   18m   v1.18.2   172.19.0.2    <none>        Ubuntu 19.10   4.15.0-112-generic   containerd://1.3.3-14-g449e9269

export NODE_ADDR=172.19.0.2
```

Install kong ingress controller:

```bash
kubectl apply --filename https://raw.githubusercontent.com/Kong/kubernetes-ingress-controller/0.10.x/deploy/single/all-in-one-dbless.yaml

kubectl get service -n kong
NAME                      TYPE           CLUSTER-IP       EXTERNAL-IP   PORT(S)                      AGE
kong-proxy                LoadBalancer   10.100.50.116    <pending>     80:32581/TCP,443:31682/TCP   9s
kong-validation-webhook   ClusterIP      10.107.249.147   <none>        443/TCP                      9s

export PROXY_PORT_HTTP=32581
export PROXY_PORT_HTTPS=31682
```

Deploy simple service:

```bash
kubectl apply -f service.yaml
```

## Check if services work via NodePort

```bash
kubectl get service
NAME             TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)          AGE
apple-service    NodePort    10.107.158.229   <none>        5678:31880/TCP   11m
banana-service   NodePort    10.109.24.128    <none>        5678:30687/TCP   11m

curl http://$NODE_ADDR:31880
apple

curl http://$NODE_ADDR:30687
banana
```

## Try to access services via ingress

### Route by path

```bash
kubectl apply -f ingress_0.yaml

curl http://localhost/apple
curl: (56) Recv failure: Connection reset by peer

curl http://localhost:$PORT_HTTP/apple
curl: (7) Failed to connect to localhost port 30493: Connection refused

curl http://localhost:80/apple
curl: (56) Recv failure: Connection reset by peer
```

Unlike nginx who makes services available at the [cluster address](../kind_echo_ingress_nginx/README.md) (which is localhost in our case), kong only exposes service at NodePort of kong-proxy service:

```bash
curl http://$NODE_ADDR:$PROXY_PORT_HTTP/apple
apple

curl http://$NODE_ADDR:$PROXY_PORT_HTTP/banana
banana
```

### Route by host

```bash
kubectl delete ingress --all
kubectl apply -f ingress_1.yaml

curl -H "Host: apple.myexample.com" http://$NODE_ADDR:$PROXY_PORT_HTTP
apple

curl -H "Host: banana.myexample.com" http://$NODE_ADDR:$PROXY_PORT_HTTP
banana
```
