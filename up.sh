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
kubectl apply -f roles/role_deployment_binding.yaml
echo "service account created"

echo "creating services..."
kubectl apply -k .

kubectl get services -n ed-system