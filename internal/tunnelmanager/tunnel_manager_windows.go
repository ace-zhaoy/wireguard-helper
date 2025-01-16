//go:build windows

package tunnelmanager

import (
	"github.com/ace-zhaoy/errors"
	"github.com/ace-zhaoy/glog/log"
	"os/exec"
	"regexp"
)

func (t *TunnelManager) isConnected() (connected bool, interfaceName string, err error) {
	defer errors.Recover(func(e error) { err = e })
	cmd := exec.Command("wg", "show")
	out, err := cmd.Output()
	errors.Check(errors.WithStack(err))
	reg := regexp.MustCompile(`interface:\s*([^\n\r]+)`)
	matches := reg.FindStringSubmatch(string(out))
	if len(matches) > 1 {
		interfaceName, connected = matches[1], true
	}
	log.Debug("interfaceName: %s, connected: %v", interfaceName, connected)
	return
}

func (t *TunnelManager) disconnect(interfaceName string) error {
	cmd := exec.Command("WireGuard.exe", "/uninstalltunnelservice", interfaceName)
	return errors.WithStack(cmd.Run())
}
