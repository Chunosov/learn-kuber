# Tunneling UPD into cluster

Kubectl doesn't support port-forwarding for UDP traffic because it requires changes in container runtime. There is the [issue](https://github.com/kubernetes/kubernetes/issues/47862) about that. And there is also the workaround described in the issue.

```bash
kind create cluster
```

Deploy the test container simulating an application dealing with UDP:

```bash
kubectl apply -f testserver.yaml
```

Attach to it an ephemeral container running `socat` to forward TCP to UDP. Port 80 in the `socat` command is the port where the UDP server is listening (see `testserver.yaml`). The port 9999 is a random port that we need to use later with port-forward:

```bash
kubectl debug testserver --image=alpine/socat -- socat tcp-listen:9999,fork,reuseaddr,max-children=1 UDP4:127.0.0.1:80
```

A new container appears in the app:

```bash
kubectl describe pod testserver

...
Containers:
  agnhost:
    Image:         registry.k8s.io/e2e-test-images/agnhost:2.39
    Port:          80/TCP
    Host Port:     0/TCP
    Args:
      netexec
      --udp-port=80
Ephemeral Containers:
  debugger-5lzzt:
    Image:         alpine/socat
    Port:          <none>
    Host Port:     <none>
    Command:
      socat
      tcp-listen:9999,fork,reuseaddr,max-children=1
      UDP4:127.0.0.1:80
...
```

Forward a port from the host to the port we indicated socat to run. Here it is the same port on the host and in the socat container 9999:

```bash
kubectl port-forward testserver 9999
Forwarding from 127.0.0.1:9999 -> 9999
Forwarding from [::1]:9999 -> 9999
```

Run another `socat` as a proxy on the host to forward from an UDP port to the forwarded TCP port:

```bash
socat UDP4-LISTEN:9010 tcp4:127.0.0.1:9999,fork
```

Now we can send UDP traffic to the local port and it will be tunneled into the cluster over TCP:

```bash
echo "hostname" | nc -u 127.0.0.1 9010
testserver
```

```
|----------- host -------------|----- cluster -----|
|                              |                   |
|  client    |   socat   | apiserver | socat | app |
|            |   proxy   |  forward  | proxy |     |
???? -----> 9010 -----> 9999 -----> 9999 -------> 80
      UDP        TCP          TCP          UDP
```
