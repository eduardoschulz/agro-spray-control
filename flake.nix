{
  description = "Flake Provisorio";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-24.11";
    # Se quiser usar uma versão estável do Go
    go.url = "github:golang/go";
    nodejs.url = "github:nixos/nixpkgs/nixos-24.11#nodejs";
  };

  outputs = { self, nixpkgs, go, nodejs, ... } @ inputs: {

    # Ambiente de desenvolvimento para o backend (Go)
    devShells.x86_64-linux.backend = nixpkgs.mkShell {
      buildInputs = [
        go 
      ];
      shellHook = ''
        export SHELL=$(which zsh)
        export GOPATH=$PWD/go
        export PATH=$GOPATH/bin:$PATH
      '';
    };

    # Ambiente de desenvolvimento para o frontend (Vue.js)
    devShells.x86_64-linux.frontend = nixpkgs.mkShell {
      buildInputs = [
        nodejs # Inclui Node.js para o Vue.js
      ];
      shellHook = ''
        export SHELL=$(which zsh)
        export PATH=$PWD/node_modules/.bin:$PATH
      '';
    };

    # Recomendação para o ambiente de produção
    packages.x86_64-linux.backend = nixpkgs.buildGoModule {
      pname = "agrospray-backend";
      version = "0.1.0";
      src = ./backend;
      vendorSha256 = "sha256-hash"; # Adicione o hash dos pacotes de dependências se necessário
    };

    packages.x86_64-linux.frontend = nodejs.buildPackages.nodePackages.createFrontend {
      pname = "agrospray-frontend";
      version = "0.1.0";
      src = ./frontend;
    };
  };
}

