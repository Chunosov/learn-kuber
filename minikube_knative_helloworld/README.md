# Run predefined helloworld with autoscale via knative in minikube

Start demo cluster:

```bash
minikube start --driver=virtualbox
```

Install knative and kong as described [here](../minikube_knative_kong_prepare/README.md).

Install predefined helloworld as knative service:

```bash
echo "
apiVersion: serving.knative.dev/v1
kind: Service
metadata:
  name: helloworld-go
  namespace: default
spec:
  template:
    spec:
      containers:
        - image: gcr.io/knative-samples/helloworld-go
          env:
            - name: TARGET
              value: Go Sample v1
" | kubectl apply -f -
```

`TARGET` is the environment variable printed out by the sample app. There is a [doc](https://knative.dev/v0.15-docs/serving/samples/hello-world/helloworld-go/) about how this image built.

Here is what we have now:

```bash
kubectl get ksvc
NAME            URL                                        LATESTCREATED         LATESTREADY   READY     REASON
helloworld-go   http://helloworld-go.default.example.com   helloworld-go-xs8j6                 Unknown   RevisionMissing

kubectl get services
NAME                          TYPE           CLUSTER-IP       EXTERNAL-IP                         PORT(S)                             AGE
helloworld-go                 ExternalName   <none>           helloworld-go.default.example.com   <none>                              26s
helloworld-go-xs8j6           ClusterIP      10.100.224.104   <none>                              80/TCP                              2m31s
helloworld-go-xs8j6-private   ClusterIP      10.102.74.151    <none>                              80/TCP,9090/TCP,9091/TCP,8022/TCP   2m31s
kubernetes                    ClusterIP      10.96.0.1        <none>                              443/TCP                             7m12s

```

`ksvc` is a knative object type, it manages `helloworld-go*` services. Use `kubectl delete ksvc helloworld-go` to delete all that stuff above, `delete service` will not work.

See that there are not pods running before we make a request:

```bash
kubectl get pods
No resources found in default namespace.
```

The service is available though proxy:

```bash
curl -v -H "Host: helloworld-go.default.example.com" http://$(minikube ip):32526
...
X-Kong-Upstream-Latency: 1752
...
Hello Go Sample v1!
```

We have to pass the Host header explicitly to help kong-proxy resolve a route to our service. In real life DNS does this job (see [this diagram](../minikube_knative_kong_prepare/README.md#Proxy-overview)).

It show a big latency because service pod gets running:

```bash
kubectl get pods
NAME                                              READY   STATUS    RESTARTS   AGE
helloworld-go-xs8j6-deployment-595ddf65f7-r6gfm   2/2     Running   0          55s
```

But it runs much faster in the next time:

```bash
curl -v -H "Host: helloworld-go.default.example.com" http://$(minikube ip):32526
...
< X-Kong-Upstream-Latency: 1
...
Hello Go Sample v1!
```

And after short time knative stops service pods if they are not used:

```bash
kubectl get pods
NAME                                              READY   STATUS    RESTARTS   AGE
helloworld-go-xs8j6-deployment-595ddf65f7-zxslt   2/2     Running   0          107s

kubectl get pods
NAME                                              READY   STATUS        RESTARTS   AGE
helloworld-go-xs8j6-deployment-595ddf65f7-zxslt   2/2     Terminating   0          111s

kubectl get pods
No resources found in default namespace.
```

So it works.
