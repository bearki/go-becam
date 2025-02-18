# 获取项目路径
$projectDir = Resolve-Path "${PSScriptRoot}\.."

# 声明要使用的库列表
$libList = @(
    "libbecam_windows_i686_dshow_mingw",
    "libbecam_windows_x86_64_dshow_mingw",
    "libbecam_windows_i686_mf_mingw",
    "libbecam_windows_x86_64_mf_mingw"
)

# 遍历库列表
foreach ($item in $libList) {
    # 声明压缩包路径
    $srcPath = "${projectDir}\download\${item}.zip"
    # 声明解压路径
    $dstPath = "${projectDir}\lib\libbecam\${item}"
    # 创建解压目录
    New-Item -Path $dstPath -ItemType Directory -Force
    # 解压
    Expand-Archive -Path $srcPath -DestinationPath $dstPath -Force
    # 处理pkg-config文件
    $libPc = (Get-Content -Path "${dstPath}\becam.pc") -creplace "ENV_LIBRARY_PATH", "${dstPath}"
    # 确保目标目录存在
    $libPcShortPath = $item -replace '^libbecam_', ""
    $libPcDstDir = Join-Path -Path $projectDir -ChildPath "lib\pkgconfig\$libPcShortPath"
    New-Item -Path $targetDir -ItemType Directory -foreach
    # 写入
    $libPc | Set-Content -Path "${libPcDstDir}\becam.pc" -Encoding UTF8
}
