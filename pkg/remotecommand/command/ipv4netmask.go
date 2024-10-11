package command

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/fabiant7t/totalos/pkg/remotecommand"
	"golang.org/x/crypto/ssh"
)

// IPv4Netmask returns the IPv4 netmask of ethernet device
func IPv4Netmask(m remotecommand.Machine, cb ssh.HostKeyCallback) (net.IP, error) {
	cmd := `
    ip -4 -j a show \
    | jq -r '.[] | select(.ifname | startswith("en") or startswith("eth")) | .addr_info[].prefixlen'
  `
	stdout, err := remotecommand.Command(m, cmd, cb)
	if err != nil {
		return nil, fmt.Errorf("Remote command IPv4Netmask failed: %w", err)
	}
	prefixlen, err := strconv.ParseInt(strings.TrimSpace(string(stdout)), 10, 32)
	if err != nil {
		return nil, fmt.Errorf("Remote command IPv4Netmask failed: %w", err)
	}
	mask := net.CIDRMask(int(prefixlen), 32)
	return net.IPv4(mask[0], mask[1], mask[2], mask[3]), nil
}
