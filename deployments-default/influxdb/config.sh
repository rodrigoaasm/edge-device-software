#!/bin/bash
set -eu -o pipefail
INIT_FLAG="/etc/influxdb2/initialized"

if [ -f "$INIT_FLAG" ]; then
  echo "InfluxDB already initialized"
  exec influxd
fi

influxd &
PID=$!

until curl --silent http://localhost:8086/health | grep -q '"status":"pass"'; do
  echo "Waiting influxdb..."
  sleep 1
done

HOST=${HOST:-"http://localhost:8086"}
DEFAULT_USER=${DEFAULT_USER:-"eds-software"}
DEFAULT_PASSWORD=${DEFAULT_PASSWORD:-"qweasd123"}
DEFAULT_TOKEN=${DEFAULT_TOKEN:-"edge-device@token"}
DEFAULT_ORGANIZATION=${DEFAULT_ORGANIZATION:-"edge-device"}
DEFAULT_BUCKET=${DEFAULT_BUCKET:-"history"}
DEFAULT_RETENTION=${DEFAULT_RETENTION:-"7d"}

echo "Initializing InfluxDB..."
influx setup \
    --force \
    --host "$HOST" \
    --username "$DEFAULT_USER" \
    --password "$DEFAULT_PASSWORD" \
    --org "$DEFAULT_ORGANIZATION" \
    --bucket "$DEFAULT_BUCKET" \
    --token "$DEFAULT_TOKEN" \
    --retention "$DEFAULT_RETENTION"

touch "$INIT_FLAG"
echo "InfluxDB initialized successfully."

wait $PID