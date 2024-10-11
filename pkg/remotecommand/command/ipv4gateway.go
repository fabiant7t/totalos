package command

import (
	"fmt"
	"net"
	"strings"

	"github.com/fabiant7t/totalos/pkg/remotecommand"
	"golang.org/x/crypto/ssh"
)

// IPv4Gateway returns the default gateway IPv4 of ethernet device
func IPv4Gateway(m remotecommand.Machine, cb ssh.HostKeyCallback) (net.IP, error) {
	cmd := `
    ip -j -4 route show \
    | jq -r '.[] | select(.dev | startswith("en") or startswith("eth")) | select(.dst == "default") | .gateway'
  `
	stdout, err := remotecommand.Command(m, cmd, cb)
	if err != nil {
		return nil, fmt.Errorf("Remote command IPv4Gateway failed: %w", err)
	}
	return net.ParseIP(strings.TrimSpace(string(stdout))), nil
}
