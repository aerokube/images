#!/bin/bash
SCREEN_RESOLUTION=${SCREEN_RESOLUTION:-"1920x1080x24"}
DISPLAY=99
fluxbox -display :$DISPLAY &
x11vnc -display ":$DISPLAY" -passwd selenoid -shared -forever -loop500 -rfbport 5900 -rfbportv6 5900 -logfile /var/log/x11vnc.log &
/usr/bin/xvfb-run -l -n "$DISPLAY" -s "-ac -screen 0 $SCREEN_RESOLUTION -noreset -listen tcp" /usr/bin/java -Xmx256m -Djava.security.egd=file:/dev/./urandom -jar /usr/share/selenium/selenium-server-standalone.jar -port 4444 -browserTimeout 120