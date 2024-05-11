#!/bin/sh

set -eu

version="0.0.0"

build_darwin() {
  GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.version=$version" -o bin/one && tar -czf release/one_"$version"_darwin_x86_64.tar.gz -C bin one && rm bin/one
  GOOS=darwin GOARCH=arm64 go build -ldflags "-X main.version=$version" -o bin/one && tar -czf release/one_"$version"_darwin_arm64.tar.gz -C bin one && rm bin/one

  cd release
  {
    shasum -a 256 one_"$version"_darwin_x86_64.tar.gz
    shasum -a 256 one_"$version"_darwin_arm64.tar.gz
  } >> one_"$version"_checksums.txt
	cd ..
}
build_linux() {
  GOOS=linux GOARCH=amd64 go build -ldflags "-X main.version=$version" -o bin/one && tar -czf release/one_"$version"_linux_x86_64.tar.gz -C bin one && rm bin/one
  GOOS=linux GOARCH=386 go build -ldflags "-X main.version=$version" -o bin/one && tar -czf release/one_"$version"_linux_x86.tar.gz -C bin one && rm bin/one
  GOOS=linux GOARCH=arm64 go build -ldflags "-X main.version=$version" -o bin/one && tar -czf release/one_"$version"_linux_arm64.tar.gz -C bin one && rm bin/one

  cd release
  {
    shasum -a 256 one_"$version"_linux_x86_64.tar.gz
    shasum -a 256 one_"$version"_linux_x86.tar.gz
    shasum -a 256 one_"$version"_linux_arm64.tar.gz
  } >> one_"$version"_checksums.txt
	cd ..
}
build_windows() {
  GOOS=windows GOARCH=amd64 go build -ldflags "-X main.version=$version" -o bin/one && tar -czf release/one_"$version"_windows_x86_64.tar.gz -C bin one && rm bin/one
  GOOS=windows GOARCH=386 go build -ldflags "-X main.version=$version" -o bin/one && tar -czf release/one_"$version"_windows_x86.tar.gz -C bin one && rm bin/one
  GOOS=windows GOARCH=arm64 go build -ldflags "-X main.version=$version" -o bin/one && tar -czf release/one_"$version"_windows_arm64.tar.gz -C bin one && rm bin/one

  cd release
  {
    shasum -a 256 one_"$version"_windows_x86_64.tar.gz
    shasum -a 256 one_"$version"_windows_x86.tar.gz
    shasum -a 256 one_"$version"_windows_arm64.tar.gz
  } >> one_"$version"_checksums.txt
  cd ..
}

printf "Building v%s...\n" "$version"

rm -rf "bin" && rm -rf "release"
mkdir -p "release"
touch release/one_"$version"_checksums.txt

echo "Building darwin..."
build_darwin
echo "Building linux..."
build_linux
echo "Building windows..."
build_windows
rm -rf "bin"

printf "\033[32;1m%s\033[0m\n" "one ${version} was successfully build."
open release