#!/usr/bin/env bash
docker build -t whistle-ui .
docker tag whistle-ui:latest 514845858982.dkr.ecr.us-west-1.amazonaws.com/whistle-ui:latest
docker push 514845858982.dkr.ecr.us-west-1.amazonaws.com/whistle-ui:latest
kubectl --kubeconfig=../kubeconfig  delete -f ../kube/whistle-ui-deployment.yaml
kubectl --kubeconfig=../kubeconfig  create -f ../kube/whistle-ui-deployment.yaml