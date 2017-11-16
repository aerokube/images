FROM selenoid/base:1.0

RUN \
    apt-get update && \
    apt-get install -y fonts-wqy-zenhei fonts-thai-tlwg-ttf fonts-ipafont-mincho fonts-sahadeva && \
    mkdir -p /var/lib/locales/supported.d/ && grep UTF-8 /usr/share/i18n/SUPPORTED > /var/lib/locales/supported.d/all && \
    locale-gen && update-locale && \
    apt-get clean && \
    rm -Rf /tmp/*
