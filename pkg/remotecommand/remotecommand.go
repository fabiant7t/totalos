package remotecommand

import (
	"errors"
	"log/slog"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

var (
	ErrEmptyCommand = errors.New("empty command")
)

type Machine interface {
	Addr() string
	Key() []byte
	Password() string
	User() string
}

func Command(m Machine, cmd string, hostKeyCallback ssh.HostKeyCallback) ([]byte, error) {
	logger := slog.With("component", "remotecommand", "host", m.Addr())
	// Refuse executing empty commands
	if cmd == "" {
		return nil, ErrEmptyCommand
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
			logger.Debug("SSH private key parsing failed", "error", err)
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
	logger.Debug("opening SSH connection")
	c, err := ssh.Dial("tcp", m.Addr(), cc)
	if err != nil {
		logger.Debug("SSH connection failed", "error", err)
		return nil, err
	}
	defer c.Close()
	sess, err := c.NewSession()
	if err != nil {
		logger.Debug("SSH session creation failed", "error", err)
		return nil, err
	}
	defer sess.Close()
	// Commands may contain secrets; only include them at debug level.
	logger.Debug("running remote command", "command", strings.Join(strings.Fields(cmd), " "))
	started := time.Now()
	stdout, err := sess.Output(cmd)
	if err != nil {
		logger.Debug("remote command failed", "duration", time.Since(started), "error", err)
	} else {
		logger.Debug("remote command completed", "duration", time.Since(started), "stdout_bytes", len(stdout))
	}
	return stdout, err
}
