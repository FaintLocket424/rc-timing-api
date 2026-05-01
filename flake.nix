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
      serverPort = "8080";

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
              "-X main.Port=${serverPort}"
            ];

            tags = [ "release" ];
          });

          debug = pkgs.buildGoModule (commonArgs // {
            pname = "${commonArgs.pname}-debug";

            dontStrip = true;

            buildFlags = [ "-gcflags=all=-N -l" ];

            ldflags = [
              "-X main.Version=${commonArgs.version}-${gitRev}-debug"
              "-X main.Port=${serverPort}"
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
              echo "Starting the server on port ${serverPort} in development mode..."
              PORT="${serverPort}" go run ./cmd/server/main.go "$@"
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

          testRateLimitingScript = pkgs.writeShellApplication {
            name = "test-rate-limiting";
            runtimeInputs = [ pkgs.vegeta pkgs.jq ];
            text = ''
              RATE="''${1:-5}"
              DURATION="''${2:-10}"

              echo "Attacking at $RATE req/s for $DURATION seconds..."

              echo "GET http://localhost:${serverPort}/api/v1/ping" | \
                  vegeta attack -rate="$RATE"/1s -duration="$DURATION"s | \
                  vegeta encode | \
                  jq -r '. | "\(.code) \(.latency / 1000000)ms \(.error)"'
            '';
          };

          lineCounterScript = pkgs.writeShellApplication {
            name = "line-count";
            runtimeInputs = [ pkgs.findutils pkgs.coreutils ];
            text = ''
              # Get file lists safely
              src_files=$(find . -name "*.go" ! -name "*_test.go" -not -path "*/vendor/*" -type f)
              test_files=$(find . -name "*_test.go" -not -path "*/vendor/*" -type f)

              # Count files
              num_src=$(echo "$src_files" | grep -c ".go" || echo 0)
              num_test=$(echo "$test_files" | grep -c "_test.go" || echo 0)

              # Count lines (ignoring empty lines)
              lines_src=$(echo "$src_files" | xargs cat 2>/dev/null | grep -cve '^$' || echo 0)
              lines_test=$(echo "$test_files" | xargs cat 2>/dev/null | grep -cve '^$' || echo 0)

              # Pretty print output
              printf "+------------------+-------+-------+\n"
              printf "| Category         | Files | Lines |\n"
              printf "+------------------+-------+-------+\n"
              printf "| Source Code      | %5d | %5d |\n" "$num_src" "$lines_src"
              printf "| Test Code        | %5d | %5d |\n" "$num_test" "$lines_test"
              printf "+------------------+-------+-------+\n"
              printf "| Total            | %5d | %5d |\n" $((num_src + num_test)) $((lines_src + lines_test))
              printf "+------------------+-------+-------+\n"
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
            testRateLimitingScript
            lineCounterScript
            listFunctions
          ];

          commandNames = builtins.concatStringsSep ", " (builtins.map (p: p.name) commandBinaries);
        in
        {
          default = pkgs.mkShell {
            packages = with pkgs; [
              go
              gopls
              delve
            ] ++ commandBinaries;

            shellHook = ''
              export PORT="${serverPort}"
              echo "---"
              echo "Go development environment loaded."
              echo "Target Port: $PORT"
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
