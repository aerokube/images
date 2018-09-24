FROM golang:1.11

COPY xseld /xseld

RUN \
    apt-get update && \
    apt-get install -y upx-ucl && \
    cd /xseld && \
    go mod download && \
    GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" && \
    upx /xseld/xseld

#FROM gcc:5
#
#COPY wm-detector /wm-detector
#
#RUN \
#    apt-get update && \
#    apt-get install -y upx-ucl && \
#    cd /wm-detector && \
#    cc wm-detector.c -o wm-detector && \
#    upx /wm-detector/wm-detector

FROM selenoid/base:3.0

RUN \
    apt-get update && \
    apt-get install -y xsel && \
    apt-get clean && \
    rm -Rf /tmp/* && rm -Rf /var/lib/apt/lists/*

COPY --from=0 /xseld/xseld /usr/bin/
#COPY --from=1 /wm-detector/wm-detector /usr/bin/
