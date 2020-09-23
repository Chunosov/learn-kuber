# Run HTTPS service with Istio gateway in minikube cluster

https://istio.io/latest/docs/tasks/traffic-management/ingress/secure-ingress/
https://programmaticponderings.com/2019/01/03/securing-your-istio-gateway-with-https/

Start a cluster and get get proxy ports:

```bash
minikube start --driver=virtualbox
istioctl install --set profile=demo
kubectl label namespace default istio-injection=enabled

export INGRESS_PORT_HTTP=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.spec.ports[?(@.name=="http2")].nodePort}')
export INGRESS_PORT_HTTPS=$(kubectl -n istio-system get service istio-ingressgateway -o jsonpath='{.spec.ports[?(@.name=="https")].nodePort}')
echo $INGRESS_PORT_HTTP,$INGRESS_PORT_HTTPS
```

Deploy sample service:

```bash
kubectl apply -f service.yaml
```

Create a root certificate and private key `example.com` to sign the certificates for services. Then Create a private key and a certificate and for `banana.fruits.com` domain:

```bash
openssl req -x509 -sha256 -nodes -days 365 -newkey rsa:2048 -subj '/O=example Inc./CN=example.com' -keyout example.com.key -out example.com.crt

openssl req -out banana.fruits.com.csr -newkey rsa:2048 -nodes -keyout banana.fruits.com.key -subj "/CN=banana.fruits.com/O=banana.fruits.com"
openssl x509 -req -days 365 -CA example.com.crt -CAkey example.com.key -set_serial 0 -in banana.fruits.com.csr -out banana.fruits.com.crt
```

Create a secret for the ingress gateway:

```bash
kubectl create -n istio-system secret tls fruits-credential --key=banana.fruits.com.key --cert=banana.fruits.com.crt
```

Deploy a gateway:

```bash
kubectl apply -f gateway.yaml
```

Try to connect to services:

```bash
curl -H "Host: apple.fruits.com" http://$(minikube ip):$INGRESS_PORT_HTTP
apple

curl -k -H "Host: banana.fruits.com" https://$(minikube ip):$INGRESS_PORT_HTTPS
curl: (35) OpenSSL SSL_connect: SSL_ERROR_SYSCALL in connection to 192.168.99.121:30851

```

```bash
curl -v -H "Host: banana.fruits.com" \
    --resolve "banana.fruits.com:$INGRESS_PORT_HTTPS:$(minikube ip)" \
    --cacert example.com.crt \
    https://$(minikube ip):$INGRESS_PORT_HTTPS/status/418

* Added banana.fruits.com:30851:192.168.99.121 to DNS cache
*   Trying 192.168.99.121...
* TCP_NODELAY set
* Connected to 192.168.99.121 (192.168.99.121) port 30851 (#0)
* ALPN, offering h2
* ALPN, offering http/1.1
* successfully set certificate verify locations:
*   CAfile: example.com.crt
  CApath: /etc/ssl/certs
* TLSv1.3 (OUT), TLS handshake, Client hello (1):
* OpenSSL SSL_connect: SSL_ERROR_SYSCALL in connection to 192.168.99.121:30851
* stopped the pause stream!
* Closing connection 0
curl: (35) OpenSSL SSL_connect: SSL_ERROR_SYSCALL in connection to 192.168.99.121:30851
```

<span style="color:red"><b>It doesn't work.</b></span>
