ARG VERSION
FROM selenoid/dev_yandex:$VERSION

COPY chromedriver /usr/bin/
COPY entrypoint.sh /

RUN chmod +x /usr/bin/chromedriver
USER selenium

EXPOSE 4444
ENTRYPOINT ["/entrypoint.sh"]
