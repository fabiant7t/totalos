package command

import (
	"slices"
	"testing"
)

func TestResolvectlDNSv4_parsePublicIPv4s(t *testing.T) {
	stdout := []byte("127.0.0.1 10.10.10.10 185.12.64.1 185.12.64.2 2a01:4ff:ff00::add:1 2a01:4ff:ff00::add:2")
	want := []string{"185.12.64.1", "185.12.64.2"}
	got := parsePublicIPv4sFromResolvectl(stdout)
	if !slices.Equal(got, want) {
		t.Errorf("Got %+v, want %+v", got, want)
	}
}
