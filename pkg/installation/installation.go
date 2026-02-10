package installation

import (
	"github.com/fabiant7t/totalos/pkg/server"
)

type Installation struct {
	Image                             string      `json:"image"`
	Rebooting                         bool        `json:"rebooting"`
	Config                            string      `json:"config"`
	StaticInitialNetworkConfiguration string      `json:"static_initial_network_configuration"`
	StorageDisk                       server.Disk `json:"storage_disk"`
	SystemDisk                        server.Disk `json:"system_disk"`
}
