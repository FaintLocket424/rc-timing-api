{ self, inputs, pkgs, system }:

let
  debugPort = 8080;

  commonArgs = {
    pname = "opengrid-bridge";
    version = "0.2.0";
    # Pointing to the flake root from inside the nix/ directory
    src = ../.;
    subPackages = [ "cmd/opengrid-bridge" ];
    vendorHash = "sha256-aMSG3Wi66LA7gEbbeNLyozggvVuestrswkBLV3WlNB8=";
  };

  treefmtEval = inputs.treefmt-nix.lib.evalModule pkgs {
    projectRootFile = "flake.nix";
    programs = {
      nixpkgs-fmt.enable = true;
      gofumpt.enable = true;
      mdformat.enable = true;
    };
    settings.global.excludes = [
      "**/testdata/**"
      "cmd/experiment/**"
    ];
  };

  golangciYaml = pkgs.writeText "golangci.yaml" ''
    version: "2"

    run:
      timeout: 5m
      skip-dirs:
        - testdata
        - cmd/experiment

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
        - godoclint

    linters-settings:
      godoclint:
        options:
          require-doc:
            ignore-unexported: false

    issues:
      exclude-use-default: false
      max-issues-per-linter: 3
      max-same-issues: 3
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

  pre-commit-check = inputs.git-hooks.lib.${system}.run {
    src = ../.;

    package = pkgs.prek;

    hooks = {
      treefmt = {
        enable = true;
        package = treefmtEval.config.build.wrapper;
      };
      golangci-lint = {
        enable = true;
        package = custom-golangci-lint;
      };
      nix-build = {
        enable = true;
        name = "Nix Build Check";
        description = "Check that the program builds successfully with Nix";
        entry = "nix build .#default --no-link";
        pass_filenames = false;
        stages = [ "push" ];
      };
    };
  };

  packages = import ./packages.nix {
    inherit pkgs self commonArgs;
  };

  checks = import ./checks.nix {
    inherit pkgs self commonArgs treefmtEval custom-golangci-lint;
  };

  devShells = import ./devshell.nix {
    inherit pkgs debugPort custom-golangci-lint pre-commit-check;
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
              exec ${packages.debug}/bin/${commonArgs.pname} -port ${toString debugPort} "$@"
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

in
{
  inherit packages apps checks devShells;
  formatter = treefmtEval.config.build.wrapper;
}
