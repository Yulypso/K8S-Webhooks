#!/usr/bin/env bash 

echo "Running Test"
cd WebhookServer/test

go test -v -coverprofile=cover.txt
#go tool cover -html=cover.txt -o cover.html