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
		return "", fmt.Errorf("Remote command Arch failed: %w", err)
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
	if _, err := remotecommand.Command(m, cmd, cb); err != nil {
		return fmt.Errorf("Remote command SoftwareRAIDNotExists failed: %w", err)
	}
	return nil
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
	if _, err := remotecommand.Command(m, cmd, cb); err != nil {
		return fmt.Errorf("Remote command WipeFileSystemSignatures failed: %w", err)
	}
	return nil
}

// Disks returns all disks
func Disks(m remotecommand.Machine, cb ssh.HostKeyCallback) ([]totalos.Disk, error) {
	cmd := `
	  lsblk -o NAME,SERIAL,SIZE,TYPE,MODEL,TRAN,WWN --json -b \
		| jq -r '.blockdevices | map(select(.type == "disk"))'
  `
	var disks []totalos.Disk
	stdout, err := remotecommand.Command(m, cmd, cb)
	if err != nil {
		return disks, fmt.Errorf("Remote command Disks failed: %w", err)
	}
	if err := json.Unmarshal(stdout, &disks); err != nil {
		return disks, fmt.Errorf("Remote command Disks failed: %w", err)
	}
	return disks, nil
}

// InstallImage downloads the raw.xz image URL and writes it to the given device.
func InstallImage(m remotecommand.Machine, isoImageURL, device string, cb ssh.HostKeyCallback) error {
	cmd := fmt.Sprintf(`
	  wget %s -O talos-metal.raw.xz \
		&& cat talos-metal.raw.xz \
		|  xz -d \
		|  dd of=%s bs=4M \
		&& sync`, isoImageURL, device)
	if _, err := remotecommand.Command(m, cmd, cb); err != nil {
		return fmt.Errorf("Remote command InstallImage failed: %w", err)
	}
	return nil
}

// FormatXFS formats the full device with one XFS partition.
// Sets disk UUID (default is "61291e61-291e-6129-1e61-291e61291e00"),
// partition label (default is "storage") and
// partition UUID (default is "61291e61-291e-6129-1e61-291e61291e01").
func FormatXFS(m remotecommand.Machine, device, diskUUID, partLabel, partUUID string, cb ssh.HostKeyCallback) error {
	if diskUUID == "" {
		diskUUID = "61291e61-291e-6129-1e61-291e61291e00"
	}
	if partLabel == "" {
		partLabel = "storage"
	}
	if partUUID == "" {
		partUUID = "61291e61-291e-6129-1e61-291e61291e01"
	}
	cmd := fmt.Sprintf(`
	  export DEVICE=%s \
		&& wipefs -af ${DEVICE} \
		&& parted ${DEVICE} --script mklabel gpt \
		&& sgdisk --disk-guid=%s ${DEVICE} \
		&& parted ${DEVICE} --script mkpart primary xfs %s %s \
		&& mkfs.xfs -f -L %s -m uuid=%s ${DEVICE}*1
	`, device, diskUUID, "0%", "100%", partLabel, partUUID)
	if _, err := remotecommand.Command(m, cmd, cb); err != nil {
		return fmt.Errorf("Remote command FormatXFS failed: %w", err)
	}
	return nil
}

// EthernetDeviceName returns the name of the ethernet device,
// like eth0 or enp0s31f6 (when kernel is configured to use
// predictable network interface names).
func EthernetDeviceName(m remotecommand.Machine, cb ssh.HostKeyCallback) (string, error) {
	cmd := `
		ip -o link show \
		| awk -F': ' '/: (en|eth)/{print $2; exit}'
	`
	stdout, err := remotecommand.Command(m, cmd, cb)
	if err != nil {
		return "", fmt.Errorf("Remote command EthernetDeviceName failed: %w", err)
	}
	return strings.TrimSpace(string(stdout)), nil
}

// IPv4 address of ethernet device
func IPv4(m remotecommand.Machine, cb ssh.HostKeyCallback) (net.IP, error) {
	cmd := `
	  ip -4 -j a show \
		| jq -r '.[] | select(.ifname | startswith("en") or startswith("eth")) | .addr_info[].local'
	`
	stdout, err := remotecommand.Command(m, cmd, cb)
	if err != nil {
		return nil, fmt.Errorf("Remote command IPv4 failed: %w", err)
	}
	return net.ParseIP(strings.TrimSpace(string(stdout))), nil
}

// IPv4Netmask returns the IPv4 netmask of ethernet device
func IPv4Netmask(m remotecommand.Machine, cb ssh.HostKeyCallback) (net.IP, error) {
	cmd := `
	  ip -4 -j a show \
		| jq -r '.[] | select(.ifname | startswith("en") or startswith("eth")) | .addr_info[].prefixlen'
	`
	stdout, err := remotecommand.Command(m, cmd, cb)
	if err != nil {
		return nil, fmt.Errorf("Remote command IPv4Netmask failed: %w", err)
	}
	prefixlen, err := strconv.ParseInt(strings.TrimSpace(string(stdout)), 10, 32)
	if err != nil {
		return nil, fmt.Errorf("Remote command IPv4Netmask failed: %w", err)
	}
	mask := net.CIDRMask(int(prefixlen), 32)
	return net.IPv4(mask[0], mask[1], mask[2], mask[3]), nil
}

