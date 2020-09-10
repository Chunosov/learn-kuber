# Start demo cluster with kind

Make sure you have kind [installed](../README.md#kind).

## Default cluster

```bash
kind create cluster

creating cluster "kind" ...
 ✓ Ensuring node image (kindest/node:v1.18.2) 🖼
 ✓ Preparing nodes 📦
 ✓ Writing configuration 📜
 ✓ Starting control-plane 🕹️
 ✓ Installing CNI 🔌
 ✓ Installing StorageClass 💾
Set kubectl context to "kind-kind"
You can now use your cluster with:

kubectl cluster-info --context kind-kind

Not sure what to do next? 😅  Check out https://kind.sigs.k8s.io/docs/user/quick-start/
```

You can see errors like:

```
ERROR: failed to create cluster: failed to ensure docker network: command "docker network create -d=bridge --ipv6 --subnet fc00:8866:27d0:bd7e::/64 kind" failed with error: exit status 1
Command Output: Error response from daemon: could not find an available, non-overlapping IPv4 address pool among the defaults to assign to the network
```

Possible [solutions](https://stackoverflow.com/questions/43720339/docker-error-could-not-find-an-available-non-overlapping-ipv4-address-pool-am):

- Check if there are too many networks already created `docker network ls` and remove all unused networks `docker network prune`

- Stop openvpn if it is running `sudo service openvpn stop`. People say the same happens with Cisco VPN. There should be a better solution to keep project working alongside with having vpn started.

There is our cluster:

```bash
docker images | grep kind
kindest/node                                    <none>              de6eb7df13da        4 months ago        1.25GB

docker ps
CONTAINER ID        IMAGE                  COMMAND                  CREATED             STATUS              PORTS                       NAMES
11608ea3fcd3        kindest/node:v1.18.2   "/usr/local/bin/entr…"   6 minutes ago       Up 6 minutes        127.0.0.1:43567->6443/tcp   kind-control-plane
```

Destroy the cluster:

```bash
kind delete cluster
Deleting cluster "kind" ...
```

## Named cluster

Make named cluster:

```bash
kind create cluster --name ch-test0
...
You can now use your cluster with:

kubectl cluster-info --context kind-ch-test0
```

It says we have to scaffold kubectl with `--context` but indeed the kind sets this new context as current so we can just call kubectl:

```bash
kubectl config get-contexts
CURRENT   NAME            CLUSTER         AUTHINFO        NAMESPACE
*         kind-ch-test0   kind-ch-test0   kind-ch-test0
          minikube        minikube        minikube

```

```bash
kubectl cluster-info
Kubernetes master is running at https://127.0.0.1:39547
KubeDNS is running at https://127.0.0.1:39547/api/v1/namespaces/kube-system/services/kube-dns:dns/proxy

To further debug and diagnose cluster problems, use 'kubectl cluster-info dump'.
```

```bash
kubectl get node
NAME                     STATUS   ROLES    AGE     VERSION
ch-test0-control-plane   Ready    master   6m24s   v1.18.2
```

## Multi-node cluster

There can be several clusters in the same time. Make another one having several nodes. Set of nodes is defined in config:

```bash
kind create cluster --name ch-test1 --config cluster.yaml
```

```bash
kubectl get node
NAME                     STATUS   ROLES    AGE     VERSION
ch-test1-control-plane   Ready    master   3m28s   v1.18.2
ch-test1-worker          Ready    <none>   2m52s   v1.18.2
ch-test1-worker2         Ready    <none>   2m51s   v1.18.2
```

```bash
docker ps
CONTAINER ID        IMAGE                  COMMAND                  CREATED             STATUS              PORTS                       NAMES
f85fa319dea6        kindest/node:v1.18.2   "/usr/local/bin/entr…"   2 minutes ago       Up 2 minutes                                    ch-test1-worker
2400c77f4f36        kindest/node:v1.18.2   "/usr/local/bin/entr…"   2 minutes ago       Up 2 minutes                                    ch-test1-worker2
b49f5ce88a18        kindest/node:v1.18.2   "/usr/local/bin/entr…"   2 minutes ago       Up 2 minutes        127.0.0.1:39803->6443/tcp   ch-test1-control-plane
899e00f3a9a0        kindest/node:v1.18.2   "/usr/local/bin/entr…"   20 minutes ago      Up 19 minutes       127.0.0.1:39547->6443/tcp   ch-test0-control-plane
```

## Use several clusters

See all existed clusters:

```bash
kind get clusters
ch-test0
ch-test1
```

Newly created cluster became current:

```bash
kubectl config get-contexts
CURRENT   NAME            CLUSTER         AUTHINFO        NAMESPACE
          kind-ch-test0   kind-ch-test0   kind-ch-test0
*         kind-ch-test1   kind-ch-test1   kind-ch-test1
          minikube        minikube        minikube
```


We can use `--context` to work with non-current cluster:

```bash
kubectl get node --context kind-ch-test0
NAME                     STATUS   ROLES    AGE   VERSION
ch-test0-control-plane   Ready    master   27m   v1.18.2
```

Or we can switch the context permanently:

```bash
kubectl config use-context kind-ch-test0
```

Then kubectl can be used without `--context`:

```bash
kubectl get node
NAME                     STATUS   ROLES    AGE   VERSION
ch-test0-control-plane   Ready    master   32m   v1.18.2
```

Named clusters should be deleted by name:

```bash
kind delete cluster
Deleting cluster "kind" ...   # actually we don't have this already

kind delete cluster --name ch-test0
Deleting cluster "ch-test0" ...

kind delete cluster --name ch-test1
Deleting cluster "ch-test1" ...
```
