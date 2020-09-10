# Start demo cluster with minikube

Make sure you have minikube [installed](../README.md#minikube). Install [kvm](https://help.ubuntu.com/community/KVM/Installation) [driver](https://minikube.sigs.k8s.io/docs/drivers/kvm2/) for minikube.

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

```bash
 minikube ssh
                         _             _
            _         _ ( )           ( )
  ___ ___  (_)  ___  (_)| |/')  _   _ | |_      __
/' _ ` _ `\| |/' _ `\| || , <  ( ) ( )| '_`\  /'__`\
| ( ) ( ) || || ( ) || || |\`\ | (_) || |_) )(  ___/
(_) (_) (_)(_)(_) (_)(_)(_) (_)`\___/'(_,__/'`\____)

$ docker ps
# shows a list of containers run on VM
```

#### Shared dirs

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