// IPv4Gateway returns the default gateway IPv4 of ethernet device
func IPv4Gateway(m remotecommand.Machine, cb ssh.HostKeyCallback) (net.IP, error) {
	cmd := `
		ip -j -4 route show \
		| jq -r '.[] | select(.dev | startswith("en") or startswith("eth")) | select(.dst == "default") | .gateway'
	`
	stdout, err := remotecommand.Command(m, cmd, cb)
	if err != nil {
		return nil, fmt.Errorf("Remote command IPv4Gateway failed: %w", err)
	}
	return net.ParseIP(strings.TrimSpace(string(stdout))), nil
}

// MAC returns the address of the ethernet interface
func MAC(m remotecommand.Machine, cb ssh.HostKeyCallback) (string, error) {
	cmd := `
	  ip -j link show \
		| jq -r '.[] | select(.ifname | startswith("en") or startswith("eth")) | .address'
	`

	stdout, err := remotecommand.Command(m, cmd, cb)
	if err != nil {
		return "", fmt.Errorf("Remote command MAC failed: %w", err)
	}
	return strings.TrimSpace(string(stdout)), nil
}

// SystemUUID as defined in SMBIOS data
func SystemUUID(m remotecommand.Machine, cb ssh.HostKeyCallback) (string, error) {
	cmd := `dmidecode -s system-uuid`
	stdout, err := remotecommand.Command(m, cmd, cb)
	if err != nil {
		return "", fmt.Errorf("Remote command SystemUUID failed: %w", err)
	}
	return strings.TrimSpace(string(stdout)), nil
}

// Storage returns a slice with the gigabytes of storage per disk
func Storage(m remotecommand.Machine, cb ssh.HostKeyCallback) ([]totalos.GigaByte, error) {
	cmd := `
	  lsblk -b --json \
		| jq -r '.blockdevices | map(select(.type =="disk")) | .[] | .size' | awk '{printf "%d\n",$1 / 1000000000}'
	`
	stdout, err := remotecommand.Command(m, cmd, cb)
	if err != nil {
		return []totalos.GigaByte{}, fmt.Errorf("Remote command Storage failed: %w", err)
	}
	var storage []totalos.GigaByte
	for _, s := range strings.Split(strings.TrimSpace(string(stdout)), "\n") {
		gb, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return []totalos.GigaByte{}, fmt.Errorf("Remote command Storage failed: %w", err)
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
		return "", fmt.Errorf("Remote command CPUName failed: %w", err)
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
		return 0, fmt.Errorf("Remote command CPUCores failed: %w", err)
	}
	count, err := strconv.ParseInt(strings.TrimSpace(string(stdout)), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("Remote command CPUCores failed: %w", err)
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
		return 0, fmt.Errorf("Remote command CPUThreads failed: %w", err)
	}
	threads, err := strconv.ParseInt(strings.TrimSpace(string(stdout)), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("Remote command CPUThreads failed: %w", err)
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
		return 0, fmt.Errorf("Remote command Memory failed: %w", err)
	}
	mem, err := strconv.ParseInt(strings.TrimSpace(string(stdout)), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("Remote command Memory failed: %w", err)
	}
	return totalos.GigaByte(mem), nil
}

// Reboot triggers a reboot. It does not return anything, since the
// machine should be offline and not be able to have an SSH chat :)
func Reboot(m remotecommand.Machine, cb ssh.HostKeyCallback) {
	_, _ = remotecommand.Command(m, `shutdown -r now`, cb)
}

// Sets talos.config in grub.cfg
func SetConfigURL(m remotecommand.Machine, configURL, device string, cb ssh.HostKeyCallback) error {
	replacer := strings.NewReplacer(
		"&", "\\&",
		"/", "\\/",
		":", "\\:",
		"{", "\\{",
		"}", "\\}",
	)
	replacement := replacer.Replace(configURL)
	cmd := fmt.Sprintf(`
    mount %sp3 /mnt \
		&& cp /mnt/grub/grub.cfg /mnt/grub/grub.cfg.orig \
		&& sed 's/talos.platform=metal/talos.platform=metal talos.config=%s/g' /mnt/grub/grub.cfg.orig \
		> /mnt/grub/grub.cfg \
		&& umount /mnt
	`, device, replacement)
	if _, err := remotecommand.Command(m, cmd, cb); err != nil {
		return fmt.Errorf("Remote command SetConfigURL failed: %w", err)
	}
	return nil
}
