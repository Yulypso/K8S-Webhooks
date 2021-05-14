#!/usr/bin/env bash

# Generate keys into a temporary directory.
echo "Generate TLS keys ..."

mkdir -p "Certificates"
chmod 0700 "Certificates"
cd "Certificates"

# Generate the CA cert and CA private key
openssl req \
        -nodes \
        -sha256 \
        -new \
        -newkey rsa:4096 \
        -days 3650 \
        -x509 \
        -subj "/C=FR/ST=Bezons/L=Bezons/O=Atos/OU=IT Department/CN=webhookserver.webhookserver-ns.svc" \
        -keyout "ca-key.pem" \
        -out "ca.pem" 

# Generate the private key and Certificate Signing Request (CSR) for the webhook server
openssl req \
        -new \
        -nodes \
        -sha256 \
        -subj "/CN=webhookserver.webhookserver-ns.svc" \
        -config "../webhookserver.csr.cnf" \
        -keyout "webhookservertls.key" \
        -out "webhookservertls.csr"


# Sign it with the private key of the CA.
openssl x509 \
        -req \
        -sha256 \
        -days 3650 \
        -CA "ca.pem" \
        -CAkey "ca-key.pem" \
        -CAcreateserial \
        -in "webhookservertls.csr" \
        -out "webhookservertls.cert" \
        -extfile "../webhookserver_v3.ext"

echo "Certificates ready!"