package command

import (
	"slices"
	"testing"
)

func TestResolveconfDNSv4_parsePublicIPv4s(t *testing.T) {
	stdout := []byte("nameserver 127.0.0.1\nnameserver 10.10.10.10\nnameserver 185.12.64.1\nnameserver 185.12.64.2\nnameserver 2a01:4ff:ff00::add:1\nnameserver 2a01:4ff:ff00::add:2")
	want := []string{"185.12.64.1", "185.12.64.2"}
	got := parsePublicIPv4sFromResolveconf(stdout)
	if !slices.Equal(got, want) {
		t.Errorf("Got %+v, want %+v", got, want)
	}
}
