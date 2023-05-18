package nginx

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/chulinx/zlxGo/log"
	"github.com/chulinx/zlxGo/net"
	"github.com/chulinx/zlxGo/stringfile"
	"github.com/lithammer/dedent"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/types"
	"os"
	"path"
	"strings"
	"text/template"
)

var (
	confdPath = path.Join(configPath, "conf.d")
)

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
	t := template.Must(template.New("text").Parse(dedent.Dedent(virtualServerTmpl)))
	err := t.Execute(&out, m.VirtualServer.Spec)
	if err != nil {
		return "", err
	}
	return out.String(), nil
}
