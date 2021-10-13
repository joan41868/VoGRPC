#! /usr/bin/env bash

 protoc --go_out . \
    --go-grpc_out .  \
proto/protobuf/*.proto

#mv proto/*.pb.go proto/genproto