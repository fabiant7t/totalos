package command

import (
	"bufio"
	"bytes"
	"fmt"
	"net/netip"
	"strings"

	"github.com/fabiant7t/totalos/pkg/remotecommand"
	"golang.org/x/crypto/ssh"
)

// ResolveconfDNSv6 queries the global DNS servers from /etc/resolv.conf
func ResolveconfDNSv6(m remotecommand.Machine, cb ssh.HostKeyCallback) ([]string, error) {
	cmd := `grep nameserver /etc/resolv.conf`
	stdout, err := remotecommand.Command(m, cmd, cb)
	if err != nil {
		return nil, fmt.Errorf("Remote command ResolvectlDNS failed: %w", err)
	}
	return parsePublicIPv6sFromResolveconf(stdout), nil
}

// stdout should contain linebreak separated IPs as a byte slice
func parsePublicIPv6sFromResolveconf(stdout []byte) []string {
	var resolvers []string
	scanner := bufio.NewScanner(bytes.NewReader(stdout))
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) == 2 && fields[0] == "nameserver" {
			addr, err := netip.ParseAddr(fields[1])
			if err == nil && addr.Is6() && !addr.IsPrivate() && !addr.IsLoopback() {
				resolvers = append(resolvers, addr.String())
			}
		}
	}
	return resolvers
}
