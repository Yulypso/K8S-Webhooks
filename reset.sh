#!/usr/bin/env bash 

echo "Reset Kubernetes Cluster..."

echo "Deleting Kubernetes Objects..."
kubectl delete -f "Deployments/Webhooks/webhookserver.yml" # delete webhookserver
kubectl delete -f "Deployments/Webhooks/mutating-webhook.yml" # delete webhookconfigurations

shopt -s nocasematch
if [[ "$1" =~ ^certificate ]]; then
    echo "Deleting Certificates..."
    rm -rf "./Certificates"
    shift
fi

if [[ "$1" =~ ^docker ]]; then
    echo "Docker build..."
    docker build -t yulypso/webhookserver:v0.0.6 . 
    #echo "Docker push..."
    #docker push yulypso/webhookserver:v0.0.6
    shift
fi

if [[ "$1" =~ ^cluster ]]; then
    kubectl delete -f "Deployments/Cluster/persistent-volume.yml"
    kubectl delete -f "Deployments/Cluster/namespace.yml"
fi

echo "Clear!"