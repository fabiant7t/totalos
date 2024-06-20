package command

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/fabiant7t/totalos/internal/totalos"
	"github.com/fabiant7t/totalos/pkg/remotecommand"
	"golang.org/x/crypto/ssh"
)

// Arch returns the machine architecture (see `uname -m`)
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

// Disks returns all disks
func Disks(m remotecommand.Machine, cb ssh.HostKeyCallback) ([]totalos.Disk, error) {
	cmd := `
	  lsblk -o NAME,SERIAL,SIZE,TYPE,MODEL --json -b \
		| jq -r '.blockdevices | map(select(.type == "disk"))'
  `
	var disks []totalos.Disk
	stdout, err := remotecommand.Command(m, cmd, cb)
	if err != nil {
		return disks, err
	}
	if err := json.Unmarshal(stdout, &disks); err != nil {
		return disks, err
	}
	return disks, nil
}

// InstallImage downloads the ISO image URL and writes it to the given device.
func InstallImage(m remotecommand.Machine, isoImageURL, device string, cb ssh.HostKeyCallback) error {
	cmd := fmt.Sprintf(`
	  wget %s -O talos-metal.iso \
		&& dd if=talos-metal.iso of=%s bs=4M \
		&& sync`, isoImageURL, device)
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

// MAC returns the address of the ethernet interface
func MAC(m remotecommand.Machine, cb ssh.HostKeyCallback) (string, error) {
	cmd := `ip -j link show eth0 | jq -r ".[0].address"`

	stdout, err := remotecommand.Command(m, cmd, cb)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(stdout)), nil
}

// SystemUUID as defined in SMBIOS data
func SystemUUID(m remotecommand.Machine, cb ssh.HostKeyCallback) (string, error) {
	cmd := `dmidecode -s system-uuid`
	stdout, err := remotecommand.Command(m, cmd, cb)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(stdout)), nil
}

// Storage returns a slice with the gigabytes storage per disk
func Storage(m remotecommand.Machine, cb ssh.HostKeyCallback) ([]totalos.GigaByte, error) {
	cmd := `
	  lsblk -b --json \
		| jq -r '.blockdevices | map(select(.type =="disk")) | .[] | .size' | awk '{printf "%d\n",$1 / 1000000000}'
	`
	stdout, err := remotecommand.Command(m, cmd, cb)
	if err != nil {
		return []totalos.GigaByte{}, err
	}
	var storage []totalos.GigaByte
	for _, s := range strings.Split(strings.TrimSpace(string(stdout)), "\n") {
		gb, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return []totalos.GigaByte{}, err
		}
		storage = append(storage, totalos.GigaByte(gb))
	}
	return storage, nil
}

// CPUName returns the CPU name
func CPUName(m remotecommand.Machine, cb ssh.HostKeyCallback) (string, error) {
	cmd := `
	  dmidecode -t processor \
		| grep Version\: \
		| cut -d ':' -f 2- \
		| awk '{print}'
	`
	stdout, err := remotecommand.Command(m, cmd, cb)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(stdout)), nil
}

// CPUCores returns the amount of real CPU cores
func CPUCores(m remotecommand.Machine, cb ssh.HostKeyCallback) (int, error) {
	cmd := `
	  dmidecode -t processor \
		| grep Core\ Count: \
		| awk '{print $3}'
	`
	stdout, err := remotecommand.Command(m, cmd, cb)
	if err != nil {
		return 0, err
	}
	count, err := strconv.ParseInt(strings.TrimSpace(string(stdout)), 10, 64)
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

// CPUThreads returns the amount of CPU threads
func CPUThreads(m remotecommand.Machine, cb ssh.HostKeyCallback) (int, error) {
	cmd := `
	  dmidecode -t processor \
		| grep Thread\ Count: \
		| awk '{print $3}'
	`
	stdout, err := remotecommand.Command(m, cmd, cb)
	if err != nil {
		return 0, err
	}
	threads, err := strconv.ParseInt(strings.TrimSpace(string(stdout)), 10, 64)
	if err != nil {
		return 0, err
	}
	return int(threads), nil
}

// Memory returns the amount of RAM
func Memory(m remotecommand.Machine, cb ssh.HostKeyCallback) (totalos.GigaByte, error) {
	cmd := `
	  dmidecode -t memory \
		| grep -i size \
		| awk '{sum += $2} END {print sum}'
	`
	stdout, err := remotecommand.Command(m, cmd, cb)
	if err != nil {
		return 0, err
	}
	mem, err := strconv.ParseInt(strings.TrimSpace(string(stdout)), 10, 64)
	if err != nil {
		return 0, err
	}
	return totalos.GigaByte(mem), nil
}

// Reboot triggers a reboot. It does not return anything, since the
// machine should be offline and not be able to have an SSH chat :)
func Reboot(m remotecommand.Machine, cb ssh.HostKeyCallback) {
	_, _ = remotecommand.Command(m, `shutdown -r now`, cb)
}