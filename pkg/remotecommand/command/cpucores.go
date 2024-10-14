package command

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fabiant7t/totalos/pkg/remotecommand"
	"golang.org/x/crypto/ssh"
)

// CPUCores returns the amount of real CPU cores
func CPUCores(m remotecommand.Machine, cb ssh.HostKeyCallback) (int, error) {
	cmd := `
    dmidecode -t processor \
    | grep Core\ Count: \
    | awk '{print $3}'
  `
	stdout, err := remotecommand.Command(m, cmd, cb)
	if err != nil {
		return 0, fmt.Errorf("Remote command CPUCores failed: %w", err)
	}
	count, err := strconv.ParseInt(strings.TrimSpace(string(stdout)), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("Remote command CPUCores failed: %w", err)
	}
	return int(count), nil
}
