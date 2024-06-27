{
    description = "Cutlink Nix Flake";

    inputs = {
        # Use stable branch of nixpkgs
        nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    };

    outputs = { self, nixpkgs, ... }: let
        pkgName = "cutlink";

        supportedSystems = [ "x86_64-linux" "aarch64-linux" "x86_64-darwin" "aarch64-darwin" ];

        # current version of project
        version = "1.0.0";

        forAllSystems = nixpkgs.lib.genAttrs supportedSystems;

        nixpkgsFor = forAllSystems (system: import nixpkgs { inherit system; });
    in {

        packages = forAllSystems (system: let
            pkgs = nixpkgsFor.${system};
        in {
            ${pkgName} = pkgs.buildGoModule {
                name = "${pkgName}";
                inherit version;
                src = ./.;

                # if `nix build` command fails because of wrong vendorHash, it will
                # report (and print) the correct vendorHash, so you MUST replace it.
                # Then re-run `nix build` command.
                vendorHash = "sha256-j5nsgQe2ViwBA757YYjK49txFxFeCh69YJK116EJJtI=";
                checkPhase = "";

                postInstall = ''
                mv $out/bin/cmd $out/bin/cutlink
                cp config.toml $out
                cp -r ./ui $out
                '';
            };
        });

        # Interactive shell for development
        devShells = forAllSystems (system: let
            pkgs = nixpkgsFor.${system};
        in { 
            default = pkgs.mkShell {
                packages = with pkgs; [
                    ## go compiler and gofmt
                    go
                    gnumake

                    ## debugging stuff
                    gdb
                ];
            };
        });


      defaultPackage = forAllSystems (system: self.packages.${system}.${pkgName});
    };
}
