#!/bin/bash
set -e
input=$1
driver_version=$2
tag=$3
channel=${4:-"default"}
test_failure_ignore=${TEST_FAILURE_IGNORE:-true}

if [ -z "$1" -o -z "$2" -o -z "$3" ]; then
    echo 'Usage: automate_opera.sh <browser_version|package_file> <operadriver_version> <tag_version> [<channel={beta|dev}>]'
    exit 1
fi
set -x

browser_version=$input
method="opera/blink/apt"
if [ -f "$input" ]; then
    cp "$input" opera/blink/local/opera.deb
    filename=$(echo "$input" | awk -F '/' '{print $NF}')
    browser_version=$(echo $filename | awk -F '_' '{print $2}' | awk -F '+' '{print $1}')
    method="opera/blink/local"
fi

./build-dev.sh $method $browser_version $channel true
if [ "$method" == "opera/blink/apt" ]; then
    ./build-dev.sh $method $browser_version $channel false
fi
pushd opera/blink
../../build.sh operadriver $browser_version $driver_version selenoid/opera:$tag
popd

test_image(){
    docker rm -f selenium || true
    docker run -d --name selenium -p 4445:4444 $1:$2
    tests_dir=../../selenoid-container-tests/
    if [ -d "$tests_dir" ]; then
        pushd "$tests_dir"
        mvn clean test -Dgrid.connection.url="http://localhost:4445/" -Dgrid.browser.name=opera -Dgrid.browser.version=$2 || $test_failure_ignore
        popd
    else
        echo "Skipping tests as $tests_dir does not exist."
    fi
}

test_image "selenoid/opera" $tag
docker tag "selenoid/opera:$tag" "selenoid/vnc_opera:$tag"
docker tag "selenoid/opera:$tag" "selenoid/vnc:opera_$tag"

read -p "Push?" yn
if [ "$yn" == "y" ]; then
	docker push "selenoid/dev_opera:"$browser_version
	if [ "$method" == "opera/blink/apt" ]; then
	    docker push "selenoid/dev_opera_full:"$browser_version
    fi
	docker push "selenoid/opera:$tag"
    docker tag "selenoid/opera:$tag" "selenoid/opera:latest"
    docker push "selenoid/opera:latest"
    docker push "selenoid/vnc:opera_"$tag
    docker push "selenoid/vnc_opera:"$tag
fi
