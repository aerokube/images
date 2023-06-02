ARG VERSION
FROM selenoid/dev_chromium:$VERSION

ENV DBUS_SESSION_BUS_ADDRESS=/dev/null
COPY entrypoint.sh /

RUN chmod +x /usr/bin/chromedriver
USER selenium

EXPOSE 4444
ENTRYPOINT ["/entrypoint.sh"]
