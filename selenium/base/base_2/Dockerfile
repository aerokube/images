FROM ubuntu:16.04

RUN \
    apt-get update && \
    apt-get install -y apt-transport-https ca-certificates tzdata locales && \
    apt-get clean && \
    rm -Rf /tmp/* && rm -Rf /var/lib/apt/lists/*
    
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
    fonts-wqy-zenhei \
    fonts-thai-tlwg-ttf \
    fonts-ipafont-mincho \
    fonts-sahadeva \
    libfontconfig1 \
    libfontenc1 \
    libfreetype6 \
    libxfont1 \
    libxft2 \
    xfonts-base \
    xfonts-encodings \
    xfonts-utils \
    flashplugin-installer \
    xvfb && \
    mkdir -p /var/lib/locales/supported.d/ && grep UTF-8 /usr/share/i18n/SUPPORTED > /var/lib/locales/supported.d/all && \
    locale-gen && update-locale && \
    fc-cache -f -v && \
    apt-get clean && \
    rm -Rf /tmp/* && rm -Rf /var/lib/apt/lists/*
