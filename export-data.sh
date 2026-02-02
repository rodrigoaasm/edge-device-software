#!/bin/bash
HOST=${1}
DEVICE=${2}
LAST="-2d"

if [[ -n "${3}" ]]; then
  LAST=${3}
fi

if [[ -z "$HOST" || -z "$DEVICE" ]]; then
  echo "Use: $0 <host> <device>"
  echo "Ex.: $0 192.168.0.8:30086 abc123 -2d"
  exit 1
fi

curl --location "${HOST}/api/v2/query?org=edge-device" \
--header 'Content-Type: application/vnd.flux' \
--header 'Authorization: Token edge-device@token' \
--data "from(bucket: \"history\")
  |> range(start: ${LAST})
  |> filter(fn: (r) => r[\"_measurement\"] == \"${DEVICE}\")" > "${DEVICE}.csv"