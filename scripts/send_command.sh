#!/bin/bash
mosquitto_pub -d \
  -h localhost -p 31883 \
  -t device/data \
  -m 'dd17dd temperatura=71 1748613118000000000'



# '{
#      "command": "update",
#      "args": { 
#        "name":"operator",
#        "image":"rodrigoasmaia/ed-operator:0.4.8"
#      }
#    }' 



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