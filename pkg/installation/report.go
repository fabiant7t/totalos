package installation

import (
	"github.com/fabiant7t/totalos/pkg/server"
)

type Report struct {
	Installation Installation   `json:"installation"`
	Machine      server.Machine `json:"machine"`
}
