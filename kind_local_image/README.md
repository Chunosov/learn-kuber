# Run helloworld Python service with nginx ingress controller

It is described [here](https://kind.sigs.k8s.io/docs/user/local-registry/) how to use kind clusters with local docker images.

Create a cluster and install nginx ingress controller:

```bash
./make_cluster.sh

kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/master/deploy/static/provider/kind/deploy.yaml

kubectl wait --namespace ingress-nginx \
  --for=condition=ready pod \
  --selector=app.kubernetes.io/component=controller \
  --timeout=120s
```

Check if the registry is available:

```bash
curl http://localhost:5000/v2/_catalog
{"repositories":[]}
```


Build custom helloworld docker image from one the previous examples and push it to the local registry:

```bash
docker build -t helloworld_py:v0 ../minikube_knative_helloworld_py
docker tag helloworld_py:v0 localhost:5000/helloworld_py:v0
docker push localhost:5000/helloworld_py:v0
```

Deploy s service:

```bash
kubectl apply -f service.yaml
```

<span style='color:red'><b>It doesn't work.</b> It can't start a pod because can't use the image. The pod finally falls to the CrashLoopBackOff state. </span>

```
  ----     ------     ----                  ----                         -------
  Normal   Scheduled  2m38s                 default-scheduler            Successfully assigned default/helloworld-py to kind-control-plane
  Normal   Pulled     61s (x5 over 2m37s)   kubelet, kind-control-plane  Container image "localhost:5000/helloworld_py:v0" already present on machine
  Normal   Created    61s (x5 over 2m37s)   kubelet, kind-control-plane  Created container helloworld-py
  Normal   Started    61s (x5 over 2m37s)   kubelet, kind-control-plane  Started container helloworld-py
  Warning  BackOff    44s (x10 over 2m35s)  kubelet, kind-control-plane  Back-off restarting failed container
```