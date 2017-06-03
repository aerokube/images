#!/bin/bash
SCREEN_RESOLUTION=${SCREEN_RESOLUTION:-"1920x1080x24"}
DISPLAY=99
XAUTHORITY=/tmp/xvfb.auth x11vnc -display ":$DISPLAY" -passwd selenoid -shared -forever -loop500 -rfbport 5900 -rfbportv6 5900 &
/usr/bin/xvfb-run -f /tmp/xvfb.auth -n "$DISPLAY" -s "-screen 0 $SCREEN_RESOLUTION -noreset -auth /tmp/xvfb.auth" /usr/bin/java -Xmx256m -Djava.security.egd=file:/dev/./urandom -Dwebdriver.gecko.driver=/usr/bin/geckodriver -jar /usr/share/selenium/selenium-server-standalone.jar -port 4444 -browserTimeout 120