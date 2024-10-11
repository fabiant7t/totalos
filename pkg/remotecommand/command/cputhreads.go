package command

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fabiant7t/totalos/pkg/remotecommand"
	"golang.org/x/crypto/ssh"
)

// CPUThreads returns the amount of CPU threads
func CPUThreads(m remotecommand.Machine, cb ssh.HostKeyCallback) (int, error) {
	cmd := `
    dmidecode -t processor \
    | grep Thread\ Count: \
    | awk '{print $3}'
  `
	stdout, err := remotecommand.Command(m, cmd, cb)
	if err != nil {
		return 0, fmt.Errorf("Remote command CPUThreads failed: %w", err)
	}
	threads, err := strconv.ParseInt(strings.TrimSpace(string(stdout)), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("Remote command CPUThreads failed: %w", err)
	}
	return int(threads), nil
}
