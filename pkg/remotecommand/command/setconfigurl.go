package command

import (
	"fmt"
	"strings"

	"github.com/fabiant7t/totalos/pkg/remotecommand"
	"golang.org/x/crypto/ssh"
)

// Sets talos.config in grub.cfg
func SetConfigURL(m remotecommand.Machine, configURL, device string, cb ssh.HostKeyCallback) error {
	replacer := strings.NewReplacer(
		"&", "\\&",
		"/", "\\/",
		":", "\\:",
		"{", "\\{",
		"}", "\\}",
	)
	replacement := replacer.Replace(configURL)
	cmd := fmt.Sprintf(`
    mount %sp3 /mnt \
    && cp /mnt/grub/grub.cfg /mnt/grub/grub.cfg.orig \
    && sed 's/talos.platform=metal/talos.platform=metal talos.config=%s/g' /mnt/grub/grub.cfg.orig \
    > /mnt/grub/grub.cfg \
    && umount /mnt
  `, device, replacement)
	if _, err := remotecommand.Command(m, cmd, cb); err != nil {
		return fmt.Errorf("Remote command SetConfigURL failed: %w", err)
	}
	return nil
}
