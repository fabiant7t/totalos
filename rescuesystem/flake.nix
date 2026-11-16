{
  description = "Minimal TotalOS rescue live ISO";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-26.05";

  outputs =
    { self, nixpkgs }:
    let
      system = "x86_64-linux";
      pkgs = nixpkgs.legacyPackages.${system};
      rescue = self.nixosConfigurations.rescuesystem;
    in
    {
      nixosConfigurations.rescuesystem = nixpkgs.lib.nixosSystem {
        inherit system;
        modules = [ ./rescuesystem/iso.nix ];
      };
      packages.${system} = {
        rescuesystem = rescue.config.system.build.isoImage;
        default = self.packages.${system}.rescuesystem;
      };
      checks.${system} = {
        config = import ./rescuesystem/tests/config.nix {
          inherit pkgs;
          inherit (nixpkgs) lib;
          config = rescue.config;
          banner = ./rescuesystem/banner.txt;
        };
        runtime = import ./rescuesystem/tests/runtime.nix {
          inherit pkgs;
          liveSystem = ./rescuesystem/live-system.nix;
          banner = ./rescuesystem/banner.txt;
        };
      };
      formatter.${system} = pkgs.nixfmt;
    };
}
