{
  description = "A very basic flake";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
  };

  outputs = { nixpkgs, ... }:
	let
		system = "x86_64-linux";
		pkgs = nixpkgs.legacyPackages.${system};
	in
	{
		devShells.${system}.default = pkgs.mkShell {
			packages = with pkgs; [
				go
				gopls
				gotools
			];

			shellHook = ''
				alias build="go build -o ./bin/hera ."
			'';
		};
  };
}
