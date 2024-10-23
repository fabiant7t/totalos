package installation

import (
	"github.com/fabiant7t/totalos/pkg/server"
)

type Installation struct {
	FormatStorageDisk bool        `json:"format_storage_disk"`
	Image             string      `json:"image"`
	Rebooting         bool        `json:"rebooting"`
	Config            string      `json:"config"`
	StorageDisk       server.Disk `json:"storage_disk"`
	SystemDisk        server.Disk `json:"system_disk"`
}
