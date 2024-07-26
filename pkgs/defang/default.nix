# This file was generated by GoReleaser. DO NOT EDIT.
# vim: set ft=nix ts=2 sw=2 sts=2 et sta
{
system ? builtins.currentSystem
, lib
, fetchurl
, installShellFiles
, stdenvNoCC
, unzip
}:
let
  shaMap = {
    x86_64-linux = "1z3n27w6vgc0xwncdnpi8lz44247k4k0c20a52mdwcwd7aj70xhj";
    aarch64-linux = "0lfv1s9hx521dc1469p5myb4rffaifk6rnmhsgln7fg96kgjjjim";
    x86_64-darwin = "18jfx4561p2ghrn74wp53xcbbc3xi3zxpdbl2f4npyip936snng7";
    aarch64-darwin = "18jfx4561p2ghrn74wp53xcbbc3xi3zxpdbl2f4npyip936snng7";
  };

  urlMap = {
    x86_64-linux = "https://github.com/DefangLabs/defang/releases/download/v0.5.36/defang_0.5.36_linux_amd64.tar.gz";
    aarch64-linux = "https://github.com/DefangLabs/defang/releases/download/v0.5.36/defang_0.5.36_linux_arm64.tar.gz";
    x86_64-darwin = "https://github.com/DefangLabs/defang/releases/download/v0.5.36/defang_0.5.36_macOS.zip";
    aarch64-darwin = "https://github.com/DefangLabs/defang/releases/download/v0.5.36/defang_0.5.36_macOS.zip";
  };
in
stdenvNoCC.mkDerivation {
  pname = "defang";
  version = "0.5.36";
  src = fetchurl {
    url = urlMap.${system};
    sha256 = shaMap.${system};
  };

  sourceRoot = ".";

  nativeBuildInputs = [ installShellFiles unzip ];

  installPhase = ''
    mkdir -p $out/bin
    cp -vr ./defang $out/bin/defang
  '';
  postInstall = ''
    installShellCompletion --cmd defang \
    --bash <($out/bin/defang completion bash) \
    --zsh <($out/bin/defang completion zsh) \
    --fish <($out/bin/defang completion fish)
  '';

  system = system;

  meta = {
    description = "Defang is the easiest way for developers to create and deploy their containerized applications";
    homepage = "https://defang.io/";
    license = lib.licenses.mit;

    sourceProvenance = [ lib.sourceTypes.binaryNativeCode ];

    platforms = [
      "aarch64-darwin"
      "aarch64-linux"
      "x86_64-darwin"
      "x86_64-linux"
    ];
  };
}
