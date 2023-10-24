ARG VERSION
FROM selenoid/dev_edge:$VERSION

ENV DBUS_SESSION_BUS_ADDRESS=/dev/null
COPY msedgedriver /usr/bin/
COPY entrypoint.sh /

RUN chmod +x /usr/bin/msedgedriver

RUN echo "cookie-file = ~/.config/pulse/cookie" >> /etc/pulse/client.conf

USER selenium

EXPOSE 4444
ENTRYPOINT ["/entrypoint.sh"]
