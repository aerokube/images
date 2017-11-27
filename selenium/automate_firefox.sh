#!/bin/bash
set -e
input=$1
server_version=$2
tag=$3
driver_version=$4

if [ -z "$1" -o -z "$2" -o -z "$3" -o -z "$4" ]; then
    echo 'Usage: automate_firefox.sh <browser_version|package_file> <selenium_version|selenoid_version> <tag_version> <geckodriver_version>'
    exit 1
fi
set -x

browser_version=$input
method="firefox/apt"
runner="selenoid"
requires_java="false"
numeric_version=$(echo "$tag" | awk -F '.' '{print $1}' )
if [ $numeric_version -lt 48 ]; then
    runner="selenium"
    requires_java="true"
fi

if [ -f $input ]; then
    filename=$(echo "$input" | awk -F '/' '{print $NF}')
    arch=$(echo "$filename" | awk -F '_' '{print $NF}' | sed -e 's|.deb||g')
    cp "$input" firefox/local/firefox_$arch.deb
    browser_version=$(echo $filename | awk -F '_' '{print $2}' | awk -F '-' '{print $1}')
    method="firefox/local"
fi

./build-dev.sh $method $browser_version false $requires_java $tag
./build-dev.sh $method $browser_version true $requires_java $tag
pushd firefox/$runner
../../build.sh $runner $tag $server_version selenoid/firefox:$tag $driver_version
popd

test_image(){
    docker rm -f selenium || true
    docker run -d --name selenium -p 4444:4444 $1:$2
    tests_dir=../../selenoid-container-tests/
    if [ -d "$tests_dir" ]; then
        pushd "$tests_dir"
        mvn clean test -Dgrid.browser.name=firefox -Dgrid.browser.version=$2 || true
        popd
    else
        echo "Skipping tests as $tests_dir does not exist."
    fi
}

test_image "selenoid/firefox" $tag
read -p "Create VNC?" vnc
if [ "$vnc" == "y" ]; then
    pushd vnc/firefox/$runner
    ../../../build-vnc.sh firefox $tag 
    popd
    test_image "selenoid/vnc_firefox" $tag
fi

read -p "Push?" yn
if [ "$yn" == "y" ]; then
	docker push "selenoid/dev_firefox:"$tag
	docker push "selenoid/dev_firefox_full:"$tag
	docker push "selenoid/firefox:$tag"
    docker tag "selenoid/firefox:$tag" "selenoid/firefox:latest"
    docker push "selenoid/firefox:latest"
    if [ "$vnc" == "y" ]; then
        docker push "selenoid/vnc:firefox_"$tag
        docker push "selenoid/vnc_firefox:"$tag
    fi    
fi
