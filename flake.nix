{
  description = "RC Timing API Development Environment";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-25.11";

  outputs = { nixpkgs, ... }:
    let
      system = "x86_64-linux";
      pkgs = import nixpkgs { system = system; };

      runServer = pkgs.writeShellApplication {
        name = "run-server";
        runtimeInputs = [ pkgs.go ];
        text = ''
          go run ./cmd/server "$@"
        '';
      };

      downloadRemotes = pkgs.writeShellApplication {
        name = "download-remotes";

        text = ''
          TARGET_DIR="./internal/scraper/bbk/testdata/$(date +%Y-%m-%d)"

          mkdir -p "$TARGET_DIR"

          urls=(
            "https://forcc.co.uk/live/"
            "http://www.rcresults.org/bingham/"
            "http://www.rcresults.org/dms/"
            "http://www.rcresults.org/rhr/"
            "http://www.rcresults.org/smcc/"
            "http://www.rcresults.org/york/"
            "http://www.rcresults.org/worksop/"
          )

          for url in "''${urls[@]}"; do
            echo "Mirroring: $url"
            go run ./cmd/mirror -url "$url" -path "$TARGET_DIR" -remove-non-htm
          done
        '';
      };

      test-bbk-scraper = pkgs.writeShellApplication {
        name = "test-bbk-scraper";
        runtimeInputs = [ pkgs.go ];
        text = ''
          go test ./internal/scraper/bbk "$@"
        '';
      };

      test-rate-limiting = pkgs.writeShellApplication {
        name = "test-rate-limiting";
        runtimeInputs = [ pkgs.vegeta pkgs.jq ];
        text = ''
          RATE="''${1:-5}"
          DURATION="''${2:-10}"

          echo "Attacking at $RATE req/s for $DURATION seconds..."

          echo "GET http://localhost:8080/api/v1/ping" | \
              vegeta attack -rate="$RATE"/1s -duration="$DURATION"s | \
              vegeta encode | \
              jq -r '. | "\(.code) \(.latency / 1000000)ms \(.error)"'
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
        runServer
        downloadRemotes
        test-bbk-scraper
        listFunctions
        test-rate-limiting
      ];

      commandNames = builtins.concatStringsSep ", " (builtins.map (p: p.name) commandBinaries);
    in
    {
      devShells.${system}.default = pkgs.mkShell {
        packages = with pkgs; [
          go
          vegeta
        ] ++ commandBinaries;

        shellHook = ''
          echo "---"
          echo "Go development environment loaded."
          echo "Available custom commands: $(lsfunc)"
          echo "---"
        '';
      };
    };
}
