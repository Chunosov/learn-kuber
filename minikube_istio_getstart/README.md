# Get statrted with Istio on minikube cluster

This is a reproduction of the official [get-started example](https://istio.io/latest/docs/setup/getting-started/) to make sure all is working with predefinied application.

Start a cluster and [install istio](https://istio.io/latest/docs/setup/getting-started/#install):

```bash
minikube start --driver=virtualbox
istioctl install --set profile=demo
kubectl label namespace default istio-injection=enabled
```

That's what we have now:

```bash
kubectl get namespace
NAME              STATUS   AGE
default           Active   8m17s
istio-system      Active   6m58s
kube-node-lease   Active   8m18s
kube-public       Active   8m18s
kube-system       Active   8m18s

kubectl get service -n istio-system
NAME                   TYPE           CLUSTER-IP       EXTERNAL-IP   PORT(S)                                                                      AGE
istio-egressgateway    ClusterIP      10.102.64.13     <none>        80/TCP,443/TCP,15443/TCP                                                     6m32s
istio-ingressgateway   LoadBalancer   10.110.146.114   <pending>     15021:30548/TCP,80:32446/TCP,443:30156/TCP,31400:30030/TCP,15443:32114/TCP   6m32s
istiod                 ClusterIP      10.111.242.3     <none>        15010/TCP,15012/TCP,443/TCP,15014/TCP,853/TCP                                7m10s

kubectl get deployment -n istion-system
No resources found in istion-system namespace.

kubectl get pod -n istion-system
No resources found in istion-system namespace.
```

Delploy sample application (it is from the Istio samples dir but we have yaml files reproduced here):

```bash
kubectl apply -f bookinfo.yaml

kubectl get service
NAME          TYPE        CLUSTER-IP      EXTERNAL-IP   PORT(S)    AGE
details       ClusterIP   10.103.131.15   <none>        9080/TCP   11s
kubernetes    ClusterIP   10.96.0.1       <none>        443/TCP    15m
productpage   ClusterIP   10.98.150.121   <none>        9080/TCP   11s
ratings       ClusterIP   10.105.229.40   <none>        9080/TCP   11s
reviews       ClusterIP   10.96.93.76     <none>        9080/TCP   11s
```

Repeat the `get pod` command several times, pods get ready after a couple of minutes. As each pod becomes ready, the Istio sidecar will be deployed along with it:

```bash
kubectl get pod
NAME                              READY   STATUS    RESTARTS   AGE
details-v1-79c697d759-96947       2/2     Running   0          4m45s
productpage-v1-65576bb7bf-cf7vr   2/2     Running   0          4m45s
ratings-v1-7d99676f7f-rhp8w       2/2     Running   0          4m44s
reviews-v1-987d495c-lvz5n         2/2     Running   0          4m45s
reviews-v2-6c5bf657cf-rm9w5       2/2     Running   0          4m45s
reviews-v3-5f7b9f4f77-tr56k       2/2     Running   0          4m45s
```

Cast some spell to check if our app is working, calling it inside of the cluster:

```bash
kubectl exec "$(kubectl get pod -l app=ratings -o jsonpath='{.items[0].metadata.name}')" -c ratings -- curl -s productpage:9080/productpage | grep -o "<title>.*</title>"
<title>Simple Bookstore App</title>
```

The guys claim it means that all ok, let's trust them as we don't have many options at this stage.

Then we need to make the app available outside of the cluster:

```bash
kubectl apply -f bookinfo-gateway.yaml

export INGRESS_PORT_HTTP=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.spec.ports[?(@.name=="http2")].nodePort}')
export INGRESS_PORT_HTTPS=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.spec.ports[?(@.name=="https")].nodePort}')
echo $INGRESS_PORT_HTTP,$INGRESS_PORT_HTTPS
32446,30156
export APP_URL=http://$(minikube ip):$INGRESS_PORT_HTTP/productpage
echo $APP_URL
http://192.168.99.114:32446/productpage
```

Then navigate to this address manually or right from terminal:

```bash
firefox $APP_URL
# or
google-chrome $APP_URL
```

It [was said](https://istio.io/latest/docs/setup/getting-started/#determining-the-ingress-ip-and-ports) we have to run `minikube tunnel` but it works without. minikube [docs](https://minikube.sigs.k8s.io/docs/handbook/accessing/#using-minikube-tunnel) say the tunnel makes LoadBalancer services accessible on some external IP. Previous examples prove we can access such services at `minikube ip` as if they were NodePort. This is exactly how we access our app in this example, so we don't need `minikube tunnel` here.

> tunnel creates a route to services deployed with type LoadBalancer and sets their Ingress to their ClusterIP
