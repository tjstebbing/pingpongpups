{ pkgs }:

let
  server1 = pkgs.stdenv.mkDerivation {
    name = "ping";
    version = "1.0";

    buildInputs = [ pkgs.go ];

    src = ./.;

    unpackPhase = "true";

    buildPhase = ''
      export GO111MODULE=off
      export GOCACHE=$(pwd)/.gocache
      mkdir -p $out/bin
      cd $src
      go build -o $out/bin/ping ping.go
    '';

    installPhase = ''
      echo "Binary built at $out/bin/ping"
    '';
  };
in
{
  inherit server1;
}
