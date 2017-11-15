@echo off

echo version = v1.0

if "%1" == "linux" (
    set GOOS=linux
) else if "%1" == "darwin" (
    set GOOS=darwin
) else (
    set GOOS=windows
)

if "%2" == "64" (
    set GOARCH=amd64
) else if "%2" == "amd64" (
    set GOARCH=amd64
) else (
    set GOARCH=386
)

set CGO_ENABLED=1

echo building %3 for %GOOS% %GOARCH%

if "%3" == "" (
    go build  -i . && echo build successfully
) else (
    go build  -i -o %3 . && echo build successfully
)

:end
