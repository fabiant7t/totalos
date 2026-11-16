# TotalOS rescue system

A minimal, RAM-backed NixOS live ISO for preparing bare-metal servers with
[TotalOS](../README.md). It provides root SSH access and the tools TotalOS needs
to inspect hardware and write a Talos Linux image. Booting the ISO does not itself
run TotalOS or install Talos.

## Build

Requirements: an `x86_64-linux` build environment, Nix with `nix-command` and
`flakes` enabled, GNU Make, and network access to fetch dependencies.

From this directory:

```sh
make
```

The image is copied to `dist/rescuesystem.iso`. The Makefile builds from a temporary
copy of the required Nix files, avoiding Git flake filtering and unrelated files
in the parent repository. Prefer this entry point over a direct `nix build .`.

Other targets:

```sh
make checksum  # Build the ISO and print its SHA-256 checksum
make up        # Update the pinned nixpkgs input in flake.lock
make clean     # Remove dist/rescuesystem.iso (not the Nix store build outputs)
```

## Boot and connect

Boot `dist/rescuesystem.iso` using a server's virtual media feature, or write it to
a USB drive. The image supports legacy BIOS and UEFI boot on x86-64 machines.

Wired interfaces matching `en*` or `eth*` request an IPv4 address through DHCP.
A DHCP server is therefore needed for the default network setup. At the local
console, log in and find the assigned address with `ip -4 addr`.

Default credentials:

- User: `root`
- Password: `totalos`
- SSH port: `22`

```sh
ssh root@<server-ip>
```

**Use an isolated or trusted network.** Root password authentication is enabled
and the default password is public. Run `passwd` after login to change it for the
current session. The SSH host key is generated under `/run/ssh` and is not
persistent across reboots; verify a changed key before updating a client's
known-hosts entry.

## Use with TotalOS

Run the TotalOS CLI from your workstation, not inside the rescue environment.
For example, after building the CLI in the repository root:

```sh
./totalos --ip <server-ip> --user root --password 'totalos'
```

Use the current password if you changed it. See the [main README](../README.md)
for image selection, configuration, networking, reporting, and reboot options.

**TotalOS is destructive:** it wipes disk signatures and writes the Talos image
to the selected disk. Verify the target server before running it.

## Live environment

- Hostname: `rescuesystem`; locale: `en_US.UTF-8`.
- Console-only system, with no X server or UDisks service.
- Hardware support and redistributable firmware enabled.
- Tools include Bash, coreutils, grep, sed, awk, jq, util-linux, dmidecode,
  iproute2, mdadm, wget, xz, zstd, systemd, and CA certificates.
- The writable root filesystem, including `/root` and `/mnt`, lives in RAM.
  Changes and downloaded files are lost on reboot. Internal disks are not used
  automatically as scratch space.
- The root tmpfs has a default size ceiling of **75% of RAM**. This is a limit,
  not a reservation. Allow enough memory for image downloads and the running
  system; the ceiling does not guarantee that this much memory is available.

## Configuration

| File | Purpose |
| --- | --- |
| `flake.nix` | Defines the x86-64 NixOS configuration and ISO package. |
| `flake.lock` | Pins the nixpkgs revision used for builds. |
| `iso.nix` | Configures boot support, compression, and the RAM-backed root. |
| `live-system.nix` | Configures packages, networking, console, and SSH access. |
| `banner.txt` | Supplies the console login banner. |
| `Makefile` | Builds, copies, checksums, updates, and removes the ISO. |

Change credentials or networking in `live-system.nix` before building if the
standard defaults do not fit your environment. The root tmpfs ceiling is
controlled by the `rescueSystem.rootTmpfsSize` option defined in `iso.nix`.
