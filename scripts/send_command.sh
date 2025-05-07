#!/bin/bash
mosquitto_pub -d \
  -h 10.42.0.126 -p 1883 \
  -t deployments/start \
  -m  '{
     "command": "update",
     "args": { 
       "name":"operator",
       "image":"rodrigoasmaia/ed-operator:0.4.8"
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

#'{
#    "command": "undeploy",
#    "args": { 
#      "name":"influxdb"
#    }
#  }'