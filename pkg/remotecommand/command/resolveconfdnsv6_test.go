package command

import (
	"slices"
	"testing"
)

func TestResolveconfDNSv6_parsePublicIPv6s(t *testing.T) {
	stdout := []byte("nameserver 185.12.64.1\nnameserver 185.12.64.2\nnameserver ::1\nnameserver fd12:3456:789a::1\nnameserver 2a01:4ff:ff00::add:1\nnameserver 2a01:4ff:ff00::add:2")
	want := []string{"2a01:4ff:ff00::add:1", "2a01:4ff:ff00::add:2"}
	got := parsePublicIPv6sFromResolveconf(stdout)
	if !slices.Equal(got, want) {
		t.Errorf("Got %+v, want %+v", got, want)
	}
}
