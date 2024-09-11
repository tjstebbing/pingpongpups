{ pkgs }:

let
  pong = pkgs.stdenv.mkDerivation {
    name = "pong";
    version = "1.0";

    buildInputs = [ pkgs.go ];

    src = ./.;

    unpackPhase = "true";

    buildPhase = ''
      export GO111MODULE=off
      export GOCACHE=$(pwd)/.gocache
      mkdir -p $out/bin
      cd $src
      go build -o $out/bin/pong pong.go
    '';

    installPhase = ''
      echo "Binary built at $out/bin/pong"
    '';
  };
in
{
  inherit pong;
}
