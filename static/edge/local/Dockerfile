FROM browsers/base:7.4.0

ARG VERSION=noop
ARG PACKAGE=microsoft-edge-stable
ARG INSTALL_DIR=msedge

LABEL browser=$PACKAGE:$VERSION

RUN \
        apt-get update && \
        apt-get -y --no-install-recommends install libatk-bridge2.0-0 \
         libatomic1 \
         libatspi2.0-0 \
         libcairo-gobject2 \
         libcolord2 \
         libepoxy0 \
         libgbm1 \
         libgtk-3-0 \
         libgtk-3-common \
         liblcms2-2 \
         librest-0.7-0 \
         libsoup-gnome2.4-1 \
         libwayland-client0 \
         libwayland-cursor0 \
         libwayland-egl1 \
         libwayland-server0 \
         libxkbcommon0 \
         xdg-utils \
         ca-certificates \
         fonts-liberation \
         libappindicator3-1 \
         libnss3 \
         lsb-base \
         libcurl4 \
         curl && \
         curl -O http://host.docker.internal:8080/microsoft-edge.deb && \
         apt-get -y purge curl && \
         dpkg -i microsoft-edge.deb && \
         (  \
           sed -i -e 's@exec -a "$0" "$HERE/msedge"@& --no-sandbox --disable-gpu@' /opt/microsoft/$INSTALL_DIR/$PACKAGE || \
           sed -i -e 's@exec -a "$0" "$HERE/msedge"@& --no-sandbox --disable-gpu@' /opt/microsoft/$INSTALL_DIR/microsoft-edge \
         ) && \
         rm microsoft-edge.deb && \
         chown root:root /opt/microsoft/$INSTALL_DIR/msedge-sandbox && \
         chmod 4755 /opt/microsoft/$INSTALL_DIR/msedge-sandbox && \
         microsoft-edge --version && \
         rm -Rf /tmp/* && rm -Rf /var/lib/apt/lists/*
