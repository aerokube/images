FROM browsers/base:7.3.6

ARG VERSION
ARG PACKAGE=microsoft-edge-stable
ARG INSTALL_DIR=msedge

LABEL browser=$PACKAGE:$VERSION

RUN \
        curl -s https://packages.microsoft.com/keys/microsoft.asc | apt-key add - && \
        echo 'deb [arch=amd64] https://packages.microsoft.com/repos/edge stable main' > /etc/apt/sources.list.d/microsoft-edge.list && \
        apt-get update && \
        apt-get -y --no-install-recommends install $PACKAGE=$VERSION && \
        (  \
          sed -i -e 's@exec -a "$0" "$HERE/msedge"@& --no-sandbox --disable-gpu@' /opt/microsoft/$INSTALL_DIR/$PACKAGE || \
          sed -i -e 's@exec -a "$0" "$HERE/msedge"@& --no-sandbox --disable-gpu@' /opt/microsoft/$INSTALL_DIR/microsoft-edge \
        ) && \
        chown root:root /opt/microsoft/$INSTALL_DIR/msedge-sandbox && \
        chmod 4755 /opt/microsoft/$INSTALL_DIR/msedge-sandbox && \
        microsoft-edge --version && \
        rm -Rf /tmp/* && rm -Rf /var/lib/apt/lists/*
