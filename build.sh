#!/bin/bash

GIT_TAG=$(git describe --long --tags | sed -e 's/-0-g.*//')

echo Buildning V$GIT_TAG

FLAGS=""
OUTPUT="loopia-dns-updater"

build() {
  env GOOS=$1 GOARCH=$2 go build -o ./$1/$OUTPUT-$1-$2 $FLAGS -ldflags "-X main.version=$GIT_TAG" cmd/main/*.go
}

build linux amd64
build linux arm
build linux arm64
build openbsd amd64
build freebsd amd64
