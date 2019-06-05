FROM selenoid/base:5.0

ARG VERSION=19.4.2.698-1
ARG CLEANUP

RUN  \
        ( [ "$CLEANUP" != "true" ] && rm -f /etc/apt/apt.conf.d/docker-clean ) || true && \
        wget -O- https://repo.yandex.ru/yandex-browser/YANDEX-BROWSER-KEY.GPG | apt-key add - && \
        echo 'deb [arch=amd64] http://repo.yandex.ru/yandex-browser/deb beta main' > /etc/apt/sources.list.d/yandex-browser.list && \
        apt-get update && \
        apt-get -y --no-install-recommends install yandex-browser-beta=$VERSION && \
        yandex-browser-beta --version && \
        ln /usr/bin/yandex-browser-beta /usr/bin/opera && \
        rm /etc/apt/sources.list.d/yandex-browser.list && \
        ($CLEANUP && rm -Rf /tmp/* && rm -Rf /var/lib/apt/lists/*) || true
