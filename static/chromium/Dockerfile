ARG VERSION
FROM selenoid/dev_chromium:$VERSION

ENV DBUS_SESSION_BUS_ADDRESS=/dev/null
COPY entrypoint.sh /

RUN chmod +x /usr/bin/chromedriver

RUN echo "cookie-file = ~/.config/pulse/cookie" >> /etc/pulse/client.conf

USER selenium

EXPOSE 4444
ENTRYPOINT ["/entrypoint.sh"]
