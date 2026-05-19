{ self }:

{
  # NixOS System Module
  nixos = { config, lib, pkgs, ... }:
    let
      cfg = config.services.opengrid-bridge;
    in
    {
      options.services.opengrid-bridge = {
        enable = lib.mkEnableOption "OpenGrid Timing Bridge Service";
        port = lib.mkOption {
          type = lib.types.port;
          default = 4998;
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
            ExecStart = "${cfg.package}/bin/opengrid-bridge -port ${toString cfg.port}";
            Restart = "always";
            RestartSec = "3";
            DynamicUser = true;
          };

          environment = {
            GIN_MODE = "release";
          };
        };

        networking.firewall.allowedTCPPorts = lib.mkIf cfg.openFirewall [ cfg.port ];
      };
    };
}
