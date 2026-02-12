package command

import (
	"fmt"
	"net/netip"
	"strings"

	"github.com/fabiant7t/totalos/pkg/remotecommand"
	"golang.org/x/crypto/ssh"
)

// ResolvectlDNSv6 queries the global DNS servers from resolvectl
func ResolvectlDNSv6(m remotecommand.Machine, cb ssh.HostKeyCallback) ([]string, error) {
	cmd := `
	  resolvectl dns \
		| grep ^Global: \
		| cut -d ' ' -f 2-
  `
	stdout, err := remotecommand.Command(m, cmd, cb)
	if err != nil {
		return nil, fmt.Errorf("Remote command ResolvectlDNS failed: %w", err)
	}
	return parsePublicIPv6sFromResolvectl(stdout), nil
}

// stdout should contain space separated IPs as a byte slice
func parsePublicIPv6sFromResolvectl(stdout []byte) []string {
	var resolvers []string
	for _, tok := range strings.Split(string(stdout), " ") {
		addr, err := netip.ParseAddr(strings.TrimSpace(tok))
		if err == nil && addr.Is6() && !addr.IsPrivate() && !addr.IsLoopback() {
			resolvers = append(resolvers, addr.String())
		}
	}
	return resolvers
}
