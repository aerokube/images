FROM ubuntu:16.04

RUN apt-get update && \
	apt-get -y --no-install-recommends install ffmpeg && \
	apt-get clean && rm -Rf /tmp/*

COPY entrypoint.sh /

ENTRYPOINT ["/entrypoint.sh"]
