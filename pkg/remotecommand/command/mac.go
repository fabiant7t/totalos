package command

import (
	"fmt"
	"strings"

	"github.com/fabiant7t/totalos/pkg/remotecommand"
	"golang.org/x/crypto/ssh"
)

// MAC returns the address of the ethernet interface
func MAC(m remotecommand.Machine, cb ssh.HostKeyCallback) (string, error) {
	cmd := `
    ip -j link show \
    | jq -r '.[] | select(.ifname | startswith("en") or startswith("eth")) | .address'
  `

	stdout, err := remotecommand.Command(m, cmd, cb)
	if err != nil {
		return "", fmt.Errorf("Remote command MAC failed: %w", err)
	}
	return strings.TrimSpace(string(stdout)), nil
}
