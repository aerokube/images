FROM selenoid/base:1.0

RUN \
    apt-get update && \
    apt-get install -y fonts-wqy-zenhei fonts-thai-tlwg-ttf fonts-ipafont-mincho fonts-ipafont-gothic && \
    ln -s /usr/share/i18n/SUPPORTED /var/lib/locales/supported.d/all && locale-gen && update-locale && \
    apt-get clean && \
    rm -Rf /tmp/*
