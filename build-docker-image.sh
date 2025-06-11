#!/bin/bash

# Function to handle tagging and pushing an image
perform_push_operations() {
    local image_to_push="$1"
    local current_container_engine="$2"

    DEFAULT_REGISTRY_HOST="greptime-registry.cn-hangzhou.cr.aliyuncs.com/tools"
    read -p "请输入镜像仓库地址 (默认为 $DEFAULT_REGISTRY_HOST): " REGISTRY_HOST
    if [ -z "$REGISTRY_HOST" ]; then
        REGISTRY_HOST="$DEFAULT_REGISTRY_HOST"
        echo "使用默认镜像仓库地址: $REGISTRY_HOST"
    fi

    REMOTE_IMAGE_NAME="$REGISTRY_HOST/$image_to_push"

    echo "正在标记镜像 $image_to_push 为 $REMOTE_IMAGE_NAME ..."
    if ! "$current_container_engine" tag "$image_to_push" "$REMOTE_IMAGE_NAME"; then
        echo "错误：镜像标记失败。"
        return 1
    fi
    echo "镜像标记成功。"

    echo "正在推送镜像 $REMOTE_IMAGE_NAME ..."
    if ! "$current_container_engine" push "$REMOTE_IMAGE_NAME"; then
        echo "错误：镜像推送失败。"
        echo "请确保您已登录到 $REGISTRY_HOST ($current_container_engine login $REGISTRY_HOST)"
        return 1
    fi
    echo "镜像 $REMOTE_IMAGE_NAME 推送成功。"
    return 0
}

# Determine script mode
SCRIPT_MODE="default"
if [[ "$1" == "only-push" ]]; then
    SCRIPT_MODE="only-push"
fi

# 检查 podman 是否可用
if command -v podman &> /dev/null; then
    CONTAINER_ENGINE="podman"
# 否则，检查 docker 是否可用
elif command -v docker &> /dev/null; then
    CONTAINER_ENGINE="docker"
else
    echo "错误：未找到 podman 或 docker。请安装其中一个。"
    exit 1
fi
echo "将使用 $CONTAINER_ENGINE。"

if [[ "$SCRIPT_MODE" == "only-push" ]]; then
    echo "--- 仅推送模式 ---"
    echo "可用的本地镜像列表:"
    "$CONTAINER_ENGINE" images

    read -p "请输入要推送的本地镜像名称 (例如, repository:tag): " LOCAL_IMAGE_TO_PUSH
    if [ -z "$LOCAL_IMAGE_TO_PUSH" ]; then
        echo "错误：镜像名称不能为空。"
        exit 1
    fi

    # 检查镜像是否存在
    if ! "$CONTAINER_ENGINE" image inspect "$LOCAL_IMAGE_TO_PUSH" &> /dev/null; then
        echo "错误：本地镜像 $LOCAL_IMAGE_TO_PUSH 未找到。"
        exit 1
    fi

    echo "准备推送镜像: $LOCAL_IMAGE_TO_PUSH"
    if perform_push_operations "$LOCAL_IMAGE_TO_PUSH" "$CONTAINER_ENGINE"; then
        echo "镜像推送流程完成。"
    else
        echo "镜像推送流程失败。"
        exit 1
    fi
else # Default mode (build and then optionally push)
    echo "--- 默认构建并推送模式 ---"

    # 以交互方式询问用户镜像标签
    read -p "请输入镜像标签 (例如, latest, 1.0): " IMAGE_TAG

    # 检查标签是否为空
    if [ -z "$IMAGE_TAG" ]; then
        echo "错误：镜像标签不能为空。"
        exit 1
    fi

    BUILT_IMAGE_BASENAME="ingester" # As implied by original script
    LOCAL_BUILT_IMAGE_NAME="${BUILT_IMAGE_BASENAME}:${IMAGE_TAG}"
    DOCKERFILE_PATH="Dockerfile" # 假设 Dockerfile 在当前目录

    # 检查 Dockerfile 是否存在
    if [ ! -f "$DOCKERFILE_PATH" ]; then
        echo "错误：在当前目录下未找到 $DOCKERFILE_PATH。"
        echo "请确保 Dockerfile 存在于脚本执行的目录中。"
        exit 1
    fi

    echo "正在构建镜像 $LOCAL_BUILT_IMAGE_NAME ..."

    # 构建镜像
    if "$CONTAINER_ENGINE" build -t "$LOCAL_BUILT_IMAGE_NAME" -f "$DOCKERFILE_PATH" .; then
        echo "镜像 $LOCAL_BUILT_IMAGE_NAME 构建成功。"
    else
        echo "错误：镜像构建失败。"
        exit 1
    fi

    # 询问是否推送到镜像仓库
    read -p "是否要将镜像 $LOCAL_BUILT_IMAGE_NAME 推送到镜像仓库? (y/n): " PUSH_CONFIRMATION

    if [[ "$PUSH_CONFIRMATION" == "y" || "$PUSH_CONFIRMATION" == "Y" ]]; then
        if perform_push_operations "$LOCAL_BUILT_IMAGE_NAME" "$CONTAINER_ENGINE"; then
            echo "镜像推送流程完成。"
        else
            echo "镜像推送流程失败。"
            exit 1
        fi
    else
        echo "镜像 $LOCAL_BUILT_IMAGE_NAME 未推送到仓库。"
    fi
