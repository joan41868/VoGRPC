#!/usr/bin/env bash

# build and start gRPC server ( with name vogrpc )
docker build -f Dockerfile -t vogrpc .
docker run  -p 50515:50515 vogrpc:latest
