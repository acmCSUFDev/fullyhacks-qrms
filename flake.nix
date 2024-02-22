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
			let
				pkgs = nixpkgs.legacyPackages.${system}.extend (self: super: {
					go = super.go_1_22;
				});

				buildGoApplication = gomod2nix.legacyPackages.${system}.buildGoApplication;
				gomod2nixTool = gomod2nix.packages.${system}.default.override {
					inherit (pkgs) go;
				};
			in
			with pkgs; {
				packages.default = buildGoApplication {
					inherit (pkgs) go;

					name = "fullyhacks-qrms";

					pwd = ./.;
					src = ./.;
					modules = ./gomod2nix.toml;
					subPackages = [ "." ];

					preBuild = "go generate ./...";

					nativeBuildInputs = [
						sqlc
						gomod2nixTool
					];
				};
				devShell = mkShell {
					packages = [
						go
						gopls
						go-tools
						sqlc
						gomod2nixTool
					];
				};
			}
		);
}
