package command

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fabiant7t/totalos/pkg/remotecommand"
	"github.com/fabiant7t/totalos/pkg/server"
	"golang.org/x/crypto/ssh"
)

// Storage returns a slice with the gigabytes of storage per disk
func Storage(m remotecommand.Machine, cb ssh.HostKeyCallback) ([]server.GigaByte, error) {
	cmd := `
    lsblk -b --json \
    | jq -r '.blockdevices | map(select(.type =="disk")) | .[] | .size' | awk '{printf "%d\n",$1 / 1000000000}'
  `
	stdout, err := remotecommand.Command(m, cmd, cb)
	if err != nil {
		return []server.GigaByte{}, fmt.Errorf("Remote command Storage failed: %w", err)
	}
	var storage []server.GigaByte
	for _, s := range strings.Split(strings.TrimSpace(string(stdout)), "\n") {
		gb, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return []server.GigaByte{}, fmt.Errorf("Remote command Storage failed: %w", err)
		}
		storage = append(storage, server.GigaByte(gb))
	}
	return storage, nil
}
