#!/bin/sh

version=1.0

CGO_ENABLED=1
GOOS=windows
GOARCH=386

echo building for $GOOS $GOARCH

GOOS=$GOOS GOARCH=$GOARCH go build -i . && echo build successfully
