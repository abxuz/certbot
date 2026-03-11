#!/bin/sh

REGISTRY=registry.doubi.fun/certbot:latest
docker build --network host -t $REGISTRY $@ . && \
docker image push $REGISTRY && \
docker image rm $REGISTRY && \
docker image prune -f