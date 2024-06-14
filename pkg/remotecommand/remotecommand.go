package remotecommand

import (
	"errors"

	"golang.org/x/crypto/ssh"
)

type Machine interface {
	Addr() string
	Key() []byte
	Password() string
	User() string
}

func Command(m Machine, cmd string, hostKeyCallback ssh.HostKeyCallback) ([]byte, error) {
	// Refuse executing empty commands
	if cmd == "" {
		return nil, errors.New("Empty command")
	}
	// hostKeyCallback nil means the host key is not being verified
	if hostKeyCallback == nil {
		hostKeyCallback = ssh.InsecureIgnoreHostKey()
	}
	// password authentication
	authMethods := []ssh.AuthMethod{}
	if pwd := m.Password(); pwd != "" {
		authMethods = append(authMethods, ssh.Password(pwd))
	}
	// key authentication
	if key := m.Key(); len(key) != 0 {
		signer, err := ssh.ParsePrivateKey(m.Key())
		if err != nil {
			return nil, err
		}
		authMethods = append(authMethods, ssh.PublicKeys(signer))
	}
	// connect to the remote machine
	cc := &ssh.ClientConfig{
		User:            m.User(),
		Auth:            authMethods,
		HostKeyCallback: hostKeyCallback,
	}
	c, err := ssh.Dial("tcp", m.Addr(), cc)
	if err != nil {
		return nil, err
	}
	sess, err := c.NewSession()
	if err != nil {
		return nil, err
	}
	defer sess.Close()
	// run the command and return stdout
	return sess.Output(cmd)
}
