#!/bin/bash
set -e

if [ -z "$1" -o -z "$2" ]; then
    echo 'Usage: build-dev.sh {firefox/official|firefox/ubuntuzilla|chrome|opera/presto|opera/blink} <browser_version> [<cleanup={true|false}>]'
    exit 1
fi
set -x

browser=$1
version=$2
cleanup=${3:-"false"}
browser_name=$(echo "$browser" | sed -e 's/\(..*\)\/..*/\1/g')
tag="selenoid/dev:"$browser_name"_"$version
if [ "$cleanup" == "false" ]; then
    tag=$tag"_full"
fi
pushd "$browser"
echo "Creating image $tag with cleanup=$cleanup..."
docker build --build-arg VERSION="$version" --build-arg CLEANUP="$cleanup" -t "$tag" .
popd
exit 0
