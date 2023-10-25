ARG VERSION
FROM selenoid/dev_opera:$VERSION

COPY operadriver_linux64/operadriver /usr/bin/
COPY entrypoint.sh /

RUN \
    chmod +x /usr/bin/operadriver && \
    echo "cookie-file = ~/.config/pulse/cookie" >> /etc/pulse/client.conf

USER selenium

EXPOSE 4444
ENTRYPOINT ["/entrypoint.sh"]
