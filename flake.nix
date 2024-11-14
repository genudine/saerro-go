{
  description = "saerro";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-parts.url = "github:hercules-ci/flake-parts";
  };

  outputs = inputs: inputs.flake-parts.lib.mkFlake { inherit inputs; } {
    systems = [ "x86_64-linux" "aarch64-linux" ];
    perSystem = { config, self', pkgs, lib, system, ... }: {
      devShells.default = import ./shell.nix { inherit pkgs; };

      packages = let 
        vendorHash = "sha256-A5hZxo0zZ3w6qryV24PjYaKQatN2G2heyuee6QaU55M=";
      in rec {
        default = saerro;
        saerro = pkgs.buildGoModule {
          inherit vendorHash;
          name = "saerro-ws";
          src = ./.;
          subPackages = [
            "cmd/ws"
            "cmd/pruner"
          ];
        };
        ws = pkgs.ociTools.buildContainer {
          args = [
            "${saerro}/ws"
          ];
        };
        pruner = pkgs.ociTools.buildContainer {
          args = [
            "${saerro}/pruner"
          ];
        };
      };
    };
  };
}