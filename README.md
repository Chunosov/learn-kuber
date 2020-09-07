# kubernetes lessons and samples

## Sequence of samples

- [minikube_helloworld](./minikube_helloworld/README.md)
- [minikube_local_image](./minikube_local_image/README.md)
- [minikube_shared_dirs](./minikube_shared_dirs/README.md)
- [minikube_knative_helloworld](./minikube_knative_helloworld/README.md)
- [minikube_knative_helloworld_py](./minikube_knative_helloworld_py/README.md)
- [minikube_local_registry](./minikube_local_registry/README.md)
- [minikube_knative_simple](./minikube_knative_simple/README.md)

## minikube

Install kubectl binary:

```bash
curl -LO "https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl"
chmod +x ./kubectl
sudo mv ./kubectl /usr/local/bin/kubectl
kubectl version --client
```

Install [minikube](https://kubernetes.io/docs/setup/learning-environment/minikube/) binary:

```bash
curl -Lo minikube https://storage.googleapis.com/minikube/releases/latest/minikube-linux-amd64
chmod +x minikube
sudo mv minikube /usr/local/bin
minikube version
```

Install [kvm](https://help.ubuntu.com/community/KVM/Installation) [driver](https://minikube.sigs.k8s.io/docs/drivers/kvm2/) for minikube.

Start demo cluster:

```bash
minikube start --driver=kvm2

😄  minikube v1.12.3 on Ubuntu 18.04
    ▪ MINIKUBE_ACTIVE_DOCKERD=minikube
✨  Using the kvm2 driver based on user configuration
👍  Starting control plane node minikube in cluster minikube
🔥  Creating kvm2 VM (CPUs=2, Memory=3900MB, Disk=20000MB) ...
🐳  Preparing Kubernetes v1.18.3 on Docker 19.03.12 ...
🔎  Verifying Kubernetes components...
🌟  Enabled addons: default-storageclass, storage-provisioner
🏄  Done! kubectl is now configured to use "minikube"

minikube status

minikube
type: Control Plane
host: Running
kubelet: Running
apiserver: Running
kubeconfig: Configured
```

You can connect to the running VM via `minikube ssh`. Or open VM in Virtual Machine Manager application, double-click VM to open a terminal and login as root without password.

### Shared dirs

With `kvm2` driver you can't use shared directories. If they needed, use VirualBox driver, it automatically bind host's `/home` directory to `/hosthome` directory inside of VM (but it is [not configurable](https://kubernetes.io/docs/setup/learning-environment/minikube/#mounted-host-folders) currently).

```bash
minikube start --driver=virtualbox
```

Check what driver minikube currently runs:

```bash
minikube profile list

|----------|------------|---------|----------------|------|---------|---------|
| Profile  | VM Driver  | Runtime |       IP       | Port | Version | Status  |
|----------|------------|---------|----------------|------|---------|---------|
| minikube | virtualbox | docker  | 192.168.99.102 | 8443 | v1.18.3 | Running |
|----------|------------|---------|----------------|------|---------|---------|
```

## knative

Install [knative](https://knative.dev/docs/install/any-kubernetes-cluster/) custom resource definitions (CRDs) and serving component:

```bash
kubectl apply --filename https://github.com/knative/serving/releases/download/v0.17.0/serving-crds.yaml
kubectl apply --filename https://github.com/knative/serving/releases/download/v0.17.0/serving-core.yaml
```

See new services:

```bash
kubectl get service -n knative-serving
NAME                TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)                           AGE
activator-service   ClusterIP   10.111.170.231   <none>        9090/TCP,8008/TCP,80/TCP,81/TCP   36s
autoscaler          ClusterIP   10.96.107.228    <none>        9090/TCP,8008/TCP,8080/TCP        36s
controller          ClusterIP   10.101.63.37     <none>        9090/TCP,8008/TCP                 36s
webhook             ClusterIP   10.103.137.140   <none>        9090/TCP,8008/TCP,443/TCP         36s
```

knative requires something called "network layer", we'll use [kong](https://docs.konghq.com/2.1.x/kong-for-kubernetes/using-kong-for-kubernetes/) for that:

```bash
kubectl apply --filename https://raw.githubusercontent.com/Kong/kubernetes-ingress-controller/0.9.x/deploy/single/all-in-one-dbless.yaml
# or
kubectl apply -f https://bit.ly/k4k8s
```

See new services:

```bash
kubectl get service -n kong
NAME                      TYPE           CLUSTER-IP      EXTERNAL-IP   PORT(S)                      AGE
kong-proxy                LoadBalancer   10.104.32.3     <pending>     80:32131/TCP,443:30891/TCP   42s
kong-validation-webhook   ClusterIP      10.97.241.236   <none>        443/TCP                      42s
```

kong-proxy service shows external ip as "pending" because LoadBalancer service type doesn't work in minikube. But the service is still available at cluster address `$(minikube ip):32131` (for HTTP, the second port is for HTTPs).

Use kong as [ingress](https://kubernetes.io/docs/concepts/services-networking/ingress/) controller for the cluster:

```bash
kubectl patch configmap/config-network --namespace knative-serving --type merge --patch '{"data":{"ingress.class":"kong"}}'
```

See also [Kong official guides on Ingress Controller](https://github.com/Kong/kubernetes-ingress-controller/tree/main/docs/guides) and [Using Kong with Knative](https://github.com/Kong/kubernetes-ingress-controller/blob/main/docs/guides/using-kong-with-knative.md) in particular.