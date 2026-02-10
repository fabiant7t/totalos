package command

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/fabiant7t/totalos/pkg/kernel"
	"github.com/fabiant7t/totalos/pkg/remotecommand"
	"golang.org/x/crypto/ssh"
)

func SetIPOptionStaticV4(m remotecommand.Machine, ipOpt *kernel.IPOptionStaticV4, device string, cb ssh.HostKeyCallback) (string, error) {
	part := fmt.Sprintf("%s3", device)
	if strings.Contains(device, "nvme") {
		part = fmt.Sprintf("%sp3", device)
	}
	cmd := fmt.Sprintf(`
    mount %s /mnt \
    && cp /mnt/grub/grub.cfg /mnt/grub/grub.cfg.orig \
    && sed 's/talos.platform=metal/talos.platform=metal ip=%s/g' /mnt/grub/grub.cfg.orig \
    > /mnt/grub/grub.cfg \
    && grep ip= /mnt/grub/grub.cfg \
    && umount /mnt
  `, part, ipOpt.String())
	stdout, err := remotecommand.Command(m, cmd, cb)
	if err != nil {
		return "", fmt.Errorf("remote command SetIPOptionStaticV4 failed: %w", err)
	}

	p := regexp.MustCompile(`ip=([^ ]+)`)
	match := p.FindStringSubmatch(string(stdout))
	if len(match) >= 2 {
		return match[1], nil
	}
	return "", nil
}
