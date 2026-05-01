{
  description = "RC Timing API";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-25.11";

    treefmt-nix = {
      url = "github:numtide/treefmt-nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs = { self, nixpkgs, treefmt-nix }:
    let
      supportedSystems = [ "x86_64-linux" "aarch64-linux" "x86_64-darwin" "aarch64-darwin" ];

      forAllSystems = nixpkgs.lib.genAttrs supportedSystems;

      nixpkgsFor = forAllSystems (system: import nixpkgs { inherit system; });

      commonArgs = {
        pname = "rc-timing-api";
        version = "0.1.0";
        src = ./.;
        subPackages = [ "cmd/server" ];
        vendorHash = "sha256-IUww4dVh6MWv7SQITZY+LKgOT1ReVgmEHY3bl2f/tM4=";

        postInstall = ''
          mv $out/bin/server $out/bin/$pname
        '';
      };
    in
    {
      packages = forAllSystems (system:
        let
          pkgs = nixpkgsFor.${system};

          gitRev = self.shortRev or self.dirtyShortRev or "unknown";
        in
        rec {
          release = pkgs.buildGoModule (commonArgs // {
            ldflags = [
              "-s"
              "-w"

              "-X main.GinMode=release"
              "-X main.Version=${commonArgs.version}-${gitRev}"
            ];

            tags = [ "release" ];
          });

          debug = pkgs.buildGoModule (commonArgs // {
            pname = "${commonArgs.pname}-debug";

            dontStrip = true;

            buildFlags = [ "-gcflags=all=-N -l" ];

            ldflags = [
              "-X main.Version=${commonArgs.version}-${gitRev}-debug"
            ];

            tags = [ "debug" ];
          });

          default = release;
        });

      checks = forAllSystems
        (system:
          let pkgs = nixpkgsFor.${system};
          in {
            go-test = pkgs.buildGoModule (commonArgs // {
              pname = "rc-timing-api-tests";
              buildPhase = "";
              installPhase = "touch $out";
            });

            golangci-lint = pkgs.buildGoModule (commonArgs // {
              pname = "rc-timing-api-lint";

              buildPhase = ''
                export GOCACHE=$TMPDIR/go-cache
                export GOLANGCI_LINT_CACHE=$TMPDIR/golangci-cache

                ${pkgs.golangci-lint}/bin/golangci-lint run --timeout 5m --skip-dirs "testdata"
              '';

              installPhase = "touch $out";
            });
          });

      devShells = forAllSystems (system:
        let
          pkgs = nixpkgsFor.${system};

          runServerScript = pkgs.writeShellApplication {
            name = "run-server";

            runtimeInputs = [ pkgs.go ];

            text = ''
              echo "Starting the server in development mode..."
              go run ./cmd/server/main.go "$@"
            '';
          };

          mirrorBBKScript = pkgs.writeShellApplication {
            name = "mirror-bbk";
            runtimeInputs = [ pkgs.wget ];
            text = ''
              if [ "$#" -ne 2 ]; then
                echo "Usage: $0 <base_url> <output_path>" >&2
                exit 1
              fi

              BASE_URL="$1"
              OUTPUT_PATH="$2"

              {
                echo "$BASE_URL/"
                echo "$BASE_URL/liveraceres.htm"
                echo "$BASE_URL/liveresults.htm"
                echo "$BASE_URL/liveschedule.htm"
                echo "$BASE_URL/livecompets.htm"
              } | wget \
                --mirror \
                --convert-links \
                --page-requisites \
                --no-parent \
                --user-agent="Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36" \
                --directory-prefix="$OUTPUT_PATH" \
                --input-file=- \
                --wait=0.1 \
                --random-wait
            '';
          };

          listFunctions = pkgs.writeShellApplication {
            name = "lsfunc";
            runtimeInputs = [ pkgs.coreutils ];
            text = ''
              echo ${commandNames}
            '';
          };

          commandBinaries = [
            runServerScript
            mirrorBBKScript
            listFunctions
          ];

          commandNames = builtins.concatStringsSep ", " (builtins.map (p: p.name) commandBinaries);
        in
        {
          default = pkgs.mkShell {
            packages = with pkgs; [
              go
              vegeta
              gopls
              delve
            ] ++ commandBinaries;

            shellHook = ''
              echo "---"
              echo "Go development environment loaded."
              echo "Available custom commands: $(lsfunc)"
              echo "---"
            '';
          };
        });

      formatter = forAllSystems (system:
        let
          pkgs = nixpkgsFor.${system};

          treefmtEval = treefmt-nix.lib.evalModule pkgs {
            projectRootFile = "flake.nix";
            programs.nixpkgs-fmt.enable = true;
            programs.gofmt.enable = true;

            settings.global.excludes = [ "**/testdata/**" ];
          };
        in
        treefmtEval.config.build.wrapper
      );
    };
}
