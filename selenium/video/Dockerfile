FROM alpine:3.6

RUN apk add -U alpine-sdk && adduser -g "" -D -G abuild selenoid
USER selenoid
RUN git clone git://git.alpinelinux.org/aports /home/selenoid/aports
ADD APKBUILD /home/selenoid/aports/community/ffmpeg/
RUN cd /home/selenoid/aports/community/ffmpeg/ && \
    abuild checksum && \
    abuild unpack && \
    abuild deps && \
    abuild build
USER root

FROM alpine:3.6

RUN apk add -U x264-libs sdl2 libxcb libbz2 xset && \
    rm -rf /var/cache/apk/*
COPY --from=0 /home/selenoid/aports/main/ffmpeg/src/ffmpeg-3.4/ffmpeg /usr/bin/ffmpeg

COPY entrypoint.sh /
ENTRYPOINT ["/entrypoint.sh"]
