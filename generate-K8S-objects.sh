#!/usr/bin/env bash

# Generate K8S objects
echo "Generate Kubernetes objects ..."

# Create Namespaces
kubectl apply -f "Deployments/Cluster/namespace.yml" # contains webhookserver-namespace and admissionwebhook-namespace

# Create Persistent-volume
data_dir="data"
hostname=`kubectl get nodes -o json | jq -r '.items[].metadata.name'`
volume_path="`pwd`/WebhookServer/$data_dir"

mkdir -p WebhookServer/$data_dir
sed -e 's/${HOSTNAME}/'"$hostname"'/g' "Deployments/Cluster/persistent-volume.yml" | 
sed -e 's*${VOLUME_PATH}*'"$volume_path"'*g' |
kubectl apply -f -

# Create Secrets to store TLS certificate
kubectl -n webhookserver-ns create secret tls webhookserver-tls --key "Certificates/webhookservertls.key" --cert "Certificates/webhookservertls.cert"

# Create Webhookserver-deployment/yml objects
kubectl apply -f "Deployments/Webhooks/webhookserver.yml"

# Get CA_B64 flag replaced by ca.cert value within webhook-deployment.yaml 
# Then create K8S validating-webhook.yml object  
ca_cert_B64=`openssl base64 -A < "Certificates/ca.pem"`
sed -e 's/${CA_B64}/'"$ca_cert_B64"'/g' "Deployments/Webhooks/mutating-webhook.yml" | kubectl apply -f -
sed -e 's/${CA_B64}/'"$ca_cert_B64"'/g' "Deployments/Webhooks/validating-webhook.yml" | kubectl apply -f -

echo "WebhookServer has been deployed..."