#!/bin/bash
VOLUME=""

while getopts "v" opt; do
  case ${opt} in
    v)
      VOLUME=true
      ;;
    \?)
      echo "Uso: $0 [-v]"
      exit 1
      ;;
  esac
done

echo "Removing operator deployment"
kubectl delete deployment operator -n ed-system

echo "Removing operator deployment"
kubectl delete service nanomq -n ed-system
kubectl delete deployment nanomq -n ed-system

echo "Removing influxdb deployment"
kubectl delete service influxdb -n ed-system
kubectl delete deployment influxdb -n ed-system

echo "Removing telegraf-bridge deployment"
kubectl delete deployment telegraf-bridge -n ed-system

echo "Closing external ports"
kubectl delete service influxdb-nodeport -n ed-system
kubectl delete service nanomq-nodeport -n ed-system

if [ "$VOLUME" = true ]; then
  echo "Removing influxdb pvc & pv..."
  kubectl delete pvc influxdb-data-pvc -n ed-system
  kubectl delete pvc influxdb-config-pvc -n ed-system
  kubectl delete pv influxdb-data
  kubectl delete pv influxdb-config

  echo "Removing system dir"
  rm -R /mnt/eds
fi

exit 0