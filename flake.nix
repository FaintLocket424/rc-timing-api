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

      # downloadRemotes = pkgs.writeScriptBin "download-remotes" ''
      #   #!/usr/bin/env fish

      #   set urls \
      #       https://forcc.co.uk/live/ \
      #       http://www.rcresults.org/bingham/ \
      #       http://www.rcresults.org/dms/ \
      #       http://www.rcresults.org/rhr/ \
      #       http://www.rcresults.org/smcc/ \
      #       http://www.rcresults.org/york/ \
      #       http://www.rcresults.org/worksop/

      #   for url in $urls
      #       echo "Mirroring: $url"
      #       go run ./cmd/mirror -url $url -path "./internal/scraper/bbk/testdata/" -remove-non-htm
      #   end
      # '';

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
      ];

      commandNames = builtins.concatStringsSep ", " (builtins.map (p: p.name) commandBinaries);
    in
    {
      devShells.${system}.default = pkgs.mkShell {
        packages = with pkgs; [
          go
          gopls
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
