FROM selenoid/base:1.0

ARG VERSION=12.16.1860
ARG CLEANUP=false

RUN  \
        wget -O- http://deb.opera.com/archive.key | apt-key add -A && \
        echo 'deb http://deb.opera.com/opera/ stable non-free' >> /etc/apt/sources.list.d/opera.list && \
        apt-get update && \
        apt-get -y --force-yes install opera=$VERSION && \
        ($CLEANUP && apt-get clean && rm -Rf /tmp/*)
