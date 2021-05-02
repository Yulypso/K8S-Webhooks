#!/usr/bin/env bash 

echo "Reset Kubernetes Cluster..."

echo "Deleting Kubernetes Objects..."
kubectl delete -f "WebhookDeployments/webhook-deployment.yaml" # delete webhook deployment
kubectl delete -f "WebhookDeployments/namespaces.yaml" # delete namespaces + secrets contained


if [ $"(tr "[:upper:]" "[:lower:]" <<< "$1")" == "certificates" ] || [ $"(tr "[:upper:]" "[:lower:]" <<< "$1")" == "certificate" ]; then
    echo "Deleting Certificates..."
    rm -rf "./Certificates"
fi

echo "Clear!"