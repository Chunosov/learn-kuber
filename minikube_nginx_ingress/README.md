# Run simple service with nginx ingress controller on minikube

https://matthewpalmer.net/kubernetes-app-developer/articles/kubernetes-ingress-guide-nginx-example.html

Run sample minikube cluster and enable `ingress` addon (minikube uses [ingress-nginx](https://kubernetes.github.io/ingress-nginx/deploy/#minikube) by default):

```bash
minikube start --driver=virtualbox
minikube addons enable ingress
```

In general, nginx ingress services are in the `ingress-nginx` namespace, but in minikube they are in the `kube-system` namespace:

```bash
kubectl get pods --all-namespaces
NAMESPACE     NAME                                       READY   STATUS      RESTARTS   AGE
kube-system   coredns-f9fd979d6-htd5p                    1/1     Running     0          4m31s
kube-system   etcd-minikube                              1/1     Running     0          4m31s
kube-system   ingress-nginx-admission-create-8zlbr       0/1     Completed   0          4m22s
kube-system   ingress-nginx-admission-patch-s7c79        0/1     Completed   0          4m22s
kube-system   ingress-nginx-controller-789d9c4dc-qnfsv   1/1     Running     0          4m22s
kube-system   kube-apiserver-minikube                    1/1     Running     0          4m31s
kube-system   kube-controller-manager-minikube           1/1     Running     0          4m31s
kube-system   kube-proxy-p7zp4                           1/1     Running     0          4m31s
kube-system   kube-scheduler-minikube                    1/1     Running     0          4m31s
kube-system   storage-provisioner                        1/1     Running     0          4m36s
```

## Route by path

Deploy simple services:

```bash
kubectl apply -f service.yaml
```

Services are available on the default HTTP port 80:

```bash
curl http://$(minikube ip)/apple
apple

curl http://$(minikube ip)/banana
banana
```

The same services also available on the default HTTPS post 443:

```bash
curl --insecure https://$(minikube ip)/banana
banana

curl --insecure https://$(minikube ip)/apple
apple
```

**Note:** curl doesn't guess if request is HTTPS only by port number, it sends HTTP request by default:

```bash
curl --insecure $(minikube ip):443/banana
<html>
<head><title>400 The plain HTTP request was sent to HTTPS port</title></head>
<body>
<center><h1>400 Bad Request</h1></center>
<center>The plain HTTP request was sent to HTTPS port</center>
<hr><center>nginx/1.19.1</center>
</body>
</html>
```

## Route by host

nginx can host many sites (services in our case) on the same IP-address. It uses the "Host" http header to decide what site should be returned to a client.

Deploy another pair of services:

```bash
kubectl apply -f service1.yaml
```

Services are available on the default HTTP port 80:

```bash
curl -H "Host: orange.fruits.org" http://$(minikube ip)
orange

curl -H "Host: lemon.fruits.org" http://$(minikube ip)
lemon

curl -k -H "Host: orange.fruits.org" https://$(minikube ip)
orange

curl -k -H "Host: lemon.fruits.org" https://$(minikube ip)
lemon
```
