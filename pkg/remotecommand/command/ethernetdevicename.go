package command

import (
	"fmt"
	"strings"

	"github.com/fabiant7t/totalos/pkg/remotecommand"
	"golang.org/x/crypto/ssh"
)

// EthernetDeviceName returns the name of the ethernet device,
// like eth0 or enp0s31f6 (when kernel is configured to use
// predictable network interface names).
func EthernetDeviceName(m remotecommand.Machine, cb ssh.HostKeyCallback) (string, error) {
	cmd := `
    ip -o link show \
    | awk -F': ' '/: (en|eth)/{print $2; exit}'
  `
	stdout, err := remotecommand.Command(m, cmd, cb)
	if err != nil {
		return "", fmt.Errorf("Remote command EthernetDeviceName failed: %w", err)
	}
	return strings.TrimSpace(string(stdout)), nil
}
