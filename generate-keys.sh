#!/usr/bin/env bash

# Generate keys into a temporary directory.
echo "Generate TLS keys ..."

mkdir -p "Certificates"
chmod 0700 "Certificates"
cd "Certificates"

# Generate the CA cert and CA private key
openssl req \
        -nodes \
        -new \
        -newkey rsa:4096 \
        -days 365 \
        -x509 \
        -subj "/C=FR/ST=Bezons/L=Bezons/O=Atos/OU= IT Department/CN=Webhook.CA" \
        -keyout "ca.key" \
        -out "ca.cert" 


# Generate the private key and Certificate Signing Request (CSR) for the webhook server
openssl req \
        -nodes \
        -newkey rsa:2048 \
        -subj "/C=FR/ST=Bezons/L=Bezons/O=Atos/OU= IT Department/CN=WebhookServer.WebhookServerNameSpace.svc" \
        -keyout "webhookservertls.key" \
        -out "webhookservertls.csr"


# Sign it with the private key of the CA.
openssl x509 \
        -req \
        -CA ca.cert \
        -CAkey ca.key \
        -CAcreateserial \
        -in "webhookservertls.csr" \
        -out "webhookservertls.cert"

echo "Certificates ready!"