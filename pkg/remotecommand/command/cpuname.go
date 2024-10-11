package command

import (
	"fmt"
	"strings"

	"github.com/fabiant7t/totalos/pkg/remotecommand"
	"golang.org/x/crypto/ssh"
)

// CPUName returns the CPU name
func CPUName(m remotecommand.Machine, cb ssh.HostKeyCallback) (string, error) {
	cmd := `
    dmidecode -t processor \
    | grep Version\: \
    | cut -d ':' -f 2- \
    | awk '{print}'
  `
	stdout, err := remotecommand.Command(m, cmd, cb)
	if err != nil {
		return "", fmt.Errorf("Remote command CPUName failed: %w", err)
	}
	return strings.TrimSpace(string(stdout)), nil
}
