FROM selenoid/base:4.0

ARG VERSION
ARG CLEANUP

RUN \
        ( [ "$CLEANUP" != "true" ] && rm -f /etc/apt/apt.conf.d/docker-clean ) || true && \
        wget -q -O - https://dl-ssl.google.com/linux/linux_signing_key.pub | apt-key add - && \
        echo 'deb [arch=amd64] http://dl.google.com/linux/chrome/deb/ stable main' > /etc/apt/sources.list.d/google.list && \
        apt-get update && \
        apt-get -y  --no-install-recommends install iproute2 google-chrome-stable=$VERSION && \
        chown root:root /opt/google/chrome/chrome-sandbox && \
        chmod 4755 /opt/google/chrome/chrome-sandbox && \
        ($CLEANUP && rm -Rf /tmp/* && rm -Rf /var/lib/apt/lists/*) || true
