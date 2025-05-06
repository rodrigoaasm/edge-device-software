#!/bin/bash
mosquitto_pub -d \
  -h 10.42.0.101 -p 1883 \
  -t deployments/start \
  -m '{
    "command": "undeploy",
    "args": { 
      "name":"influxdb"
    }
  }' 

# '{
#     "command": "deploy",
#     "args": { 
#       "name":"influxdb",
#       "image":"rodrigoasmaia/influxdb:2.7.11-alpine"
#     }
#   }'   
  
# '{
#     "command": "deploy",
#     "args": { 
#       "name":"telegraf",
#       "image":"rodrigoasmaia/telegraf:0.1.0"
#     }
#   }' 