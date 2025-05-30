#!/bin/bash

curl --location 'http://10.43.10.191:8086/api/v2/query?org=edge-device' \
--header 'Content-Type: application/vnd.flux' \
--header 'Authorization: Token edge-device@token' \
--data 'from(bucket: "history") 
    |> range(start: -7d)'

  