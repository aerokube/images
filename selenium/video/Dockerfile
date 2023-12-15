FROM alpine:3

RUN apk add -U x264-libs sdl2 libxcb libbz2 xset pulseaudio-utils ffmpeg && \
    rm -rf /var/cache/apk/*

COPY entrypoint.sh /
ENTRYPOINT ["/entrypoint.sh"]
