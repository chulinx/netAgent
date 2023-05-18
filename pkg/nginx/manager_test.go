package nginx

import (
	"fmt"
	virtualserverv1alpha1 "github.com/chulinx/netAgent/api/v1alpha1"
	"testing"
)

var (
	result = `server
      {
          listen       8282;
          server_name web.chulinx.com;
          location / {
                proxy_pass https://webserver:2003;
                proxy_redirect     off;
                proxy_set_header   Host             $host;
                proxy_set_header   X-Real-IP        $remote_addr;
                proxy_set_header   X-Forwarded-For  $proxy_add_x_forwarded_for;
          }
          access_log    /var/log/nginx/access.log;
      }`
	vs = virtualserverv1alpha1.VirtualServer{
		Spec: virtualserverv1alpha1.VirtualServerSpec(struct {
			ListenPort int32
			ServerName string
			Proxys     []virtualserverv1alpha1.Location
		}{
			ListenPort: 8282,
			ServerName: "web.chulinx.com",
			Proxys: []virtualserverv1alpha1.Location{
				{
					Scheme: "https",
					//NameSpace: "default",
					Service: "webserver",
					Port:    2003,
				},
			},
		}),
	}
)

func Test_generateConfig(t *testing.T) {
	newManger := NewManager(vs)
	c, err := newManger.generateConfig()
	if err != nil {
		t.Errorf("Error,%s", err.Error())
	}
	if c == result {
		t.Log("Success")
	} else {
		fmt.Println(c)
		fmt.Println(result)
	}
}
