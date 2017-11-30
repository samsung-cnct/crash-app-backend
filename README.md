# Crash Application Backend
[![pipeline status](https://git.cnct.io/common-tools/samsung-cnct_crash-app-backend/badges/master/pipeline.svg)](https://git.cnct.io/common-tools/samsung-cnct_crash-app-backend/commits/master)

This is a reverse proxy validator that validates input from the kraken crash application destined for elasticsearch.

The crashbackend image is running in the common tools kubernetes production cluster along with elasticsearch.

## Prerequisites
Docker

##  Getting Started
    git clone https://github.com/samsung-cnct/crash-app-backend
    cd crash-app-backend

    // build the app, the binary is output to the _containerize dir 
    make clean build-app

    // run all linters and tests
    make test

    // To build a new container
    export DOCKER_REPO='quay.io/yourrepo'
    ./build.sh --kube -- clean build-image

    // push it to your docker repo
    ./build.sh --kube — push

## Cobra startup
The _containerize/Dockerfile shows how the crashbackend server is started.

    crashbackend serve --target http://elasticsearch:9200 --ratelimit 60

  





