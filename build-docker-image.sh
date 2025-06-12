#!/bin/bash

# Function to handle tagging and pushing an image
perform_push_operations() {
    local image_to_push="$1"
    local current_container_engine="$2"

    DEFAULT_REGISTRY_HOST="greptime-registry.cn-hangzhou.cr.aliyuncs.com/tools"
    read -p "Please enter the image registry address (default is $DEFAULT_REGISTRY_HOST): " REGISTRY_HOST
    if [ -z "$REGISTRY_HOST" ]; then
        REGISTRY_HOST="$DEFAULT_REGISTRY_HOST"
        echo "Using default image registry address: $REGISTRY_HOST"
    fi

    REMOTE_IMAGE_NAME="$REGISTRY_HOST/$image_to_push"

    echo "Tagging image $image_to_push as $REMOTE_IMAGE_NAME ..."
    if ! "$current_container_engine" tag "$image_to_push" "$REMOTE_IMAGE_NAME"; then
        echo "Error: Image tagging failed."
        return 1
    fi
    echo "Image tagged successfully."

    echo "Pushing image $REMOTE_IMAGE_NAME ..."
    if ! "$current_container_engine" push "$REMOTE_IMAGE_NAME"; then
        echo "Error: Image push failed."
        echo "Please ensure you are logged in to $REGISTRY_HOST ($current_container_engine login $REGISTRY_HOST)"
        return 1
    fi
    echo "Image $REMOTE_IMAGE_NAME pushed successfully."
    return 0
}

# Determine script mode
SCRIPT_MODE="default"
if [[ "$1" == "only-push" ]]; then
    SCRIPT_MODE="only-push"
fi

# Check if podman is available
if command -v podman &> /dev/null; then
    CONTAINER_ENGINE="podman"
# Otherwise, check if docker is available
elif command -v docker &> /dev/null; then
    CONTAINER_ENGINE="docker"
else
    echo "Error: Neither podman nor docker found. Please install one of them."
    exit 1
fi
echo "Will use $CONTAINER_ENGINE."

if [[ "$SCRIPT_MODE" == "only-push" ]]; then
    echo "--- Push-Only Mode ---"
    echo "Available local images:"
    "$CONTAINER_ENGINE" images

    read -p "Please enter the name of the local image to push (e.g., repository:tag): " LOCAL_IMAGE_TO_PUSH
    if [ -z "$LOCAL_IMAGE_TO_PUSH" ]; then
        echo "Error: Image name cannot be empty."
        exit 1
    fi

    # Check if the image exists
    if ! "$CONTAINER_ENGINE" image inspect "$LOCAL_IMAGE_TO_PUSH" &> /dev/null; then
        echo "Error: Local image $LOCAL_IMAGE_TO_PUSH not found."
        exit 1
    fi

    echo "Preparing to push image: $LOCAL_IMAGE_TO_PUSH"
    if perform_push_operations "$LOCAL_IMAGE_TO_PUSH" "$CONTAINER_ENGINE"; then
        echo "Image push process completed."
    else
        echo "Image push process failed."
        exit 1
    fi
else # Default mode (build and then optionally push)
    echo "--- Default Build and Push Mode ---"

    # Interactively ask the user for the image tag
    read -p "Please enter the image tag (e.g., latest, 1.0): " IMAGE_TAG

    # Check if the tag is empty
    if [ -z "$IMAGE_TAG" ]; then
        echo "Error: Image tag cannot be empty."
        exit 1
    fi

    BUILT_IMAGE_BASENAME="ingester" # As implied by original script
    LOCAL_BUILT_IMAGE_NAME="${BUILT_IMAGE_BASENAME}:${IMAGE_TAG}"
    DOCKERFILE_PATH="Dockerfile" # Assume Dockerfile is in the current directory

    # Check if Dockerfile exists
    if [ ! -f "$DOCKERFILE_PATH" ]; then
        echo "Error: $DOCKERFILE_PATH not found in the current directory."
        echo "Please ensure Dockerfile exists in the directory where the script is executed."
        exit 1
    fi

    echo "Building image $LOCAL_BUILT_IMAGE_NAME ..."

    # Build the image
    if "$CONTAINER_ENGINE" build -t "$LOCAL_BUILT_IMAGE_NAME" -f "$DOCKERFILE_PATH" .; then
        echo "Image $LOCAL_BUILT_IMAGE_NAME built successfully."
    else
        echo "Error: Image build failed."
        exit 1
    fi

    # Ask whether to push to the image registry
    read -p "Do you want to push the image $LOCAL_BUILT_IMAGE_NAME to the image registry? (y/n): " PUSH_CONFIRMATION

    if [[ "$PUSH_CONFIRMATION" == "y" || "$PUSH_CONFIRMATION" == "Y" ]]; then
        if perform_push_operations "$LOCAL_BUILT_IMAGE_NAME" "$CONTAINER_ENGINE"; then
            echo "Image push process completed."
        else
            echo "Image push process failed."
            exit 1
        fi
    else
        echo "Image $LOCAL_BUILT_IMAGE_NAME was not pushed to the registry."
    fi
