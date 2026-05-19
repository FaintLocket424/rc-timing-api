{ pkgs, debugPort, custom-golangci-lint, pre-commit-check }:

let
  runDevServerScript = pkgs.writeShellApplication {
    name = "run-dev-server";

    runtimeInputs = [ pkgs.go ];

    text = ''
      echo "Starting the server on port ${toString debugPort} in development mode..."
      go run ./cmd/opengrid-bridge/main.go -port ${toString debugPort} "$@"
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
      prek
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
}
