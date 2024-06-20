package command

import (
	"errors"
	"sort"

	"github.com/fabiant7t/totalos/internal/totalos"
)

// SelectSystemDisk takes all disks, sorts them by their serial
// alphabetically and returns the first one. The result is
// deterministic.
func SelectSystemDisk(disks []totalos.Disk) (totalos.Disk, error) {
	var serials []string
	disksBySerial := make(map[string]totalos.Disk, len(disks))

	for _, disk := range disks {
		serials = append(serials, disk.Serial)
		disksBySerial[disk.Serial] = disk
	}
	sort.Strings(serials)

	if len(serials) > 0 {
		return disksBySerial[serials[0]], nil
	}
	return totalos.Disk{}, errors.New("Cannot select system disk")
}

// SelectStorageDisk finds and returns the biggest available disk which
// is not the system disk.
func SelectStorageDisk(disks []totalos.Disk, systemDisk totalos.Disk) (totalos.Disk, error) {
	var storageDisk totalos.Disk

	for _, disk := range disks {
		// Ignore system disk
		if disk.Name == systemDisk.Name {
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
