{
	inputs = {
		nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
		flake-utils.url = "github:numtide/flake-utils";
		gomod2nix.url = "github:nix-community/gomod2nix";
		gomod2nix.inputs.nixpkgs.follows = "nixpkgs";
		gomod2nix.inputs.flake-utils.follows = "flake-utils";
	};

	outputs =
		{ self, nixpkgs, flake-utils, gomod2nix }:

		flake-utils.lib.eachDefaultSystem (system:
			with gomod2nix.legacyPackages.${system};
			with nixpkgs.legacyPackages.${system}.extend (self: super: {
				go = super.go_1_22;
			});
			{
				devShell = mkShell {
					packages = [
						go
						gopls
						go-tools
						sqlc
					];
				};
			}
		);
}
