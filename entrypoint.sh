#!/usr/sh
#netAgentLog="/var/log/nginx/netAgent.log"
#[ -d /var/log/nginx/ ] mkdir -p /var/log/nginx/
#[ -f $ingressNgLog ] || touch $netAgentLog
nohup /netagent --health-probe-bind-address=:8081 --metrics-bind-address=127.0.0.1:8080 --leader-elect >> /dev/stdout 2>/dev/stderr &
sleep 3
/usr/sbin/nginx -g "daemon off;"