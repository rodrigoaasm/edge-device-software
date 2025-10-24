#!/bin/bash
echo "create namespace"
kubectl create namespace ed-system

if [ -f "./.certs/ca.key" ]; then
  echo "ca already created"
else
  echo "creating ca certificate..."
  openssl genrsa -out ./.certs/ca.key 2048
  openssl req -x509 -new -nodes \
    -key ./.certs/ca.key -sha256 -days 3650 \
    -out ./.certs/ca.pem -config ./.certs/ca.cnf
  echo "ca created"
fi

echo "creating server key"
openssl genrsa -out ./.certs/server.key 2048
echo "creating server csr"
openssl req -new -key ./.certs/server.key \
  -out ./.certs/server.csr -config ./.certs/cert.cnf
echo "creating server certificate"
openssl x509 -req -in ./.certs/server.csr \
-CA ./.certs/ca.pem \
-CAkey ./.certs/ca.key -CAcreateserial \
-out ./.certs/server.pem -days 3650
echo "certificates created"

echo "adding certificate in cluster"
kubectl create secret generic certs-secret -n ed-system \
  --from-file=ca.pem=./.certs/ca.pem \
  --from-file=server.key=./.certs/server.key \
  --from-file=server.pem=./.certs/server.pem
echo "certificate added"

rm ./.certs/server.key
rm ./.certs/server.pem
rm ./.certs/server.csr

echo "creating roles"
kubectl apply -f roles/role_manager.yaml
kubectl apply -f roles/role_deployment.yaml
kubectl apply -f roles/role_service.yaml
kubectl apply -f roles/role_pod_reader.yaml
kubectl apply -f roles/role_events.yaml
kubectl apply -f roles/role_monitoring.yaml
echo "roles created"

echo "creating service account for operator"
kubectl apply -f roles/service_account.yaml
kubectl apply -f roles/role_manager_binding.yaml
kubectl apply -f roles/role_deployment_binding.yaml
kubectl apply -f roles/role_service_binding.yaml
kubectl apply -f roles/role_pod_reader_binding.yaml
kubectl apply -f roles/role_events_binding.yaml
kubectl apply -f roles/role_monitoring_binding.yaml
echo "service account created"

echo "creating profiles"
kubectl apply -f priority-profiles/profiles.yaml

echo "creating services..."
kubectl apply -k .

kubectl get services -n ed-system