#!/bin/bash
mosquitto_sub -d \
  -h 10.42.0.101 -p 1883 \
  -t deployments/results
