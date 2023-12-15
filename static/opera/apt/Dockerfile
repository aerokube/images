FROM browsers/base:7.4.1

ARG VERSION
ARG PACKAGE=opera-stable

LABEL browser=$PACKAGE:$VERSION

RUN  \
        curl -s https://deb.opera.com/archive.key | apt-key add - && \
        echo 'deb https://deb.opera.com/opera-stable/ stable non-free' > /etc/apt/sources.list.d/opera-blink.list && \
        apt-get update && \
        apt-get -y --no-install-recommends install libgtk-3-0 $PACKAGE=$VERSION && \
        ( [ "$PACKAGE" != "opera-stable" ] && ln /usr/bin/$PACKAGE /usr/bin/opera ) || true && \
        opera --version && \
        rm /etc/apt/sources.list.d/$PACKAGE.list && \
        rm -Rf /tmp/* && rm -Rf /var/lib/apt/lists/*
