#!/bin/bash
set -e
input=$1
driver_version=$2
tag=$3
channel=${4:-"default"}
test_failure_ignore=${TEST_FAILURE_IGNORE:-true}

if [ -z "$1" -o -z "$2" -o -z "$3" ]; then
    echo 'Usage: ./automate_chrome.sh <browser_version|package_file> <chromedriver_version|latest> <tag_version> [<channel={beta|dev}>]'
    exit 1
fi
set -x

browser_version=$input
method="chrome/apt"
if [ -f "$input" ]; then
    cp "$input" chrome/local/google-chrome.deb
    filename=$(echo "$input" | awk -F '/' '{print $NF}')
    browser_version=$(echo $filename | awk -F '_' '{print $2}' | awk -F '-' '{print $1}')
    method="chrome/local"
fi

get_latest_chromedriver() {
    chrome_version=$(echo "$1" | awk -F '-' '{print $1}')
    chrome_channel=$2
    base_url="https://chromedriver.storage.googleapis.com"
    if [ "$chrome_channel" == "dev" ]; then
        chrome_major_version=$(echo "$chrome_version" | cut -d. -f1)
        status_code=$(wget --spider -S $base_url/LATEST_RELEASE_$chrome_major_version 2>&1 | awk '/HTTP\// {print $2}')
        if [ "$status_code" == "404" ]; then
            let chrome_major_version--
        fi
        echo "$(wget -qO- $base_url/LATEST_RELEASE_$chrome_major_version)"
    else
        chrome_build_version=$(echo "$chrome_version" | cut -d. -f1-3)
        echo "$(wget -qO- $base_url/LATEST_RELEASE_$chrome_build_version)"
    fi
}

if [ "$driver_version" == "latest" ]; then
    driver_version=$(get_latest_chromedriver $browser_version $channel)
fi

./build-dev.sh $method $browser_version $channel true
if [ "$method" == "chrome/apt" ]; then
    ./build-dev.sh $method $browser_version $channel false
fi
pushd chrome
../build.sh chromedriver $browser_version $driver_version selenoid/chrome:$tag
popd

test_image(){
    docker rm -f selenium || true
    docker run -d --name selenium -p 4445:4444 $1:$2
    tests_dir=../../selenoid-container-tests/
    if [ -d "$tests_dir" ]; then
        pushd "$tests_dir"
        mvn clean test -Dgrid.connection.url="http://localhost:4445/" -Dgrid.browser.name=chrome -Dgrid.browser.version=$2 || $test_failure_ignore
        popd
    else
        echo "Skipping tests as $tests_dir does not exist."
    fi
}

test_image "selenoid/chrome" $tag
docker tag "selenoid/chrome:$tag" "selenoid/vnc_chrome:$tag"
docker tag "selenoid/chrome:$tag" "selenoid/vnc:chrome_$tag"

read -p "Push?" yn
if [ "$yn" == "y" ]; then
	docker push "selenoid/dev_chrome:"$browser_version
	if [ "$method" == "chrome/apt" ]; then
    	docker push "selenoid/dev_chrome_full:"$browser_version
    fi
	docker push "selenoid/chrome:$tag"
    docker tag "selenoid/chrome:$tag" "selenoid/chrome:latest"
    docker push "selenoid/chrome:latest"
    docker push "selenoid/vnc:chrome_"$tag
    docker push "selenoid/vnc_chrome:"$tag
fi
