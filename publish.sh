#!/bin/bash

# --- 配置区域 ---
DOCKERHUB_USER="7as0nch"
PLATFORM="linux/amd64" # 统一指定构建平台

# 各个服务的版本号独立管理
VERSION_BACKEND="v1.0.4-beta.1"
VERSION_BACKWEB="v1.0.4-beta.1"
VERSION_LITECHAT="v1.0.4-beta.1"

# --- 脚本逻辑 ---
SERVICE=$1  # 接收第一个参数作为指定服务名

build_and_push() {
    local name=$1
    local dir=$2
    local dockerfile=$3
    local context=$4
    local version=$5

    echo "=== 正在构建服务: $name (版本: $version, 平台: $PLATFORM) ==="
    
    # 使用 docker buildx 构建以支持跨平台
    docker buildx build --platform $PLATFORM \
        -t $DOCKERHUB_USER/aichat-$name:$version \
        -t $DOCKERHUB_USER/aichat-$name:latest \
        -f "$dir/$dockerfile" "$dir/$context" --push

    if [ $? -eq 0 ]; then
        echo "Successfully pushed $name"
        # 自动同步更新 k8s-deployment.yaml 中的版本号
        if [[ "$OSTYPE" == "darwin"* ]]; then
            sed -i '' "s|7as0nch/aichat-$name:v[0-9.]*|7as0nch/aichat-$name:$version|g" k8s-deployment.yaml
        else
            sed -i "s|7as0nch/aichat-$name:v[0-9.]*|7as0nch/aichat-$name:$version|g" k8s-deployment.yaml
        fi
        echo "K8s deployment sync done for $name"
    else
        echo "Failed to build/push $name"
        exit 1
    fi
}

# 确保 buildx 已初始化
docker buildx create --use --name mybuilder 2>/dev/null || docker buildx use mybuilder

if [ -z "$SERVICE" ] || [ "$SERVICE" == "backend" ]; then
    build_and_push "backend" "backend" "Dockerfile" "." "$VERSION_BACKEND"
fi

if [ -z "$SERVICE" ] || [ "$SERVICE" == "backweb" ]; then
    build_and_push "backweb" "backweb" "scripts/deploy/Dockerfile" "." "$VERSION_BACKWEB"
fi

if [ -z "$SERVICE" ] || [ "$SERVICE" == "litechat" ]; then
    build_and_push "litechat" "litechat" "dockerfile" "." "$VERSION_LITECHAT"
fi

echo "Done! 所有选定服务已发布，K8s 配置已更新。"
