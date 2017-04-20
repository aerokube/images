FROM selenoid/base:1.0

ARG VERSION
ARG CLEANUP

RUN  \
        ( [ "$CLEANUP" != "true" ] && rm -f /etc/apt/apt.conf.d/docker-clean ) || true && \
        apt-get update && \
        apt-get -y --no-install-recommends install firefox=$VERSION && \
        ($CLEANUP && rm -Rf /tmp/*) || true
