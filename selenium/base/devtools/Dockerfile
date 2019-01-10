FROM golang:1.11 as go

COPY devtools /devtools

RUN \
    apt-get update && \
    apt-get install -y upx-ucl && \
    cd /devtools && \
    GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" && \
    upx /devtools/devtools

FROM selenoid/base:4.0

COPY --from=go /devtools/devtools /usr/bin/
