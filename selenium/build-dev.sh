#!/bin/bash
set -e

if [ -z "$1" -o -z "$2" ]; then
    echo 'Usage: build-dev.sh {firefox/apt|firefox/ubuntuzilla|chrome/apt|chrome/local|opera/presto|opera/blink/local|opera/blink/apt} <browser_version> [<cleanup={true|false}>] [<requires_java={true|false}>]'
    exit 1
fi
set -x

browser=$1
version=$2
cleanup=${3:-"false"}
requires_java=${4:-"false"}
tag_version="$version"
if [ -n "$5" ]; then
    tag_version=$5
fi
browser_name=$(echo "$browser" | sed -e 's/\(\/..*\)\+//g')
tag="selenoid/dev:"$browser_name"_"$tag_version
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
    if [ "$browser" == "chrome/local" -o "$browser" == "opera/blink/local" ]; then
        debWildcard="$browser/*.deb"
        cp $debWildcard "$dir_name"
    fi
    cp "$browser/Dockerfile" "$dir_name"
fi
pushd "$dir_name"
echo "Creating image $tag with cleanup=$cleanup..."
docker build --build-arg VERSION="$version" --build-arg CLEANUP="$cleanup" -t "$tag" .
popd
rm -Rf "$dir_name"
exit 0
