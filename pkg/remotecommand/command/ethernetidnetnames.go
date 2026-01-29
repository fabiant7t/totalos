package command

import (
	"fmt"
	"strings"

	"github.com/fabiant7t/totalos/pkg/remotecommand"
	"golang.org/x/crypto/ssh"
)

// EthernetIDNetNames
func EthernetIDNetNames(m remotecommand.Machine, cb ssh.HostKeyCallback) (map[string]string, error) {
	cmd := `
    udevadm info /sys/class/net/e* \
    | grep ID_NET_NAME_ \
    | cut -d " " -f 2-
  `
	names := make(map[string]string)
	stdout, err := remotecommand.Command(m, cmd, cb)
	if err != nil {
		return names, fmt.Errorf("Remote command EthernetIDNetNames failed: %w", err)
	}
	for _, line := range strings.Split(strings.TrimSpace(string(stdout)), "\n") {
		k, v, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		names[k] = v
	}
	return names, nil
}
