package command

import (
	"fmt"
	"strings"

	"github.com/fabiant7t/totalos/pkg/remotecommand"
	"golang.org/x/crypto/ssh"
)

// SystemManufacturer as defined in SMBIOS data
func SystemManufacturer(m remotecommand.Machine, cb ssh.HostKeyCallback) (string, error) {
	cmd := `dmidecode -s system-manufacturer`
	stdout, err := remotecommand.Command(m, cmd, cb)
	if err != nil {
		return "", fmt.Errorf("Remote command SystemManufacturer failed: %w", err)
	}
	return strings.TrimSpace(string(stdout)), nil
}
