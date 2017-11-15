@echo off

echo version = v1.0

set GOOS=windows
set GOARCH=386
set CGO_ENABLED=1

echo building for %GOOS% %GOARCH%

go build  -i . && echo build successfully

:end
