package nginx

var (
	virtualServerTmpl = `server
      {
          {{ if .TlsSecret }}listen       {{.ListenPort}} ssl;
          server_name {{.ServerName}};
          ssl_certificate     /opt/{{ .NameSpace}}-{{ .SecretName }}.crt;
          ssl_certificate_key /opt/{{ .NameSpace}}-{{ .SecretName }}.key;
		  {{ else }}listen       {{.ListenPort}};
		  server_name {{.ServerName}};{{ end }}
          {{ range $Proxy := .Proxys }}location / {
             {{ if $Proxy.ProxyPass }}proxy_pass   {{ $Proxy.ProxyPass }};{{ end }}
             {{ if not $Proxy.ProxyRedirect }}proxy_redirect     off;{{ end }}
             {{ if $Proxy.ProxyHttpVersion }}proxy_http_version {{ $Proxy.ProxyHttpVersion }};{{ end }}
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
