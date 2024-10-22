# Run helloworld service in WSL KIND cluster

Make sure you have a WSL instance [prepared](../wsl_start_kind/README.md).

In WSL instance terminal:

```bash
kind create cluster --name demo

kubectl create deployment hello --image=k8s.gcr.io/echoserver:1.10

kubectl expose deployment hello --type=NodePort --port=8080
```

**WSL not always works as honest Linux**

On kind in normal Linux installation we could use node address to connect to the service:

```bash
kubectl get service
NAME          TYPE           CLUSTER-IP      EXTERNAL-IP   PORT(S)          AGE
hello         NodePort       10.96.193.27    <none>        8080:31903/TCP   6m6s

NAME                 STATUS   ROLES           AGE     VERSION   INTERNAL-IP   EXTERNAL-IP   OS-IMAGE                         KERNEL-VERSION                       CONTAINER-RUNTIME
demo-control-plane   Ready    control-plane   4h31m   v1.31.0   172.18.0.2    <none>        Debian GNU/Linux 12 (bookworm)   5.15.153.1-microsoft-standard-WSL2   containerd://1.7.18

curl 172.18.0.2:31903
# it works
```

This doesn't work on WSL instance from both Linux and Windows terminals:

```
curl 172.18.0.2:31903
curl: (28) Failed to connect to 172.18.0.2 port 31903 after 134162 ms: Couldn't connect to server
```

Use port-forwarding to access the service:

```bash
kubectl port-forward service/hello 8077:8080
```

Now it can be accessed:

```bash
curl 127.0.0.1:8077

Hostname: hello-7fd68756fd-hvpp7

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
        request_uri=http://127.0.0.1:8080/

Request Headers:
        accept=*/*
        host=127.0.0.1:8077
        user-agent=curl/8.5.0

Request Body:
        -no body in request- 
```

When port forwarding is started in WSL terminal, the service can be accessed from both Windows and WSL. When port forwarding is started in Windows terminal, the service can only be accessed from Windows.