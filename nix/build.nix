# { buildGoModule, lib }:
with import <nixpkgs> {};

buildGoModule rec {
  pname = "direnv-gc";
  version = "0.1.0";

  src = ./..;

  vendorSha256 = null;
}
