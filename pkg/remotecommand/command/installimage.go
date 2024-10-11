package command

import (
	"fmt"

	"github.com/fabiant7t/totalos/pkg/remotecommand"
	"golang.org/x/crypto/ssh"
)

// InstallImage downloads the raw.xz image URL and writes it to the given device.
func InstallImage(m remotecommand.Machine, isoImageURL, device string, cb ssh.HostKeyCallback) error {
	cmd := fmt.Sprintf(`
    wget %s -O talos-metal.raw.xz \
    && cat talos-metal.raw.xz \
    |  xz -d \
    |  dd of=%s bs=4M \
    && sync`, isoImageURL, device)
	if _, err := remotecommand.Command(m, cmd, cb); err != nil {
		return fmt.Errorf("Remote command InstallImage failed: %w", err)
	}
	return nil
}
