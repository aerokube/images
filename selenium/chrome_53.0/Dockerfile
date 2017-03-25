FROM selenoid/base:latest

RUN     wget -q -O - https://dl-ssl.google.com/linux/linux_signing_key.pub | apt-key add - && \
        echo 'deb [arch=amd64] http://dl.google.com/linux/chrome/deb/ stable main' > /etc/apt/sources.list.d/google.list && \
        apt-get update && \
        apt-get -y --force-yes install google-chrome-stable=53.0.2785.101-1 unzip && \
        chown root:root /opt/google/chrome/chrome-sandbox && \
        chmod 4755 /opt/google/chrome/chrome-sandbox && \
        wget -O chromedriver.zip http://chromedriver.storage.googleapis.com/2.24/chromedriver_linux64.zip && \
        unzip chromedriver.zip -d /usr/bin && \
        apt-get clean && \
        rm -Rf /tmp/*

COPY entrypoint.sh /
ENTRYPOINT /entrypoint.sh
