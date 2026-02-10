package server_test

import (
	"testing"

	"github.com/fabiant7t/totalos/pkg/server"
)

func TestInterfaceName(t *testing.T) {
	for _, tc := range []struct {
		name string
		idnn *server.IDNetNames
		want string
	}{
		{"from database is being ignored", &server.IDNetNames{FromDatabase: "d"}, ""},
		{"1st choice is onboard", &server.IDNetNames{Onboard: "o", Slot: "s", Path: "p", MAC: "m"}, "o"},
		{"2nd choice is slot", &server.IDNetNames{Slot: "s", Path: "p", MAC: "m"}, "s"},
		{"3rd choice is path", &server.IDNetNames{Path: "p", MAC: "m"}, "p"},
		{"4th choice is mac", &server.IDNetNames{MAC: "m"}, "m"},
		{"returns empty string if there are no names to choose from", &server.IDNetNames{}, ""},
	} {
		if got, want := tc.idnn.InterfaceName(), tc.want; got != want {
			t.Errorf("%s: got %s, want %s", tc.name, got, want)
		}
	}
}
