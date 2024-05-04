#!/bin/bash

ARCH=$(uname -m)
case $ARCH in
    i386|i686) ARCH=x86 ;;
    aarch64*) ARCH=arm64 ;;
esac

LATEST_VERSION=$(curl -L -s -H 'Accept: application/json' https://github.com/exelban/one/releases/latest | sed -e 's/.*"tag_name":"\([^"]*\)".*/\1/')
FILE="one_${LATEST_VERSION//v/}_$(uname -s)_${ARCH}.tar.gz"

curl -L -o one.tar.gz "https://github.com/exelban/one/releases/download/${LATEST_VERSION}/${FILE}"
tar xzvf one.tar.gz one
install -m 755 one -t "/usr/local/bin"
rm one.tar.gz