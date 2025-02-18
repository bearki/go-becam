#!/bin/bash

export PKG_CONFIG_PATH=$HOME/Desktop/Project/go-becam/lib/pkgconfig/linux_x86_64_v4l2_gnu

export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=1

go clean -cache
go build -v -ldflags="-s -w" cmd/main.go
