#!/bin/bash
set -e

if [ -z $1 -o -z $2 ]; then
    echo 'Usage: build-dev.sh {firefox/official|firefox/ubuntuzilla|chrome|opera/presto|opera/blink} <browser_version> [<cleanup={true|false}>]'
    exit 1
fi

browser=$1
version=$2
cleanup=${3:-"false"}
browser_name=$(echo "$browser" | sed -e 's/\(..*\)\/..*/\1/g')
tag="selenoid/$browser_name:$version"
pushd "$browser"
echo "Creating image $tag..."
docker build --build-arg VERSION="$version" --build-arg CLEANUP="$cleanup" -t "$tag" .
popd
exit 0
