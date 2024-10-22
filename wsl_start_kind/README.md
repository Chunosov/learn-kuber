# Start demo cluster with WSL and kind

Make sure you have WSL and Docker [prepared](../wsl_start_k3s/README.md). 

Install kind on WSL instance:

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

**NB:** Networking on WSL cluster is somewhat different than that on true Linux installation, see the [example](../wsl_echo_kind/README.md). The load balancer [example](../kind_lb/README.md) also doesn't work on WSL instance.

*TODO:* How to access services properly on WSL?
