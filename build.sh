#!/bin/bash

GIT_TAG=$(git describe --long --tags | sed -e 's/-0-g.*//')

echo Buildning V$GIT_TAG

FLAGS=""
OUTPUT="loopia-dns-updater"
env GOOS=linux GOARCH=amd64 go build -o linux/$OUTPUT $FLAGS -ldflags "-X main.version=$GIT_TAG" cmd/main/*.go
env GOOS=openbsd GOARCH=amd64 go build -o openbsd/$OUTPUT $FLAGS -ldflags "-X main.version=$GIT_TAG" cmd/main/*.go
env GOOS=freebsd GOARCH=amd64 go build -o freebsd/$OUTPUT $FLAGS -ldflags "-X main.version=$GIT_TAG" cmd/main/*.go
