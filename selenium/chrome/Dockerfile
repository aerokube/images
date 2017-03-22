FROM selenoid/base:latest

ARG VERSION
ARG CLEANUP=false

RUN \
        wget -q -O - https://dl-ssl.google.com/linux/linux_signing_key.pub | sudo apt-key add - && \
        echo 'deb [arch=amd64] http://dl.google.com/linux/chrome/deb/ stable main' > /etc/apt/sources.list.d/google.list && \
        apt-get update && \
        apt-get -y --force-yes install google-chrome-stable=$VERSION unzip && \
        chown root:root /opt/google/chrome/chrome-sandbox && \
        chmod 4755 /opt/google/chrome/chrome-sandbox && \
        ($CLEANUP && apt-get clean && rm -Rf /tmp/*)
