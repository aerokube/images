FROM ubuntu:18.04 as build

ARG WEBKITGTK_VERSION="2.30.5"

RUN \
    apt-get update && \
    apt-get -y install xz-utils wget && \
    wget https://webkitgtk.org/releases/webkitgtk-$WEBKITGTK_VERSION.tar.xz && \
    tar xf webkitgtk-$WEBKITGTK_VERSION.tar.xz && \
    mkdir -p /opt/webkit && \
    cd webkitgtk-$WEBKITGTK_VERSION && \
    yes | DEBIAN_FRONTEND=noninteractive Tools/gtk/install-dependencies && \
    cmake -DPORT=GTK -DCMAKE_BUILD_TYPE=Release -DCMAKE_INSTALL_PREFIX=/opt/webkit -DUSE_WPE_RENDERER=OFF -DENABLE_MINIBROWSER=ON -DENABLE_BUBBLEWRAP_SANDBOX=OFF -DENABLE_SPELLCHECK=OFF -DENABLE_WAYLAND_TARGET=OFF -GNinja && \
    ninja && \
    ninja install

FROM golang:1.15 as go

COPY cmd/prism /prism

RUN \
    apt-get update && \
    apt-get install -y upx-ucl libx11-dev && \
    cd /prism && \
    GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" && \
    upx /prism/prism

FROM selenoid/base:7.0

COPY --from=build /opt/webkit /opt/webkit
COPY --from=go /prism/prism /usr/bin/

ENV LD_LIBRARY_PATH /opt/webkit/lib/:${LD_LIBRARY_PATH}

RUN \
    apt-get update && \
    apt-get -y install \
        libsoup2.4-1 \
        libgtk-3-0 \
        libwebp6 \
        libwebpdemux2 \
        libsecret-1-0 \
        libhyphen0 \
        libwoff1 \
        libharfbuzz-icu0 \
        libgstreamer-gl1.0-0 \
        libopenjp2-7 \
        libnotify4 \
        libxslt1.1 \
        libegl1 && \
    ldconfig && \
    apt-get clean && \
    rm -Rf /tmp/* && rm -Rf /var/lib/apt/lists/*

COPY entrypoint.sh /

USER selenium

EXPOSE 4444
ENTRYPOINT ["/entrypoint.sh"]
