# Start demo cluster on WSL

Make sure you have WSL and Docker [prepared](../wsl/README.md). 

Install kind on the WSL instance:

```bash
curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.24.0/kind-linux-amd64
chmod +x kind
sudo mv kind /usr/local/bin
```

It works [the same way](../kind_empty_cluster/README.md) as on true Linux installation:

```bash
kind create cluster
```

Copy kubeconfig from `~/.kube/config` on WSL instance to `%USERPROFILE%\.kube\config` on the Windows host. Now the cluster can be accessed from both Linux and Windows terminals.

It's said [here](https://kind.sigs.k8s.io/docs/user/known-issues/#docker-desktop-for-macos-and-windows):

> Docker containers cannot be executed natively on macOS and Windows, therefore Docker Desktop runs them in a Linux VM. As a consequence, the container networks are not exposed to the host and you cannot reach the kind nodes via IP. You may be able to work around this limitation by configuring [extra port mappings](https://kind.sigs.k8s.io/docs/user/configuration/#extra-port-mappings), leveraging [cloud-provider-kind](https://github.com/kubernetes-sigs/cloud-provider-kind), using a network proxy, or other solution specific to your environment.

They say about "container networks are not exposed to the host" but actually they are not available even on the WSL instance, see the [example](../kind_echo_wsl/README.md). But `cloud-provider-kind` helps indeed, see the [example](../kind_lb/README.md).
