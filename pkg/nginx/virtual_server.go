package nginx

import (
	"errors"
	"fmt"
	virtualserverv1alpha1 "github.com/chulinx/netAgent/api/v1alpha1"
	"github.com/chulinx/zlxGo/log"
	"github.com/chulinx/zlxGo/net"
	"github.com/chulinx/zlxGo/stringfile"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/types"
	"path"
	"strings"
)

// VirtualServerManager 管理nginx配置的生命周期
type VirtualServerManager struct {
	VirtualServer *virtualserverv1alpha1.VirtualServer
}

func NewVirtualServerManager(v virtualserverv1alpha1.VirtualServer) *VirtualServerManager {
	return &VirtualServerManager{
		VirtualServer: &v,
	}
}

// CreateOrUpdate write Or Update http server conf to confdPath
func (m *VirtualServerManager) CreateOrUpdate() error {
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
	confFile := path.Join(confdPath, fmt.Sprintf("%s-%s.conf", m.VirtualServer.Name,
		m.VirtualServer.Namespace))
	if strings.Contains(confFile, "--0") {
		return errors.New("vs config error")
	}
	err = stringfile.RewriteFile(content, confFile)
	if err != nil {
		return err
	}
	return reloadNginx()
}

// RemoveVirtualServerManager remove conf to confdPath
func (m *VirtualServerManager) RemoveVirtualServerManager(namespacedName types.NamespacedName) error {
	confFile := path.Join(confdPath, fmt.Sprintf("%s-%s.conf", namespacedName.Name,
		namespacedName.Namespace))
	log.Info("remove virtual host config file", zap.Any("config file", confFile))
	err := removeConf(confFile)
	if err != nil {
		return err
	}
	return reloadNginx()
}

func (m *VirtualServerManager) generateConfig() (string, error) {
	log.Info("start generate nginx config")
	return renderConf(virtualServerTmpl, m.VirtualServer.Spec)
}
