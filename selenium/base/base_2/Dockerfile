FROM selenoid/base:1.0

RUN \
    apt-get update && \
    apt-get install -y fonts-wqy-zenhei fonts-thai-tlwg-ttf fonts-ipafont-mincho fonts-ipafont-gothic && \
    locale-gen ru_RU.UTF-8 && update-locale && \
    apt-get clean && \
    rm -Rf /tmp/*
