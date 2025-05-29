#!/bin/bash
echo "creating roles"
kubectl apply -f roles/role_manager.yaml
kubectl apply -f roles/role_deployment.yaml
echo "roles created"

echo "create namespace"
kubectl create namespace ed-system

echo "creating service account for operator"
kubectl apply -f roles/service_account.yaml
kubectl apply -f roles/role_manager_binding.yaml
kubectl apply -f roles/role_manager_deployment.yaml
echo "service account created"

echo "create operator deployment"
kubectl apply -f ed-operator/deployment.yaml

echo "create services default:"
echo "- nanomq"
kubectl apply -f deployments-default/nanomq/deployment.yaml
echo "- Influxdb"
kubectl apply -f deployments-default/influxdb/deployment.yaml
echo "- telegraf-bridge"
kubectl apply -f deployments-default/telegraf-bridge/deployment.yaml
