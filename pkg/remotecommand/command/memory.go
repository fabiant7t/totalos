package command

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fabiant7t/totalos/pkg/remotecommand"
	"github.com/fabiant7t/totalos/pkg/server"
	"golang.org/x/crypto/ssh"
)

// Memory returns the amount of RAM
func Memory(m remotecommand.Machine, cb ssh.HostKeyCallback) (server.GigaByte, error) {
	cmd := `
    dmidecode -t memory \
    | grep -i size \
    | awk '{sum += $2} END {print sum}'
  `
	stdout, err := remotecommand.Command(m, cmd, cb)
	if err != nil {
		return 0, fmt.Errorf("Remote command Memory failed: %w", err)
	}
	mem, err := strconv.ParseInt(strings.TrimSpace(string(stdout)), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("Remote command Memory failed: %w", err)
	}
	return server.GigaByte(mem), nil
}
