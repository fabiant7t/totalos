package command

import (
	"fmt"
	"net"
	"strings"

	"github.com/fabiant7t/totalos/pkg/remotecommand"
	"golang.org/x/crypto/ssh"
)

// IPv4 address of ethernet device
func IPv4(m remotecommand.Machine, cb ssh.HostKeyCallback) (net.IP, error) {
	cmd := `
    ip -4 -j a show \
    | jq -r '.[] | select(.ifname | startswith("en") or startswith("eth")) | .addr_info[].local'
  `
	stdout, err := remotecommand.Command(m, cmd, cb)
	if err != nil {
		return nil, fmt.Errorf("Remote command IPv4 failed: %w", err)
	}
	return net.ParseIP(strings.TrimSpace(string(stdout))), nil
}
