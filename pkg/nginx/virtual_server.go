package nginx

import (
	"errors"
	"fmt"
	virtualserverv1alpha1 "github.com/chulinx/netAgent/api/v1alpha1"
	"github.com/chulinx/zlxGo/log"
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

type VirtualServerManagerTemplate struct {
	virtualserverv1alpha1.VirtualServerSpec
	NameSpace  string
	SecretName string
}

func NewVirtualServerManager(v virtualserverv1alpha1.VirtualServer) *VirtualServerManager {
	return &VirtualServerManager{
		VirtualServer: &v,
	}
}

// CreateOrUpdate write Or Update http server conf to confdPath
func (m *VirtualServerManager) CreateOrUpdate() error {
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
	var (
		nameSpace, secretName string
	)

	secretSplit := strings.Split(m.VirtualServer.Spec.TlsSecret, "/")
	if len(secretSplit) == 2 {
		nameSpace = secretSplit[0]
		secretName = secretSplit[1]
	} else {
		nameSpace = m.VirtualServer.Namespace
		secretName = m.VirtualServer.Spec.TlsSecret
	}
	vsTemple := VirtualServerManagerTemplate{
		VirtualServerSpec: m.VirtualServer.Spec,
		NameSpace:         nameSpace,
		SecretName:        secretName,
	}
	return renderConf(virtualServerTmpl, vsTemple)
}
