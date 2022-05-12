# Sample controller for built-in resource

An example of kubernetes controller watching pods.

Start demo cluster with kind, kubeconfig should be available for the controller:

```bash
kind create cluster
```

Start our demo pod controller:

```bash
go run .
```

You should see that some system pods were already catched by the controller:

```log
Pod added: kube-system/etcd-kind-control-plane
...
```

Let's add another pod:

```bash
kubectl apply -f pod.yaml
```

We should see that controller has handled the event:

```log
Pod added: default/echo
Pod updated: default/echo
...
```

There is also a bunch of update events caused by the system.

Lets delete the pod:

```bash
kubectl delete pod echo
```

After a couple of updates the pod gets deleted:

```log
...
Pod updated: default/echo
Pod updated: default/echo
Pod deleted: default/echo
```
