#!/bin/sh
# Installs hat for Linux and WSL.
# curl -fsSL https://raw.githubusercontent.com/baodq97/hatrack/main/install.sh | sh
set -eu
[ "$(uname -s)" = Linux ] || { echo "hatrack runs on Linux/WSL and Windows only" >&2; exit 1; }
case "$(uname -m)" in
  x86_64|amd64) arch=amd64 ;;
  aarch64|arm64) arch=arm64 ;;
  *) echo "unsupported CPU: $(uname -m)" >&2; exit 1 ;;
esac
dir="${HAT_INSTALL_DIR:-$HOME/.local/bin}"
mkdir -p "$dir"
curl -fsSL "https://github.com/baodq97/hatrack/releases/latest/download/hat-linux-$arch" -o "$dir/hat.tmp"
chmod +x "$dir/hat.tmp" && mv "$dir/hat.tmp" "$dir/hat"
echo "installed $dir/hat"
case ":$PATH:" in *":$dir:"*) ;; *) echo "add $dir to your PATH" ;; esac
