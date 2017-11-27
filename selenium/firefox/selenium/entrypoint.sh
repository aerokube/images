#!/bin/bash
SCREEN_RESOLUTION=${SCREEN_RESOLUTION:-"1920x1080x24"}
DISPLAY=99
/usr/bin/xvfb-run -l -n "$DISPLAY" -s "-ac -screen 0 $SCREEN_RESOLUTION -noreset -listen tcp" /usr/bin/java -Xmx256m -Djava.security.egd=file:/dev/./urandom -jar /usr/share/selenium/selenium-server-standalone.jar -port 4444 -browserTimeout 120