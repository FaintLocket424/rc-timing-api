{ pkgs, self, commonArgs, treefmtEval, custom-golangci-lint }:

{
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
}
