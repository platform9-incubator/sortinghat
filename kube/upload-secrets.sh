#!/usr/bin/env bash
aws ecr get-login --region us-west-1 | sh -
cat > /tmp/image-pull-secret.yaml << EOF
apiVersion: v1
kind: Secret
metadata:
  name: rparikhregistrykey
data:
  .dockerconfigjson: $(cat ~/.docker/config.json | base64 -w 0)
type: kubernetes.io/dockerconfigjson
EOF
kubectl --kubeconfig=../kubeconfig replace -f /tmp/image-pull-secret.yaml