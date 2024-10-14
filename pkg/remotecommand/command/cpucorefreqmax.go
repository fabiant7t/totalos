package command

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fabiant7t/totalos/pkg/remotecommand"
	"github.com/fabiant7t/totalos/pkg/server"
	"golang.org/x/crypto/ssh"
)

// CPUCoreFreqMax returns the maximum frequency of a core
func CPUCoreFreqMax(m remotecommand.Machine, cb ssh.HostKeyCallback) (server.MHz, error) {
	cmd := `
    lscpu -e+MHZ -J \
    | jq '.cpus[].maxmhz' \
    | sort -nu \
    | tail -n 1
  `
	stdout, err := remotecommand.Command(m, cmd, cb)
	if err != nil {
		return 0, fmt.Errorf("Remote command CPUCoreFreqMax failed: %w", err)
	}
	freq, err := strconv.ParseInt(strings.TrimSpace(string(stdout)), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("Remote command CPUCoreFreqMax failed: %w", err)
	}
	return server.MHz(freq), nil
}
