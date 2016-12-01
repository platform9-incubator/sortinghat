#!/usr/bin/env bash
docker build -t sortinghat-ui .
docker tag sortinghat-ui:latest 514845858982.dkr.ecr.us-west-1.amazonaws.com/sortinghat-ui:latest
docker push 514845858982.dkr.ecr.us-west-1.amazonaws.com/sortinghat-ui:latest
kubectl --kubeconfig=../kubeconfig  delete -f ../kube/sortinghat-ui-deployment.yaml
kubectl --kubeconfig=../kubeconfig  create -f ../kube/sortinghat-ui-deployment.yaml