# Use Istio as ingress controller for knative on minikube cluster

Start a cluster:

```bash
minikube start --driver=virtualbox
```

Install [knative](https://knative.dev/docs/install/any-kubernetes-cluster/) custom resource definitions and core components:

```bash
kubectl apply --filename https://github.com/knative/serving/releases/download/v0.17.0/serving-crds.yaml
kubectl apply --filename https://github.com/knative/serving/releases/download/v0.17.0/serving-core.yaml
```

Install Istio for knative as described [here](https://knative.dev/docs/install/installing-istio/) (make sure you have istioctl [installed](../README.md#istioctl)):

```bash
istioctl install -f istio-minimal-operator.yaml
kubectl apply --filename https://github.com/knative/net-istio/releases/download/v0.17.0/release.yaml
```

Extract service ports:

```bash
kubectl --namespace istio-system get service istio-ingressgateway
NAME                   TYPE           CLUSTER-IP      EXTERNAL-IP   PORT(S)                                                      AGE
istio-ingressgateway   LoadBalancer   10.106.178.24   <pending>     15021:30779/TCP,80:32317/TCP,443:30794/TCP,15443:31847/TCP   2m42s

export PROXY_PORT_HTTP=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.spec.ports[?(@.name=="http2")].nodePort}')
export PROXY_PORT_HTTPS=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.spec.ports[?(@.name=="https")].nodePort}')
echo $PROXY_PORT_HTTP,$PROXY_PORT_HTTPS
32317,30794
```

Deploy [sample service](https://knative.dev/docs/serving/getting-started-knative-app/):

```bash
kubectl apply -f service.yaml
```

See there arev no pods before we call the service:

```bash
kubectl get pod
No resources found in default namespace.
```

Actually there can be some pods right after we have deployed the service, but they terminate after some time:

```bash
kubectl get pod
NAME                                              READY   STATUS              RESTARTS   AGE
helloworld-go-qnph8-deployment-5646c79d88-4fl48   0/2     ContainerCreating   0          68s
helloworld-go-qnph8-deployment-b77995578-9pb44    0/2     ContainerCreating   0          69s

kubectl get pod
NAME                                              READY   STATUS        RESTARTS   AGE
helloworld-go-qnph8-deployment-5646c79d88-4fl48   1/2     Terminating   0          3m18s

kubectl get pod
No resources found in default namespace.
```

Call the service and see how autoscaler starts new pods:

```bash
kubectl get ksvc
NAME            URL                                        LATESTCREATED         LATESTREADY           READY   REASON
helloworld-go   http://helloworld-go.default.example.com   helloworld-go-qnph8   helloworld-go-qnph8   True

curl -H "Host: helloworld-go.default.example.com" http://$(minikube ip):$PROXY_PORT_HTTP
Hello Go Sample v1!

kubectl get pod
NAME                                              READY   STATUS    RESTARTS   AGE
helloworld-go-qnph8-deployment-5646c79d88-s4wxb   2/2     Running   0          8s
```

After some time all pods terminate:

```bash
kubectl get pod
No resources found in default namespace.
```
