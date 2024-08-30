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

## How to build the ISO?

```sh
sudo pacman -S archiso
sudo mkarchiso -v -w /tmp/totalosrescue rescue
sudo dd if=out/totalosrescue-2024.08.30-x86_64.iso of=/dev/<usb-device>
sync
```

Now you can unplug the USB stick and boot a server. There is no need
to interact, you should be able to ssh into the box:

```sh
$ ssh root@192.168.178.206
Warning: Permanently added '192.168.178.206' (ED25519) to the list of known hosts.
root@192.168.178.206's password:

   █████              █████              ████
  ░░███              ░░███              ░░███
  ███████    ██████  ███████    ██████   ░███   ██████   █████
 ░░░███░    ███░░███░░░███░    ░░░░░███  ░███  ███░░███ ███░░
   ░███    ░███ ░███  ░███      ███████  ░███ ░███ ░███░░█████
   ░███ ███░███ ░███  ░███ ███ ███░░███  ░███ ░███ ░███ ░░░░███
   ░░█████ ░░██████   ░░█████ ░░████████ █████░░██████  ██████
    ░░░░░   ░░░░░░     ░░░░░   ░░░░░░░░ ░░░░░  ░░░░░░  ░░░░░░


   ████████   ██████   █████   ██████  █████ ████  ██████
  ░░███░░███ ███░░███ ███░░   ███░░███░░███ ░███  ███░░███
   ░███ ░░░ ░███████ ░░█████ ░███ ░░░  ░███ ░███ ░███████
   ░███     ░███░░░   ░░░░███░███  ███ ░███ ░███ ░███░░░
   █████    ░░██████  ██████ ░░██████  ░░████████░░██████
  ░░░░░      ░░░░░░  ░░░░░░   ░░░░░░    ░░░░░░░░  ░░░░░░

  Fabian Topfstedt | 2024 | https://github.com/fabiant7t/totalos/

The password for user `root` is `totalos`. SSH login is permitted.
Quickly use `totalos` to install Talos on this machine and NEVER
leave the server running this system unattended.

Last login: Fri Aug 30 19:01:45 2024
Connected Network Interface: enp0s31f6
  MAC Address: 6c:4b:90:01:af:fe
  IPv4 Address: 192.168.178.206
  IPv6 Address: 2001:871:265:1a5c:6e4b:90ff:fe0f:6018/64
```
