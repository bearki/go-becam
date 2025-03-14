#!/usr/bin/env python3
import argparse
import json
import os
from pathlib import Path
import platform
import shutil
import sys
import tarfile
import urllib.request
import zipfile
import glob


def parse_arguments(projectDir: Path):
    """参数解析"""

    # 声明解析器
    parser = argparse.ArgumentParser(description="依赖库处理器")
    # 声明参数
    parser.add_argument(
        "--libs-root",
        type=str,
        default=os.path.join(projectDir, "libs"),
        metavar="libs",
        help="三方库根目录",
    )
    parser.add_argument(
        "--enable-pkgconfig", type=bool, default=True, help="是否启用pkg-config配置文件"
    )
    parser.add_argument(
        "--pkgconfig-root",
        type=str,
        default=os.path.join(projectDir, "libs", "pkgconfig"),
        metavar="libs/pkgconfig",
        help="pkg-config配置文件存放根目录",
    )
    parser.add_argument(
        "--prioritylib-type",
        type=str,
        default="static",
        choices=("static", "shared"),
        metavar="static",
        help="优先使用的库类型",
    )
    parser.add_argument(
        "--limit-os",
        type=str,
        default=platform.system(),
        choices=("Windows", "Linux", "Android"),
        metavar="Windows",
        help="要是使用的库平台",
    )
    parser.add_argument(
        "--limit-archs",
        type=str,
        nargs="*",
        metavar="i686 x86_64 arm aarch64",
        help="限制要使用的Arch范围, 空格分隔",
    )
    parser.add_argument(
        "--token",
        type=str,
        default=os.environ.get("GITEA_TOKEN"),
        metavar="xxxxxxxxxxxxxxxxxxxxxxxxxx",
        help="限制要使用的Arch范围, 空格分隔",
    )
    # 执行解析
    return parser.parse_args()


def load_dependencies(projectDir: Path):
    """加载依赖配置"""

    deps_path = os.path.join(projectDir, "dependencies.json")
    with open(deps_path, "r", encoding="utf-8") as f:
        return json.load(f)


def create_directories(directories: list):
    """创建必要目录结构"""

    for d in directories:
        os.makedirs(d, exist_ok=True)


def generate_pkg_config(
    libName: str, libPath: str, pkgconfigDir: str, system: str, arch: str
):
    """生成pkg-config文件"""

    # 获取pc文件名
    pcFile = f"{libName.removeprefix("lib")}.pc"
    # 读取原始pc文件内容
    with open(os.path.join(libPath, pcFile), "r", encoding="utf-8") as f:
        pcData = f.read()

    # 赋值库pkg-config中的额库路径
    pcData = pcData.replace("${Replace:Library:Dir}", libPath)
    pcData = pcData.replace("ENV_LIBRARY_PATH", libPath)

    # 确定存储目录
    if "i686" in arch:
        pcDir = os.path.join(pkgconfigDir, f"{system}_i686")
    elif "x86_64" in arch:
        pcDir = os.path.join(pkgconfigDir, f"{system}_x86_64")
    elif "aarch64" in arch or "arm64" in arch or "armv8" in arch:
        pcDir = os.path.join(pkgconfigDir, f"{system}_aarch64")
    elif "arm" in arch or "armv7" in arch:  # 注意该判断一定要在aarch64之后
        pcDir = os.path.join(pkgconfigDir, f"{system}_arm")
    else:
        pcDir = os.path.join(pkgconfigDir, f"{system}_{arch}")

    # 创建目录
    os.makedirs(pcDir, exist_ok=True)
    # 写入新的PC文件到指定位置
    with open(os.path.join(pcDir, pcFile), "w", encoding="utf-8") as f:
        f.write(pcData)


def prioritize_libraries(libPath: str, priorityType: str):
    """处理库文件优先级"""

    # 匹配.a静态库
    aStaticFiles = glob.glob(os.path.join(libPath, "lib", "*.a"))
    # 匹配.lib静态库
    libStaticFiles = glob.glob(os.path.join(libPath, "lib", "*.lib"))
    # 匹配SO动态库
    soDynamicFiles = glob.glob(os.path.join(libPath, "lib", "*.so"))
    # 匹配DLL动态库
    dllDynamicFiles = glob.glob(os.path.join(libPath, "lib", "*.dll"))

    # 两者是否同时存在
    if (aStaticFiles or libStaticFiles) and (soDynamicFiles or dllDynamicFiles):
        # 检查保留优先级
        if priorityType == "static":
            for f in soDynamicFiles:
                os.remove(f)
            for f in dllDynamicFiles:
                os.remove(f)
        else:
            for f in aStaticFiles:
                os.remove(f)
            for f in libStaticFiles:
                os.remove(f)


