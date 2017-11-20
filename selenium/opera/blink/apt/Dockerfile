FROM selenoid/base:2.0

ARG VERSION
ARG CLEANUP

RUN  \
        ( [ "$CLEANUP" != "true" ] && rm -f /etc/apt/apt.conf.d/docker-clean ) || true && \
        wget -O- http://deb.opera.com/archive.key | apt-key add - && \
        echo 'deb [arch=amd64] http://deb.opera.com/opera-stable/ stable non-free' > /etc/apt/sources.list.d/opera-blink.list && \
        apt-get update && \
        apt-get -y --no-install-recommends install opera-stable=$VERSION && \
        opera --version && \
        ($CLEANUP && rm -Rf /tmp/*) || true
