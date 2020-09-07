# Run service with persistent volumes in minikube cluster

Start demo cluster (`kvm2` vm-driver doesn't support shared dirs so we have to use `virtualbox`):

```bash
minikube start --driver=virtualbox
```

Connect to docker inside of minilube VM:

```bash
eval $(minikube -p minikube docker-env)
docker ps # should show containers run in VM (not local ones)
```

Build docker image to make it available to the cluster:

```bash
docker build -t kuber_learn__home_lister:v0 .
```

Deploy the service:

```bash
kubectl apply -f kubeconfig.yaml

deployment.apps/kuber-learn--home-lister-1-deployment created
service/kuber-learn--home-lister-1-service created
persistentvolumeclaim/hosthome-pv-claim created
persistentvolume/hosthome-pv-volume created
```

Call the service:

```bash
export SERVICE_ADDR=$(minikube service kuber-learn--home-lister-1-service --url)
echo $SERVICE_ADDR
curl $SERVICE_ADDR

Content of hosthome: datasets,nikolay
```

It can read content of a dir on the physical machine.
