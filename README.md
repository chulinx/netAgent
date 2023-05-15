# netagent
> kubernetes ingress like, develop by operator-sdk.

## Deploy 

```shell
kubectl apply -f https://
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

## Create VirtualHost

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
  serverName: webserver.chulinx.com
```

## Access webserver 

```shell
curl --resolve webserver.chulinx.com:8585:192.168.0.230  webserver.chulinx.com:8585
webserver
```