fi

echo "脚本执行完毕。"
exit 0

# 检查 podman 是否可用
if command -v podman &> /dev/null; then
    CONTAINER_ENGINE="podman"
# 否则，检查 docker 是否可用
elif command -v docker &> /dev/null; then
    CONTAINER_ENGINE="docker"
else
    echo "错误：未找到 podman 或 docker。请安装其中一个。"
    exit 1
fi

echo "将使用 $CONTAINER_ENGINE 来构建镜像。"

# 以交互方式询问用户镜像标签
read -p "请输入镜像标签 (例如, latest, 1.0): " IMAGE_TAG

# 检查标签是否为空
if [ -z "$IMAGE_TAG" ]; then
    echo "错误：镜像标签不能为空。"
    exit 1
fi

LOCAL_IMAGE_NAME="ingester:$IMAGE_TAG"
DOCKERFILE_PATH="Dockerfile" # 假设 Dockerfile 在当前目录

# 检查 Dockerfile 是否存在
if [ ! -f "$DOCKERFILE_PATH" ]; then
    echo "错误：在当前目录下未找到 $DOCKERFILE_PATH。"
    echo "请确保 Dockerfile 存在于脚本执行的目录中。"
    exit 1
fi

echo "正在构建镜像 $LOCAL_IMAGE_NAME ..."

# 构建镜像
if "$CONTAINER_ENGINE" build -t "$LOCAL_IMAGE_NAME" -f "$DOCKERFILE_PATH" .; then
    echo "镜像 $LOCAL_IMAGE_NAME 构建成功。"
else
    echo "错误：镜像构建失败。"
    exit 1
fi

# 询问是否推送到镜像仓库
read -p "是否要将镜像推送到镜像仓库? (y/n): " PUSH_TO_REGISTRY

if [[ "$PUSH_TO_REGISTRY" == "y" || "$PUSH_TO_REGISTRY" == "Y" ]]; then
    DEFAULT_REGISTRY_HOST="greptime-registry.cn-hangzhou.cr.aliyuncs.com/tools"
    read -p "请输入镜像仓库地址 (默认为 $DEFAULT_REGISTRY_HOST): " REGISTRY_HOST
    if [ -z "$REGISTRY_HOST" ]; then
        REGISTRY_HOST="$DEFAULT_REGISTRY_HOST"
        echo "使用默认镜像仓库地址: $REGISTRY_HOST"
    fi

    REMOTE_IMAGE_NAME="$REGISTRY_HOST/$LOCAL_IMAGE_NAME"

    echo "正在标记镜像 $LOCAL_IMAGE_NAME 为 $REMOTE_IMAGE_NAME ..."
    if "$CONTAINER_ENGINE" tag "$LOCAL_IMAGE_NAME" "$REMOTE_IMAGE_NAME"; then
        echo "镜像标记成功。"
    else
        echo "错误：镜像标记失败。"
        exit 1
    fi

    echo "正在推送镜像 $REMOTE_IMAGE_NAME ..."
    if "$CONTAINER_ENGINE" push "$REMOTE_IMAGE_NAME"; then
        echo "镜像 $REMOTE_IMAGE_NAME 推送成功。"
    else
        echo "错误：镜像推送失败。"
        echo "请确保您已登录到 $REGISTRY_HOST ($CONTAINER_ENGINE login $REGISTRY_HOST)"
        exit 1
    fi
else
    echo "镜像未推送到仓库。"
fi