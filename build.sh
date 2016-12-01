#!/usr/bin/env bash
CGO_ENABLED=0 go build -a -installsuffix cgo
docker build -t sortinghat-logingest .
docker tag sortinghat-logingest:latest 514845858982.dkr.ecr.us-west-1.amazonaws.com/sortinghat-logingest:latest
docker push 514845858982.dkr.ecr.us-west-1.amazonaws.com/sortinghat-logingest:latest
kubectl --kubeconfig=./kubeconfig delete -f ./kube/sortinghat-logingest-deployment.yaml
kubectl --kubeconfig=./kubeconfig create -f ./kube/sortinghat-logingest-deployment.yaml