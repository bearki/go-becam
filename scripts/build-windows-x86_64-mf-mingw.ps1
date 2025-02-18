$Env:PKG_CONFIG_PATH="${PSScriptRoot}\..\lib\pkgconfig\windows_x86_64_mf_mingw"

$Env:GOOS="windows"
$Env:GOARCH="amd64"
$Env:CGO_ENABLED="1"

go clean -cache
go build -v -ldflags="-s -w" cmd/main.go
