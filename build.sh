#!/bin/bash
set -e

SERVICE=${1}
VERSION=${2}
BUILD_DIR="deployments-default"

if [[ -n "${3}" ]]; then
  BUILD_DIR=${3}
fi

if [[ -z "$SERVICE" || -z "$VERSION" ]]; then
  echo "Use: $0 <serviço> <versão> <build-dir>"
  echo "Ex.: $0 service v1.0.0"
  exit 1
fi

if [ $SERVICE = 'ed-operator' ]; then
  DOCKERFILE="ed-operator/edge-device-operator/Dockerfile"
  BUILD_DIR="ed-operator/edge-device-operator/"
else
  DOCKERFILE="${BUILD_DIR}/${SERVICE}/Dockerfile"
  BUILD_DIR="${BUILD_DIR}/${SERVICE}/"
fi

if [[ ! -f "$DOCKERFILE" ]]; then
  echo "Erro: Dockerfile not found in: $DOCKERFILE"
  exit 1
fi

echo "🔧 building rodrigoasmaia/${SERVICE}:${VERSION}..."
docker build --rm \
  -t rodrigoasmaia/${SERVICE}:${VERSION} \
  -f "$DOCKERFILE" \
  "$BUILD_DIR"

echo "✅ Build done!"
echo "🔧 Importing  rodrigoasmaia/${SERVICE}:${VERSION} to contained..."
mkdir -p .images
docker save rodrigoasmaia/${SERVICE}:${VERSION} \
  -o .images/rodrigoasmaia-${SERVICE}:${VERSION}.tar

k3s ctr images import .images/rodrigoasmaia-${SERVICE}:${VERSION}.tar
echo "✅ Import done!"

