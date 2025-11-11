#!/usr/bin/env bash
set -euxo pipefail

#TAG=${VERSION}-${PLATFORM}
TAG=${VERSION}
IMAGE=ghcr.io/${ORG}/${COMPONENT}:${TAG}

docker build -f perforator/deploy/docker/Dockerfile.prebuilt --platform ${PLATFORM} --target ${COMPONENT} --push --tag ${IMAGE} ./dist
if [[ -z "${EXTRA_CR}" ]]; then
  echo "Additional container registry not configured"
else
  IMAGE2=${EXTRA_CR}/${COMPONENT}:${TAG}
  docker tag ${IMAGE} ${IMAGE2}
  docker push ${IMAGE2}
fi
