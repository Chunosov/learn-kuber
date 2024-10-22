# Run helloworld service in WSL k3s cluster

Make sure you have a WSL cluster [running](../wsl_start/README.md).

Make a deployment. This command runs container in a pod, but it is unavailable from outside of cluster:

```bash
kubectl apply -f echo.yaml
service/echo created
deployment.apps/echo created

kubectl get deployments
NAME   READY   UP-TO-DATE   AVAILABLE   AGE
echo   1/1     1            1           13s

kubectl get pods
NAME                       READY   STATUS    RESTARTS   AGE
echo-67446c6864-f2dg8      1/1     Running   0          35s

kubectl get service
NAME         TYPE        CLUSTER-IP   EXTERNAL-IP   PORT(S)          AGE
echo         NodePort    10.43.7.6    <none>        8080:32620/TCP   55s
```

Service type NodePort means that a port from a node will be exposed at the cluster address. Here 8080 is port on the node, 32620 is respective port on the cluster address.

*TODO:* How to expose a cluster address?

Use port forwarding to get access into the cluster:

```
kubectl port-forward service/echo 8077:8080
Forwarding from 127.0.0.1:8077 -> 8080
Forwarding from [::1]:8077 -> 8080
```

And connect to the echo app:

```bash
curl localhost:8077

Hostname: echo-67446c6864-f2dg8

Pod Information:
        -no pod information available-

Server values:
        server_version=nginx: 1.13.3 - lua: 10008

Request Information:
        client_address=127.0.0.1
        method=GET
        real path=/
        query=
        request_version=1.1
        request_scheme=http
        request_uri=http://localhost:8080/

Request Headers:
        accept=*/*
        host=localhost:8077
        user-agent=curl/8.0.1

Request Body:
        -no body in request-
```

It works.
