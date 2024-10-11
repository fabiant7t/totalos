package command

import (
	"fmt"

	"github.com/fabiant7t/totalos/pkg/remotecommand"
	"golang.org/x/crypto/ssh"
)

// FormatXFS formats the full device with one XFS partition.
// Sets disk UUID (default is "61291e61-291e-6129-1e61-291e61291e00"),
// partition label (default is "storage") and
// partition UUID (default is "61291e61-291e-6129-1e61-291e61291e01").
func FormatXFS(m remotecommand.Machine, device, diskUUID, partLabel, partUUID string, cb ssh.HostKeyCallback) error {
	if diskUUID == "" {
		diskUUID = "61291e61-291e-6129-1e61-291e61291e00"
	}
	if partLabel == "" {
		partLabel = "storage"
	}
	if partUUID == "" {
		partUUID = "61291e61-291e-6129-1e61-291e61291e01"
	}
	cmd := fmt.Sprintf(`
    export DEVICE=%s \
    && wipefs -af ${DEVICE} \
    && parted ${DEVICE} --script mklabel gpt \
    && sgdisk --disk-guid=%s ${DEVICE} \
    && parted ${DEVICE} --script mkpart primary xfs %s %s \
    && mkfs.xfs -f -L %s -m uuid=%s ${DEVICE}*1
  `, device, diskUUID, "0%", "100%", partLabel, partUUID)
	if _, err := remotecommand.Command(m, cmd, cb); err != nil {
		return fmt.Errorf("Remote command FormatXFS failed: %w", err)
	}
	return nil
}
