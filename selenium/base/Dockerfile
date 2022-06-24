ARG UBUNTU_VERSION=20.04

FROM golang:1.17 as go

COPY xseld /xseld

COPY fileserver /fileserver

RUN \
    if [ `uname -m` = "aarch64" ]; then ARCH="arm64"; else ARCH="amd64"; fi && \
    apt-get update && \
    apt-get install -y upx-ucl libx11-dev && \
    cd /xseld && \
    GOOS=linux GOARCH=$ARCH go build -ldflags="-s -w" && \
    upx /xseld/xseld && \
    cd /fileserver && \
    go test -race && \
    GOOS=linux GOARCH=$ARCH go build -ldflags="-s -w" && \
    upx /fileserver/fileserver

# For M1 Chromium images it's required to override a version to 18.04 as latest Ubuntu distributions don't ship updates
FROM ubuntu:$UBUNTU_VERSION

RUN \
    apt update && \
    apt remove -y libcurl4 && \
    apt install -y apt-transport-https ca-certificates tzdata locales libcurl4 curl gnupg && \
    DEBIAN_FRONTEND=noninteractive apt -y upgrade && \
    echo ttf-mscorefonts-installer msttcorefonts/accepted-mscorefonts-eula select true | debconf-set-selections && \
    echo 'UTC' | tee /etc/timezone && \
    dpkg-reconfigure -f noninteractive tzdata && \
    echo "gtk-cursor-blink=0" > /root/.gtkrc-2.0 && \
    apt update && \
    apt install -y ttf-mscorefonts-installer \
    ttf-dejavu-core \
    fontconfig \
    fontconfig-config \
    fonts-dejavu-core \
    fonts-liberation \
    fonts-ubuntu-font-family-console \
    fonts-wqy-zenhei \
    fonts-thai-tlwg-ttf \
    fonts-ipafont-mincho \
    fonts-sahadeva \
    fonts-noto-unhinted \
    fonts-noto-color-emoji \
    libfontconfig1 \
    libfontenc1 \
    libfreetype6 \
    libxfont2 \
    libxft2 \
    libnss3-tools \
    xfonts-base \
    xfonts-encodings \
    xfonts-utils \
    xvfb \
    pulseaudio \
    fluxbox \
    x11vnc \
    feh \
    wmctrl \
    libnss-wrapper \
    xsel && \
    if [ `uname -m` = "amd64" ]; then apt install -y flashplugin-installer; fi && \
    mkdir -p /var/lib/locales/supported.d/ && grep UTF-8 /usr/share/i18n/SUPPORTED > /var/lib/locales/supported.d/all && \
    locale-gen && update-locale && \
    fc-cache -f -v && \
    adduser --system --home /home/selenium --uid 4096 \
    --ingroup root --disabled-password --shell /bin/bash selenium && \
    mkdir -p /home/selenium/Downloads && \
    mkdir -p /home/selenium/.fluxbox && \
    chgrp -R 0 /home/selenium && \
    chmod -R g=u /home/selenium && \
    ln -sf /bin/true /usr/bin/xdg-open && \
    apt-get clean && \
    rm -Rf /tmp/* && rm -Rf /var/lib/apt/lists/*

COPY fluxbox /usr/share/fluxbox/styles/
COPY --chown=selenium:root fluxbox /home/selenium/.fluxbox/
COPY aerokube.png /usr/share/images/fluxbox/

COPY --from=go /fileserver/fileserver /usr/bin/
COPY --from=go /xseld/xseld /usr/bin/
