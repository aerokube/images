#!/bin/bash

set -e

die(){
    echo $1
    return 1
}

require_command(){
    cmd_name=$1
    if [ -z $(command -v $1) ]; then
        die "$1 command required for this script to run"
    fi
}

request_answer(){
    prompt=$1
    default_value=$2
    if [ -n "$default_value" ]; then
        prompt="$prompt [$default_value]"
    fi
    read -p "$prompt " value
    if [ -z "$value" -a -n "$default_value" ]; then
        value="$default_value"
    fi
    echo "$value"
}

validate_android_version(){
    version="$1"
    avd_name="android$version-1"
    build_tools="build-tools;28.0.3"
    case "$version" in
	4.4)
        platform="android-19"
        emulator_image="system-images;android-19;default;x86"
		;;
	5.0)
        platform="android-21"
        emulator_image="system-images;android-21;default;x86"
		;;
	5.1)
        platform="android-22"
        emulator_image="system-images;android-22;default;x86"
		;;
	6.0)
        platform="android-23"
        emulator_image="system-images;android-23;default;x86"
		;;
	7.0)
        platform="android-24"
        emulator_image="system-images;android-24;default;x86"
		;;
	7.1)
        platform="android-25"
        emulator_image="system-images;android-25;google_apis;x86"
		;;
	8.0)
        platform="android-26"
        emulator_image="system-images;android-26;google_apis;x86"
		;;
	8.1)
        platform="android-27"
        emulator_image="system-images;android-27;google_apis;x86"
		;;
	9.0)
        platform="android-28"
        emulator_image="system-images;android-28;google_apis;x86"
		;;
	*)
		echo "Unsupported Android version"
		false
		;;
    esac
}

download_chromedriver() {
    pushd "$TMP_DIR"
    wget -O chromedriver.zip http://chromedriver.storage.googleapis.com/$1/chromedriver_linux64.zip
    unzip chromedriver.zip
    rm chromedriver.zip
    popd
}

test_image(){
    tests_dir=../../selenoid-container-tests/
    if [ -d "$tests_dir" ]; then
        echo "Running test suite on image."
        docker rm -f selenium || true
        docker run -d --privileged --name selenium -p 4445:4444 $1
        echo "Waiting for image to start..."
        sleep 20
        pushd "$tests_dir"
        mvn clean test -Dgrid.connection.url="http://localhost:4445/wd/hub" -Dgrid.browser.name=chrome || true
        popd
        docker rm -f selenium || true
    else
        echo "Skipping tests as $tests_dir does not exist."
    fi
}

require_command "docker"
require_command "sed"
require_command "true"
require_command "false"
require_command "wget"
require_command "unzip"
require_command "cut"

TMP_DIR="android/tmp"
rm -Rf ./"$TMP_DIR" || true
mkdir -p "$TMP_DIR"
cp android/entrypoint.sh "$TMP_DIR/entrypoint.sh"

appium_version=$(request_answer "Specify Appium version:" "1.13.0")

until [ "$?" -ne 0 ]; do
    android_version=$(request_answer "Specify Android version:" "8.1")
    if validate_android_version "$android_version"; then
        break
    fi
done
IFS=';' read -ra emulator_image_info <<< "$emulator_image"
emulator_image_type=${emulator_image_info[2]}
sed -i.bak "s|@AVD_NAME@|$avd_name|g" "$TMP_DIR/entrypoint.sh"
sed -i.bak "s|@PLATFORM@|$platform|g" "$TMP_DIR/entrypoint.sh"

android_device=$(request_answer "Specify device preset name if needed (e.g. \"Nexus 4\"):")
sdcard_size=$(request_answer "Specify SD card size, Mb:" 500)
userdata_size=$(request_answer "Specify userdata.img size, Mb:" 500)

image_name="android"
default_tag="$android_version"
chrome_mobile=$(request_answer "Are you building a Chrome Mobile image (for mobile web testing):" "n")
if [ "y" == "$chrome_mobile" ]; then
    sed -i.bak 's|@CHROME_MOBILE@|yes|g' "$TMP_DIR/entrypoint.sh"
    image_name="chrome-mobile"
    default_tag="$android_version"
else
    sed -i.bak 's|@CHROME_MOBILE@||g' "$TMP_DIR/entrypoint.sh"
fi

chromedriver_version=$(request_answer "Specify Chromedriver version if needed (required for Chrome Mobile):")
chrome_major_version="$(cut -d'.' -f1 <<<${chromedriver_version})"
chrome_minor_version="$(cut -d'.' -f2 <<<${chromedriver_version})"
chrome_version="$chrome_major_version.$chrome_minor_version"
if [ -n ${chrome_version} ]; then
    default_tag="$chrome_version"
fi

tag=$(request_answer "Specify image tag:" "selenoid/$image_name:$default_tag")
need_quickboot=$(request_answer "Add Android quick boot snapshot?" "y")

if [ -n "$chromedriver_version" ]; then
    download_chromedriver "$chromedriver_version"
fi

rm -Rf *.bak || true
set -x

tmp_tag="$tag"_tmp
docker build -t "$tmp_tag" \
    --build-arg APPIUM_VERSION="$appium_version" \
    --build-arg ANDROID_DEVICE="$android_device" \
    --build-arg AVD_NAME="$avd_name" \
    --build-arg BUILD_TOOLS="$build_tools" \
    --build-arg PLATFORM="$platform" \
    --build-arg EMULATOR_IMAGE="$emulator_image" \
    --build-arg EMULATOR_IMAGE_TYPE="$emulator_image_type" \
    --build-arg SDCARD_SIZE="$sdcard_size" \
    --build-arg USERDATA_SIZE="$userdata_size" android

if [ "$need_quickboot" == "y" ]; then
    id=$(docker run -e CHROME_MOBILE="$chrome_mobile" -d --privileged "$tmp_tag")
    sleep 60
    docker exec "$id" "/usr/bin/emulator-snapshot.sh"
    sleep 30 # Wait for snapshot to save
    docker commit "$id" "$tag"
    docker rm -f "$id" || true
else
    docker tag "$tmp_tag" "$tag"
fi
docker rmi -f "$tmp_tag" || true
set +x

if [ "y" == "$chrome_mobile" ]; then
    test_image "$tag"
fi

read -p "Push?" yn
if [ "$yn" == "y" ]; then
    docker push "$tag"
fi
