{
  description = "RC Timing API";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-25.11";

    flake-utils.url = "github:numtide/flake-utils";

    treefmt-nix = {
      url = "github:numtide/treefmt-nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs = { self, nixpkgs, flake-utils, treefmt-nix }:
    let
      defaultPort = 4998;
    in
    flake-utils.lib.eachDefaultSystem
      (system:
        let
          debugPort = 8080;

          commonArgs = {
            pname = "rc-timing-api";
            version = "0.1.0";
            src = ./.;
            subPackages = [ "cmd/rc-timing-api" ];
            vendorHash = "sha256-IUww4dVh6MWv7SQITZY+LKgOT1ReVgmEHY3bl2f/tM4=";
          };

          pkgs = import nixpkgs { inherit system; };
        in
        rec {
          packages =
            let
              gitRev = self.shortRev or self.dirtyShortRev or "unknown";
            in
            rec {
              release = pkgs.buildGoModule (commonArgs // {
                ldflags = [
                  "-s"
                  "-w"

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
            };

          apps = {
            debug = {
              type = "app";
              program =
                let
                  runner = pkgs.writeShellApplication {
                    name = "rc-timing-api-debug-runner";
                    text = ''
                      echo "Running Debug Build on Port ${toString debugPort}..."
                      export PORT=${toString debugPort}
                      exec ${packages.debug}/bin/${commonArgs.pname} "$@"
                    '';
                  };
                in
                "${runner}/bin/rc-timing-api-debug-runner";
            };

            default = {
              type = "app";
              program = "${packages.default}/bin/${commonArgs.pname}";
            };
          };

          checks = {
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
          };

          formatter =
            let
              treefmtEval = treefmt-nix.lib.evalModule pkgs {
                projectRootFile = "flake.nix";
                programs.nixpkgs-fmt.enable = true;
                programs.gofmt.enable = true;

                settings.global.excludes = [ "**/testdata/**" ];
              };
            in
            treefmtEval.config.build.wrapper
          ;

          devShells =
            let

              runDevServerScript = pkgs.writeShellApplication {
                name = "run-dev-server";

                runtimeInputs = [ pkgs.go ];

                text = ''
                  echo "Starting the server on port ${toString debugPort} in development mode..."
                  PORT="${toString debugPort}" go run ./cmd/rc-timing-api/main.go "$@"
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

                  echo "GET http://localhost:${toString debugPort}/api/v1/ping" | \
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
                runDevServerScript
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
                  echo "---"
                  echo "Go development environment loaded."
                  echo "Default Dev Port: ${toString debugPort}"
                  echo "Available custom commands: $(lsfunc)"
                  echo "---"
                '';
              };
            };

        }
      ) // {
      nixosModules.rc-timing-api = { config, lib, pkgs, ... }:
        let cfg = config.services.rc-timing-api; in {
          options.services.rc-timing-api = {
            enable = lib.mkEnableOption "RC Timing API Service";
            port = lib.mkOption {
              type = lib.types.port;
              default = defaultPort;
              description = "The port the API should listen on.";
            };
            openFirewall = lib.mkOption {
              type = lib.types.bool;
              default = false;
              description = "Whether to automatically open the firewall for the API port.";
            };
            package = lib.mkOption {
              type = lib.types.package;
              default = self.packages.${pkgs.system}.default;
              description = "The rc-timing-api package to use.";
            };
          };

          config = lib.mkIf cfg.enable {
            environment.systemPackages = [ cfg.package ];

            systemd.services.rc-timing-api = {
              description = "RC Timing API Server";
              after = [ "network.target" ];
              wantedBy = [ "multi-user.target" ];

              serviceConfig = {
                ExecStart = "${cfg.package}/bin/rc-timing-api";
                Restart = "always";
                RestartSec = "3";
                DynamicUser = true;
              };

              environment = {
                GIN_MODE = "release";
                PORT = toString cfg.port;
              };
            };

            networking.firewall.allowedTCPPorts = lib.mkIf cfg.openFirewall [ cfg.port ];
          };
        };
    };

}
