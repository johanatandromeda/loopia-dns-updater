#!/bin/bash

GIT_TAG=$(git describe --long --tags | sed -e 's/-0-g.*//')

echo Buildning V$GIT_TAG

FLAGS=""
OUTPUT="loopia-ipv6-updater"
env GOOS=linux GOARCH=amd64 go build -o $OUTPUT $FLAGS -ldflags "-X main.version=$GIT_TAG" cmd/main/*.go
