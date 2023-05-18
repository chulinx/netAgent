#!/usr/sh
[ -f /etc/nginx/conf.d/default.conf ] && rm -f /etc/nginx/conf.d/default.conf
ln -s /var/log/nginx/tcp-access.log /dev/stdout
nohup /netagent --health-probe-bind-address=:8081 --metrics-bind-address=127.0.0.1:8080 --leader-elect >> /dev/stdout 2>/dev/stderr &
sleep 3
/usr/sbin/nginx -g "daemon off;"