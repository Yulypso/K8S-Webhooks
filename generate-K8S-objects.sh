#!/usr/bin/env bash

# Generate K8S objects
echo "Generate Kubernetes objects ..."

# Create Namespaces
kubectl apply -f "WebhookDeployments/namespaces.yml" # contains webhookserver-namespace and admissionwebhook-namespace

# Create Secrets to store TLS certificate
kubectl -n webhookserver-ns create secret tls webhookserver-tls --key "Certificates/webhookservertls.key" --cert "Certificates/webhookservertls.cert"

# Create Webhookserver-deployment/yml objects
kubectl apply -f "WebhookDeployments/webhookserver-deployment.yml"

# Get CA_B64 flag replaced by ca.cert value within webhook-deployment.yaml 
# Then create K8S validating-webhook.yml object  
ca_cert_B64=`openssl base64 -A < "Certificates/ca.pem"`
sed -e 's/${CA_B64}/'"$ca_cert_B64"'/g' "WebhookDeployments/mutating-webhook.yml" | kubectl apply -f -
sed -e 's/${CA_B64}/'"$ca_cert_B64"'/g' "WebhookDeployments/validating-webhook.yml" | kubectl apply -f -

echo "WebhookServer has been deployed..."