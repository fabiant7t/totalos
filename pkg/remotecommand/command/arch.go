package command

import (
	"fmt"
	"strings"

	"github.com/fabiant7t/totalos/pkg/remotecommand"
	"golang.org/x/crypto/ssh"
)

// Arch returns the machine architecture (see `uname -m`)
func Arch(m remotecommand.Machine, cb ssh.HostKeyCallback) (string, error) {
	stdout, err := remotecommand.Command(m, "uname -m", cb)
	if err != nil {
		return "", fmt.Errorf("Remote command Arch failed: %w", err)
	}
	return strings.TrimSpace(string(stdout)), nil
}
