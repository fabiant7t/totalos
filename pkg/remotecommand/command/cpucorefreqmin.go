package command

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fabiant7t/totalos/pkg/remotecommand"
	"github.com/fabiant7t/totalos/pkg/server"
	"golang.org/x/crypto/ssh"
)

// CPUCoreFreqMin returns the minimum frequency of a core
func CPUCoreFreqMin(m remotecommand.Machine, cb ssh.HostKeyCallback) (server.MHz, error) {
	cmd := `
    lscpu -e+MHZ -J \
    | jq -r '.cpus[].minmhz' \
    | sort -nu \
    | head -n 1
  `
	stdout, err := remotecommand.Command(m, cmd, cb)
	if err != nil {
		return 0, fmt.Errorf("Remote command CPUCoreFreqMin failed: %w", err)
	}
	freq, err := strconv.ParseFloat(strings.TrimSpace(string(stdout)), 64)
	if err != nil {
		return 0, fmt.Errorf("Remote command CPUCoreFreqMin failed: %w", err)
	}
	return server.MHz(freq), nil
}
