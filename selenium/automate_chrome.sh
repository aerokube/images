#!/bin/bash
set -e
set -x
browser_version=$1
driver_version=$2
tag=$3
./build-dev.sh chrome/apt $browser_version true true
./build-dev.sh chrome/apt $browser_version false true
pushd chrome
../build.sh chromedriver $browser_version $driver_version selenoid/chrome:$tag
popd
docker rm -f selenium || true
docker run -d --name selenium -p 4444:4444 selenoid/chrome:$tag
tests_dir=../../selenoid-container-tests/
if [ -d "$tests_dir" ]; then
    pushd "$tests_dir"
    mvn clean test -Dgrid.connection.url="http://localhost:4444/" -Dgrid.browser.version=$tag || true
    popd
else
    echo "Skipping tests as $tests_dir does not exist."
fi
read -p "Push?" yn
if [ "$yn" == "y" ]; then
	docker push "selenoid/dev:chrome_"$browser_version
	docker push "selenoid/dev:chrome_"$browser_version"_full"
	docker push "selenoid/chrome:$tag"
    docker tag "selenoid/chrome:$tag" "selenoid/chrome:latest"
    docker push "selenoid/chrome:latest"
fi
