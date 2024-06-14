{ pkgs ? import <nixpkgs> {} }: pkgs.mkShell {
  buildInputs = with pkgs; [
    go
    just
    docker-compose
    sqlite
  ];
}
