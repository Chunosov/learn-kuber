# Using load balancer in kind cluster

Reproduction of the official [example](./https://kind.sigs.k8s.io/docs/user/loadbalancer/).

Make a kind cluster:

```
kind create cluster --name demo
```

Download and run the kind load balancer in a separate terminal:

```bash
wget https://github.com/kubernetes-sigs/cloud-provider-kind/releases/download/v0.4.0/cloud-provider-kind_0.4.0_linux_amd64.tar.gz
tar -xf cloud-provider-kind_0.4.0_linux_amd64.tar.gz
./cloud-provider-kind
```

It will show several errors like:

```log
I1022 08:03:42.663675   11246 controller.go:174] probe HTTP address https://demo-control-plane:6443
I1022 08:03:42.669141   11246 controller.go:177] Failed to connect to HTTP address https://demo-control-plane:6443: Get "https://demo-control-plane:6443": dial tcp: lookup demo-control-plane on 127.0.0.53:53: server misbehaving
```

This seems ok, just wait a bit for lines appear:

```log
I1022 08:03:53.156518   11246 controller.go:84] Creating new cloud provider for cluster demo
I1022 08:03:53.174664   11246 controller.go:91] Starting cloud controller for cluster demo
```

**NB:** `cloud-provider-kind` works even on WSL instance, but there it must be run under `sudo`. On regular VM it works without sudo.

Deploy a LoadBalancer service with two backend pods:

```bash
# provided that the target system is being interacted over ssh, do in the host terminal:
scp echo.yaml {REMOTE_USER}@{REMOTE_ADDRESS}:echo.yaml

# then on the target system in the home of the {REMOTE_USER}:
kubectl apply -f echo.yaml

kubectl get pods
NAME       READY   STATUS    RESTARTS   AGE
echo-bar   1/1     Running   0          7s
echo-foo   1/1     Running   0          8s
```

Load balancer monitors services and gives each service having the `LoadBalancer` type an IP address :

```bash
LB_IP=$(kubectl get svc/echo -o=jsonpath='{.status.loadBalancer.ingress[0].ip}')

echo $LB_IP
172.18.0.3
```

Note that it is an IP in the same range as cluster nodes have:

```bash
kubectl get nodes -owide
NAME                 STATUS   ROLES           AGE   VERSION   INTERNAL-IP   EXTERNAL-IP   OS-IMAGE                         KERNEL-VERSION     CONTAINER-RUNTIME
demo-control-plane   Ready    control-plane   16h   v1.31.0   172.18.0.2    <none>        Debian GNU/Linux 12 (bookworm)   6.8.0-47-generic   containerd://1.7.18
```

We'd run the service as NodePort and accessed it on 172.18.0.2.

With load balancer we access in on 172.18.0.3:

```bash
for _ in {1..10}; do
  curl ${LB_IP}:5678
done

bar
foo
bar
foo
bar
foo
bar
bar
bar
foo
```

Should output "foo" and "bar" on separate lines.
