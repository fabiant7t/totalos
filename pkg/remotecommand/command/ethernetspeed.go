package command

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fabiant7t/totalos/pkg/remotecommand"
	"github.com/fabiant7t/totalos/pkg/server"
	"golang.org/x/crypto/ssh"
)

// EthernetSpeed returns the speed of the ethernet device in Mbps.
func EthernetSpeed(m remotecommand.Machine, cb ssh.HostKeyCallback) (server.Mbps, error) {
	cmd := `
    cat /sys/class/net/$(ip -o link show | awk -F': ' '/: (en|eth)/{print $2; exit}')/speed
  `
	stdout, err := remotecommand.Command(m, cmd, cb)
	if err != nil {
		return server.Mbps(0), fmt.Errorf("Remote command EthernetSpeed failed: %w", err)
	}
	mbps_int, err := strconv.ParseInt(strings.TrimSpace(string(stdout)), 10, 64)
	if err != nil {
		return server.Mbps(0), fmt.Errorf("Remote command EthernetSpeed failed converting value: %w", err)
	}
	return server.Mbps(mbps_int), nil
}
