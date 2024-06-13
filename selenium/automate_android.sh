#!/bin/bash

set -e

die(){
    echo "$1"
    return 1
}

require_command(){
    if [ -z "$(command -v $1)" ]; then
        die "$1 command required for this script to run"
    fi
}

request_answer(){
    prompt=$1
    default_value=$2
    if [ -n "$default_value" ]; then
        prompt="$prompt [$default_value]"
    fi
    read -r -p "$prompt " value
    if [ -z "$value" ] && [ -n "$default_value" ]; then
        value="$default_value"
    fi
    echo "$value"
}

validate_android_version(){
    version="$1"
    type=${2:-"default"}
    abi=${3:-"x86"}
    avd_name="android$version-1"
    build_tools="build-tools;29.0.2"
    replace_img="y"
    case "$version" in
	4.4)
        platform="android-19"
        emulator_image="system-images;android-19;$type;$abi"
		;;
	5.0)
        platform="android-21"
        emulator_image="system-images;android-21;$type;$abi"
		;;
	5.1)
        platform="android-22"
        emulator_image="system-images;android-22;$type;$abi"
		;;
	6.0)
        platform="android-23"
        emulator_image="system-images;android-23;$type;$abi"
		;;
	7.0)
        platform="android-24"
        emulator_image="system-images;android-24;$type;$abi"
		;;
	7.1)
        platform="android-25"
        emulator_image="system-images;android-25;$type;$abi"
		;;
	8.0)
        platform="android-26"
        emulator_image="system-images;android-26;$type;$abi"
		;;
	8.1)
        platform="android-27"
        emulator_image="system-images;android-27;$type;$abi"
		;;
	9.0)
        platform="android-28"
        emulator_image="system-images;android-28;$type;$abi"
        replace_img=""
		;;
	10.0)
        platform="android-29"
        emulator_image="system-images;android-29;$type;$abi"
        replace_img=""
    ;;
	11.0)
        platform="android-30"
        emulator_image="system-images;android-30;$type;$abi"
        replace_img=""
		;;
    12.0)
        platform="android-31"
        emulator_image="system-images;android-31;$type;$abi"
        replace_img=""
        ;;
    13.0)
        platform="android-33"
        emulator_image="system-images;android-33;$type;$abi"
        replace_img=""
        ;;
    14.0)
        platform="android-34"
        emulator_image="system-images;android-34;$type;$abi"
        replace_img=""
        ;;
	*)
		echo "Unsupported Android version"
		false
		;;
    esac
}

validate_android_image_type(){
    type="$1"
    case "$type" in
	"default" | "google_apis" | "google_apis_playstore" | "android-wear" | "android-tv" )
		;;
	*)
		echo "Unsupported Android image type"
		false
		;;
    esac
}

validate_android_abi(){
    abi="$1"
    case "$abi" in
	"armeabi-v7a" | "arm64-v8a" | "x86" | "x86_64" )
		;;
	*)
		echo "Unsupported Application Binary Interface"
		false
		;;
    esac
}

download_chromedriver() {
    pushd "$TMP_DIR"
    wget -O chromedriver.zip https://chromedriver.storage.googleapis.com/$1/chromedriver_linux64.zip
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
cp -r ../static/chrome/devtools "$TMP_DIR/devtools"

appium_version=$(request_answer "Specify Appium version:" "2.0.1")

until [ "$?" -ne 0 ]; do
    android_image_type=$(request_answer "Specify Android image type (possible values: \"default\", \"google_apis\", \"google_apis_playstore\", \"android-tv\", \"android-wear\"):" "default")
    if validate_android_image_type "$android_image_type"; then
        break
    fi
done
until [ "$?" -ne 0 ]; do
    android_abi=$(request_answer "Specify Application Binary Interface (possible values: \"armeabi-v7a\", \"arm64-v8a\", \"x86\", \"x86_64\"):" "x86")
    if validate_android_abi "$android_abi"; then
        break
    fi
done
until [ "$?" -ne 0 ]; do
    android_version=$(request_answer "Specify Android version:" "8.1")
    if validate_android_version "$android_version" "$android_image_type" "$android_abi"; then
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
else
    sed -i.bak 's|@CHROME_MOBILE@||g' "$TMP_DIR/entrypoint.sh"
fi

chromedriver_version=$(request_answer "Specify Chromedriver version if needed (required for Chrome Mobile):")
if [ -n "$chromedriver_version" ]; then
    chrome_major_version="$(cut -d'.' -f1 <<<${chromedriver_version})"
    chrome_minor_version="$(cut -d'.' -f2 <<<${chromedriver_version})"
    if [ -n "$chrome_major_version" ] &&  [ -n "$chrome_minor_version" ]; then
        default_tag="$chrome_major_version.$chrome_minor_version"
    fi
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
    --build-arg REPLACE_IMG="$replace_img" \
    --build-arg AVD_NAME="$avd_name" \
    --build-arg BUILD_TOOLS="$build_tools" \
    --build-arg PLATFORM="$platform" \
    --build-arg EMULATOR_IMAGE="$emulator_image" \
    --build-arg EMULATOR_IMAGE_TYPE="$emulator_image_type" \
    --build-arg ANDROID_ABI="$android_abi" \
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

read -r -p "Push?" yn
if [ "$yn" == "y" ]; then
    docker push "$tag"
fi
