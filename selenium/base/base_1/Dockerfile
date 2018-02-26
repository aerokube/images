FROM ubuntu:16.04

RUN \
    apt-get update && \
    apt-get install -y apt-transport-https ca-certificates tzdata && \
    apt-get clean && \
    rm -Rf /tmp/*
    
COPY sources.list /etc/apt/

RUN \
    echo ttf-mscorefonts-installer msttcorefonts/accepted-mscorefonts-eula select true | debconf-set-selections && \
    echo 'UTC' | tee /etc/timezone && \
    dpkg-reconfigure -f noninteractive tzdata && \
    echo "gtk-cursor-blink=0" > /root/.gtkrc-2.0

RUN apt-get update && apt-get install -y ttf-mscorefonts-installer \
    ttf-dejavu-core \
    fontconfig \
    fontconfig-config \
    fonts-dejavu-core \
    fonts-liberation \
    fonts-ubuntu-font-family-console \
    libfontconfig1 \
    libfontenc1 \
    libfreetype6 \
    libxfont1 \
    libxft2 \
    xfonts-base \
    xfonts-encodings \
    xfonts-utils \
    flashplugin-installer \
    x11-utils \
    fluxbox \
    xvfb && \
    fc-cache -f -v && \
    apt-get clean && \
    rm -Rf /tmp/*
