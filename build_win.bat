@ECHO OFF

set PATH=C:\Program Files\mingw-w64\x86_64-8.1.0-posix-seh-rt_v6-rev0\mingw64\bin;%GOROOT%\bin;%GOPATH%\bin;%PATH%

set GOOS=windows
set GOARCH=amd64
set CGO_ENABLED=1

go clean
go get
go fmt
go build --ldflags "-s -w -extldflags '-static' -H windowsgui" -o UpdateChecker.exe
