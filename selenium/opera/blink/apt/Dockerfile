FROM selenoid/base:5.0

ARG VERSION
ARG CLEANUP
ARG PACKAGE=opera-stable

LABEL browser=$PACKAGE:$VERSION

RUN  \
        ( [ "$CLEANUP" != "true" ] && rm -f /etc/apt/apt.conf.d/docker-clean ) || true && \
        wget -O- http://deb.opera.com/archive.key | apt-key add - && \
        echo 'deb http://deb.opera.com/opera-stable/ stable non-free' > /etc/apt/sources.list.d/opera-blink.list && \
        apt-get update && \
        apt-get -y --no-install-recommends install $PACKAGE=$VERSION && \
        ( [ "$PACKAGE" != "opera-stable" ] && ln /usr/bin/$PACKAGE /usr/bin/opera ) || true && \
        opera --version && \
        rm /etc/apt/sources.list.d/$PACKAGE.list && \
        ($CLEANUP && rm -Rf /tmp/* && rm -Rf /var/lib/apt/lists/*) || true
