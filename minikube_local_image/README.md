# Run simple service in minikube cluster from a local docker image

Start demo cluster:

```bash
minikube start --driver=kvm2
```

Connect to docker [inside](https://stackoverflow.com/questions/52310599/what-does-minikube-docker-env-mean) of minilube VM:

```bash
eval $(minikube -p minikube docker-env)
docker ps # should show containers run in VM (not local ones)
```

Build docker image to make it available to the cluster:

```bash
docker build -t kuber_learn__simple_service_1:v0 .
```

Deploy the service:

```bash
kubectl apply -f kubeconfig.yaml

deployment.apps/kuber-learn--simple-service-1-deployment created
service/kuber-learn--simple-service-1-service created
```

Call the service:

```bash
export SERVICE_ADDR=$(minikube service kuber-learn--simple-service-1-service --url)
echo $SERVICE_ADDR
curl $SERVICE_ADDR

Hello World! I'm on internal port 5000
```

It works!

---

We can check on the VM if the image is built and runnable:

```bash
minikube ssh

docker images
REPOSITORY                                TAG                 IMAGE ID            CREATED             SIZE
kuber_learn__simple_service_1             v0                  22c831c3298f        17 minutes ago      122MB

docker run --rm -d -p 5000:5000  kuber_learn__simple_service_1:v0
ee30c5d57e6b51d44a50d4b263d2bcc8bcb843b5a522cfb1550b17b85def27b7

docker ps
CONTAINER ID        IMAGE                              COMMAND                  CREATED             STATUS              PORTS                    NAMES
ee30c5d57e6b        kuber_learn__simple_service_1:v0   "python server.py"       3 seconds ago       Up 2 seconds        0.0.0.0:5000->5000/tcp   suspicious_haibt

curl localhost:5000
Hello World!

docker stop ee30c5d57e6b
ee30c5d57e6b
```

---

If something wrong happens

```bash
kubectl rollout status deployment.v1.apps/kuber-learn--simple-service-1-deployment

error: deployment "kuber-learn--simple-service-1-deployment" exceeded its progress deadline


kubectl get pods

NAME                                                        READY   STATUS             RESTARTS   AGE
kuber-learn--simple-service-1-deployment-5858c69b6c-2lp29   0/1     ImagePullBackOff   0          27m
kuber-learn--simple-service-1-deployment-5858c69b6c-k55cb   0/1     ImagePullBackOff   0          27m
kuber-learn--simple-service-1-deployment-5858c69b6c-wkjsc   0/1     ImagePullBackOff   0          27m
kuber-learn--simple-service-1-deployment-76bbc8496f-xsw69   0/1     ImagePullBackOff   0          23m
```

then check the pod status:

```bash
kubectl describe pod kuber-learn--simple-service-1-deployment-5858c69b6c-2lp29

...

Events:
  Type     Reason     Age                From               Message
  ----     ------     ----               ----               -------
  Normal   Scheduled  32s                default-scheduler  Successfully assigned default/kuber-learn--simple-service-1-deployment-6dffc4579f-687vf to minikube
  Normal   BackOff    25s                kubelet, minikube  Back-off pulling image "kuber_learn__simple_service_1:v0"
  Warning  Failed     25s                kubelet, minikube  Error: ImagePullBackOff
  Normal   Pulling    13s (x2 over 31s)  kubelet, minikube  Pulling image "kuber_learn__simple_service_1:v0"
  Warning  Failed     10s (x2 over 26s)  kubelet, minikube  Failed to pull image "kuber_learn__simple_service_1:v0": rpc error: code = Unknown desc = Error response from daemon: pull access denied for kuber_learn__simple_service_1, repository does not exist or may require 'docker login': denied: requested access to the resource is denied
  Warning  Failed     10s (x2 over 26s)  kubelet, minikube  Error: ErrImagePull
```

In this case the solution was to set `imagePullPolicy: Never` in `kubeconfig.yaml`. Seems some issues with `kvm2` driver. If we run cluster as `--driver=virtualbox` then all works ok even without this option. UPD: but sometimes the option is required even with the virtualbox driver ...weird stuff.
