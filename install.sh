#!/bin/sh
set -eu

repo="ragnarok22/timer-go"
binary="timer"
version="${VERSION:-latest}"
install_dir="${INSTALL_DIR:-$HOME/.local/bin}"

detect_os() {
  case "$(uname -s)" in
    Darwin) printf '%s' darwin ;;
    Linux) printf '%s' linux ;;
    *)
      printf 'Unsupported OS: %s\n' "$(uname -s)" >&2
      exit 1
      ;;
  esac
}

detect_arch() {
  case "$(uname -m)" in
    x86_64 | amd64) printf '%s' amd64 ;;
    arm64 | aarch64) printf '%s' arm64 ;;
    *)
      printf 'Unsupported architecture: %s\n' "$(uname -m)" >&2
      exit 1
      ;;
  esac
}

download() {
  url="$1"
  output="$2"

  if command -v curl >/dev/null 2>&1; then
    curl -fsSL "$url" -o "$output"
    return
  fi

  if command -v wget >/dev/null 2>&1; then
    wget -q "$url" -O "$output"
    return
  fi

  printf 'Either curl or wget is required.\n' >&2
  exit 1
}

os="$(detect_os)"
arch="$(detect_arch)"
asset="${binary}_${os}_${arch}.tar.gz"

if [ "$version" = "latest" ]; then
  url="https://github.com/${repo}/releases/latest/download/${asset}"
else
  url="https://github.com/${repo}/releases/download/${version}/${asset}"
fi

tmp_dir="$(mktemp -d)"
trap 'rm -rf "$tmp_dir"' EXIT INT TERM

download "$url" "$tmp_dir/$asset"
tar -xzf "$tmp_dir/$asset" -C "$tmp_dir"

mkdir -p "$install_dir"
mv "$tmp_dir/$binary" "$install_dir/$binary"
chmod 755 "$install_dir/$binary"

printf 'Installed %s to %s\n' "$binary" "$install_dir/$binary"
printf 'Make sure %s is in your PATH.\n' "$install_dir"
