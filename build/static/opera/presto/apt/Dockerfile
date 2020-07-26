FROM selenoid/base:6.0

ARG VERSION=12.16.1860

RUN  \
        wget -O- http://deb.opera.com/archive.key | apt-key add - && \
        echo 'deb http://deb.opera.com/opera/ stable non-free' >> /etc/apt/sources.list.d/opera.list && \
        apt-get update && \
        apt-get -y install opera=$VERSION && \
        rm -Rf /tmp/* && rm -Rf /var/lib/apt/lists/*
