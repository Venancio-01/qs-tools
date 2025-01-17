#!/bin/bash

# 设置版本号（从 version.go 中获取）
VERSION=$(grep -oP 'var Version = "\K[^"]+' internal/cmd/version.go)

# 设置构建目录
BUILD_DIR="build"
BINARY_NAME="qs-tools"

# 设置远程主机信息
REMOTE_HOST="192.168.70.150"
REMOTE_USER="liqingshan"
REMOTE_PATH="/home/liqingshan"
REMOTE_PASS="lqs988910"

# 设置第二个远程主机信息
REMOTE_HOST2="107.173.165.209"
REMOTE_USER2="root"
REMOTE_PATH2="/root/upload"
REMOTE_PASS2="#6aL*k5d&2Lg*V"

# 清理构建目录
rm -rf $BUILD_DIR
mkdir -p $BUILD_DIR

# 构建函数
build() {
  os=$1
  output="${BUILD_DIR}/${BINARY_NAME}"

  # Windows 需要加上 .exe 后缀
  if [ "$os" = "windows" ]; then
    output="${output}.exe"
  fi

  echo "正在构建 $os/amd64..."

  # 设置编译参数
  # CGO_ENABLED=0: 禁用 CGO，使用纯 Go 实现
  # -trimpath: 删除编译路径信息
  # -ldflags:
  #   -s: 删除符号表
  #   -w: 删除 DWARF 调试信息
  #   -extldflags "-static": 静态链接
  CGO_ENABLED=0 GOOS=$os GOARCH=amd64 \
    go build -trimpath \
    -ldflags="-s -w -extldflags '-static'" \
    -o "$output" \
    ./cmd/qs-tools

  if [ $? -eq 0 ]; then
    echo "✅ 构建成功: $output"

    # 显示文件信息
    if [ "$os" = "linux" ]; then
      file "$output"
    fi

    # 计算文件大小
    if [ -f "$output" ]; then
      size=$(ls -lh "$output" | awk '{print $5}')
      echo "   文件大小: $size"
    fi
  else
    echo "❌ 构建失败: $os/amd64"
    exit 1
  fi
}

echo "开始构建 qs-tools $VERSION..."

# 构建 Linux 64位版本
build "linux"

# 构建 Windows 64位版本
build "windows"

echo "构建完成！"
echo "构建文件位于 $BUILD_DIR 目录"

# 显示构建文件列表
echo -e "\n构建文件列表："
ls -lh $BUILD_DIR

# # 传输 Linux 版本到远程主机
# if [ -f "$BUILD_DIR/qs-tools" ]; then
#   echo -e "\n开始传输文件到远程主机..."

#   # 检查是否可以连接到远程主机
#   if ping -c 1 $REMOTE_HOST &>/dev/null; then
#     echo "正在传输到 $REMOTE_USER@$REMOTE_HOST:$REMOTE_PATH ..."

#     # 检查 sshpass 是否安装
#     if ! command -v sshpass &>/dev/null; then
#       echo "正在安装 sshpass..."
#       sudo apt-get update && sudo apt-get install -y sshpass
#     fi

#     # 使用 sshpass 传输文件
#     if sshpass -p "$REMOTE_PASS" scp -o StrictHostKeyChecking=no "$BUILD_DIR/qs-tools" "$REMOTE_USER@$REMOTE_HOST:$REMOTE_PATH/"; then
#       echo "✅ 文件传输成功！"

#       # 设置远程文件权限
#       sshpass -p "$REMOTE_PASS" ssh -o StrictHostKeyChecking=no "$REMOTE_USER@$REMOTE_HOST" "chmod +x $REMOTE_PATH/qs-tools"
#       echo "✅ 已设置可执行权限"
#     else
#       echo "❌ 文件传输失败"
#     fi
#   else
#     echo "❌ 无法连接到远程主机 $REMOTE_HOST"
#   fi
# fi

# 传输 Linux 版本到第二个远程主机
if [ -f "$BUILD_DIR/qs-tools" ] || [ -f "$BUILD_DIR/qs-tools.exe" ]; then
  echo -e "\n开始传输文件到第二个远程主机..."

  # 检查是否可以连接到远程主机
  if ping -c 1 $REMOTE_HOST2 &>/dev/null; then
    echo "正在传输到 $REMOTE_USER2@$REMOTE_HOST2:$REMOTE_PATH2 ..."

    # 检查 sshpass 是否安装
    if ! command -v sshpass &>/dev/null; then
      echo "正在安装 sshpass..."
      sudo apt-get update && sudo apt-get install -y sshpass
    fi

    # 传输 Linux 版本
    if [ -f "$BUILD_DIR/qs-tools" ]; then
      if sshpass -p "$REMOTE_PASS2" scp -o StrictHostKeyChecking=no "$BUILD_DIR/qs-tools" "$REMOTE_USER2@$REMOTE_HOST2:$REMOTE_PATH2/"; then
        echo "✅ Linux 版本传输成功！"
      else
        echo "❌ Linux 版本传输失败"
      fi
    fi

    # 传输 Windows 版本
    if [ -f "$BUILD_DIR/qs-tools.exe" ]; then
      if sshpass -p "$REMOTE_PASS2" scp -o StrictHostKeyChecking=no "$BUILD_DIR/qs-tools.exe" "$REMOTE_USER2@$REMOTE_HOST2:$REMOTE_PATH2/"; then
        echo "✅ Windows 版本传输成功！"
      else
        echo "❌ Windows 版本传输失败"
      fi
    fi

  else
    echo "❌ 无法连接到远程主机 $REMOTE_HOST2"
  fi
fi
