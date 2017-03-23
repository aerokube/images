#!/bin/bash
set -e
if [ -z $1 -o -z $2 ]; then
    echo 'Usage: build.sh <chrome_version> <driver_version> [<screen_resolution in form 1280x1600x24>]'
    exit 1
fi
version=$1
screenResolution=${3:-"1280x1600x24"}
dirName="/tmp/$(uuidgen | sed -e 's|-||g')"
mkdir -p "$dirName"
cat Docker.driver.tmpl | sed -e "s|@@VERSION@@|$1|g" > "$dirName/Dockerfile"
pushd "$dirName"
docker build --build-arg DRIVER_VERSION="$2" --build-arg SCREEN_RESOLUTION="$screenResolution" .
popd
rm -Rf "$dirName"
