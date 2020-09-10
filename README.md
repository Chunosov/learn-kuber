# kubernetes lessons and samples

## Sequence of samples (minikube)

- [Start demo cluster with minikube](./minikube_empty_cluster/README.md) (minikube)
- [Run helloworld service in minikube cluster via CLI](./minikube_helloworld/README.md) (minikube)
- [Run simple service in minikube cluster from a local docker image](./minikube_local_image/README.md) (minikube)
- [Run service with persistent volumes in minikube cluster](./minikube_shared_dirs/README.md) (minikube)
- [Setup local docker registry with minikube](./minikube_local_registry/README.md) (minikube)
- [Run example echo server behind kong-proxy on minikube](./minikube_kong_echo/README.md) (minikube, kong)
- [Run example gRPC service behind kong-proxy on minikube](./minikube_kong_grpc/README.md) (minikube, kong)
- [Run simple Python gRPC service behind kong-proxy on minikube](./minikube_kong_grpc_py/README.md) (minikube, kong)
- [Run simple Go gRPC service behind kong-proxy on minikube](./minikube_kong_grpc_go/README.md) (minikube, kong, kong)
- [Prepare knative and kong on minikube](./minikube_knative_kong_prepare/README.md) (minikube, knative)
- [Run predefined helloworld with autoscale via knative in minikube](./minikube_knative_helloworld/README.md) (minikube, knative, kong)
- [Helloworld application in python using knative on minikube](./minikube_knative_helloworld_py/README.md) (minikube, knative, kong)
- [Run simple service from local image via knative in minikube](./minikube_knative_simple/README.md) (minikube, knative, kong)
- [Example of gRPC ping app on Go with minikube and knative](./minikube_knative_grpc_go/README.md) (minikube, knative, kong)
- [Example of gRPC service on Python with minikube and knative](./minikube_knative_grpc_py/README.md) (minikube, knative, kong)

## Sequence of samples (kind)

- [Start demo cluster with kind](./kind_empty_cluster) (kind)

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
