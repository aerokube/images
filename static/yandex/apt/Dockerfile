FROM browsers/base:7.3.5

ARG VERSION=19.4.2.698-1

LABEL browser=yandex-browser-beta:$VERSION

RUN  \
        curl -s https://repo.yandex.ru/yandex-browser/YANDEX-BROWSER-KEY.GPG | apt-key add - && \
        echo 'deb [arch=amd64] https://repo.yandex.ru/yandex-browser/deb beta main' > /etc/apt/sources.list.d/yandex-browser.list && \
        apt-get update && \
        apt-get -y --no-install-recommends install yandex-browser-beta=$VERSION && \
        yandex-browser-beta --version && \
        rm /etc/apt/sources.list.d/yandex-browser.list && \
        rm -Rf /tmp/* && rm -Rf /var/lib/apt/lists/*
