package services

import (
	"fmt"
	"strings"

	"github.com/fabiant7t/totalos/pkg/remotecommand"
	"golang.org/x/crypto/ssh"
)

func Arch(m remotecommand.Machine, hostKeyCallback ssh.HostKeyCallback) (string, error) {
	stdout, err := remotecommand.Command(m, "uname -m", hostKeyCallback)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(stdout)), nil
}

// SoftwareRAIDNotExists ensures that there are no software RAID partitions (/dev/md*)
func SoftwareRAIDNotExists(m remotecommand.Machine, hostKeyCallback ssh.HostKeyCallback) error {
	cmd := `
	  fdisk -l \
		|  grep -E '^Disk /dev/md[0-9]+' \
		&& mdadm --stop /dev/md/* \
		|| echo Software RAID already missing
	`
	_, err := remotecommand.Command(m, cmd, hostKeyCallback)
	return err
}

// WipeFileSystemSignatures erases all available signatures of all NVMe and SATA drives.
// It does not shred the data, though!
func WipeFileSystemSignatures(m remotecommand.Machine, hostKeyCallback ssh.HostKeyCallback) error {
	cmd := `
		for satadisk in $(ls /dev/sd*); do
		    wipefs -fa ${satadisk};
		done;
		for nvmedisk in $(ls /dev/nvme*n1); do
		    wipefs -fa ${nvmedisk};
		done;
  `
	_, err := remotecommand.Command(m, cmd, hostKeyCallback)
	return err
}

// NominateInstallDisk queries all bulk devices, sorts them alphabetically
// by their serial number and returns the first disk.
// The result should be deterministic.
func NominateInstallDisk(m remotecommand.Machine, hostKeyCallback ssh.HostKeyCallback) (string, error) {
	cmd := `
		name=$(
		    lsblk --json -o NAME,SERIAL,TYPE \
				| jq -r '.blockdevices | map(select(.type == "disk")) | sort_by(.serial) | .[0].name'
    );
		echo /dev/${name}
	`
	stdout, err := remotecommand.Command(m, cmd, hostKeyCallback)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(stdout)), nil
}

// InstallImage downloads the ISO image URL and writes it to the given disk.
func InstallImage(m remotecommand.Machine, isoImageURL string, disk string, hostKeyCallback ssh.HostKeyCallback) error {
	cmd := fmt.Sprintf(`
	  wget %s -O talos-metal.iso \
		&& dd if=talos-metal.iso of=%s bs=4M \
		&& sync`, isoImageURL, disk)
	_, err := remotecommand.Command(m, cmd, hostKeyCallback)
	if err != nil {
		return err
	}
	return nil
}

// Reboot triggers a reboot. It does not return anything, since the
// machine should be offline and not be able to have an SSH chat :)
func Reboot(m remotecommand.Machine, hostKeyCallback ssh.HostKeyCallback) {
	_, _ = remotecommand.Command(m, "reboot", hostKeyCallback)
}
