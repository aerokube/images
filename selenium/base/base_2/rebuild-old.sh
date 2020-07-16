#!/bin/bash

set -e

if [ -z "$1" ]; then
    echo 'Usage: rebuild-old.sh <image_ref>'
    exit 1
fi
set -x

image=$1

dir_name="/tmp/$(uuidgen | sed -e 's|-||g')"
mkdir -p "$dir_name"

cat "Dockerfile" | sed -e "s|selenoid/base:1.0|$image|g" > "$dir_name/Dockerfile"
pushd "$dir_name"

additional_docker_args=""
if [ -n "$http_proxy" -a -n "$https_proxy" ]; then
    additional_docker_args+="--build-arg http_proxy=$http_proxy --build-arg https_proxy=$https_proxy "
fi
if [ -n "$HTTP_PROXY" -a -n "$HTTPS_PROXY" ]; then
    additional_docker_args+="--build-arg http_proxy=$HTTP_PROXY --build-arg https_proxy=$HTTPS_PROXY "
fi

echo "Adding fonts and encodings to $image..."
docker build $additional_docker_args -t "$image" .
popd
rm -Rf "$dir_name"
exit 0
