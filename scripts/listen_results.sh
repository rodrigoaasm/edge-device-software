#!/bin/bash
mosquitto_sub -d \
  -h 10.42.0.126 -p 1883 \
  -t deployments/#
  
