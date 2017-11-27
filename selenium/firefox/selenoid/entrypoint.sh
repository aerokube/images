#!/bin/bash
SCREEN_RESOLUTION=${SCREEN_RESOLUTION:-"1920x1080x24"}
DISPLAY=99
/usr/bin/xvfb-run -l -n "$DISPLAY" -s "-ac -screen 0 $SCREEN_RESOLUTION -noreset -listen tcp" /usr/bin/selenoid -conf /etc/selenoid/browsers.json -disable-docker -timeout 1h -enable-file-upload -capture-driver-logs
