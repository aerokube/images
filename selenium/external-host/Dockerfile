FROM ubuntu:18.04
COPY entrypoint.sh /
COPY choose /usr/local/bin/
COPY rproxy /usr/local/bin/
RUN apt update && \
  apt -y upgrade && \
  apt -y install xvfb && \
  apt -y install x11-utils && \
  apt -y install x11vnc && \
  apt -y install tigervnc-common && \
  apt -y install tigervnc-viewer && \
  chmod +x /entrypoint.sh /usr/local/bin/choose
ENTRYPOINT ["/entrypoint.sh"]

