package command

import (
	"fmt"

	"github.com/fabiant7t/totalos/pkg/remotecommand"
	"golang.org/x/crypto/ssh"
)

// WipeFileSystemSignatures erases all available signatures of all NVMe and SATA drives.
// It does not shred the data, though!
func WipeFileSystemSignatures(m remotecommand.Machine, cb ssh.HostKeyCallback) error {
	cmd := `
    for satadisk in $(ls /dev/sd*); do
        wipefs -fa ${satadisk};
    done;
    for nvmedisk in $(ls /dev/nvme*n1); do
        wipefs -fa ${nvmedisk};
    done;
  `
	if _, err := remotecommand.Command(m, cmd, cb); err != nil {
		return fmt.Errorf("Remote command WipeFileSystemSignatures failed: %w", err)
	}
	return nil
}
