# Wait for service to be ready

https://vadosware.io/post/so-you-need-to-wait-for-some-kubernetes-resources/
https://stackoverflow.com/questions/51079849/kubernetes-wait-for-other-pod-to-be-ready


Start example cluster:

```bash
minikube start --driver=kvm2
```

Deploy demo app:

```bash
kubectl apply -f deploy.yaml
```

Wait for client pod gets to completed state and see the logs:

```bash
minikube ssh
cat /tmp/echo/log
```

```log
1647344955: SERVICE starting
1647344963: CLIENT waiting for service...
1647344969: CLIENT waiting for service...
1647344975: CLIENT waiting for service...
1647344981: CLIENT waiting for service...
1647344985: SERVICE started
1647344991: CLIENT
Hello world!
1647345004: CLIENT
Hello world!
1647345022: CLIENT
Hello world!
1647345058: CLIENT
Hello world!
1647345129: CLIENT
Hello world!
1647345231: CLIENT
Hello world!
```
