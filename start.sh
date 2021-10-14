#!/usr/bin/env bash

# build and start gRPC server ( with name vogrpc )
docker build -f Dockerfile -t vogrpc .
docker run  -p 50515:50515 vogrpc:latest

# build and start envoy proxy & - link it with the grpc container
docker build -f envoy.dockerfile -t envoy .
docker run -p 8080:8080 envoy:latest --link vogrpc:vogrpc

# TODO: frontend ?