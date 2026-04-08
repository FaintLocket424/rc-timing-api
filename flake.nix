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

      listFunctions = pkgs.writeShellApplication {
        name = "lsfunc";
        runtimeInputs = [ pkgs.coreutils ];
        text = ''
          echo ${commandNames}
        '';
      };

      commandBinaries = [
        runServer
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
