package totalos

import (
	"fmt"
	"os"
)

// server contains information about a Hetzner dedicated server
// and implements the remotecommand.Machine interface.
type server struct {
	ip       string
	port     uint16
	user     string
	password string
	key      []byte
}

func (s *server) Addr() string {
	if s.port != 0 {
		return fmt.Sprintf("%s:%d", s.ip, s.port)
	}
	return fmt.Sprintf("%s:22", s.ip)
}

func (s *server) User() string {
	if s.user == "" {
		return "root"
	}
	return s.user
}

func (s *server) Password() string {
	return s.password
}

func (s *server) Key() []byte {
	return s.key
}

func (s *server) SetKeyFromFile(path string) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	s.key = b
	return nil
}

func NewServer(ip, user string, port uint16, password string, key []byte) *server {
	return &server{
		ip:       ip,
		user:     user,
		password: password,
		port:     port,
		key:      key,
	}
}
