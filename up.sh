#!/bin/bash
ED_OPERATOR_VERSION="0.3.0"
NANOMQ_VERSION=latest
INFLUXDB_VERSION="0.3.0"
TELEGRAF_BRIDGE_VERSION="0.3.0"
TELEGRAF_AGG_VERSION="0.3.0". 

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

echo "create operator deployment"
pushd "ed-operator" > /dev/null
kustomize edit set image rodrigoasmaia/ed-operator=rodrigoasmaia/ed-operator:$ED_OPERATOR_VERSION
kustomize build . | kubectl apply -f -

echo "create services default:"

echo "- nanomq"
pushd "../deployments-default/nanomq" > /dev/null
kustomize edit set image emqx/nanomq=emqx/nanomq:$NANOMQ_VERSION
kustomize build . | kubectl apply -f -

echo "- Influxdb"
pushd "../influxdb" > /dev/null
kubectl apply -f ../../volumes-default/influxdb_pv.yaml
kustomize edit set image influxdb:2.7.11-alpine=rodrigoasmaia/influxdb:$INFLUXDB_VERSION
kustomize build . | kubectl apply -f -

echo "- telegraf-bridge"
pushd "../telegraf-bridge" > /dev/null
kustomize edit set image rodrigoasmaia/telegraf-bridge=rodrigoasmaia/telegraf-bridge:$TELEGRAF_BRIDGE_VERSION
kustomize build . | kubectl apply -f -

echo "- telegraf-agg"
pushd "../telegraf-agg" > /dev/null
kustomize edit set image rodrigoasmaia/telegraf-agg=rodrigoasmaia/telegraf-agg:$TELEGRAF_AGG_VERSION
kustomize build . | kubectl apply -f -

echo "Opening external ports..."
pushd "../.."
kubectl apply -f nodeport-default/nodeports.yaml

kubectl get services -n ed-system