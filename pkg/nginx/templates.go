package nginx

var (
	virtualServerTmpl = `server
      {
          listen       {{.ListenPort}};
          server_name {{.ServerName}};
          {{ range $Proxy := .Proxys }}location / {
                {{ if $Proxy.NameSpace }}proxy_pass {{$Proxy.Scheme}}://{{$Proxy.Service}}.{{$Proxy.NameSpace}}:{{$Proxy.Port}};{{ else }}proxy_pass {{$Proxy.Scheme}}://{{$Proxy.Service}}:{{$Proxy.Port}};{{ end }}
                {{ if not $Proxy.ProxyRedirect }}proxy_redirect     off;{{ end }}
				{{ if $Proxy.ProxyRedirect }}proxy_http_version {{ $Proxy.ProxyHttpVersion }};{{ end }}
				{{ range $key,$value := $Proxy.ProxyHeaders }}proxy_set_header   {{ $key }} {{ $value }};
				{{ end }}
          }
		  access_log    /var/log/nginx/access.log;{{ end }}
      }`

	streamServerTmpl = `server {
                listen {{.ListenPort}};
                proxy_pass {{ .Proxy.Service }}.{{ .Proxy.NameSpace}}:{{ .Proxy.Port }};
                access_log   /var/log/nginx/tcp-access.log stream_proxy;
        }`
)
