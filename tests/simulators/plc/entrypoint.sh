if [ "$TLS_ENABLE" = "true" ]; then  
  cat > opc-ua-agent.cnf <<EOF
[ req ]
default_bits       = 2048
prompt             = no
default_md         = sha256
req_extensions     = req_ext
distinguished_name = dn

[ dn ]
C  = BR
ST = MG
L  = Itajuba
O  = UNIFEI
OU = ED-SYSTEM
CN = opc-ua-agent

[ req_ext ]
subjectAltName = @alt_names
basicConstraints = CA:FALSE
keyUsage = nonRepudiation, digitalSignature, keyEncipherment, dataEncipherment
extendedKeyUsage = clientAuth, serverAuth

[ alt_names ]
DNS.1 = opc-ua-agent
IP.1  = $ALTERNATIVE_DOMAIN
IP.2  = 0.0.0.0
URI.1 = urn:gopcua:server:test
EOF

  echo "creating $DEVICE_ID key"
  openssl genrsa -out /certs/$DEVICE_ID.key 2048
  echo "creating $DEVICE_ID csr"
  openssl req -new -x509 -sha256 -key /certs/$DEVICE_ID.key -out /certs/$DEVICE_ID.pem -days 3650 -config opc-ua-agent.cnf -extensions req_ext
fi

./plc

# openssl genrsa -out opc-ua-agent.key 2048
#  openssl req -new -key opc-ua-agent.key \
#     -out opc-ua-agent.csr
#   openssl x509 -req -in opc-ua-agent.csr \
#   -CA ca.pem \
#   -CAkey ca.key -CAcreateserial \
#   -out opc-ua-agent.pem -days 3650