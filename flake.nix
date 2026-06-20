{
  description = "🐙 Less 3";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
  };

  outputs = { self, nixpkgs, ... }:
    let
      systems = [
        "x86_64-linux"
        "x86_64-darwin"
        "aarch64-darwin"
      ];
      forAllSystems = nixpkgs.lib.genAttrs systems;
    in
    {
      packages = forAllSystems (system:
        let
          pkgs = nixpkgs.legacyPackages.${system};
	  version = "0.1.0";
          commit = self.rev or "dirty";
        in
        {
          default = pkgs.buildGoModule {
            pname = "hera";
            inherit version;
            src = self;
            modules = ./gomod2nix.toml;

            ldflags = [
              "-X main.version=${version}"
              "-X main.commit=${commit}"
	      "-s"
	      "-w"
            ];

            vendorHash = "sha256-GdV+7ccktqbsDwfNLBU8fEsOjtHXBKiqjn7m1lMFvUU=";

            meta = {
              description = "A tactical turn-based game. Made with ♡";
              homepage = "https://github.com/IwnuplyNotTyan/hera";
              mainProgram = "hera";
            };
          };
	});
      devShells = forAllSystems (system:
        let
          pkgs = nixpkgs.legacyPackages.${system};
        in
        {
          default = pkgs.mkShell {
            packages = with pkgs; [
              go
              gopls
              gotools
              golangci-lint
            ];
          };
        });
    };
}
