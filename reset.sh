#!/usr/bin/env bash 

echo "Reset Kubernetes Cluster..."

echo "Deleting Kubernetes Objects..."
kubectl delete -f "Deployments/Webhooks/webhookserver.yml" # delete webhookserver
kubectl delete -f "Deployments/Webhooks/mutating-webhook.yml" # delete webhookconfigurations
kubectl delete -f "Deployments/Webhooks/validating-webhook.yml" # delete webhookconfigurations
kubectl delete -f "Deployments/Cluster/persistent-volume.yml" # delete persistent volumes
kubectl delete -f "Deployments/Cluster/namespace.yml" # delete namespaces

shopt -s nocasematch
if [[ "$1" =~ ^certificate ]]; then
    echo "Deleting Certificates..."
    rm -rf "./Certificates"
    shift
fi

if [[ "$1" =~ ^docker ]]; then
    echo "Go build..."
    rm WebhookServer/cmd/server/webhookserver
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o WebhookServer/cmd/server/webhookserver WebhookServer/cmd/server/*.go
    echo "Docker build..."
    docker build -t yulypso/webhookserver:v0.0.6 . 
    echo "Docker push..."
    docker push yulypso/webhookserver:v0.0.6
    shift
fi

echo "Clear!"