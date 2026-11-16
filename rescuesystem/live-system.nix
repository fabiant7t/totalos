{ pkgs, modulesPath, ... }:
{
  imports = [ (modulesPath + "/profiles/minimal.nix") ];

  system.stateVersion = "26.05";
  networking.hostName = "rescuesystem";
  i18n.defaultLocale = "en_US.UTF-8";
  console.font = "${pkgs.terminus_font}/share/consolefonts/ter-v16n.psf.gz";
  services.getty.greetingLine = builtins.readFile ./banner.txt;
  environment.etc.issue.source = ./banner.txt;

  environment.systemPackages = with pkgs; [
    bash
    coreutils
    gnugrep
    gnused
    gawk
    jq
    util-linux
    dmidecode
    iproute2
    mdadm
    wget
    xz
    zstd
    systemd
    cacert
  ];
  security.pki.installCACerts = true;
  hardware.enableRedistributableFirmware = true;

  networking.useDHCP = false;
  networking.useNetworkd = true;
  systemd.network.networks."20-wired" = {
    matchConfig.Name = "en* eth*";
    networkConfig.DHCP = "ipv4";
  };
  services.resolved.enable = true;
  networking.firewall.allowedTCPPorts = [ 22 ];

  users.users.root = {
    shell = pkgs.bashInteractive;
    initialPassword = "totalos";
  };
  services.openssh = {
    enable = true;
    settings = {
      PermitRootLogin = "yes";
      PasswordAuthentication = true;
      KbdInteractiveAuthentication = false;
      PermitEmptyPasswords = false;
      MaxStartups = "100:30:200";
    };
    hostKeys = [
      {
        path = "/run/ssh/ssh_host_ed25519_key";
        type = "ed25519";
      }
    ];
  };
  systemd.services.sshd.unitConfig.StartLimitIntervalSec = 0;
  systemd.services."sshd@".unitConfig.StartLimitIntervalSec = 0;
  systemd.tmpfiles.rules = [ "d /mnt 0755 root root -" ];

  documentation.enable = false;
  services.xserver.enable = false;
  services.udisks2.enable = false;
}
