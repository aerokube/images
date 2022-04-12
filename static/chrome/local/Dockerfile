FROM browsers/base:7.3.6

ARG VERSION=noop
ARG PACKAGE=google-chrome-stable
ARG INSTALL_DIR=chrome

LABEL browser=$PACKAGE:$VERSION

RUN \
        apt-get update && \
        apt-get -y --no-install-recommends install gconf-service \
         libasound2 \
         libatk1.0-0 \
         libc6 \
         libcairo2 \
         libcups2 \
         libdbus-1-3 \
         libexpat1 \
         libfontconfig1 \
         libfreetype6 \
         libgcc1 \
         libgconf-2-4 \
         libgdk-pixbuf2.0-0 \
         libglib2.0-0 \
         libgtk2.0-0 \
         libgtk-3-0 \
         libnspr4 \
         libnss3 \
         libpango1.0-0 \
         libstdc++6 \
         libx11-6 \
         libx11-xcb1 \
         libxcb1 \
         libxcomposite1 \
         libxcursor1 \
         libxdamage1 \
         libxext6 \
         libxfixes3 \
         libxi6 \
         libxrandr2 \
         libxrender1 \
         libxss1 \
         libxtst6 \
         ca-certificates \
         fonts-liberation \
         libappindicator3-1 \
         libnss3 \
         lsb-base \
         xdg-utils \
         libcurl4 \
         iproute2 \
         curl && \
         curl -O http://host.docker.internal:8080/google-chrome.deb && \
         apt-get -y purge curl && \
         dpkg -i google-chrome.deb && \
         sed -i -e 's@exec -a "$0" "$HERE/chrome"@& --no-sandbox --disable-gpu@' /opt/google/$INSTALL_DIR/$PACKAGE && \
         rm google-chrome.deb && \
         chown root:root /opt/google/$INSTALL_DIR/chrome-sandbox && \
         chmod 4755 /opt/google/$INSTALL_DIR/chrome-sandbox && \
         google-chrome --version && \
         rm -Rf /tmp/* && \
         rm -Rf /var/lib/apt/lists/*
