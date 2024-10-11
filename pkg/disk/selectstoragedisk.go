package disk

import (
	"errors"

	"github.com/fabiant7t/totalos/pkg/server"
)

// SelectStorageDisk finds and returns the biggest available disk which
// is not the system disk.
func SelectStorageDisk(disks []server.Disk, systemDisk server.Disk, pref *Preference) (server.Disk, error) {
	if pref == nil {
		pref = &Preference{}
	}

	var storageDisk server.Disk

	for _, disk := range disks {
		// Ignore system disk
		if disk.Name == systemDisk.Name {
			continue
		}
		// May ignore USB device
		if disk.Transport == "usb" && pref.IgnoreUSB {
			continue
		}
		// Take the biggest available disk
		if disk.Size > storageDisk.Size {
			storageDisk = disk
		}
	}
	if storageDisk.Size == 0 {
		return storageDisk, errors.New("Cannot select storage disk")
	}
	return storageDisk, nil
}
