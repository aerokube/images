#!/bin/bash
set -e

require_command(){
    cmd_name=$1
    if [ -z $(command -v $1) ]; then
        echo "$1 command required for this script to run"
        exit 1
    fi
}

require_command "awk"
require_command "cut"
require_command "docker"
require_command "ifconfig"
require_command "jq"
require_command "sed"
require_command "true"
require_command "unzip"
require_command "uuidgen"
require_command "wget"

input=$1
driver_version=$2
tag=$3
channel=${4:-"default"}
test_failure_ignore=${TEST_FAILURE_IGNORE:-true}

if [ -z "$1" -o -z "$2" -o -z "$3" ]; then
    echo 'Usage: automate_opera.sh <browser_version|package_file> <operadriver_version|latest> <tag_version> [<channel={beta|dev}>]'
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

get_latest_operadriver() {
    echo "$(wget -qO- "https://api.github.com/repos/operasoftware/operachromiumdriver/releases/latest" | jq -r '.tag_name' | awk -F 'v.' '{print $2}')"
}

if [ "$driver_version" == "latest" ]; then
    driver_version=$(get_latest_operadriver)
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
