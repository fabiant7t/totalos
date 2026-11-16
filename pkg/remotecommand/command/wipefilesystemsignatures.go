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
		for disk in $(
			lsblk --json --tree --paths --output NAME,TYPE,TRAN,PKNAME \
			| jq -r '.blockdevices[] | select(.type == "disk" and (.tran == "nvme" or .tran == "sata")) | .name'
		); \
		do
			wipefs -fa "${disk}";
		done;
	`
	if _, err := remotecommand.Command(m, cmd, cb); err != nil {
		return fmt.Errorf("remote command WipeFileSystemSignatures failed: %w", err)
	}
	return nil
}
