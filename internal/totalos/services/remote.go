package services

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/fabiant7t/totalos/pkg/remotecommand"
	"golang.org/x/crypto/ssh"
)

func Arch(m remotecommand.Machine, cb ssh.HostKeyCallback) (string, error) {
	stdout, err := remotecommand.Command(m, "uname -m", cb)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(stdout)), nil
}

// SoftwareRAIDNotExists ensures that there are no software RAID partitions (/dev/md*)
func SoftwareRAIDNotExists(m remotecommand.Machine, cb ssh.HostKeyCallback) error {
	cmd := `
	  fdisk -l \
		|  grep -E '^Disk /dev/md[0-9]+' \
		&& mdadm --stop /dev/md/* \
		|| echo Software RAID already missing
	`
	_, err := remotecommand.Command(m, cmd, cb)
	return err
}

// WipeFileSystemSignatures erases all available signatures of all NVMe and SATA drives.
// It does not shred the data, though!
func WipeFileSystemSignatures(m remotecommand.Machine, cb ssh.HostKeyCallback) error {
	cmd := `
		for satadisk in $(ls /dev/sd*); do
		    wipefs -fa ${satadisk};
		done;
		for nvmedisk in $(ls /dev/nvme*n1); do
		    wipefs -fa ${nvmedisk};
		done;
  `
	_, err := remotecommand.Command(m, cmd, cb)
	return err
}

// NominateInstallDisk queries all bulk devices, sorts them alphabetically
// by their serial number and returns the first disk.
// The result should be deterministic.
func NominateInstallDisk(m remotecommand.Machine, cb ssh.HostKeyCallback) (string, error) {
	cmd := `
		name=$(
		    lsblk --json -o NAME,SERIAL,TYPE \
				| jq -r '.blockdevices | map(select(.type == "disk")) | sort_by(.serial) | .[0].name'
    );
		echo /dev/${name}
	`
	stdout, err := remotecommand.Command(m, cmd, cb)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(stdout)), nil
}

// InstallImage downloads the ISO image URL and writes it to the given disk.
func InstallImage(m remotecommand.Machine, isoImageURL string, disk string, cb ssh.HostKeyCallback) error {
	cmd := fmt.Sprintf(`
	  wget %s -O talos-metal.iso \
		&& dd if=talos-metal.iso of=%s bs=4M \
		&& sync`, isoImageURL, disk)
	_, err := remotecommand.Command(m, cmd, cb)
	return err
}

// IPv4 address of eth0
func IPv4(m remotecommand.Machine, cb ssh.HostKeyCallback) (net.IP, error) {
	cmd := `
	  ip -4 -j a show eth0 \
		| jq -r '.[0].addr_info[].local'
	`
	stdout, err := remotecommand.Command(m, cmd, cb)
	if err != nil {
		return nil, err
	}
	return net.ParseIP(strings.TrimSpace(string(stdout))), nil
}

// IPv4Netmask returns the IPv4 netmask of eth0
func IPv4Netmask(m remotecommand.Machine, cb ssh.HostKeyCallback) (net.IP, error) {
	cmd := `
	  ip -4 -j a show eth0 \
		| jq -r '.[0].addr_info[].prefixlen'
	`
	stdout, err := remotecommand.Command(m, cmd, cb)
	if err != nil {
		return nil, err
	}
	prefixlen, err := strconv.ParseInt(strings.TrimSpace(string(stdout)), 10, 32)
	if err != nil {
		return nil, err
	}
	mask := net.CIDRMask(int(prefixlen), 32)
	return net.IPv4(mask[0], mask[1], mask[2], mask[3]), nil
}

// IPv4Gateway returns the default gateway IPv4 of eth0
func IPv4Gateway(m remotecommand.Machine, cb ssh.HostKeyCallback) (net.IP, error) {
	cmd := `
		ip -j -4 route show \
		| jq -r '.[] | select(.dev == "eth0" and .dst == "default") | .gateway'
	`
	stdout, err := remotecommand.Command(m, cmd, cb)
	if err != nil {
		return nil, err
	}
	return net.ParseIP(strings.TrimSpace(string(stdout))), nil
}

// Reboot triggers a reboot. It does not return anything, since the
// machine should be offline and not be able to have an SSH chat :)
func Reboot(m remotecommand.Machine, cb ssh.HostKeyCallback) {
	_, _ = remotecommand.Command(m, "reboot", cb)
}
