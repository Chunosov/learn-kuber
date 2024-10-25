# Start demo k3s cluster on WSL

Make sure you have [WSL](../wsl/README.md) prepared.

Check if the Linux instance has the [required settings](https://ranchermanager.docs.rancher.com/v2.5/getting-started/installation-and-upgrade/installation-requirements) set to 1:

```bash
sysctl net.bridge.bridge-nf-call-iptables
net.bridge.bridge-nf-call-iptables = 1
```

If needed, can be changed via:

```bash
sudo sysctl -w <variable name>=<value>
```

Install Rancher [k3s](https://www.rancher.com/products/k3s) on the WSL instance:

```bash
sudo docker run --privileged -d --restart=unless-stopped -p 9080:80 -p 9443:443 rancher/rancher
```

Find a temp admin password of the cluster manager

```bash
docker ps # Find id of the rancher/rancher container
docker logs {container-id} 2>&1 | grep "Bootstrap Password:"
```

```log
2024/10/18 09:21:25 [INFO] Bootstrap Password: z2bdfh8n6lpgf8g6gmhqfmhwzvwbmhl6qpm5zbxj4vxbph2t86nc9k
```

Open a browser, navigate to `localhost:9080`, and login with the temp password. Provide a new password for secure connections.

Go to the Cluster Management page, download cluster's kubeconfig `local.yaml`, put it into `%USERPROFILE%\.kube` on the Windows host (create the dir if it doesn't exists), and rename to `config` (no yaml ext). Now the Windows host can manage the cluster:

```bash
kubectl version
Client Version: v1.30.2
Kustomize Version: v5.0.4-0.20230601165947-6ce0bf390ce3
Server Version: v1.30.2+k3s2

kubectl get nodes -o wide
NAME         STATUS   ROLES                       AGE     VERSION        INTERNAL-IP   EXTERNAL-IP   OS-IMAGE                              KERNEL-VERSION                       CONTAINER-RUNTIME
local-node   Ready    control-plane,etcd,master   2d23h   v1.30.2+k3s2   172.17.0.2    <none>        SUSE Linux Enterprise Server 15 SP6   5.15.153.1-microsoft-standard-WSL2   containerd://1.7.17-k3s1
```
