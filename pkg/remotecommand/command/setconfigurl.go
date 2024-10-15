package command

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/fabiant7t/totalos/pkg/remotecommand"
	"golang.org/x/crypto/ssh"
)

// Sets talos.config in grub.cfg and returns the URL
func SetConfigURL(m remotecommand.Machine, configURL, device string, cb ssh.HostKeyCallback) (string, error) {
	replacer := strings.NewReplacer(
		"&", "\\&",
		"/", "\\/",
		":", "\\:",
		"{", "\\{",
		"}", "\\}",
	)
	part := fmt.Sprintf("%s3", device)
	if strings.Contains(device, "nvme") {
		part = fmt.Sprintf("%sp3", device)
	}
	replacement := replacer.Replace(configURL)
	cmd := fmt.Sprintf(`
    mount %s /mnt \
    && cp /mnt/grub/grub.cfg /mnt/grub/grub.cfg.orig \
    && sed 's/talos.platform=metal/talos.platform=metal talos.config=%s/g' /mnt/grub/grub.cfg.orig \
    > /mnt/grub/grub.cfg \
    && grep talos.config /mnt/grub/grub.cfg \
    && umount /mnt
  `, part, replacement)
	stdout, err := remotecommand.Command(m, cmd, cb)
	if err != nil {
		return "", fmt.Errorf("Remote command SetConfigURL failed: %w", err)
	}

	p := regexp.MustCompile(`talos.config=([^ ]+)`)
	match := p.FindStringSubmatch(string(stdout))
	if len(match) >= 2 {
		return match[1], nil
	}
	return "", nil
}
