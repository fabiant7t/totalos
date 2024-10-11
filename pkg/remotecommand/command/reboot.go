package command

import (
	"github.com/fabiant7t/totalos/pkg/remotecommand"
	"golang.org/x/crypto/ssh"
)

// Reboot triggers a reboot. It does not return anything, since the
// machine should be offline and not be able to have an SSH chat :)
func Reboot(m remotecommand.Machine, cb ssh.HostKeyCallback) {
	_, _ = remotecommand.Command(m, `shutdown -r now`, cb)
}
