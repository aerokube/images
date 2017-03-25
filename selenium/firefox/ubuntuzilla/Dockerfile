FROM selenoid/base:1.0

ARG VERSION
ARG CLEANUP=false

RUN \
        apt-key adv --recv-keys --keyserver keyserver.ubuntu.com B7B9C16F2667CA5C && \
        echo 'deb http://downloads.sourceforge.net/project/ubuntuzilla/mozilla/apt all main' > /etc/apt/sources.list.d/ubuntuzilla.list && \
        dpkg --add-architecture i386 && \
        apt-get update && \
        ( \
            apt-get -y install firefox-mozilla-build=$VERSION-0ubuntu1 || \
            apt-get -y install firefox-mozilla-build=$VERSION-0ubuntu1:i386 \
        ) && \
        ($CLEANUP && apt-get clean && rm -Rf /tmp/*)