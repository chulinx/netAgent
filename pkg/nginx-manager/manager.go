package nginx_manager

import (
	"bytes"
	"errors"
	"fmt"
	virtualserverv1alpha1 "github.com/chulinx/netAgent/api/v1alpha1"
	"github.com/chulinx/zlxGo/log"
	"github.com/chulinx/zlxGo/net"
	"github.com/chulinx/zlxGo/stringfile"
	"github.com/lithammer/dedent"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/types"
	"os"
	"os/exec"
	"path"
	"strings"
	"text/template"
)

const (
	binaryFilename = "/usr/sbin/nginx"
	configPath     = "/etc/nginx/"
)

var (
	confdPath = path.Join(configPath, "conf.d")
	proxyTmpl = `server
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
)

// Manager 管理nginx配置的生命周期
type Manager struct {
	VirtualServer *virtualserverv1alpha1.VirtualServer
	// confdPath default "/etc/nginx/conf.d"
	confdPath string
}

func NewManager(v virtualserverv1alpha1.VirtualServer) *Manager {
	return &Manager{
		VirtualServer: &v,
		confdPath:     path.Join(configPath, "conf.d"),
	}
}

// CreateOrUpdateVirtualServerManager write conf to confdPath
func (m *Manager) CreateOrUpdateVirtualServerManager() error {
	for _, location := range m.VirtualServer.Spec.Proxys {
		domain := fmt.Sprintf("%s.%s", location.Service, location.NameSpace)
		_, err := net.LookAddrIPFromDomain(domain)
		if err != nil {
			return err
		}
	}
	content, err := m.generateConfig()
	if err != nil {
		return err
	}
	fileName := fmt.Sprintf("%s-%s.conf", m.VirtualServer.Name,
		m.VirtualServer.Namespace,
	)
	if strings.Contains(fileName, "--0") {
		return errors.New("vs config error")
	}
	err = stringfile.RewriteFile(content, path.Join(confdPath, fileName))
	if err != nil {
		return err
	}
	return m.reloadNginx()
}

// RemoveVirtualServerManager write conf to confdPath
func (m *Manager) RemoveVirtualServerManager(namespacedName types.NamespacedName) error {
	fileName := fmt.Sprintf("%s-%s.conf", namespacedName.Name,
		namespacedName.Namespace,
	)
	confFile := path.Join(confdPath, fileName)
	if _, err := os.Stat(confFile); os.IsNotExist(err) {
		return nil
	}
	log.Info("remove virtual host config file", zap.Any("config file", confFile))
	err := os.Remove(confFile)
	if err != nil {
		return err
	}
	return m.reloadNginx()
}

func (m *Manager) reloadNginx() error {
	if err := shellOut(fmt.Sprintf("%v -s %v -e stderr", binaryFilename, "reload")); err != nil {
		return fmt.Errorf("nginx reload failed: %w", err)
	}
	return nil
}

func (m *Manager) generateConfig() (string, error) {
	var out bytes.Buffer
	log.Info("start generate nginx config")
	t := template.Must(template.New("text").Parse(dedent.Dedent(proxyTmpl)))
	err := t.Execute(&out, m.VirtualServer.Spec)
	if err != nil {
		return "", err
	}
	return out.String(), nil
}

func shellOut(cmd string) (err error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	log.Info("start run cmd", zap.Any("cmd", cmd))
	command := exec.Command("sh", "-c", cmd)
	command.Stdout = &stdout
	command.Stderr = &stderr

	err = command.Start()
	if err != nil {
		return fmt.Errorf("failed to execute %v, err: %w", cmd, err)
	}

	err = command.Wait()
	if err != nil {
		return fmt.Errorf("command %v stdout: %q\nstderr: %q\nfinished with error: %w", cmd,
			stdout.String(), stderr.String(), err)
	}
	log.Info("nginx reload success")
	return nil
}
