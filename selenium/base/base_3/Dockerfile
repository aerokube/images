FROM golang:1.10

COPY fileserver /fileserver

RUN \
    apt-get update && \
    apt-get install -y upx-ucl && \
    cd /fileserver && \
    GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" && \
    upx /fileserver/fileserver

FROM selenoid/base:2.0

RUN \
    adduser --system --home /home/selenium \
    --ingroup nogroup --disabled-password --shell /bin/bash selenium && \
    mkdir -p /home/selenium/Downloads && \
    chown selenium:nogroup /home/selenium/Downloads && \
    mkdir -p /home/selenium/.fluxbox && \
    chown selenium:nogroup /home/selenium/.fluxbox && \
    ln -s /home/selenium/Downloads /static && \
    apt-get update && \
    apt-get install -y pulseaudio fluxbox x11-utils feh x11vnc && \
    apt-get clean && \
    rm -Rf /tmp/* && rm -Rf /var/lib/apt/lists/*

COPY fluxbox/aerokube /usr/share/fluxbox/styles/
COPY fluxbox/init /home/selenium/.fluxbox/
COPY aerokube.png /usr/share/images/fluxbox/
COPY --from=0 /fileserver/fileserver /usr/bin/
