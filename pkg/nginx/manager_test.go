package nginx

import (
	"fmt"
	virtualserverv1alpha1 "github.com/chulinx/netAgent/api/v1alpha1"
	"testing"
)

var (
	vs1 = virtualserverv1alpha1.VirtualServer{
		Spec: virtualserverv1alpha1.VirtualServerSpec{
			ListenPort: 8282,
			ServerName: "web.chulinx.com",
			TlsSecret:  "default/tls-secret",
			Proxys: []virtualserverv1alpha1.Location{
				{
					ProxyPass:        "http://webserver:2003",
					ProxyRedirect:    false,
					ProxyHttpVersion: "1.1",
					ProxyHeaders: map[string]string{
						"Host":            "$host",
						"X-Real-IP":       "$remote_addr",
						"X-Forwarded-For": "$proxy_add_x_forwarded_for",
					},
				},
			},
		},
	}
	vs2 = virtualserverv1alpha1.VirtualServer{
		Spec: virtualserverv1alpha1.VirtualServerSpec{
			ListenPort: 8283,
			ServerName: "web.chulinx.com",
			Proxys: []virtualserverv1alpha1.Location{
				{
					ProxyPass:        "http://webserver:2003",
					ProxyRedirect:    false,
					ProxyHttpVersion: "1.1",
					ProxyHeaders: map[string]string{
						"Host":            "$host",
						"X-Real-IP":       "$remote_addr",
						"X-Forwarded-For": "$proxy_add_x_forwarded_for",
					},
				},
			},
		},
	}
	s = virtualserverv1alpha1.StreamServer{Spec: virtualserverv1alpha1.StreamServerSpec{
		ListenPort: 80,
		Name:       "test",
		Proxy: virtualserverv1alpha1.Proxy{
			NameSpace: "default",
			Service:   "webserver",
			Port:      80,
		},
	}}
)

func Test_generateConfig(t *testing.T) {
	newManger := NewVirtualServerManager(vs1)
	c, err := newManger.generateConfig()
	if err != nil {
		t.Errorf("Error,%s", err.Error())
	}
	fmt.Println(c)
	newManger = NewVirtualServerManager(vs2)
	c, err = newManger.generateConfig()
	if err != nil {
		t.Errorf("Error,%s", err.Error())
	}
	fmt.Println(c)
	//streamServer := NewStreamServerManager(s)
	//fmt.Println(streamServer.generateConfig())
}
