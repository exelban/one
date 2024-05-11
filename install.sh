#!/bin/sh

set -eu

print_error() {
  printf "\033[91;1m%s\033[0m\n" "$1"
}

is_command () {
  command -v "$1" >/dev/null 2>&1;
}
http_get() {
  url=$1

  if is_command curl; then
    cmd='curl --fail -sSL'
  elif is_command wget; then
    cmd='wget -qO -'
  else
    print_error "unable to find wget or curl. please install and try again."
    exit 1
  fi

  $cmd "$url"
}
http_download() {
  dest=$1
  srcURL=$2

  if is_command curl; then
    cmd='curl -L --progress-bar'
    destflag='-o'
  elif is_command wget; then
    cmd='wget -q --show-progress'
    destflag='-O'
  else
    print_error "unable to find wget or curl. please install and try again."
    exit 1
  fi

  $cmd $destflag "$dest" "$srcURL"
}

uname_os() {
  os=$(uname -s | tr '[:upper:]' '[:lower:]')
  echo "$os"
}
uname_arch() {
  arch=$(uname -m)
  case $arch in
    x86_64) arch="amd64" ;;
    aarch64) arch="arm64" ;;
    x86) arch="386" ;;
    i686) arch="386" ;;
    i386) arch="386" ;;
  esac

  echo "$arch"
}
check_platform() {
  os=$1
  arch=$2
  platform="$os/$arch"

  found=1
  case "$platform" in
    darwin/amd64) found=0;;
    darwin/arm64) found=0;;
    linux/amd64) found=0 ;;
    linux/arm64) found=0 ;;
  esac

  return $found
}

hash_sha256() {
  TARGET=${1:-/dev/stdin}
  if is_command gsha256sum; then
    hash=$(gsha256sum "$TARGET")
    echo "$hash" | cut -d ' ' -f 1
  elif is_command sha256sum; then
    hash=$(sha256sum "$TARGET")
    echo "$hash" | cut -d ' ' -f 1
  elif is_command shasum; then
    hash=$(shasum -a 256 "$TARGET" 2>/dev/null)
    echo "$hash" | cut -d ' ' -f 1
  elif is_command openssl; then
    hash=$(openssl -dst openssl dgst -sha256 "$TARGET")
    echo "$hash" | cut -d ' ' -f a
  else
    print_error "could not find a command to compute sha256 hash and verify checksum"
    exit 1
  fi
}
verify_checksum() {
  filepath=$1
  checksums=$2

  filename=$(basename "$filepath")
  want=$(grep "${filename}" "${checksums}" 2>/dev/null | cut -d ' ' -f 1)
  if [ -z "$want" ]; then
    print_error "unable to find checksum for '${filename}' in '${checksums}'"
    exit 1
  fi
  got=$(hash_sha256 "$filepath")
  if [ "$want" != "$got" ]; then
    print_error "checksum for '$filepath' did not verify ${want} vs $got"
    exit 1
  fi
}

install_one() {
  sudo_cmd=""

  os=$(uname_os)
  arch=$(uname_arch)

  if ! check_platform "$os" "$arch"; then
    print_error "System not supported: $os/$arch"
    print_error "Please compile manually from https://github.com/exelban/one"
    exit 1
  fi

  binary=one
  github_download="https://github.com/exelban/one/releases/download"
  tmpdir="$(mktemp -d)"
  bindir=${bindir:-/usr/local/bin}

  if hash sudo 2>/dev/null; then
    sudo_cmd="sudo"
    echo "You need a root privilege to install one binary to $bindir"

    if ! is_command "$sudo_cmd"; then
      print_error "command not found: sudo. You need a root privilege to continue the installation."
      exit 1;
    fi
  fi

  $sudo_cmd mkdir -p "$bindir"

  resp=$(http_get "https://api.github.com/repos/exelban/one/releases")
  version=$(echo "$resp" | tr ',' '\n' | grep -m 1 "\"tag_name\": " | cut -f4 -d'"')
  if [ -z "$version" ]; then
    print_error "Error fetching latest version. Please try again."
    exit 1
  fi

  version="${version#v}"

  checksum=${binary}_${version}_checksums.txt
  filename=${binary}_${version}_${os}_${arch}
  tarball="${filename}.tar.gz"
  binary_url="${github_download}/v${version}/${tarball}"
  checksum_url="${github_download}/v${version}/${checksum}"

  echo "Latest release version is v$version."

  echo "Downloading $binary_url."
  http_download "$tmpdir/$tarball" "$binary_url"

  echo "Downloading the checksum file for v$version"
  http_download "$tmpdir/$checksum" "$checksum_url"

  echo "Comparing checksums for binaries."
  verify_checksum "$tmpdir/$tarball" "$tmpdir/$checksum"

  echo "Inflating the binary."
  (cd "${tmpdir}" && tar -xzf "${tarball}")

  $sudo_cmd install -d "${bindir}"
  $sudo_cmd install "${tmpdir}/${binary}" "${bindir}/"

  printf "\033[32;1m%s\033[0m\n" "one v${version} was successfully installed in $bindir."
}

exit_error() {
  # shellcheck disable=SC2181
  if [ "$?" -ne 0 ]; then
    print_error "A problem occurred while installing one. Please report it on https://github.com/exelban/one/issues so that we can help you."
  fi
}

trap exit_error EXIT
install_one