#!/bin/bash

# 获取项目路径
projectDir=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)/..

# 声明要使用的库列表
libList=(
    "libbecam_linux_arm_v4l2_rv1126"
    "libbecam_linux_i686_v4l2_gnu"
    "libbecam_linux_x86_64_v4l2_gnu"
)

# 遍历库列表
for item in "${libList[@]}"; do
    # 声明压缩包路径
    srcPath="${projectDir}/download/${item}.tar.gz"
    # 声明解压路径
    dstPath="${projectDir}/lib/libbecam/${item}"
    # 创建解压目录（强制创建）
    mkdir -p "$dstPath"
    # 解压
    tar -xzf "$srcPath" -C "$dstPath"
    # 确保pkg-config对应目录存在
    libPcDir="${projectDir}/lib/pkgconfig/${item//libbecam_/}"
    mkdir -p $libPcDir
    # 处理pkg-config文件
    sed -e "s@ENV_LIBRARY_PATH@$dstPath@g" "${dstPath}/becam.pc" > "${libPcDir}/becam.pc"
done