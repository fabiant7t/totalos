# totalos rescue distro

totalos was supposed to install Talos on servers running the Hetzner
Rescue System, a Debian based distro that bare metal machines network
boot into when started in rescue mode.

While totalos is limited to use binaries available on this rescue
system, any bare metal machine should be an installation target.
`totalos rescue` is a custom rescue system, which can be disk dumped
on a USB stick used to live boot any bare metal server (like all
those Lenovo ThinkCentres I could not resist buying for my homelab).

The user is `root`, the password is `totalos` and SSH login is permitted.
This means you must be careful, boot it and install talos right away.

It is based on archiso with these additional packages:
- coreutils
- dmidecode
- gawk
- gptfdisk
- grep
- iproute2
- jq
- mdadm
- parted
- sed
- systemd-sysvcompat
- util-linux
- wget
- xfsprogs
- xz
