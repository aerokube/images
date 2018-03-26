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
    apt-get update && \
    apt-get install -y pulseaudio fluxbox x11-utils xloadimage && \
    apt-get clean && \
    rm -Rf /tmp/* && rm -Rf /var/lib/apt/lists/*
    
COPY wallpaper.png /usr/share/images/fluxbox/
COPY --from=0 /fileserver/fileserver /usr/bin/
