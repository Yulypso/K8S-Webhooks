#!/usr/bin/env bash

# Generate keys into a temporary directory.
echo "Generating TLS keys ..."

mkdir -p "Certificates"
chmod 0700 "Certificates"
cd "Certificates"

# Generate the CA cert and private key
openssl req \
        -nodes \
        -new \
        -newkey rsa:4096 \
        -days 365 \
        -x509 \
        -subj "/C=FR/ST=Bezons/L=Bezons/O=Atos/OU= IT Department/CN=Webhook.CA" \
        -keyout ca.key \
        -out ca.cert 


# Generate a Certificate Signing Request (CSR) for the webhook server private key
openssl req \
        -nodes \
        -newkey rsa:2048 \
        -subj "/CN=WebhookServer.WebhookServerNameSpace.svc" \
        -keyout webhookserver.key \
        -out webhookserver.csr


# Sign it with the private key of the CA.
openssl x509 \
        -req \
        -CA ca.cert \
        -CAkey ca.key \
        -CAcreateserial \
        -in webhookserver.csr \
        -out webhookserver.cert


