ARG VERSION
FROM selenoid/dev_firefox:$VERSION

COPY geckodriver /usr/bin/
COPY selenoid /usr/bin/
COPY --chown=selenium:root browsers.json /home/selenium/
COPY entrypoint.sh /

USER selenium

RUN echo "cookie-file = ~/.config/pulse/cookie" >> /etc/pulse/client.conf

EXPOSE 4444
ENTRYPOINT ["/entrypoint.sh"]
