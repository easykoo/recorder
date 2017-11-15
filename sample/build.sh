#!/bin/sh

#eg: ./build.sh
#	./build.sh darwin 64
#	./build.sh darwin 64 wechat1


version=1.0

CGO_ENABLED=0

if [ "$1" == "windows" ]; then
    GOOS=windows
elif [ "$1" == "linux" ]; then
    GOOS=linux
else
	GOOS=darwin
fi

if [ "$2" == "64" ]; then
	GOARCH=amd64
elif [ "$2" == "amd64" ]; then
	GOARCH=amd64
else
	GOARCH=386
fi

if ["$3" == ""]; then
    TARGET=""
else
    TARGET="-o $3"
fi


if [ "$GOOS" == "windows" ]; then
    HIDE="-H windowsgui"
else
    HIDE=""
fi
    HIDE=""

echo building $3 for $GOOS $GOARCH

GOOS=$GOOS GOARCH=$GOARCH go build -ldflags "-X main.BUILD_VERSION=$version -X 'main.BUILD_TIME=`date +%Y-%m-%d_%H:%M:%S`' -X 'main.GO_VERSION=`go version`' -s -w $HIDE" -i . $TARGET && echo build successfully
