#!/bin/bash
###
 # @Author: chengjiang
 # @Date: 2025-10-13 12:31:09
 # @Description: 
### 

set -e

# 检查环境参数
if [ -z "$1" ]; then
  echo "用法: $0 <环境> [dev|prod]"
  echo "示例: $0 dev    # 部署开发环境"
  echo "示例: $0 prod   # 部署生产环境"
  exit 1
fi

ENV=$1

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
LOG_FILE=${SCRIPT_DIR}/build-local-docker-image.log
ERROR=""
IMAGE_NAME="litechat_${ENV}"

function stop_and_remove_container() {
    # judge if the container is running
    IS_RUNNING=$(docker ps -q -f name=${IMAGE_NAME}) || ERROR="stop_and_remove_container failed"
    if [ "$IS_RUNNING" ]; then
      # Stop and remove the existing container
      echo "Info: Stopping running container ${IMAGE_NAME}"
      docker stop ${IMAGE_NAME} >/dev/null 2>&1 || ERROR="stop_and_remove_container failed"
      docker rm ${IMAGE_NAME} >/dev/null 2>&1 || ERROR="stop_and_remove_container failed"
    fi
    echo "Info: No running container named ${IMAGE_NAME} found"
}

function remove_image() {
    # Remove the existing image
    docker rmi ${IMAGE_NAME} >/dev/null 2>&1 || echo "Info: No existing image named ${IMAGE_NAME} found"
}

function install_dependencies() {
    # Install all dependencies
    cd ${SCRIPT_DIR}
    pnpm install || ERROR="install_dependencies failed"
}

function build_image() {
    # build docker
    echo "pnpm run build"
    pnpm run build || ERROR="build_image failed by pnpm run build"
    docker build --platform linux/amd64 -f dockerfile_${ENV} . -t ${IMAGE_NAME} || ERROR="build_image failed"
    docker save ${IMAGE_NAME} -o ${IMAGE_NAME}.tar
    echo "Info: Saving docker image to ${IMAGE_NAME}.tar.gz"
    tar czf ${IMAGE_NAME}.tar.gz ${IMAGE_NAME}.tar
}

function log_message() {
    if [[ ${ERROR} != "" ]];
    then
        >&2 echo "build failed, Please check build-local-docker-image.log for more details"
        >&2 echo "ERROR: ${ERROR}"
        exit 1
    else
        echo "docker image with tag '${IMAGE_NAME}' built sussessfully. Use below sample command to run the container"
        echo ""
        echo "docker run -d -p 5666:5666 --name ${IMAGE_NAME} ${IMAGE_NAME}"
    fi
}

echo "Info: Stopping and removing existing container and image" | tee ${LOG_FILE}
stop_and_remove_container
remove_image

echo "Info: Installing dependencies" | tee -a ${LOG_FILE}
install_dependencies 1>> ${LOG_FILE} 2>> ${LOG_FILE}

if [[ ${ERROR} == "" ]]; then
    echo "Info: Building docker image" | tee -a ${LOG_FILE}
    build_image 1>> ${LOG_FILE} 2>> ${LOG_FILE}
fi

log_message | tee -a ${LOG_FILE}
