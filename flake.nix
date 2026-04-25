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

      count-lines = pkgs.writeShellApplication {
        name = "count-lines";
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

      coverage-stats = pkgs.writeShellApplication {
        name = "coverage-stats";
        runtimeInputs = [ pkgs.go ];
        text = ''
          EXCLUDE_DIRS="cmd|examples"
          COVERAGE_THRESHOLD=100.0

          TARGET_PKG="''${1:-./...}"

          # We redirect standard output of go test to /dev/null,
          # but let standard error (where error messages go) through.
          if [ "$TARGET_PKG" == "./..." ]; then
            mapfile -t PKGS < <(go list ./... | grep -E -v "$EXCLUDE_DIRS")
            if ! go test -coverprofile=coverage.out "''${PKGS[@]}" > /dev/null; then
              echo -e "\n[!] Tests failed. Coverage report will not be generated."
              exit 1
            fi
          else
            if ! go test -coverprofile=coverage.out "$TARGET_PKG" > /dev/null; then
              echo -e "\n[!] Tests failed. Coverage report will not be generated."
              exit 1
            fi
          fi

          # Now the output starts cleanly here
          echo -e "=== Analysis for: $TARGET_PKG ==="

          echo -e "\n=== 5 Worst Covered Functions ==="
          go tool cover -func=coverage.out | grep -v "total:" | sort -k3 -n | head -n 5

          echo -e "\n=== 5 Functions below $COVERAGE_THRESHOLD% ==="
          go tool cover -func=coverage.out | awk -v limit="$COVERAGE_THRESHOLD" '
            $1 != "total:" {
              pct = $NF; sub(/%/, "", pct);
              if (pct + 0 < limit) { print $0 }
            }
          ' | head -n 5

          echo -e "\n=== 5 Functions with 0.0% coverage ==="
          go tool cover -func=coverage.out | awk '$1 != "total:" && $NF == "0.0%"' | head -n 5

          rm coverage.out
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
        test-rate-limiting
        count-lines
        coverage-stats
        listFunctions
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
