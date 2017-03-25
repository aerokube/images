FROM selenoid/base:1.0

ARG VERSION
ARG CLEANUP

RUN  \
        apt-get update && \
        apt-get -y install firefox=$VERSION && \
        ($CLEANUP && apt-get clean && rm -Rf /tmp/*) || true
