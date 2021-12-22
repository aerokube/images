FROM golang:1.17 as go

COPY tmp/devtools /devtools

RUN \
    apt-get update && \
    apt-get install -y upx-ucl libx11-dev && \
    cd /devtools && \
    GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" && \
    upx /devtools/devtools

FROM ubuntu:18.04

ARG APPIUM_VERSION="1.8.1"

RUN \
    apt update && \
    apt remove -y libcurl4 && \
    apt install -y apt-transport-https ca-certificates tzdata locales libcurl4 curl gnupg && \
	curl --silent --location https://deb.nodesource.com/setup_12.x | bash - && \
	apt install -y --no-install-recommends \
	    curl \
	    iproute2 \
	    nodejs \
	    openjdk-8-jdk-headless \
	    unzip \
	    xvfb \
	    libpulse0 \
		libxcomposite1 \
		libxcursor1 \
		libxi6 \
		libasound2 \
        fluxbox \
        x11vnc \
        feh \
        wmctrl \
	    libglib2.0-0  \
        lib32stdc++6 && \
    apt-get clean && \
    rm -Rf /tmp/* && rm -Rf /var/lib/apt/lists/*

RUN cd / && npm install --prefix ./opt/ appium@$APPIUM_VERSION

COPY android.conf /etc/ld.so.conf.d/
COPY fluxbox/aerokube /usr/share/fluxbox/styles/
COPY fluxbox/init /root/.fluxbox/
COPY fluxbox/aerokube.png /usr/share/images/fluxbox/
COPY --from=go /devtools/devtools /usr/bin/

# Android SDK
ENV ANDROID_SDK_ROOT /opt
ENV PATH $ANDROID_SDK_ROOT/cmdline-tools:$ANDROID_SDK_ROOT/cmdline-tools/tools/bin:$ANDROID_SDK_ROOT/tools:$ANDROID_SDK_ROOT/platform-tools:$ANDROID_SDK_ROOT/emulator:$PATH
ENV LD_LIBRARY_PATH ${ANDROID_SDK_ROOT}/emulator/lib64:${ANDROID_SDK_ROOT}/emulator/lib64/gles_swiftshader:${ANDROID_SDK_ROOT}/emulator/lib64/qt/lib:${ANDROID_SDK_ROOT}/emulator/lib64/vulkan:${LD_LIBRARY_PATH}
ENV JAVA_HOME /usr/lib/jvm/java-8-openjdk-amd64

ARG ANDROID_DEVICE=""
ARG REPLACE_IMG="y"
ARG AVD_NAME="android6.0-1"
ARG BUILD_TOOLS="build-tools;23.0.1"
ARG PLATFORM="android-23"
ARG EMULATOR_IMAGE="system-images;android-23;default;x86"
ARG EMULATOR_IMAGE_TYPE="default"
ARG ANDROID_ABI="x86"
ARG SDCARD_SIZE="500"
ARG USERDATA_SIZE="500"

RUN \
	curl -o sdk-tools.zip https://dl.google.com/android/repository/commandlinetools-linux-7583922_latest.zip && \
	mkdir -p /opt/cmdline-tools /root/.config/Android\ Open\ Source\ Project/ && \
	unzip -q sdk-tools.zip -d /opt/cmdline-tools && \
    mv /opt/cmdline-tools/cmdline-tools /opt/cmdline-tools/tools &&\
	rm sdk-tools.zip && \
	yes | sdkmanager --licenses

RUN \
	sdkmanager "tools" "platform-tools" "$BUILD_TOOLS" "platforms;$PLATFORM" "$EMULATOR_IMAGE" && \
	mksdcard "$SDCARD_SIZE"M sdcard.img && \
	echo "no" | ( \
	    ([ -n "$ANDROID_DEVICE" ] && avdmanager create avd -n "$AVD_NAME" -k "$EMULATOR_IMAGE" --abi "$ANDROID_ABI" --device "$ANDROID_DEVICE" --sdcard sdcard.img ) || \
	    avdmanager create avd -n "$AVD_NAME" -k "$EMULATOR_IMAGE" --abi "$ANDROID_ABI" --sdcard sdcard.img \
    ) && \
	ldconfig && \
	( \
	    resize2fs /root/.android/avd/$AVD_NAME.avd/userdata.img "$USERDATA_SIZE"M || \
	    /opt/emulator/qemu-img resize -f raw /root/.android/avd/$AVD_NAME.avd/userdata.img "$USERDATA_SIZE"M \
    ) && \
	([ -z "$REPLACE_IMG" ] || mv /root/.android/avd/$AVD_NAME.avd/userdata.img /root/.android/avd/$AVD_NAME.avd/userdata-qemu.img ) && \
	rm /opt/system-images/$PLATFORM/$EMULATOR_IMAGE_TYPE/"$ANDROID_ABI"/userdata.img

COPY ["Emulator.conf", "/root/.config/Android Open Source Project/Emulator.conf"]
COPY emulator-snapshot.sh tmp/chromedriver* *.apk /usr/bin/

# Entrypoint
COPY tmp/entrypoint.sh /
ENTRYPOINT ["/entrypoint.sh"]
