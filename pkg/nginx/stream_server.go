package nginx

import (
	"errors"
	"fmt"
	streamserverv1alpha1 "github.com/chulinx/netAgent/api/v1alpha1"
	"github.com/chulinx/zlxGo/log"
	"github.com/chulinx/zlxGo/stringfile"
	"go.uber.org/zap"
	"k8s.io/apimachinery/pkg/types"
	"path"
	"strings"
)

// StreamServerManager 管理nginx配置的生命周期
type StreamServerManager struct {
	StreamServer *streamserverv1alpha1.StreamServer
}

func NewStreamServerManager(s streamserverv1alpha1.StreamServer) *StreamServerManager {
	return &StreamServerManager{
		StreamServer: &s,
	}
}

// CreateOrUpdate write Or Update http server conf to confdPath
func (s *StreamServerManager) CreateOrUpdate() error {
	content, err := s.generateConfig()
	if err != nil {
		return err
	}
	confFile := path.Join(confdPath, fmt.Sprintf("%s-%s.stream", s.StreamServer.Name,
		s.StreamServer.Namespace))
	if strings.Contains(confFile, "--0") {
		return errors.New("stream config error")
	}
	err = stringfile.RewriteFile(content, confFile)
	if err != nil {
		return err
	}
	return reloadNginx()
}

// RemoveVirtualServerManager remove conf to confdPath
func (s *StreamServerManager) RemoveVirtualServerManager(namespacedName types.NamespacedName) error {
	confFile := path.Join(confdPath, fmt.Sprintf("%s-%s.stream", namespacedName.Name,
		namespacedName.Namespace))
	log.Info("remove stream server host config file", zap.Any("config file", confFile))
	err := removeConf(confFile)
	if err != nil {
		return err
	}
	return reloadNginx()
}

func (s *StreamServerManager) generateConfig() (string, error) {
	log.Info("start generate nginx config")
	return renderConf(streamServerTmpl, s.StreamServer.Spec)
}
