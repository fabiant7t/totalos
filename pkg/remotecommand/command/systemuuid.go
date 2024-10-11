package command

import (
	"fmt"
	"strings"

	"github.com/fabiant7t/totalos/pkg/remotecommand"
	"golang.org/x/crypto/ssh"
)

// SystemUUID as defined in SMBIOS data
func SystemUUID(m remotecommand.Machine, cb ssh.HostKeyCallback) (string, error) {
	cmd := `dmidecode -s system-uuid`
	stdout, err := remotecommand.Command(m, cmd, cb)
	if err != nil {
		return "", fmt.Errorf("Remote command SystemUUID failed: %w", err)
	}
	return strings.TrimSpace(string(stdout)), nil
}
