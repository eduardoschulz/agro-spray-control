{
  description = "Painel de Pulverização - Desenvolvimento";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-24.11";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils, ... }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs {
          inherit system;
        };
      in {
        devShells = {
          backend = pkgs.mkShell {
            name = "backend-shell";
            packages = with pkgs; [
              go
            ];
          };

          frontend = pkgs.mkShell {
            name = "frontend-shell";
            packages = with pkgs; [
              nodejs_20
              yarn
            ];
          };
        };
      }
    );
}

