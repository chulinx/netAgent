package nginx

var (
	virtualServerTmpl = `server
      {
          listen       {{.ListenPort}};
          server_name {{.ServerName}};
          {{ range $Proxy := .Proxys }}location / {
                {{ if $Proxy.NameSpace }}proxy_pass {{$Proxy.Scheme}}://{{$Proxy.Service}}.{{$Proxy.NameSpace}}:{{$Proxy.Port}};{{ else }}proxy_pass {{$Proxy.Scheme}}://{{$Proxy.Service}}:{{$Proxy.Port}};{{ end }}
                proxy_redirect     off;
                proxy_set_header   Host             $host;
                proxy_set_header   X-Real-IP        $remote_addr;
                proxy_set_header   X-Forwarded-For  $proxy_add_x_forwarded_for;
          }
		  access_log    /var/log/nginx/access.log;{{ end }}
      }`

	streamServerTmpl = `server {
                listen {{.ListenPort}};
                proxy_pass {{ .Proxy.Service }}.{{ .Proxy.NameSpace}}:{{ .Proxy.Port }};
                access_log   /var/log/nginx/tcp-access.log stream_proxy;
        }`
)
