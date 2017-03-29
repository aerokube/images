#!/bin/bash
SCREEN_RESOLUTION=${SCREEN_RESOLUTION:"1920x1080x24"}
/usr/bin/xvfb-run -a -s "-screen 0 $SCREEN_RESOLUTION -noreset" /usr/bin/java -Xmx256m -Djava.security.egd=file:/dev/./urandom -jar /usr/share/selenium/selenium-server-standalone.jar -port 4444 -timeout 60 -browserTimeout 120 -newSessionWaitTimeout 10000