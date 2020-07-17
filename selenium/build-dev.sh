#!/bin/bash
set -e

if [ -z "$1" -o -z "$2" -o -z "$3" ]; then
    echo 'Usage: build-dev.sh {firefox/apt|firefox/local|chrome/apt|chrome/local|opera/presto|opera/blink/local|opera/blink/apt|yandex/local|yandex/apt} <browser_version> {default|beta|dev|esr} [<cleanup={true|false}>] [<requires_java={true|false}>]'
    exit 1
fi
set -x

browser=$1
version=$2
channel=$3
cleanup=${4:-"false"}
requires_java=${5:-"false"}
tag_version="$version"
if [ -n "$6" ]; then
    tag_version=$6
fi
disable_docker_cache=${DISABLE_DOCKER_CACHE:-false}
browser_name="${browser%%/*}"
image_name="dev_"$browser_name
if [ "$cleanup" == "false" ]; then
    image_name=$image_name"_full"
fi
tag="selenoid/"$image_name":"$tag_version
dir_name="/tmp/$(uuidgen | sed -e 's|-||g')"
mkdir -p "$dir_name"
additional_docker_args=""
if [ -n "$http_proxy" -a -n "$https_proxy" ]; then
    additional_docker_args+="--build-arg http_proxy=$http_proxy --build-arg https_proxy=$https_proxy "
fi
if [ -n "$HTTP_PROXY" -a -n "$HTTPS_PROXY" ]; then
    additional_docker_args+="--build-arg http_proxy=$HTTP_PROXY --build-arg https_proxy=$HTTPS_PROXY "
fi
if [ "$browser" == "firefox/local" -o "$browser" == "firefox/apt" ]; then
    requires_java_value=""
    if [ "$requires_java" == "true" ]; then
        requires_java_value="_java"
    fi
    cat "$browser/Dockerfile.tmpl" | sed -e "s|@@REQUIRES_JAVA@@|$requires_java_value|g" > "$dir_name/Dockerfile"
else
    cp "$browser/Dockerfile" "$dir_name"
fi
if [ "$browser" == "chrome/local" -o "$browser" == "opera/blink/local" -o "$browser" == "firefox/local" -o "$browser" == "yandex/local" ]; then
    debWildcard="$browser/*.deb"
    cp $debWildcard "$dir_name"
    docker rm -f static-server || true
    docker run -d -p 8080:8043 -v "$dir_name":/srv/http --name static-server pierrezemb/gostatic
    if [ $(uname) == "Darwin" ]; then
        host_ip=$(ifconfig | grep -E "([0-9]{1,3}\.){3}[0-9]{1,3}" | grep -v 127.0.0.1 | awk '{ print $2 }' | cut -f2 -d: | head -n1)
    else
        host_ip=$(ifconfig docker0 | grep inet | grep -v inet6 | awk '{print $2;}' | sed -e 's|addr:||g')
    fi
    if [ -z "$host_ip" ]; then
        echo "Failed to determine host machine IP..."
        exit 1
    fi
    additional_docker_args="--add-host apt-repo:$host_ip"
fi
if [ "$channel" != "default" ]; then
    case $browser_name in
        firefox)
            case $channel in
                beta)
                    additional_docker_args+=" --build-arg PPA=ppa:mozillateam/firefox-next"
                    ;;
                dev)
                    additional_docker_args+=" --build-arg PACKAGE=firefox-trunk --build-arg PPA=ppa:ubuntu-mozilla-daily/ppa"
                    ;;
                esr)
                    additional_docker_args+=" --build-arg PACKAGE=firefox-esr --build-arg PPA=ppa:mozillateam/ppa"
                    ;;
            esac
            ;;
        chrome)
            case $channel in
                beta)
                    additional_docker_args+=" --build-arg PACKAGE=google-chrome-beta --build-arg INSTALL_DIR=chrome-beta"
                    ;;
                dev)
                    additional_docker_args+=" --build-arg PACKAGE=google-chrome-unstable --build-arg INSTALL_DIR=chrome-unstable"
                    ;;
            esac
            ;;
        opera)
            case $channel in
                beta)
                    additional_docker_args+=" --build-arg PACKAGE=opera-beta"
                    ;;
                dev)
                    additional_docker_args+=" --build-arg PACKAGE=opera-developer"
                    ;;
            esac
            ;;
    esac
fi
if [ "$disable_docker_cache" == "true" ]; then
    additional_docker_args+=" --no-cache"
fi
pushd "$dir_name"
echo "Creating image $tag with cleanup=$cleanup..."
docker build $additional_docker_args --build-arg VERSION="$version" --build-arg CLEANUP="$cleanup" -t "$tag" .
popd
rm -Rf "$dir_name"
exit 0
