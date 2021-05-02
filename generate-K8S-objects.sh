#!/usr/bin/env bash

# Generate K8S objects
echo "Generate Kubernetes objects ..."

# Create Namespaces
#kubectl create namespace webhookserver-namespace # contains WebhookServer
#kubectl create namespace admissionwebhook-namespace # to test deployments
kubectl apply -f "WebhookDeployments/namespaces.yaml" # contains webhookserver-namespace and admissionwebhook-namespace

# Create Secrets to store TLS certificate
kubectl -n webhookserver-namespace create secret tls webhookserver-tls --key "Certificates/webhookservertls.key" --cert "Certificates/webhookservertls.cert"

# TEST: default namespace
kubectl create secret tls webhookserver-tls --key "Certificates/webhookservertls.key" --cert "Certificates/webhookservertls.cert"


# Get CA_B64 flag replaced by ca.cert value within webhook-deployment.yaml 
# Then create K8S webhook-deployment.yaml objects  
ca_cert_B64=`openssl base64 -A < "Certificates/ca.cert"`
sed -e 's@''${CA_B64}@'"$ca_cert_B64"'@g' < "WebhookDeployments/webhook-deployment.yaml" \
    | kubectl create -f -

echo "WebhookServer has been deployed..."