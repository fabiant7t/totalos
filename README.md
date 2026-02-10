# totalos

`totalos` is a small CLI that installs Talos Linux on bare-metal servers that are booted into a rescue environment (Hetzner Rescue System or the custom `totalos rescue` ISO). It connects over SSH, inventories the machine, wipes disk signatures, writes a Talos metal image, optionally injects a config URL and a static IP kernel option, and emits a JSON report (optionally POSTed to a webhook).

**This tool is destructive.** It wipes filesystem signatures and writes a disk image with `dd`.

**What It Does**
- Connects via SSH to a rescue-booted server.
- Collects hardware and network details (CPU, memory, disks, NIC, DMI info, IPv4).
- Picks a system disk deterministically (lowest serial, ignoring USB) and writes a Talos raw image to it.
- Optionally injects `talos.config=<url>` into `grub.cfg`.
- Optionally injects `ip=<...>` static network config into `grub.cfg`.
- Selects a storage disk (largest non-system disk) for reporting.
- Prints a JSON report and optionally POSTs it to a webhook.
- Optionally reboots the server.

**How It Chooses the Image**
- If `--image` is provided, it is used as-is.
- Otherwise, it queries the Talos GitHub releases API and picks the latest non-draft, non-prerelease `metal-*.raw.zst` matching the machine architecture.

**Disk Selection Rules**
- System disk: smallest serial (alphabetical) among non-USB disks.
- Storage disk: largest non-system disk (USB not excluded here).

**Requirements**
Remote rescue system must provide:
- `ssh` access as root
- `lsblk`, `jq`, `dmidecode`, `ip`, `udevadm`, `fdisk`, `mdadm`, `wipefs`, `wget`, `xz` or `zstd`, `dd`, `mount`, `umount`

Local build environment:
- Go toolchain matching `go.mod` (`go1.24.2` toolchain)

**Build**
```sh
go build -o totalos ./cmd/totalos
```

**Usage**
```sh
./totalos \
  --ip 203.0.113.10 \
  --user root \
  --password 'secret' \
  --config https://example.com/talos-config.yaml \
  --webhook https://example.com/hook \
  --static \
  --reboot
```

**Flags**
- `--ip` (required) target server IP
- `--port` SSH port (default `22`)
- `--user` SSH user (default `root`)
- `--password` SSH password (required unless `--key` is set)
- `--key` path to SSH private key (required unless `--password` is set)
- `--image` URL to `raw.xz`, `raw.zst`, or `iso` image (optional)
- `--config` URL to Talos machine config (optional, injected as `talos.config=...`)
- `--webhook` URL to receive JSON report via HTTP POST (optional)
- `--static` set static initial network configuration (adds `ip=...` kernel option)
- `--reboot` reboot server after install
- `--version` print version and exit

**Static Network Option Details**
When `--static` is set, the tool builds an `ip=` kernel command-line entry using the current IPv4 address, netmask, gateway, and interface name. It also sets DNS and NTP:
- DNS: `86.54.11.100` (DNS4EU) and `9.9.9.9` (Quad9)
- NTP: `162.159.200.1` (Cloudflare)

**Report Output**
The tool prints a JSON report to stdout and optionally POSTs it to `--webhook`.

Example structure:
```json
{
  "installation": {
    "image": "https://.../metal-amd64.raw.zst",
    "rebooting": true,
    "config": "https://example.com/talos-config.yaml",
    "static_initial_network_configuration": "...",
    "storage_disk": { "name": "sdb", "size": 2000398934016, "serial": "..." },
    "system_disk": { "name": "sda", "size": 500107862016, "serial": "..." }
  },
  "machine": {
    "arch": "x86_64",
    "hostname": "talos-203-0-113-10",
    "ipv4_network": { "ip": "203.0.113.10", "netmask": "255.255.255.0", "gateway": "203.0.113.1", "cidr": "203.0.113.10/24" },
    "cpu": { "name": "...", "cores": 8, "threads": 16 },
    "memory": { "size_gb": 64 },
    "system": { "manufacturer": "...", "product_name": "...", "uuid": "..." },
    "ethernet": { "device": "enp0s31f6", "mac": "...", "speed_mbps": 1000 }
  }
}
```

**Rescue ISO**
For the custom rescue environment, see `rescue/README.md`.
