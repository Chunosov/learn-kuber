# Helloworld application in python using knative on minikube

This is reproduction of the [official example](https://knative.dev/docs/serving/samples/hello-world/helloworld-python/).

## 1. Start

Start minikube cluster and install knative:

```bash
minikube start --driver=virtualbox
../knative_install.sh
```

## 2. Build and test the image

Connect to cluster's docker and build an image:

```bash
eval $(minikube -p minikube docker-env)
docker build -t dev.local/helloworld_py:v0 .
```

Try to run image directly:

```bash
minikube ssh
docker run --rm -d -p 8080:8080 -e PORT=8080 --name helloworld_py dev.local/helloworld_py:v0
curl localhost:8080
Hello World!  # <-- It works
docker stop helloworld_py
```

## 3. Delploy knative service

Edit knative configuration to disable [tag resolution](https://knative.dev/docs/serving/tag-resolution/) for `dev.local`

```bash
KUBE_EDITOR="nano" kubectl -n knative-serving edit configmap config-deployment

# Then copy the "registriesSkippingTagResolving" line
# from the "_example" section to its parent "data" section.
# It should look similar to:
# registriesSkippingTagResolving: dev.local,ko.local
```

Deploy the service:

```bash
kubectl apply -f kubeconfig.yaml
```

## 4. Test knative service

Get kong-proxy port and try to call the service:

```bash
kubectl get services
NAME                          TYPE           CLUSTER-IP       EXTERNAL-IP                         PORT(S)                             AGE
helloworld-py                 ExternalName   <none>           helloworld-py.default.example.com   <none>                              5m15s
helloworld-py-j8nl4           ClusterIP      10.108.115.135   <none>                              80/TCP                              5m25s
helloworld-py-j8nl4-private   ClusterIP      10.101.172.65    <none>                              80/TCP,9090/TCP,9091/TCP,8022/TCP   5m25s
kubernetes                    ClusterIP      10.96.0.1        <none>                              443/TCP                             3d22h

kubectl get service -n kong
NAME                      TYPE           CLUSTER-IP       EXTERNAL-IP   PORT(S)                      AGE
kong-proxy                LoadBalancer   10.98.127.254    <pending>     80:32188/TCP,443:32637/TCP   3d22h

curl -H "Host: helloworld-py.default.example.com" http://$(minikube ip):32188
Hello Python Sample v0!  # <--  It works
```
