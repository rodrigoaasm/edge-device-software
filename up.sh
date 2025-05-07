#!/bin/bash
## create roles
kubectl apply -f roles/role_manager.yaml
kubectl apply -f roles/role_deployment.yaml

## create namespace
kubectl create namespace ed-system

## create service account for operator
kubectl apply -f roles/service_account.yaml
kubectl apply -f roles/role_manager_binding.yaml
kubectl apply -f roles/role_manager_deployment.yaml

## create operator deployment
kubectl apply -f ed-operator/deployment.yaml

## create services default
kubectl apply -f deployments-default/nanomq/deployment.yaml

