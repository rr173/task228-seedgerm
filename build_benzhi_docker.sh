#!/usr/bin/env bash
# 评测构建脚本：双架构镜像构建（评测用 benzhi.Dockerfile）
# usage: bash build_benzhi_docker.sh <镜像名> <平台>
#   平台示例: linux/amd64  linux/arm64  linux/amd64,linux/arm64
set -euo pipefail

IMAGE_NAME="${1:-my-project}"
PLATFORM="${2:-linux/amd64}"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

echo "==> building $IMAGE_NAME for platform(s): $PLATFORM using benzhi.Dockerfile"

docker buildx build \
  --platform "$PLATFORM" \
  --file benzhi.Dockerfile \
  --tag "$IMAGE_NAME" \
  --load \
  .

echo "==> done: $IMAGE_NAME [$PLATFORM]"
