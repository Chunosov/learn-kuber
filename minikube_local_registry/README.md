# Setup local docker registry with minikube

The magic is taken from [here](https://github.com/triggermesh/knative-local-registry).

Start demo cluster:

```bash
minikube start --driver=virtualbox
```

Make local registry

```bash
kubectl create namespace registry
kubectl apply -f local-registry.yaml
```

It should be available from inside of the cluster at `knative.registry.svc.cluster.local:80` addess:

```bash
minikube ssh
curl knative.registry.svc.cluster.local:80/v2/_catalog

{"repositories":[]}
```

Connect to the cluster's docker and try to build and push an image:

```bash
eval $(minikube -p minikube docker-env)
docker build -t my-new-image:v0 ../minikube_local_image
docker tag my-new-image:v0 knative.registry.svc.cluster.local/my-new-image:v0
docker push knative.registry.svc.cluster.local/my-new-image:v0
```

Try pulling manually:

```bash
# remove the image for a clean experiment
docker rmi knative.registry.svc.cluster.local/my-new-image:v0

docker pull knative.registry.svc.cluster.local/my-new-image:v0
docker run --rm -d -p 5000:5000 --name my-new-service-v0 knative.registry.svc.cluster.local/my-new-image:v0
minikube ssh
curl localhost:5000

Hello World!  # <-- it works!

docker stop my-new-service-v0
```

Try as kubernetes service:

```bash
# remove image for a clean experiment
docker rmi knative.registry.svc.cluster.local/my-new-image:v0

kubectl apply -f service.yaml

curl $(minikube service my-new-service-v0-service --url)

Hello World!  # <-- it works!
```

kubernetes can pull the image and run the service from it.
