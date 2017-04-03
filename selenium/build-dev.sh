#!/bin/bash
set -e

if [ -z "$1" -o -z "$2" ]; then
    echo 'Usage: build-dev.sh {firefox/official|firefox/ubuntuzilla|chrome|opera/presto|opera/blink} <browser_version> [<cleanup={true|false}>] [<requires_java={true|false}>]'
    exit 1
fi
set -x

browser=$1
version=$2
cleanup=${3:-"false"}
requires_java=${4:-"false"}
browser_name=$(echo "$browser" | sed -e 's/\(..*\)\/..*/\1/g')
tag="selenoid/dev:"$browser_name"_"$version
if [ "$cleanup" == "false" ]; then
    tag=$tag"_full"
fi
dir_name="/tmp/$(uuidgen | sed -e 's|-||g')"
mkdir -p "$dir_name"
if [ "$browser" == "firefox/ubuntuzilla" ]; then
    requires_java_value=""
    if [ "$requires_java" == "true" ]; then
        requires_java_value="_java"
    fi
    cat "$browser/Dockerfile.tmpl" | sed -e "s|@@REQUIRES_JAVA@@|$requires_java_value|g" > "$dir_name/Dockerfile"
else
    cp "$browser/Dockerfile" "$dir_name"
fi
pushd "$dir_name"
echo "Creating image $tag with cleanup=$cleanup..."
docker build --build-arg VERSION="$version" --build-arg CLEANUP="$cleanup" -t "$tag" --no-cache .
popd
rm -Rf "$dir_name"
exit 0
