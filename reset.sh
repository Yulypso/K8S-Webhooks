#!/usr/bin/env bash 

echo "Reset Kubernetes Cluster..."

echo "Deleting Kubernetes Objects..."
kubectl delete -f "WebhookDeployments/webhookserver-deployment.yml" # delete webhookserver deployment
kubectl delete -f "WebhookDeployments/validating-webhook.yml" # delete webhookconfigurations
kubectl delete -f "WebhookDeployments/namespaces.yml" # delete namespaces + secrets contained
kubectl delete -f "TestDeployments/pod-1.yml"

shopt -s nocasematch # ignore case sensitive
if [[ "$1" =~ ^certificate ]]; then
    echo "Deleting Certificates..."
    rm -rf "./Certificates"
    shift
fi

if [[ "$1" =~ ^docker ]]; then
    echo "Docker build..."
    docker build -t yulypso/webhookserver:v0.0.2 . 
    echo "Docker push..."
    docker push yulypso/webhookserver:v0.0.2
fi

echo "Clear!"