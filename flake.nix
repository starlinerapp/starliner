{
  description = "Starliner development environment";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs { inherit system; };

        darwinDeps = pkgs.lib.optionals pkgs.stdenv.isDarwin [
          pkgs.apple-sdk
          (pkgs.darwinMinVersionHook "11.0")
        ];

        commonPkgs = with pkgs; [
          git
          curl
          jq
          gnumake
        ];

        serverPkgs =
          with pkgs;
          [
            go
            golangci-lint
            protobuf
            protoc-gen-go
            protoc-gen-go-grpc
            sqlc
            pgformatter
          ]
          ++ darwinDeps;

        serverShellHook = ''
          unset GOROOT
          export GOPATH="''${GOPATH:-$HOME/go}"
          export GOBIN="$GOPATH/bin"
          export PATH="$PATH:$GOBIN"

          _nix_install_go_tool() {
            local name="$1" pkg="$2"
            if ! command -v "$name" &>/dev/null; then
              echo "starliner: installing $name via go install ..."
              go install "$pkg" || echo "starliner: warning — failed to install $name"
            fi
          }
          _nix_install_go_tool swag    "github.com/swaggo/swag/cmd/swag@v1.16.6"
          _nix_install_go_tool arch-go "github.com/arch-go/arch-go/v2@v2.1.2"
          unset -f _nix_install_go_tool
        '';

      in
      {
        formatter = pkgs.nixfmt-tree;
        devShells = {
          default = pkgs.mkShell {
            name = "starliner";
            packages = commonPkgs ++ serverPkgs;
            shellHook = serverShellHook;
          };
        };
      }
    );
}
