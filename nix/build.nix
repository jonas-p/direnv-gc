# { buildGoModule, lib }:
with import <nixpkgs> {};

buildGoModule rec {
  pname = "direnv-gc";
  version = "0.1.2";

  src = ./..;

  vendorSha256 = null;
}
