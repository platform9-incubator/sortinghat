#!/usr/bin/env bash
export KUBECTL="kubectl --kubeconfig=../kubeconfig"
$KUBECTL create -f mongo-volumes.yaml
$KUBECTL create -f mongo-deployment.yaml
$KUBECTL create -f mongo-service.yaml
