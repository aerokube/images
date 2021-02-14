FROM browsers/base:7.2

RUN  \
        wget -O- http://deb.opera.com/archive.key | apt-key add - && \
        echo 'deb http://deb.opera.com/opera/ stable non-free' >> /etc/apt/sources.list.d/opera.list && \
        apt-get update && \
        apt-get -y install opera=12.16.1860 openjdk-8-jre-headless && \
        rm -Rf /tmp/* && rm -Rf /var/lib/apt/lists/*

COPY selenium-server-standalone.jar /usr/share/selenium/
COPY entrypoint.sh /

USER selenium

EXPOSE 4444
ENTRYPOINT ["/entrypoint.sh"]
