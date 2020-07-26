ARG VERSION
FROM selenoid/dev_firefox:$VERSION

RUN \
    apt-get update && \
    apt-get install -y openjdk-8-jre-headless && \
    apt-get clean && \
    rm -Rf /tmp/*

COPY selenium-server-standalone.jar /usr/share/selenium/
COPY entrypoint.sh /

USER selenium

EXPOSE 4444
ENTRYPOINT ["/entrypoint.sh"]
