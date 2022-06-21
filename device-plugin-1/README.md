# Device Plugin Example

Simplest example for kubernetes [Device Plugin](https://kubernetes.io/docs/concepts/extend-kubernetes/compute-storage-net/device-plugins/). The example provides a plugin for imaginary devices. It shows basic plugin implementation and demonstrates how kubernetes uses plugins for scheduling pods.

Use multinode [kind](https://kind.sigs.k8s.io/) cluster for this demo:

```bash
make cluster
```

Build and deploy the driver:

```bash
make build
make start
```

Check that nodes advertise custom devices:

```bash
kubectl get node kind-worker -o yaml
...
status:
  allocatable:
    myplugin.io/mydevice-1: "1"
    myplugin.io/mydevice-2: "0"
    myplugin.io/mydevice-3: "0"
  capacity:
    myplugin.io/mydevice-1: "1"
    myplugin.io/mydevice-2: "0"
    myplugin.io/mydevice-3: "0"
```

```bash
kubectl get node kind-worker2 -o yaml
...
status:
  allocatable:
    myplugin.io/mydevice-1: "0"
    myplugin.io/mydevice-2: "1"
    myplugin.io/mydevice-3: "3"
  capacity:
    myplugin.io/mydevice-1: "0"
    myplugin.io/mydevice-2: "1"
    myplugin.io/mydevice-3: "3"
```

Deploy example pods and check how they scheduled to nodes:

```bash
kubectl apply -f pods.yaml
```

Pod `dp1-1` is scheduled on the `kind-worker` node because that provides `mydevice-1`. Pod `dp1-2` is scheduled on `kind-worker2` because that provides `mydevice-2`:

```bash
kubectl get pods -o wide

NAME        READY   STATUS    RESTARTS   AGE     IP           NODE           NOMINATED NODE   READINESS GATES
dp1-1       1/1     Running   0          3m4s    10.244.1.6   kind-worker    <none>           <none>
dp1-2       1/1     Running   0          82s     10.244.3.6   kind-worker2   <none>           <none>
dp1-3       0/1     Pending   0          82s     <none>       <none>         <none>           <none>
```

```bash
kubectl logs dp1-1

Device vars:
DEV_1_ID_0='/node1/dev1/1'
```

```bash
kubectl logs dp1-2

Device vars:
DEV_2_ID_0='/node2/dev2/1'
```

Pod `dp1-3` is not scheduled because it requires two kind of devices and `kind-worker2` provides them both but one of them is already allocated by `dp1-2`:

```bash
kubectl describe pod dp1-3
...
0/3 nodes are available: 1 Insufficient myplugin.io/mydevice-3, 1 node(s) had taint {node-role.kubernetes.io/master: }, that the pod didn't tolerate, 2 Insufficient myplugin.io/mydevice-2.
```

So if we stop the pod `dp1-2` to free the device, the pod `dp1-3` get finally scheduled:

```bash
kubectl delete pod dp1-2
```

```bash
kubectl get pods -o wide

NAME        READY   STATUS    RESTARTS   AGE     IP           NODE           NOMINATED NODE   READINESS GATES
dp1-1       1/1     Running   0          11m     10.244.1.6   kind-worker    <none>           <none>
dp1-3       1/1     Running   0          9m18s   10.244.3.7   kind-worker2   <none>           <none>
```

```bash
kubectl logs dp1-3

Device vars:
DEV_2_ID_0='/node2/dev2/1'
DEV_3_ID_0='/node2/dev3/2'
DEV_3_ID_1='/node2/dev3/3'
```
