#!/bin/bash
###
 # @Author: chengjiang
 # @Date: 2025-04-24 09:34:32
 # @Description: 运行脚本
###

# 设置脚本运行出错时立即退出
set -e

# 检查环境参数
if [ -z "$1" ]; then
  echo "用法: $0 <环境> [dev|prod]"
  echo "示例: $0 dev    # 部署开发环境"
  echo "示例: $0 prod   # 部署生产环境"
  exit 1
fi

ENV=$1

# 验证环境参数
if [ "$ENV" != "dev" ] && [ "$ENV" != "prod" ]; then
  echo "错误: 环境参数只能是 'dev' 或 'prod'"
  exit 1
fi

# 打印分隔线
divider="===================="

echo "$divider 开始部署 $ENV 环境 $divider"

# cd /root/workspace/git/inoutweb

# 第一部分：拉取最新代码
echo "$divider 拉取最新代码 $divider"
git pull
if [ $? -ne 0 ]; then
  echo "git pull 失败，请检查远程仓库连接或权限！"
  exit 1
fi
echo "git pull 成功！"

# 第二部分：Pnpm build / docker image
echo "$divider 构建项目 $divider"
sh ./build-local-docker-image.sh $ENV
if [ $? -ne 0 ]; then
  echo "Pnpm build / docker image 构建失败，请检查项目配置或错误日志！"
  exit 1
fi
echo "Pnpm build / docker image 构建成功！"

# 第三部分：启动 Docker 容器
echo "$divider 启动 Docker 容器 ($ENV 环境) $divider"

# 根据环境选择不同的docker-compose文件
if [ "$ENV" = "dev" ]; then
  COMPOSE_FILE="docker-compose.yml"
  echo "使用开发环境配置文件: $COMPOSE_FILE (端口映射: 5666:5666)"
else
  COMPOSE_FILE="docker-compose-prod.yml"
  echo "使用生产环境配置文件: $COMPOSE_FILE (端口映射: 80:5666)"
fi

docker-compose -f $COMPOSE_FILE up --build -d
if [ $? -ne 0 ]; then
  echo "Docker Compose 启动失败，请检查 $COMPOSE_FILE 配置或错误日志！"
  exit 1
fi
echo "Docker 容器启动成功！"

# 脚本完成
echo "$divider $ENV 环境部署完成！ $divider"
