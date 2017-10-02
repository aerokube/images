FROM ubuntu:14.04

ARG DEVICE=""
ARG SDCARD_SIZE=500
ARG ANDROID_SDK_URL="https://dl.google.com/android/repository/sdk-tools-linux-3859397.zip"
ARG APPIUM_VERSION="1.7.1"

RUN \
	dpkg --add-architecture i386 && \
	apt-key adv --keyserver keyserver.ubuntu.com --recv-keys 7FCD11186050CD1A && \
	apt-get update && \
	apt-get install -y curl && \
	curl --silent --location https://deb.nodesource.com/setup_6.x | bash - && \
	apt-get install -y --no-install-recommends \
	    ant \
	    autossh \
	    curl \
	    git \
	    gettext-base \
	    libncurses5:i386 \
	    libstdc++6:i386 \
	    lib32stdc++6 \
	    lib32z1 \
	    libgconf2-4 \
	    maven \
	    nodejs \
	    java8-runtime \
	    python \
	    qemu-kvm \
	    qemu-system-x86 \
	    screen \
	    unzip \
	    zlib1g:i386 \
	    xserver-xorg \
	    xserver-xorg-video-dummy \
	    x11-xserver-utils \
	    x11vnc \
	    libgl1-mesa-dev \
	    libx11-6 \
	    libx11-dev \
	    libx11-6:i386 \
	    libx11-dev:i386 && \
	cd / && npm install --prefix ./opt/ appium@"$APPIUM_VERSION" wd && \
	apt-get clean && \
	rm -rf /var/lib/apt/lists/* && \
	rm -rf /tmp/*

# Static files
COPY xorg.conf /etc/X11/
COPY android.conf /etc/ld.so.conf.d/

# You may want to optionally copy some APK, e.g. Chrome Mobile version
#COPY chrome.apk /

# Android SDK
ENV ANDROID_HOME /opt/android-sdk-linux
ENV PATH /opt/android-sdk-linux/platform-tools:/opt/android-sdk-linux/tools:/opt/android-sdk-linux/tools/bin:/opt/android-sdk-linux/emulator:$PATH
ENV LD_LIBRARY_PATH ${ANDROID_HOME}/emulator/lib64/:${ANDROID_HOME}/emulator/lib64/gles_mesa:${ANDROID_HOME}/emulator/lib64/qt/lib/:${LD_LIBRARY_PATH}
RUN \
	curl -o sdk-tools.zip "$ANDROID_SDK_URL"  && \
	mkdir -p /opt/android-sdk-linux && \
	unzip sdk-tools.zip -d /opt/android-sdk-linux && \
	rm sdk-tools.zip && \
	yes | sdkmanager --licenses && \
	sdkmanager "tools" "platform-tools" "build-tools;23.0.1" "platforms;android-23" "system-images;android-23;default;x86" && \
	mksdcard "$SDCARD_SIZE"M sdcard.img && \
	echo "no" | ( \ 
	    ([ -n "$DEVICE" ] && avdmanager create avd -n android6.0-1 -k "system-images;android-23;default;x86" --abi x86 --device "$DEVICE" --sdcard /sdcard.img ) || \
	    avdmanager create avd -n android6.0-1 -k "system-images;android-23;default;x86" --abi x86 --sdcard /sdcard.img \
    ) && \
	ldconfig && \
	resize2fs /root/.android/avd/android6.0-1.avd/userdata.img 350M && \
	mv /root/.android/avd/android6.0-1.avd/userdata.img /root/.android/avd/android6.0-1.avd/userdata-qemu.img && \
	rm /opt/android-sdk-linux/system-images/android-23/default/x86/userdata.img

# Entrypoint
COPY entrypoint.sh /
ENTRYPOINT /entrypoint.sh
