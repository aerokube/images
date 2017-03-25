FROM selenoid/base:1.0

ARG VERSION
ARG CLEANUP

RUN  \
        wget -O- http://deb.opera.com/archive.key | apt-key add - && \
        echo 'deb http://deb.opera.com/opera-stable/ stable non-free' > /etc/apt/sources.list.d/opera-blink.list && \
        apt-get update && \
        apt-get -y --force-yes install opera-stable=$VERSION && \
        ($CLEANUP && apt-get clean && rm -Rf /tmp/*) || true
