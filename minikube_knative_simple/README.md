# Run simple service from local image via knative in minikube

Start demo cluster and install knative:

```bash
minikube start --driver=virtualbox
../knative_install.sh
```

We'll use simple helloworld service from the previous "minikube_local_image" step:

```bash
eval $(minikube -p minikube docker-env)
docker build -t kuber_learn__simple_service_1:v0 ../minikube_local_image
```

## Without docker registry

knative use special rules for dealing with `dev.local/*` images, this case is describe [here](https://github.com/knative/serving/issues/6101):

```bash
docker tag kuber_learn__simple_service_1:v0 dev.local/kuber_learn__simple_service_1:v0
```

Then edit knative configuration to disable [tag resolution](https://knative.dev/docs/serving/tag-resolution/) for `dev.local`:

```bash
KUBE_EDITOR="nano" kubectl -n knative-serving edit configmap config-deployment

# then copy the "registriesSkippingTagResolving" line
# from the "_example" section to its parent "data" section.
```

Then deploy the service:

```bash
kubectl apply -f kubeconfig_local_images.yaml
```

Get kong-proxy port and try to call the service:

```bash
kubectl get ksvc
NAME                            URL                                                        LATESTCREATED                         LATESTREADY                           READY   REASON
kuber-learn--simple-service-1   http://kuber-learn--simple-service-1.default.example.com   kuber-learn--simple-service-1-7r6rv   kuber-learn--simple-service-1-7r6rv   True

kubectl get service -n kong
NAME                      TYPE           CLUSTER-IP       EXTERNAL-IP   PORT(S)                      AGE
kong-proxy                LoadBalancer   10.98.127.254    <pending>     80:32188/TCP,443:32637/TCP   4d

curl -H "Host: kuber-learn--simple-service-1.default.example.com" http://$(minikube ip):32188
Hello World! I'm on internal port 8080  # <-- It works
```

## Using local docker registry

There is some magic given [here](https://github.com/triggermesh/knative-local-registry) and tested [locally](../minikube_local_registry/README.md).

Install local cluster registry, push the image into it, and remove existed services and images (for clear experiment):

```bash
# remove service from the previous step, if exists
kubectl delete -f kubeconfig_local_images.yaml

# create registry
kubectl create namespace registry
kubectl apply -f ../minikube_local_registry/local-registry.yaml

# push the image
docker tag kuber_learn__simple_service_1:v0 knative.registry.svc.cluster.local/kuber_learn__simple_service_1:v0
docker push knative.registry.svc.cluster.local/kuber_learn__simple_service_1:v0

# remove local copies
docker rmi knative.registry.svc.cluster.local/kuber_learn__simple_service_1:v0
docker rmi kuber_learn__simple_service_1:v0
```

Then edit knative configuration to disable tag resolution for our local registry:

```bash
KUBE_EDITOR="nano" kubectl -n knative-serving edit configmap config-deployment

# add local registry address to the "registriesSkippingTagResolving" option:
# registriesSkippingTagResolving: dev.local,knative.registry.svc.cluster.local
```

Deploy the service:

```bash
kubectl apply -f kubeconfig_local_registry.yaml
```

Get kong-proxy port and try to call the service:

```bash
kubectl get ksvc
NAME                            URL                                                        LATESTCREATED                         LATESTREADY                           READY   REASON
kuber-learn--simple-service-1   http://kuber-learn--simple-service-1.default.example.com   kuber-learn--simple-service-1-2z4t6   kuber-learn--simple-service-1-2z4t6   True

kubectl get service -n kong
NAME                      TYPE           CLUSTER-IP       EXTERNAL-IP   PORT(S)                      AGE
kong-proxy                LoadBalancer   10.98.127.254    <pending>     80:32188/TCP,443:32637/TCP   4d

curl -H "Host: kuber-learn--simple-service-1.default.example.com" http://$(minikube ip):32188
Hello World! I'm on internal port 8080  # <-- It works
```
