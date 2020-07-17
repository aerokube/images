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
echo "Adding fonts and encodings to $image..."
docker build -t "$image" .
popd
rm -Rf "$dir_name"
exit 0
