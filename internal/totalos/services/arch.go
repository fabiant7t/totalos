package services

import (
	"strings"

	"github.com/fabiant7t/totalos/pkg/remotecommand"
	"golang.org/x/crypto/ssh"
)

func Arch(m remotecommand.Machine, hostKeyCallback ssh.HostKeyCallback) (string, error) {
	cmd := "uname -m"
	stdout, err := remotecommand.Command(m, cmd, hostKeyCallback)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(stdout)), nil
}
