# Usage of minikube for examples

## Installation
Install kubectl if you don't have it:

```bash
curl -LO "https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl"
chmod +x ./kubectl
sudo mv ./kubectl /usr/local/bin/kubectl
```

Install minikube:

```bash
curl -Lo minikube https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
chmod +x minikube
sudo mv minikube /usr/local/bin
```

## Create clusters

Start minikube cluster:

```bash
minikube start --driver=kvm2
```

or you can use virtualbox VM:

```bash
minikube start --driver=virtualbox
```

## Use context

Minikube overrides active kubectl context

```bash
kubectl config get-contexts
CURRENT   NAME        CLUSTER     AUTHINFO    NAMESPACE
          kind-kind   kind-kind   kind-kind
*         minikube    minikube    minikube    default
```

so it will be used by default:

```bash
kubectl get nodes
NAME       STATUS   ROLES                  AGE    VERSION
minikube   Ready    control-plane,master   104s   v1.23.3
```

You can switch back to kind context:

```bash
kubectl config use-context kind-kind
```

Then all the next calls to kubectl will be directed to kind cluster:

```bash
kubectl get nodes
NAME                 STATUS   ROLES                  AGE   VERSION
kind-control-plane   Ready    control-plane,master   28h   v1.21.1
```

## Minikube vs kind

Important difference between minikube and kind that minikube uses VM and hence it supports node restart:

```bash
minikube stop

....

minikube start
# cluster restored
```

This can be important for testing some issues related to node restarts. While kind is docker based and it doesn't offcially supports node restarts. Even though kind cluster survives host machine reboots, it doesn't promice to work properly after that :)

## Use dev docker images in minikube

While in kind cluster you can use `kind load docker-image` for using some custom docker images in deployments, there is another way in minikube.

You have to export minikube docker env into your host machine:

```bash
eval $(minikube -p minikube docker-env)
```

Then you can see what images used in minikube VM:

```bash
docker images

# they are in minikube VM, not in your host machine
REPOSITORY                                TAG       IMAGE ID       CREATED         SIZE
k8s.gcr.io/kube-apiserver                 v1.23.3   f40be0088a83   2 months ago    135MB
k8s.gcr.io/kube-scheduler                 v1.23.3   99a3486be4f2   2 months ago    53.5MB
....
```

You can see the same if you ssh into minikube node:

```bash
minikube ssh
                         _             _
            _         _ ( )           ( )
  ___ ___  (_)  ___  (_)| |/')  _   _ | |_      __
/' _ ` _ `\| |/' _ `\| || , <  ( ) ( )| '_`\  /'__`\
| ( ) ( ) || || ( ) || || |\`\ | (_) || |_) )(  ___/
(_) (_) (_)(_)(_) (_)(_)(_) (_)`\___/'(_,__/'`\____)

$ docker images
REPOSITORY                                TAG       IMAGE ID       CREATED         SIZE
k8s.gcr.io/kube-apiserver                 v1.23.3   f40be0088a83   2 months ago    135MB
k8s.gcr.io/kube-scheduler                 v1.23.3   99a3486be4f2   2 months ago    53.5MB
...
```

In such exported env you can build dev images that will be used in test deployments:

```
docker build -t my-image:0.1 .
```

Then any deployment using `my-image:0.1` will use that image.

## Delete cluster

```bash
minikube delete
```

You *have* to switch to another context manually after you deleted minikube cluster. Because minikube doesn't restore context and there will be no active context after minikube cluster deleted:

```bash
kubectl config get-contexts
CURRENT   NAME        CLUSTER     AUTHINFO    NAMESPACE
          kind-kind   kind-kind   kind-kind 

kubectl config use-context kind-kind
```
