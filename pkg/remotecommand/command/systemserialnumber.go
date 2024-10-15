package command

import (
	"fmt"
	"strings"

	"github.com/fabiant7t/totalos/pkg/remotecommand"
	"golang.org/x/crypto/ssh"
)

// SystemSerialNumber as defined in SMBIOS data
func SystemSerialNumber(m remotecommand.Machine, cb ssh.HostKeyCallback) (string, error) {
	cmd := `dmidecode -s system-serial-number`
	stdout, err := remotecommand.Command(m, cmd, cb)
	if err != nil {
		return "", fmt.Errorf("Remote command SystemSerialNumber failed: %w", err)
	}
	return strings.TrimSpace(string(stdout)), nil
}
