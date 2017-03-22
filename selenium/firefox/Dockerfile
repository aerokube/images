FROM selenoid/base:latest

ARG VERSION
ARG CLEANUP=false

RUN  \
        apt-get update && \
        apt-get -y --force-yes install firefox=$VERSION && \
        ($CLEANUP && apt-get clean && rm -Rf /tmp/*)
