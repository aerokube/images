FROM browsers/base:7.3.6

ARG VERSION=19.4.2.698-1

LABEL browser=yandex-browser-stable:$VERSION

RUN  \
        curl -s https://repo.yandex.ru/yandex-browser/YANDEX-BROWSER-KEY.GPG | apt-key add - && \
        echo 'deb [arch=amd64] https://repo.yandex.ru/yandex-browser/deb stable main' > /etc/apt/sources.list.d/yandex-browser.list && \
        apt-get update && \
        apt-get -y --no-install-recommends install yandex-browser-stable=$VERSION && \
        yandex-browser --version && \
        rm /etc/apt/sources.list.d/yandex-browser.list && \
        rm -Rf /tmp/* && rm -Rf /var/lib/apt/lists/*
