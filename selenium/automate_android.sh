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
    case "$version" in
	4.4)
        build_tools="build-tools;19.1.0"
        platform="android-19"
        emulator_image="system-images;android-19;default;x86"
		break
		;;
	5.0)
        build_tools="build-tools;21.1.2"
        platform="android-21"
        emulator_image="system-images;android-21;default;x86"
		break
		;;
	5.1)
        build_tools="build-tools;22.0.1"
        platform="android-22"
        emulator_image="system-images;android-22;default;x86"
		break
		;;
	6.0)
        build_tools="build-tools;23.0.3"
        platform="android-23"
        emulator_image="system-images;android-23;default;x86"
		break
		;;
	7.0)
        build_tools="build-tools;24.0.3"
        platform="android-24"
        emulator_image="system-images;android-24;default;x86"
		break
		;;
	7.1)
        build_tools="build-tools;25.0.3"
        platform="android-25"
        emulator_image="system-images;android-25;google_apis_playstore;x86"
		break
		;;
	8.0)
        build_tools="build-tools;26.0.3"
        platform="android-26"
        emulator_image="system-images;android-26;google_apis_playstore;x86"
		break
		;;
	8.1)
        build_tools="build-tools;27.0.3"
        platform="android-27"
        emulator_image="system-images;android-27;google_apis_playstore;x86"
		break
		;;
	*)
		echo "Unsupported Android version"
		false
		;;
    esac
}

download_chromedriver() {
    pushd "android"
    wget -O chromedriver.zip http://chromedriver.storage.googleapis.com/$1/chromedriver_linux64.zip
    unzip chromedriver.zip
    rm chromedriver.zip
    popd
}

require_command "docker"
require_command "sed"
require_command "true"
require_command "false"
require_command "wget"

TMP_DIR="android/tmp"
rm -Rf ./"$TMP_DIR" || true
mkdir -p "$TMP_DIR"
cp android/entrypoint.sh "$TMP_DIR/entrypoint.sh"

appium_version=$(request_answer "Specify Appium version:" "1.8.0")

until [ "$?" -ne 0 ]; do
    android_version=$(request_answer "Specify Android version:" "6.0")
    validate_android_version "$android_version"
done
IFS=';' read -ra emulator_image_info <<< "$emulator_image"
emulator_image_type=${emulator_image_info[2]}
sed -i "s|@AVD_NAME@|$avd_name|g" "$TMP_DIR/entrypoint.sh"
sed -i "s|@PLATFORM@|$platform|g" "$TMP_DIR/entrypoint.sh"

android_device=$(request_answer "Specify device preset name if needed (e.g. \"Nexus 4\"):")
sdcard_size=$(request_answer "Specify SD card size, Mb:" 500)
userdata_size=$(request_answer "Specify userdata.img size, Mb:" 500)

image_name="android"
chrome_mobile=$(request_answer "Are you building a Chrome Mobile image (for mobile web testing):" "n")
if [ "y" == "$chrome_mobile" ]; then
    sed -i 's|@CHROME_MOBILE@|yes|g' "$TMP_DIR/entrypoint.sh"
    image_name="chrome-mobile"
else
    sed -i 's|@CHROME_MOBILE@||g' "$TMP_DIR/entrypoint.sh"
fi

chromedriver_version=$(request_answer "Specify Chromedriver version if needed (required for Chrome Mobile):")

tag=$(request_answer "Specify image tag:" "selenoid/$image_name:$android_version")
need_quickboot=$(request_answer "Add Android quick boot snapshot?" "y")

if [ -n "$chromedriver_version" ]; then
    download_chromedriver "$chromedriver_version"
fi 

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
    id=$(docker run -d --privileged "$tmp_tag")
    sleep 60
    docker exec "$id" "/usr/bin/emulator-snapshot.sh"
    sleep 30 # Wait for snapshot to save
    docker commit "$id" "$tag"
    docker rm -f "$id" || true
    docker rmi -f "$tmp_tag" || true
else
    docker tag "$tmp_tag" "$tag"
fi
docker rmi -f "$tmp_tag" || true
