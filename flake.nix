{
  description = "OpenGrid Timing Bridge development flake";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-26.05";

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

  outputs = { self, nixpkgs, flake-utils, treefmt-nix, git-hooks }@inputs:
    flake-utils.lib.eachDefaultSystem
      (system:
        let
          pkgs = import nixpkgs { inherit system; };
        in
        import ./nix/default.nix { inherit self inputs pkgs system; }
      ) // {
      nixosModules = {
        opengrid-bridge = (import ./nix/module.nix { inherit self; }).nixos;
        default = self.nixosModules.opengrid-bridge;
      };
    };
}
