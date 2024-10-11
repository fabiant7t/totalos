package server

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
	return fmt.Sprintf("%s:%d", s.ip, s.port)
}

func (s *server) User() string {
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

type Args struct {
	Port     uint16
	Password string
	Key      []byte
}

func New(ip, user string, args *Args) *server {
	srv := &server{
		ip:       ip,
		user:     user,
		password: args.Password,
		port:     args.Port,
		key:      args.Key,
	}
	if srv.user == "" {
		srv.user = "root"
	}
	if srv.port == 0 {
		srv.port = 22
	}
	return srv
}
