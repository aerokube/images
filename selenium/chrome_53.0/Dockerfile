FROM selenoid/base:latest

RUN     apt-get update && \
        apt-get -y --force-yes install google-chrome-stable=53.0.2785.89-1 && \
        chown root:root /opt/google/chrome/chrome-sandbox && \
        chmod 4755 /opt/google/chrome/chrome-sandbox && \
        wget -O chromedriver.zip http://chromedriver.storage.googleapis.com/2.24/chromedriver_linux64.zip && \
        unzip chromedriver.zip -d /usr/bin && \
        apt-get clean && \
        rm -Rf /tmp/*

COPY entrypoint.sh /
ENTRYPOINT /entrypoint.sh
