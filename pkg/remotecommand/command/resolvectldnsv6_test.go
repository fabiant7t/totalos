package command

import (
	"slices"
	"testing"
)

func TestResolvectlDNSv6_parsePublicIPv4s(t *testing.T) {
	stdout := []byte("185.12.64.1 185.12.64.2 ::1 fd12:3456:789a::1 2a01:4ff:ff00::add:1 2a01:4ff:ff00::add:2")
	want := []string{"2a01:4ff:ff00::add:1", "2a01:4ff:ff00::add:2"}
	got := parsePublicIPv6sFromResolvectl(stdout)
	if !slices.Equal(got, want) {
		t.Errorf("Got %+v, want %+v", got, want)
	}
}
