package kernel_test

import (
	"testing"

	"github.com/fabiant7t/totalos/pkg/kernel"
)

func TestIPOptionStaticV4String(t *testing.T) {
	opt := kernel.IPOptionStaticV4{
		ClientIP:  "192.168.0.42",
		GatewayIP: "192.168.0.1",
		Netmask:   "255.255.255.0",
		Device:    "eno1",
		DNS0IP:    "86.54.11.100",
		DNS1IP:    "9.9.9.9",
		NTP0IP:    "162.159.200.123",
	}
	want := "192.168.0.42::192.168.0.1:255.255.255.0::eno1:off:86.54.11.100:9.9.9.9:162.159.200.123"
	if got := opt.String(); got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}
