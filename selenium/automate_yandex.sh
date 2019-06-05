#!/bin/bash
set -e
input=$1
driver_version=$2
tag=$3

if [ -z "$1" -o -z "$2" -o -z "$3" ]; then
    echo 'Usage: automate_yandex.sh <browser_version|package_file> <operadriver_version> <tag_version>'
    exit 1
fi
set -x

browser_version=$input
method="yandex/apt"
if [ -f "$input" ]; then
    cp "$input" yandex/local/yandex-browser.deb
    filename=$(echo "$input" | awk -F '/' '{print $NF}')
    browser_version=$(echo $filename | awk -F '_' '{print $2}' | awk -F '+' '{print $1}')
    method="yandex/local"
fi

./build-dev.sh $method $browser_version true
if [ "$method" == "yandex/apt" ]; then
    ./build-dev.sh $method $browser_version false
fi
pushd yandex
../build.sh operadriver $browser_version $driver_version selenoid/yandex:$tag
popd

test_image(){
    docker rm -f selenium || true
    docker run -d --name selenium -p 4445:4444  $1:$2
    tests_dir=../../selenoid-container-tests/
    if [ -d "$tests_dir" ]; then
        pushd "$tests_dir"
        mvn clean test -Dgrid.connection.url="http://localhost:4445/" -Dgrid.browser.name=yandex -Dgrid.browser.version=$2 || true
        popd
    else
        echo "Skipping tests as $tests_dir does not exist."
    fi
}

test_image "selenoid/yandex" $tag
docker tag "selenoid/yandex:$tag" "selenoid/vnc_yandex:$tag"
docker tag "selenoid/yandex:$tag" "selenoid/vnc:yandex_$tag"

read -p "Push?" yn
if [ "$yn" == "y" ]; then
	docker push "selenoid/dev_yandex:"$browser_version
	if [ "$method" == "yandex/apt" ]; then
	    docker push "selenoid/dev_yandex_full:"$browser_version
    fi
	docker push "selenoid/yandex:$tag"
    docker tag "selenoid/yandex:$tag" "selenoid/yandex:latest"
    docker push "selenoid/yandex:latest"
    docker push "selenoid/vnc:yandex_"$tag
    docker push "selenoid/vnc_yandex:"$tag
fi
