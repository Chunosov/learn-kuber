# kubernetes lessons and samples

## Sequence of samples

💀 - not working

### minikube cluster

- [Start demo cluster](./minikube_empty_cluster/README.md) (minikube)
- [Run helloworld service via CLI](./minikube_helloworld/README.md) (minikube)
- [Run simple service from a local docker image](./minikube_local_image/README.md) (minikube)
- [Run service with persistent volumes](./minikube_shared_dirs/README.md) (minikube)
- [Setup local docker registry](./minikube_local_registry/README.md) (minikube)
- [Run example echo server behind kong-proxy](./minikube_kong_echo/README.md) (minikube, kong)
- 💀 [Run example gRPC service behind kong-proxy](./minikube_kong_grpc/README.md) (minikube, kong)
- 💀 [Run simple Python gRPC service behind kong-proxy](./minikube_kong_grpc_py/README.md) (minikube, kong)
- 💀 [Run simple Go gRPC service behind kong-proxy](./minikube_kong_grpc_go/README.md) (minikube, kong)
- [Prepare knative and kong](./minikube_knative_kong_prepare/README.md) (minikube, knative)
- [Run predefined helloworld with autoscale via knative](./minikube_knative_helloworld/README.md) (minikube, knative, kong)
- [Helloworld application in python using knative](./minikube_knative_helloworld_py/README.md) (minikube, knative, kong)
- [Run simple service from local image via knative](./minikube_knative_simple/README.md) (minikube, knative, kong)
- 💀 [Example of gRPC ping app on Go with knative](./minikube_knative_grpc_go/README.md) (minikube, knative, kong)
- 💀 [Example of gRPC service on Python with knative](./minikube_knative_grpc_py/README.md) (minikube, knative, kong)
- [Run simple service with nginx ingress controller](./minikube_nginx_ingress/README.md) (minikube, nginx)
- [Run gRPC service with nginx ingress controller](./minikube_nginx_grpc/README.md) (minikube, nginx)
- [Get statrted with Istio](./minikube_istio_getstart/README.md) (minikube, istio)
- [Run simple service with Istio](./minikube_istio_helloworld/README.md) (minikube, istio)
- 💀 [Run simple gRPC service with Istio](./minikube_istio_grpc/README.md) (minikube, istio)
- [Use Istio as ingress controller for knative](./minikube_istio_knative/README.md) (minikube, knative, istio)
- 💀 [Run HTTPS service with Istio gateway](./minikube_istio_https/README.md) (minikube, istio)
- [Run get-started example with Ambassador ingress](./minikube_ambassador_getstart/README.md) (minikube, ambassador)
- [Run simple echo service with Ambassador ingress](./minikube_ambassador_echo/README.md) (minikube, ambassador)
- 💀 [Run gRPC service with Ambassador ingress](./minikube_ambassador_grpc/README.md) (minikube, ambassador)

### kind cluster

- [Start demo cluster](./kind_empty_cluster) (kind)
- [Run simple echo service as NodePort](./kind_echo/README.md) (kind)
- [Run echo service and nginx ingress controller](./kind_echo_ingress_nginx/README.md) (kind, nginx)
- 💀 [Run helloworld from local image with nginx ingress controller](./kind_local_image/README.md) (kind, nginx)
- [Start sample gRPC service with nginx ingress controller](./kind_nginx_grpc/README.md) (kind, nginx)
- [Run simple service with kong ingress controller](./kind_kong_ingress/README/md) (kind, kong)
- 💀 [Run gRPC service with kong ingress controller](./kind_kong_grpc/README.md) (kind, kong)

### Various examples

- [Basic example for kustomize](./kustomize_0/README.md) (kustomize, shell)
- [Show pod lifecycle](./pod_lifecycle/README.md) (minikube, init-container)
- [Wait for service](./wait_for_service/README.md) (minikube, init-container)

## Tools

### kubectl

Install [kubectl](https://kubernetes.io/docs/reference/kubectl/overview/) binary:

```bash
curl -LO "https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl"
chmod +x ./kubectl
sudo mv ./kubectl /usr/local/bin/kubectl

kubectl version
Client Version: version.Info{Major:"1", Minor:"19", GitVersion:"v1.19.1", GitCommit:"206bcadf021e76c27513500ca24182692aabd17e", GitTreeState:"clean", BuildDate:"2020-09-09T11:26:42Z", GoVersion:"go1.15", Compiler:"gc", Platform:"linux/amd64"}
Server Version: version.Info{Major:"1", Minor:"19", GitVersion:"v1.19.0", GitCommit:"e19964183377d0ec2052d1f1fa930c4d7575bd50", GitTreeState:"clean", BuildDate:"2020-08-26T14:23:04Z", GoVersion:"go1.15", Compiler:"gc", Platform:"linux/amd64"}
```

### minikube

Install [minikube](https://kubernetes.io/docs/setup/learning-environment/minikube/) binary for trying minikube examples:

```bash
curl -Lo minikube https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
chmod +x minikube
sudo mv minikube /usr/local/bin

minikube version
minikube version: v1.13.0
commit: 0c5e9de4ca6f9c55147ae7f90af97eff5befef5f-dirty
```

### kind

Install [kind](https://kind.sigs.k8s.io/) binary for trying kind examples.

```bash
curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.8.1/kind-linux-amd64
chmod +x kind
sudo mv kind /usr/local/bin

kind version
kind v0.8.1 go1.14.2 linux/amd64
```

### grpcurl

We use [grpcurl](https://github.com/fullstorydev/grpcurl) for testing gRPC services.

```bash
curl -LO https://github.com/fullstorydev/grpcurl/releases/download/v1.7.0/grpcurl_1.7.0_linux_x86_64.tar.gz
tar -zxvf grpcurl_1.7.0_linux_x86_64.tar.gz -C .
chmod +x ./grpcurl
sudo mv grpcurl /usr/local/bin
rm grpcurl_1.7.0_linux_x86_64.tar.gz

grpcurl -version
grpcurl v1.7.0
```

### istio

Installation steps for Istio are described [here](https://istio.io/latest/docs/setup/getting-started/#download):

```bash
curl -L https://istio.io/downloadIstio | sh -
sudo mv istio-1.7.1/bin/istioctl /usr/local/bin

istioctl version
no running Istio pods in "istio-system"
1.7.1
```

## kubernetes spells

Get pod name by application name:

```bash
kubectl get pod -l app=$APP -o jsonpath='{.items[0].metadata.name}'
```

Get service node port:

```bash
kubectl get service $SERVICE --output='jsonpath="{.spec.ports[0].nodePort}"'
```

Show logs for a pod (`-p` show log for the previos crashed instance, it's very useful when pod gets stuck in the CrashLoopBackOff state):

```bash
kubectl logs $PODNAME -p
```
