#!/bin/bash
SCREEN_RESOLUTION=${SCREEN_RESOLUTION:-"1920x1080x24"}
DISPLAY=99
x11vnc -display ":$DISPLAY" -passwd selenoid -shared -viewonly -forever -loop500 -rfbport 5900 -rfbportv6 5900 -logfile /var/log/x11vnc.log &
/usr/bin/xvfb-run -l -n "$DISPLAY" -s "-ac -screen 0 $SCREEN_RESOLUTION -noreset -listen tcp" /usr/bin/selenoid -conf /etc/selenoid/browsers.json -disable-docker -timeout 1h -enable-file-upload -capture-driver-logs
