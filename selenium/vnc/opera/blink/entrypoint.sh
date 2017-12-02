#!/bin/bash
SCREEN_RESOLUTION=${SCREEN_RESOLUTION:-"1920x1080x24"}
DISPLAY=99
x11vnc -display ":$DISPLAY" -passwd selenoid -shared -viewonly -forever -loop500 -rfbport 5900 -rfbportv6 5900 -logfile /home/selenium/x11vnc.log &
/usr/bin/xvfb-run -l -n "$DISPLAY" -s "-ac -screen 0 $SCREEN_RESOLUTION -noreset -listen tcp" /usr/bin/operadriver --port=4444 --whitelisted-ips='' --verbose