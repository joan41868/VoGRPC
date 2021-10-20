#!/usr/bin/env bash

# build local image
docker build -t joan41868/vogrpc:latest .

# push to the heroku docker registry
heroku container:push web -a vogrpc-test-be

# start the app
heroku container:release web