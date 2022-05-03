#!/bin/bash

VERSION=${1:-101.0.4951.64-0ubuntu0.18.04.1}
TAG=${2:-chromium_101.0}
BASE_TAG=${3:-7.3.6}

# Cleanup stuff
export BUILDKIT_PROGRESS=plain
docker rmi -f selenoid/vnc:$TAG browsers/base:$BASE_TAG $(docker images -q selenoid/dev_chromium:*)
rm -rf ../selenoid-container-tests

# Prepare for building images
go get github.com/markbates/pkger/cmd/pkger
go generate github.com/aerokube/images
go build

# Forked tests with a bugfix
git clone -b add-missing-dependency https://github.com/sskorol/selenoid-container-tests.git ../selenoid-container-tests

# Force build browsers/base image as it has arm64-specific updates
cd ./selenium/base && docker build --no-cache --build-arg UBUNTU_VERSION=18.04 -t browsers/base:$BASE_TAG . && docker system prune -f

# Build chromium image
cd ../../ && ./images chromium -b $VERSION -t selenoid/vnc:$TAG --test && docker system prune -f
