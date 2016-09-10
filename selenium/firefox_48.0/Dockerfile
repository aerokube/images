FROM selenoid/base:latest

RUN wget -O geckodriver-v0.10.0-linux64.tar.gz https://github.com/mozilla/geckodriver/releases/download/v0.10.0/geckodriver-v0.10.0-linux64.tar.gz && \
        tar xvzf geckodriver-v0.10.0-linux64.tar.gz -C /usr/bin && \
        apt-key adv --recv-keys --keyserver keyserver.ubuntu.com C1289A29 && \
        echo 'deb http://downloads.sourceforge.net/project/ubuntuzilla/mozilla/apt all main' > /etc/apt/sources.list.d/ubuntuzilla.list && \
        apt-get update && \
        apt-get -y --force-yes install curl firefox-mozilla-build=48.0.1-0ubuntu1 && \
        apt-get clean && \
        rm -Rf /tmp/*

COPY entrypoint.sh /
ENTRYPOINT /entrypoint.sh
