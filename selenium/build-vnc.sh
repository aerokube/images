#!/bin/bash

set -e

if [ -z "$1" -o -z "$2" ]; then
    echo 'Usage: build-vnc.sh <browser_name> <browser_version>'
    exit 1
fi
set -x

browser=$1
version=$2

dir_name="/tmp/$(uuidgen | sed -e 's|-||g')"
mkdir -p "$dir_name"

cat "Dockerfile.tmpl" | sed -e "s|@@VERSION@@|$version|g" > "$dir_name/Dockerfile"
cp "entrypoint.sh" "$dir_name"
pushd "$dir_name"
browser_name=$(echo "$browser" | sed -e 's/\(\/..*\)\+//g')
tag="selenoid/vnc_"$browser_name":"$version
old_tag="selenoid/vnc:"$browser_name"_"$version
echo "Creating VNC image $tag..."
docker build -t "$tag" .
docker tag "$tag" "$old_tag"
popd
rm -Rf "$dir_name"
exit 0
