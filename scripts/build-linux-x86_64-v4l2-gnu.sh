#!/bin/bash

# 声明有异常时立即终止
set -e


# 项目目录
projectDir=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)

# 设置C工具
Toolchain="$HOME/build-tools/x86-64--glibc--stable-2021.11-5"
export PKG_CONFIG_PATH="${projectDir}/lib/pkgconfig/linux_x86_64_v4l2_gnu"
export AR="${Toolchain}/bin/x86_64-buildroot-linux-gnu-ar"
export CC="${Toolchain}/bin/x86_64-buildroot-linux-gnu-gcc"
export CXX="${Toolchain}/bin/x86_64-buildroot-linux-gnu-g++"

# 配置Go工具
export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=1

# 检查pkg-config
pkg-config --cflags --libs becam

# 执行编译
go clean -cache
go build -v -ldflags="-s -w" "${projectDir}/cmd/main.go"
