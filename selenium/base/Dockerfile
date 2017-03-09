FROM ubuntu:14.04.1

RUN apt-get update && apt-get install -y --force-yes  apt-transport-https
COPY sources.list /etc/apt/
RUN echo ttf-mscorefonts-installer msttcorefonts/accepted-mscorefonts-eula select true | debconf-set-selections

RUN apt-get update && apt-get install -y --force-yes ttf-mscorefonts-installer \
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
    xvfb

RUN echo 'Europe/Moscow' | sudo tee /etc/timezone && dpkg-reconfigure -f noninteractive tzdata
RUN echo "gtk-cursor-blink=0" > /root/.gtkrc-2.0
