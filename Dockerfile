FROM ubuntu:20.04 as build

ARG WEBKIT_VERSION="610.4.3.1.7"

RUN \
    apt-get update && \
    apt-get -y install --no-install-recommends ca-certificates subversion && \
    svn checkout https://svn.webkit.org/repository/webkit/tags/Safari-$WEBKIT_VERSION webkit && \
    mkdir -p /opt/webkit && \
    cd webkit && \
    yes | DEBIAN_FRONTEND=noninteractive Tools/gtk/install-dependencies && \
    cmake -DPORT=GTK -DCMAKE_BUILD_TYPE=Release -DCMAKE_INSTALL_PREFIX=/opt/webkit -DUSE_WPE_RENDERER=OFF -DENABLE_MINIBROWSER=ON -DENABLE_BUBBLEWRAP_SANDBOX=OFF -DENABLE_GAMEPAD=OFF -DENABLE_SPELLCHECK=OFF -DENABLE_WAYLAND_TARGET=OFF -DUSE_OPENJPEG=OFF -GNinja && \
    ninja && \
    ninja install && \
    rm -Rf /var/lib/apt/lists/*

FROM golang:1.16 as go

COPY cmd/prism /prism

RUN \
    apt-get update && \
    apt-get install --no-install-recommends -y upx-ucl libx11-dev && \
    cd /prism && \
    GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" && \
    upx /prism/prism && \
    rm -Rf /var/lib/apt/lists/*

FROM browsers/base:7.2

COPY --from=build /opt/webkit /opt/webkit
COPY --from=go /prism/prism /usr/bin/

ENV LD_LIBRARY_PATH /opt/webkit/lib/:${LD_LIBRARY_PATH}

RUN \
    apt-get update && \
    apt-get -y install --no-install-recommends \
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