def process_library(downloadDir: str, args: argparse.Namespace, lib):
    """处理单个库"""

    # 处理传入参数
    libsDir = args.libs_root
    pkgconfigDir = args.pkgconfig_root
    limitOS = args.limit_os
    limitArchs = args.limit_archs
    token = args.token
    priorityLibType = args.prioritylib_type
    enabledPkgconfig = args.enable_pkgconfig

    # 检查目标是否存在
    if (
        ("Targets" not in lib)
        or (limitOS not in lib["Targets"])
        or ("Archs" not in lib["Targets"][limitOS])
    ):
        # 不存在目标，跳过
        return

    # 转换系统为小写
    for arch in lib["Targets"][limitOS]["Archs"]:
        # 是否需要限制下载目标架构
        if limitArchs and arch not in limitArchs:
            continue

        print("--------------------------")

        # 提取库信息
        libName = lib["Name"]
        libVersion = lib["Version"]
        libSystem = limitOS.lower()
        libRepository = lib["Repository"]
        libTarget = f"{libSystem}_{arch}"
        libFullTarget = f"{libName}_{libTarget}"
        libCompressType = lib["Targets"][limitOS]["CompressType"]
        downloadUrl = f"{libRepository}/releases/download/{libVersion}/{libFullTarget}.{libCompressType}"
        downloadPath = os.path.join(downloadDir, f"{libFullTarget}.{libCompressType}")

        # 打印库信息
        print(f"libName:         {libName}")
        print(f"libVersion:      {libVersion}")
        print(f"libSystem:       {libSystem}")
        print(f"libRepository:   {libRepository}")
        print(f"libTarget:       {libTarget}")
        print(f"libFullTarget:   {libFullTarget}")
        print(f"libCompressType: {libCompressType}")
        print(f"downloadUrl:     {downloadUrl}")
        print(f"downloadPath:    {downloadPath}")

        # 是否需要拼接授权
        if token:
            downloadUrl += f"?token={token}"
        # 执行下载
        urllib.request.urlretrieve(downloadUrl, downloadPath)

        # 拼接库目录
        libPath = os.path.join(libsDir, libName, libFullTarget)
        # 创建库目录
        os.makedirs(libPath, exist_ok=True)
        # 解压处理（区分不同压缩类型zip、tar.gz）
        if libCompressType == "zip":
            with zipfile.ZipFile(downloadPath, "r") as zf:
                zf.extractall(libPath)
        elif libCompressType in ["tar.gz", "tgz", "tar.xz", "txz", "tar.bz2", "tbz2"]:
            with tarfile.open(downloadPath, "r:*") as tf:
                tf.extractall(libPath, filter="tar")
        else:
            print(
                f"Decompression failure: Unsupported compression format",
                file=sys.stderr,
            )

        # 库文件优先级处理
        prioritize_libraries(libPath, priorityLibType)
        # 是否启用pkg-config配置文件
        if enabledPkgconfig:
            # 生成pkg-config文件
            generate_pkg_config(libName, libPath, pkgconfigDir, libSystem, arch)


def main():
    # Python版本检查
    if sys.version_info < (3, 7):
        raise RuntimeError("Python 3.7+ required")

    # 获取脚本所在路径
    scriptPath = Path(sys.argv[0]).resolve()
    # 获取项目根路径
    projectDir = scriptPath.parent.parent
    # 获取下载根路径
    downloadDir = os.path.join(projectDir, "download")
    # 参数解析
    args = parse_arguments(projectDir)

    # 依赖处理
    try:
        # 加载依赖项
        dependencies = load_dependencies(projectDir)
    except FileNotFoundError as e:
        # 抛出异常
        print(f"Error loading dependencies.json: {e}", file=sys.stderr)
        sys.exit(1)

    # 创建必要目录
    create_directories([args.libs_root, downloadDir])

    # 下载和解压处理
    for lib in dependencies:
        process_library(downloadDir, args, lib)

    # 清理下载目录
    shutil.rmtree(downloadDir)


if __name__ == "__main__":
    print(
        "------------------------------- 开始处理依赖库 -------------------------------"
    )
    main()
    print(
        "------------------------------- 依赖库处理结束 -------------------------------"
    )
