package command

import (
	"encoding/json"
	"fmt"

	"github.com/fabiant7t/totalos/pkg/remotecommand"
	"github.com/fabiant7t/totalos/pkg/server"
	"golang.org/x/crypto/ssh"
)

// Disks returns all disks
func Disks(m remotecommand.Machine, cb ssh.HostKeyCallback) ([]server.Disk, error) {
	cmd := `
    lsblk -o NAME,SERIAL,SIZE,TYPE,MODEL,TRAN,WWN --json -b \
    | jq -r '.blockdevices | map(select(.type == "disk"))'
  `
	var disks []server.Disk
	stdout, err := remotecommand.Command(m, cmd, cb)
	if err != nil {
		return disks, fmt.Errorf("Remote command Disks failed: %w", err)
	}
	if err := json.Unmarshal(stdout, &disks); err != nil {
		return disks, fmt.Errorf("Remote command Disks failed: %w", err)
	}
	return disks, nil
}
