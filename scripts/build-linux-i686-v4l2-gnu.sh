#!/bin/bash

export PKG_CONFIG_PATH=$HOME/Desktop/Project/go-becam/lib/pkgconfig/linux_i686_v4l2_gnu

export GOOS=linux
export GOARCH=386
export CGO_ENABLED=1

go clean -cache
go build -v -ldflags="-s -w" cmd/main.go
