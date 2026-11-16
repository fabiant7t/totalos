package disk

import (
	"errors"
	"sort"

	"github.com/fabiant7t/totalos/pkg/server"
)

// SelectSystemDisk takes all disks, sorts them by their serial
// alphabetically and returns the first one. The result is
// deterministic.
func SelectSystemDisk(disks []server.Disk, pref *Preference) (server.Disk, error) {
	if pref == nil {
		pref = &Preference{}
	}

	var serials []string
	disksBySerial := make(map[string]server.Disk, len(disks))

	for _, disk := range disks {
		// Skip USB devices unless explicitly allowed.
		if disk.Transport == "usb" && !pref.AllowUSB {
			continue
		}
		serials = append(serials, disk.Serial)
		disksBySerial[disk.Serial] = disk
	}
	sort.Strings(serials)

	if len(serials) > 0 {
		return disksBySerial[serials[0]], nil
	}
	return server.Disk{}, errors.New("cannot select system disk")
}
