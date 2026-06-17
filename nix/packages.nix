{ pkgs, self, commonArgs }:

let
  gitRev = self.shortRev or self.dirtyShortRev or "unknown";
in
rec {
  release = pkgs.buildGoModule (commonArgs // {
    ldflags = [
      "-s"
      "-w"
      "-X main.Version=${commonArgs.version}"
      "-X main.Commit=${gitRev}"
    ];

    tags = [ "release" ];
  });

  debug = pkgs.buildGoModule (commonArgs // {
    pname = "${commonArgs.pname}-debug";
    dontStrip = true;
    buildFlags = [ "-gcflags=all=-N -l" ];
    ldflags = [
      "-X main.Version=${commonArgs.version}-debug"
      "-X main.Commit=${gitRev}"
    ];
    tags = [ "debug" ];
  });

  default = release;
}
