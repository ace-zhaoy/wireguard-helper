//go:build windows

package tunnel

import (
	"github.com/ace-zhaoy/errors"
	"os"
	"os/exec"
)

func (t *Tunnel) connect(configFile string) error {
	defer os.Remove(configFile)
	cmd := exec.Command("WireGuard.exe", "/installtunnelservice", configFile)
	return errors.WithStack(cmd.Run())
}
