{
  description = "OpenGrid Timing Bridge development flake";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-25.11";

    flake-utils.url = "github:numtide/flake-utils";

    treefmt-nix = {
      url = "github:numtide/treefmt-nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };

    git-hooks = {
      url = "github:cachix/git-hooks.nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs = { self, nixpkgs, flake-utils, treefmt-nix, git-hooks }:
    let
      defaultPort = 4998;
    in
    flake-utils.lib.eachDefaultSystem
      (system:
        let
          debugPort = 8080;

          commonArgs = {
            pname = "opengrid-bridge";
            version = "0.1.0";
            src = ./.;
            subPackages = [ "cmd/opengrid-bridge" ];
            vendorHash = "sha256-IUww4dVh6MWv7SQITZY+LKgOT1ReVgmEHY3bl2f/tM4=";
          };

          pkgs = import nixpkgs { inherit system; };

          treefmtEval = treefmt-nix.lib.evalModule pkgs {
            projectRootFile = "flake.nix";
            programs = {
              nixpkgs-fmt.enable = true;
              gofumpt.enable = true;
              mdformat.enable = true;
            };
            settings.global.excludes = [ "**/testdata/**" ];
          };

          golangciYaml = pkgs.writeText "golangci.yaml" ''
            version: "2"

            run:
              timeout: 5m
              skip-dirs:
                - testdata

            linters:
              disable-all: false
              enable:
                - revive
                - gosec
                - gocritic
                - nilerr
                - errcheck
                - godot
                - goconst
          '';

          custom-golangci-lint = pkgs.writeShellApplication {
            name = "golangci-lint";
            runtimeInputs = [ pkgs.golangci-lint ];
            text = ''
              # If the first argument is 'run', inject our Nix-managed config
              if [ "''${1:-}" = "run" ]; then
                exec golangci-lint "$1" --config="${golangciYaml}" "''${@:2}"
              else
                # For commands like 'version' or 'cache', just pass normally
                exec golangci-lint "$@"
              fi
            '';
          };

          pre-commit-check = git-hooks.lib.${system}.run {
            src = ./.;
            hooks = {
              treefmt = {
                enable = true;
                package = treefmtEval.config.build.wrapper;
              };
              golangci-lint = {
                enable = true;
                package = custom-golangci-lint;
              };
            };
          };
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
                    name = "opengrid-bridge-debug-runner";
                    text = ''
                      echo "Running Debug Build on Port ${toString debugPort}..."
                      export PORT=${toString debugPort}
                      exec ${packages.debug}/bin/${commonArgs.pname} "$@"
                    '';
                  };
                in
                "${runner}/bin/opengrid-bridge-debug-runner";

              meta.description = "Run the API in debug mode";
            };

            default = {
              type = "app";
              program = "${packages.default}/bin/${commonArgs.pname}";
              meta.description = "Run the API in release mode";
            };
          };

          checks = {
            formatting = treefmtEval.config.build.check self;

            go-test = pkgs.buildGoModule (commonArgs // {
              pname = "opengrid-bridge-tests";
              buildPhase = "";
              installPhase = "touch $out";
            });

            golangci-lint = pkgs.buildGoModule (commonArgs // {
              pname = "opengrid-bridge-lint";
              doCheck = false;

              buildPhase = ''
                export GOCACHE=$TMPDIR/go-cache
                export GOLANGCI_LINT_CACHE=$TMPDIR/golangci-cache

                ${custom-golangci-lint}/bin/golangci-lint run
              '';

              installPhase = "touch $out";
            });
          };

          formatter = treefmtEval.config.build.wrapper;

          devShells =
            let
              runDevServerScript = pkgs.writeShellApplication {
                name = "run-dev-server";

                runtimeInputs = [ pkgs.go ];

                text = ''
                  echo "Starting the server on port ${toString debugPort} in development mode..."
                  PORT="${toString debugPort}" go run ./cmd/opengrid-bridge/main.go "$@"
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
                  custom-golangci-lint
                ] ++ commandBinaries;

                shellHook = ''
                  ${pre-commit-check.shellHook}
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
      nixosModules.opengrid-bridge = { config, lib, pkgs, ... }:
        let cfg = config.services.opengrid-bridge; in {
          options.services.opengrid-bridge = {
            enable = lib.mkEnableOption "OpenGrid Timing Bridge Service";
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
              description = "The opengrid-bridge package to use.";
            };
          };

          config = lib.mkIf cfg.enable {
            environment.systemPackages = [ cfg.package ];

            systemd.services.opengrid-bridge = {
              description = "OpenGrid Timing Bridge Server";
              after = [ "network.target" ];
              wantedBy = [ "multi-user.target" ];

              serviceConfig = {
                ExecStart = "${cfg.package}/bin/opengrid-bridge";
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
