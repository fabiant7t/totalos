package command

import (
	"fmt"

	"github.com/fabiant7t/totalos/pkg/remotecommand"
	"golang.org/x/crypto/ssh"
)

// SoftwareRAIDNotExists ensures that there are no software RAID partitions (/dev/md*)
func SoftwareRAIDNotExists(m remotecommand.Machine, cb ssh.HostKeyCallback) error {
	cmd := `
    fdisk -l \
    |  grep -E '^Disk /dev/md[0-9]+' \
    && mdadm --stop /dev/md/* \
    || echo Software RAID already missing
  `
	if _, err := remotecommand.Command(m, cmd, cb); err != nil {
		return fmt.Errorf("Remote command SoftwareRAIDNotExists failed: %w", err)
	}
	return nil
}
