FROM ubuntu:18.04 as build

ARG WEBKITGTK_VERSION="2.28.2"
ARG PREFIX="/opt/webkit"

RUN \
    apt-get update && \
    apt-get -y install xz-utils wget && \
    wget https://webkitgtk.org/releases/webkitgtk-$WEBKITGTK_VERSION.tar.xz && \
    tar xf webkitgtk-$WEBKITGTK_VERSION.tar.xz && \
    mkdir -p $PREFIX && \
    cd webkitgtk-$WEBKITGTK_VERSION && \
    yes | DEBIAN_FRONTEND=noninteractive Tools/gtk/install-dependencies && \
    cmake -DPORT=GTK -DCMAKE_BUILD_TYPE=Release -DCMAKE_INSTALL_PREFIX=$PREFIX -DUSE_WPE_RENDERER=OFF -DENABLE_MINIBROWSER=ON -DENABLE_BUBBLEWRAP_SANDBOX=OFF -DENABLE_SPELLCHECK=OFF -DENABLE_WAYLAND_TARGET=OFF -GNinja && \
    ninja && \
    ninja install

FROM selenoid/base:6.0

COPY --from=build $PREFIX $PREFIX

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
    apt-get clean && \
    rm -Rf /tmp/* && rm -Rf /var/lib/apt/lists/*

COPY entrypoint.sh /

USER selenium

EXPOSE 4444
ENTRYPOINT ["/entrypoint.sh"]
