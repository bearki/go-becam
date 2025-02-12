#!/bin/bash

export PKG_CONFIG_PATH=$HOME/Desktop/Project/go-becam/lib/pkgconfig/linux_arm
export CC=$HOME/build-tools/RV1126_RV1109_LINUX_SDK_V2.2.4/prebuilts/gcc/linux-x86/arm/gcc-arm-8.3-2019.03-x86_64-arm-linux-gnueabihf/bin/arm-linux-gnueabihf-gcc
export CXX=$HOME/build-tools/RV1126_RV1109_LINUX_SDK_V2.2.4/prebuilts/gcc/linux-x86/arm/gcc-arm-8.3-2019.03-x86_64-arm-linux-gnueabihf/bin/arm-linux-gnueabihf-g++

export GOOS=linux
export GOARCH=arm
export CGO_ENABLED=1

go clean -cache
go build -v cmd/main.go
