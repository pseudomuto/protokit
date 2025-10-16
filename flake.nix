{
  description = "Protokit - A starter kit for building protoc-plugins";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
    }:
    flake-utils.lib.eachSystem
      [
        "x86_64-linux"
        "aarch64-linux"
        "x86_64-darwin"
        "aarch64-darwin"
      ]
      (
        system:
        let
          pkgs = import nixpkgs {
            inherit system;
            config.allowUnfree = true;
          };

          # Pick language/tool versions here (adjust as you like)
          go = pkgs.go_1_25;

          # Common build utils
          buildUtils = with pkgs; [
            go-task
            golangci-lint
            goreleaser
            protobuf
          ];
        in
        {
          # `nix develop` drops you into this shell
          devShells.default = pkgs.mkShell {
            packages = [
              go
              buildUtils
            ];

            CGO_ENABLED = "0";

            # Helpful prompt when you enter the shell
            shellHook = ''
              echo "▶ Dev shell ready on ${system}"
              echo "   Go:    $(${go}/bin/go version)"
            '';
          };
        }
      );
}
