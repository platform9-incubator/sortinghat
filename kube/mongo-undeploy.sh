#!/usr/bin/env bash
export KUBECTL="kubectl --kubeconfig=../kubeconfig"
$KUBECTL delete -f mongo-service.yaml
$KUBECTL delete -f mongo-deployment.yaml
$KUBECTL delete -f mongo-volumes.yaml
