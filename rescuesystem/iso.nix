{
  config,
  lib,
  modulesPath,
  ...
}:
{
  imports = [
    ./live-system.nix
    (modulesPath + "/installer/cd-dvd/iso-image.nix")
  ];

  options.rescueSystem.rootTmpfsSize = lib.mkOption {
    type = lib.types.str;
    default = "75%";
    description = ''
      Maximum size of the ephemeral RAM-backed root, including /root and /mnt.
      Talos image downloads consume real RAM in addition to the running system.
      This is a ceiling, not reserved memory: provision enough RAM for both the
      compressed image and the system. No internal disk is used for scratch space.
    '';
  };

  config = {
    # The ISO builder derives its actual filename from baseName, not fileName.
    image.baseName = lib.mkForce "rescuesystem";
    isoImage = {
      makeBiosBootable = true;
      makeEfiBootable = true;
      makeUsbBootable = true;
      squashfsCompression = "zstd -Xcompression-level 6";
    };
    hardware.enableAllHardware = true;
    lib.isoFileSystems."/" = lib.mkImageMediaOverride {
      fsType = "tmpfs";
      options = [
        "mode=0755"
        "size=${config.rescueSystem.rootTmpfsSize}"
      ];
    };
  };
}
