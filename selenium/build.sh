#!/bin/bash
set -e

download_selenium() {
    selenium_version=$1
    url=""
    case "$selenium_version" in
        "2.15.0" | "2.19.0" | "2.20.0" | "2.21.0" | "2.25.0" | "2.32.0" | "2.35.0" | "2.37.0" | "2.39.0" | "2.40.0" | "2.41.0" | "2.43.1" | "2.44.0" | "2.45.0" | "2.48.2")
            url="https://repo.jenkins-ci.org/releases/org/seleniumhq/selenium/selenium-server-standalone/$selenium_version/selenium-server-standalone-$selenium_version.jar"
            ;;
        "2.47.1")
            url="http://selenium-release.storage.googleapis.com/2.47/selenium-server-standalone-2.47.1.jar"
            ;;
        "2.53.1")
            url="http://selenium-release.storage.googleapis.com/2.53/selenium-server-standalone-2.53.1.jar"
            ;;
        "3.2.0")
            url="http://selenium-release.storage.googleapis.com/3.2/selenium-server-standalone-3.2.0.jar"
            ;;
        "3.3.1")
            url="http://selenium-release.storage.googleapis.com/3.3/selenium-server-standalone-3.3.1.jar"
            ;;
        "3.4.0")
            url="https://selenium-release.storage.googleapis.com/3.4/selenium-server-standalone-3.4.0.jar"
            ;;
        *)
            echo "Unsupported Selenium version: $selenium_version"
            exit 1
            ;;
    esac
    wget -O selenium-server-standalone.jar "$url"
}

download_geckodriver() {
    local download_url=""
    if [ "$1" == "latest" ]; then
        download_url=$(wget -qO- "https://api.github.com/repos/mozilla/geckodriver/releases/$1" | jq -r '.assets[].browser_download_url | select(contains("linux64"))')
    else
        download_url="https://github.com/mozilla/geckodriver/releases/download/v$1/geckodriver-v$1-linux64.tar.gz"
    fi
    wget -O geckodriver.tar.gz "$download_url"
    tar xvzf geckodriver.tar.gz
    rm -Rf geckodriver.tar.gz
}

download_chromedriver() {
    wget -O chromedriver.zip http://chromedriver.storage.googleapis.com/$1/chromedriver_linux64.zip
    unzip chromedriver.zip
    rm chromedriver.zip
}

download_operadriver() {
    local download_url=""
    if [ "$1" == "latest" ]; then
        download_url=$(wget -qO- "https://api.github.com/repos/operasoftware/operachromiumdriver/releases/$1" | jq -r '.assets[].browser_download_url | select(contains("linux64"))')
    else
        download_url="https://github.com/operasoftware/operachromiumdriver/releases/download/v.$1/operadriver_linux64.zip"
    fi
    wget -O operadriver.zip "$download_url"
    unzip operadriver.zip
    if [ -d operadriver_linux64 ]; then
        cp operadriver_linux64/operadriver ./operadriver
    fi
    chmod +x operadriver
    rm -Rf operadriver.zip operadriver_linux64
}

download_yandexdriver() {
    local download_url=""
    if [ "$1" == "latest" ]; then
        download_url=$(wget -qO- "https://api.github.com/repos/yandex/YandexDriver/releases" | jq -r 'first(.[].assets[].browser_download_url | select(contains("linux")))')
    else
        download_url=$(wget -qO- "https://api.github.com/repos/yandex/YandexDriver/releases/tags/v$1-stable" | jq -r '.assets[].browser_download_url | select(contains("linux"))')
        if [ -z "$download_url" ]; then
            echo "Unsupported Yandexdriver version: $1"
            exit 1
        fi
    fi
    wget -O yandexdriver.zip "$download_url"
    unzip yandexdriver.zip
    chmod +x yandexdriver
    rm yandexdriver.zip
}

download_selenoid() {
    local download_url=""
    if [ "$1" == "latest" ]; then
        download_url=$(wget -qO- "https://api.github.com/repos/aerokube/selenoid/releases/$1" | jq -r '.assets[].browser_download_url | select(contains("linux_amd64"))')
    else
        download_url="https://github.com/aerokube/selenoid/releases/download/$1/selenoid_linux_amd64"
    fi
    wget -O selenoid "$download_url"
    chmod +x selenoid
}

if [ -z "$1" -o -z "$2" -o -z "$3" -o -z "$4" ]; then
    echo 'Usage: build.sh {chromedriver|operadriver|yandexdriver|selenoid|selenium} <browser_version> <driver_or_selenium_version> <tag> [<supplementary_version>]'
    exit 1
fi
set -x

mode=$1
version=$2
tag=$4
disable_docker_cache=${DISABLE_DOCKER_CACHE:-false}
dir_name="/tmp/$(uuidgen | sed -e 's|-||g')"
mkdir -p "$dir_name"
pushd "$dir_name"
template_file="Dockerfile.driver.tmpl"
additional_docker_args=""
if [ "$mode" == "chromedriver" ]; then
    download_chromedriver "$3"
    additional_docker_args+="--label driver=chromedriver:$3"
elif [ "$mode" == "operadriver" ]; then
    download_operadriver "$3"
    additional_docker_args+="--label driver=operadriver:$3"
elif [ "$mode" == "yandexdriver" ]; then
    download_yandexdriver "$3"
    additional_docker_args+="--label driver=yandexdriver:$3"
elif [ "$mode" == "selenoid" ]; then
    download_selenoid "$3"
    additional_docker_args+="--label selenoid=$3"
    download_geckodriver "$5"
    additional_docker_args+=" --label driver=geckodriver:$5"
elif [ "$mode" == "selenium" ]; then
    download_selenium "$3"
    additional_docker_args+="--label selenium=$3"
    template_file="Dockerfile.server.tmpl"
else
    echo "Unsupported mode: will do nothing. Exiting."
    exit 1
fi
popd
cat "$template_file" | sed -e "s|@@VERSION@@|$version|g" > "$dir_name/Dockerfile"
if [ "$mode" == "chromedriver" ]; then
    cp -R devtools "$dir_name/devtools"
fi
if [ -f "browsers.json.tmpl" ]; then
    cat browsers.json.tmpl | sed -e "s|@@VERSION@@|$version|g" > "$dir_name/browsers.json"
fi
if [ -f "entrypoint.sh" ]; then
    cp entrypoint.sh "$dir_name/entrypoint.sh"
fi
if [ "$disable_docker_cache" == "true" ]; then
    additional_docker_args+=" --no-cache"
fi
pushd "$dir_name"
docker build $additional_docker_args -t "$tag" .
popd
rm -Rf "$dir_name"
exit 0