fi

echo "Script execution finished."
exit 0

# Check if podman is available
if command -v podman &> /dev/null; then
    CONTAINER_ENGINE="podman"
# Otherwise, check if docker is available
elif command -v docker &> /dev/null; then
    CONTAINER_ENGINE="docker"
else
    echo "Error: Neither podman nor docker found. Please install one of them."
    exit 1
fi

echo "Will use $CONTAINER_ENGINE to build the image."

# Interactively ask the user for the image tag
read -p "Please enter the image tag (e.g., latest, 1.0): " IMAGE_TAG

# Check if the tag is empty
if [ -z "$IMAGE_TAG" ]; then
    echo "Error: Image tag cannot be empty."
    exit 1
fi

LOCAL_IMAGE_NAME="ingester:$IMAGE_TAG"
DOCKERFILE_PATH="Dockerfile" # Assume Dockerfile is in the current directory

# Check if Dockerfile exists
if [ ! -f "$DOCKERFILE_PATH" ]; then
    echo "Error: $DOCKERFILE_PATH not found in the current directory."
    echo "Please ensure Dockerfile exists in the directory where the script is executed."
    exit 1
fi

echo "Building image $LOCAL_IMAGE_NAME ..."

# Build the image
if "$CONTAINER_ENGINE" build -t "$LOCAL_IMAGE_NAME" -f "$DOCKERFILE_PATH" .; then
    echo "Image $LOCAL_IMAGE_NAME built successfully."
else
    echo "Error: Image build failed."
    exit 1
fi

# Ask whether to push to the image registry
read -p "Do you want to push the image to the image registry? (y/n): " PUSH_TO_REGISTRY

if [[ "$PUSH_TO_REGISTRY" == "y" || "$PUSH_TO_REGISTRY" == "Y" ]]; then
    DEFAULT_REGISTRY_HOST="greptime-registry.cn-hangzhou.cr.aliyuncs.com/tools"
    read -p "Please enter the image registry address (default is $DEFAULT_REGISTRY_HOST): " REGISTRY_HOST
    if [ -z "$REGISTRY_HOST" ]; then
        REGISTRY_HOST="$DEFAULT_REGISTRY_HOST"
        echo "Using default image registry address: $REGISTRY_HOST"
    fi

    REMOTE_IMAGE_NAME="$REGISTRY_HOST/$LOCAL_IMAGE_NAME"

    echo "Tagging image $LOCAL_IMAGE_NAME as $REMOTE_IMAGE_NAME ..."
    if "$CONTAINER_ENGINE" tag "$LOCAL_IMAGE_NAME" "$REMOTE_IMAGE_NAME"; then
        echo "Image tagged successfully."
    else
        echo "Error: Image tagging failed."
        exit 1
    fi

    echo "Pushing image $REMOTE_IMAGE_NAME ..."
    if "$CONTAINER_ENGINE" push "$REMOTE_IMAGE_NAME"; then
        echo "Image $REMOTE_IMAGE_NAME pushed successfully."
    else
        echo "Error: Image push failed."
        echo "Please ensure you are logged in to $REGISTRY_HOST ($CONTAINER_ENGINE login $REGISTRY_HOST)"
        exit 1
    fi
else
    echo "Image was not pushed to the registry."
fi