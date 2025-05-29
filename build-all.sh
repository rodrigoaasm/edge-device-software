#!/bin/bash
set -e

VERSION=${1}

if [[ -z "$VERSION" ]]; then
  echo "Use: $0 <versão>"
  echo "Ex.: $0 1.0.0"
  exit 1
fi

./build.sh telegraf-bridge $VERSION
./build.sh telegraf-agg $VERSION
./build.sh influxdb $VERSION
./build.sh ed-operator $VERSION