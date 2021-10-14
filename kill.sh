#!/usr/bin/env bash

docker kill $(docker ps -a -q)
