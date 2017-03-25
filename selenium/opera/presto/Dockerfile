FROM selenoid/base:1.0

ARG VERSION=12.16.1860
ARG CLEANUP

RUN  \
        ( [ "$CLEANUP" != "true" ] && rm -f /etc/apt/apt.conf.d/docker-clean ) || true && \
        wget -O- http://deb.opera.com/archive.key | apt-key add -A && \
        echo 'deb http://deb.opera.com/opera/ stable non-free' >> /etc/apt/sources.list.d/opera.list && \
        apt-get update && \
        apt-get -y install opera=$VERSION && \
        ($CLEANUP && rm -Rf /tmp/*) || true
