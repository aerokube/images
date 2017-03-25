FROM selenoid/base:1.0

ARG VERSION
ARG CLEANUP

RUN  \
        ( [ "$CLEANUP" != "true" ] && rm -f /etc/apt/apt.conf.d/docker-clean ) || true && \
        apt-get update && \
        apt-get -y install firefox=$VERSION && \
        ($CLEANUP && rm -Rf /tmp/*) || true
