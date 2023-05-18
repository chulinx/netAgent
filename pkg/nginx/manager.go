package nginx

import (
	"bytes"
	"fmt"
	"github.com/lithammer/dedent"
	"k8s.io/apimachinery/pkg/types"
	"os"
	"os/exec"
	"path"
	"text/template"
)

const (
	binaryFilename = "/usr/sbin/nginx"
	configPath     = "/etc/nginx/"
)

var (
	confdPath = path.Join(configPath, "conf.d")
)

type Manager interface {
	CreateOrUpdate() error
	RemoveVirtualServerManager(namespacedName types.NamespacedName) error
}

func removeConf(confFile string) error {
	if _, err := os.Stat(confFile); os.IsNotExist(err) {
		return nil
	}
	err := os.Remove(confFile)
	if err != nil {
		return err
	}
	return nil
}

func renderConf(tmpl string, data interface{}) (string, error) {
	var out bytes.Buffer
	t := template.Must(template.New("text").Parse(dedent.Dedent(tmpl)))
	err := t.Execute(&out, data)
	if err != nil {
		return "", err
	}
	return out.String(), nil
}

func reloadNginx() error {
	if err := shellOut(fmt.Sprintf("%v -s %v -e stderr", binaryFilename, "reload")); err != nil {
		return fmt.Errorf("nginx reload failed: %w", err)
	}
	return nil
}

func shellOut(cmd string) (err error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
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
	return nil
}
