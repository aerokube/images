#!/bin/bash
set -e
input=$1
server_version=$2
tag=$3
test_failure_ignore=${TEST_FAILURE_IGNORE:-true}

if [ -z "$1" -o -z "$2" -o -z "$3" ]; then
    echo 'Usage: automate_firefox.sh <browser_version|package_file> <selenium_version|selenoid_version> <tag_version> [<geckodriver_version>] [<channel={beta|dev}>]'
    exit 1
fi
set -x

browser_version=$input
method="firefox/apt"
runner="selenoid"
requires_java="false"
numeric_version=$(echo "$tag" | awk -F '.' '{print $1}' )
driver_version=""
if [ $numeric_version -lt 48 ]; then
    channel=${4:-"default"}
    runner="selenium"
    requires_java="true"
else
    driver_version=$4
    channel=${5:-"default"}
    if [ -z "$driver_version" ]; then
        echo 'Driver version is required for Firefox 48 and above'
        exit 1
    fi
fi

if [ -f "$input" ]; then
    filename=$(echo "$input" | awk -F '/' '{print $NF}')
    arch=$(echo "$filename" | awk -F '_' '{print $NF}' | sed -e 's|.deb||g')
    rm -f firefox/local/firefox*.deb
    cp "$input" firefox/local/firefox_$arch.deb
    browser_version=$(echo $filename | sed 's/^[^_]*_//g' | awk -F '-' '{print $1}')
    method="firefox/local"
fi

./build-dev.sh $method $browser_version $channel true $requires_java $tag
if [ "$method" == "firefox/apt" ]; then
    ./build-dev.sh $method $browser_version $channel false $requires_java $tag
fi
pushd firefox/$runner
../../build.sh $runner $tag $server_version selenoid/firefox:$tag "$driver_version"
popd

test_image(){
    docker rm -f selenium || true
    docker run -d --name selenium -p 4445:4444 $1:$2
    tests_dir=../../selenoid-container-tests/
    if [ -d "$tests_dir" ]; then
        pushd "$tests_dir"
        mvn clean test -Dgrid.connection.url="http://localhost:4445/wd/hub" -Dgrid.browser.name=firefox -Dgrid.browser.version=$2 || $test_failure_ignore
        popd
    else
        echo "Skipping tests as $tests_dir does not exist."
    fi
}

test_image "selenoid/firefox" $tag
docker tag "selenoid/firefox:$tag" "selenoid/vnc_firefox:$tag"
docker tag "selenoid/firefox:$tag" "selenoid/vnc:firefox_$tag"

read -p "Push?" yn
if [ "$yn" == "y" ]; then
	docker push "selenoid/dev_firefox:"$tag
	if [ "$method" == "firefox/apt" ]; then
	    docker push "selenoid/dev_firefox_full:"$tag
    fi
	docker push "selenoid/firefox:$tag"
    docker tag "selenoid/firefox:$tag" "selenoid/firefox:latest"
    docker push "selenoid/firefox:latest"
    docker push "selenoid/vnc:firefox_"$tag
    docker push "selenoid/vnc_firefox:"$tag
fi
