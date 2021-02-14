FROM browsers/base:7.2

ARG VERSION
ARG PACKAGE=microsoft-edge-beta
ARG INSTALL_DIR=msedge-beta

LABEL browser=$PACKAGE:$VERSION

RUN \
        wget -q -O - https://packages.microsoft.com/keys/microsoft.asc | apt-key add - && \
        echo 'deb [arch=amd64] https://packages.microsoft.com/repos/edge stable main' > /etc/apt/sources.list.d/microsoft-edge.list && \
        apt-get update && \
        apt-get -y --no-install-recommends install $PACKAGE=$VERSION && \
        chown root:root /opt/microsoft/$INSTALL_DIR/msedge-sandbox && \
        chmod 4755 /opt/microsoft/$INSTALL_DIR/msedge-sandbox && \
        microsoft-edge --version && \
        rm -Rf /tmp/* && rm -Rf /var/lib/apt/lists/*
