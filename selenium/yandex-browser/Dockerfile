FROM selenoid/base:latest

ARG VERSION
ARG CLEANUP=false

RUN  \
        wget -q -O - http://repo.yandex.ru/yandex-browser/YANDEX-BROWSER-KEY.GPG | apt-key add - && \
        echo 'deb http://repo.yandex.ru/yandex-browser/deb/ beta main' > /etc/apt/sources.list.d/yandex-browser.list && \
        apt-get update && \
        apt-get -y --force-yes install yandex-browser-beta=$VERSION && \
        ($CLEANUP && apt-get clean && rm -Rf /tmp/*)
