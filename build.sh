#!/usr/bin/env bash
CGO_ENABLED=0 go build -a -installsuffix cgo
docker build -t whistle-log .
docker tag whistle-log:latest 514845858982.dkr.ecr.us-west-1.amazonaws.com/whistle-log:latest
docker push 514845858982.dkr.ecr.us-west-1.amazonaws.com/whistle-log:latest