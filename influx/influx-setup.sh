influx setup \
  --force \
  --host "http://influxdb:8086"  \
  --username "admin" \
  --password "qweasd123" \
  --org "edge-device" \
  --bucket "telemetry" \
  --token  "edge-device@token" \
  --retention "7d"