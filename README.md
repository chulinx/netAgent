# netagent
> kubernetes ingress like, develop by operator-sdk.

## Deploy 

```shell
kubectl apply -f https://raw.githubusercontent.com/chulinx/netAgent/main/deploy/netagent.yaml
kubectl expose -n netagent-system deployment netagent-controller-manager --port 8585 --type NodePort
```

## Create deployment webserver

```shell
kubectl create deployment webserver --port 80 --image nginx:alpine
kubectl get pod -l app=webserver
# wait pod is running
kubectl exec -it `kubectl get pod -l app=webserver -o=jsonpath='{.items[0].metadata.name}'` -- bash -c "echo webserver >//usr/share/nginx/html/index.html"
kubectl expose deployment webserver --port 80 --type ClusterIP
```

## Create VirtualHost For Http Proxy

```yaml
apiVersion: crd.chulinx/v1alpha1
kind: VirtualServer
metadata:
  name: webserver
spec:
  listenPort: 8585
  proxys:
  - name: web
    nameSpace: default
    path: /
    port: 80
    scheme: http
    service: webserver
    proxyRedirect: false
    proxyHeaders:
      Host: "$host"
      X-Real-IP:       "$remote_addr"
      X-Forwarded-For: "$proxy_add_x_forwarded_for"
  serverName: webserver.chulinx.com
```

## Create VirtualHost With TLS
```yaml
apiVersion: crd.chulinx/v1alpha1
kind: VirtualServer
metadata:
  name: webserver
spec:
  listenPort: 8585
  tls: true
  tlsMountPath: "/opt"
  proxys:
  - name: web
    nameSpace: default
    path: /
    port: 80
    scheme: http
    service: webserver
    proxyRedirect: false
    proxyHeaders:
      Host: "$host"
      X-Real-IP:       "$remote_addr"
      X-Forwarded-For: "$proxy_add_x_forwarded_for"
  serverName: webserver.chulinx.com
```

## Create VirtualHost For Http Websocket

```yaml
apiVersion: crd.chulinx/v1alpha1
kind: VirtualServer
metadata:
  name: webserver-ws
spec:
  listenPort: 8585
  proxys:
  - name: web
    nameSpace: default
    path: /
    port: 80
    scheme: http
    service: webserver
    proxyHttpVersion: "1.1"
    proxyHeaders:
      Host: "$host"
      X-Real-IP:       "$remote_addr"
      X-Forwarded-For: "$proxy_add_x_forwarded_for"
      Upgrade: "$http_upgrade"
      Connection: "Upgrade"  
  serverName: webserver.chulinx.com
```

## Create streamServer
```yaml
apiVersion: crd.chulinx/v1alpha1
kind: StreamServer
metadata:
  name: streamserver-sample
spec:
  listenPort: 8585
  proxy:
    nameSpace: default
    service: webserver
    port: 80
```

## Access webserver 

```shell
curl --resolve webserver.chulinx.com:8585:192.168.0.230  webserver.chulinx.com:8585
webserver
```