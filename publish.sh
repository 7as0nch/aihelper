#!/bin/bash

# --- 配置区域 ---
DOCKERHUB_USER="7as0nch"
VERSION="v1.0.2"  # 统一版本号，后续手动在此修改
PLATFORM="linux/amd64" # 统一指定构建平台，例如 linux/amd64, linux/arm64

# --- 脚本逻辑 ---
SERVICE=$1  # 接收第一个参数作为指定服务名

build_and_push() {
    local name=$1
    local dir=$2
    local dockerfile=$3
    local context=$4

    echo "=== 正在构建服务: $name (版本: $VERSION, 平台: $PLATFORM) ==="
    
    # 使用 docker buildx 构建以支持跨平台
    docker buildx build --platform $PLATFORM \
        -t $DOCKERHUB_USER/aichat-$name:$VERSION \
        -t $DOCKERHUB_USER/aichat-$name:latest \
        -f "$dir/$dockerfile" "$dir/$context" --push

    if [ $? -eq 0 ]; then
        echo "Successfully pushed $name"
    else
        echo "Failed to build/push $name"
        exit 1
    fi
}

# 确保 buildx 已初始化
docker buildx create --use --name mybuilder 2>/dev/null || docker buildx use mybuilder

if [ -z "$SERVICE" ] || [ "$SERVICE" == "backend" ]; then
    build_and_push "backend" "backend" "Dockerfile" "."
fi

if [ -z "$SERVICE" ] || [ "$SERVICE" == "backweb" ]; then
    build_and_push "backweb" "backweb" "scripts/deploy/Dockerfile" "."
fi

if [ -z "$SERVICE" ] || [ "$SERVICE" == "litechat" ]; then
    build_and_push "litechat" "litechat" "dockerfile" "."
fi

echo "Done!"
