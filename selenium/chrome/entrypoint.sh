#!/bin/bash
SCREEN_RESOLUTION=${SCREEN_RESOLUTION:-"1920x1080x24"}
/usr/bin/xvfb-run -a -s "-screen 0 $SCREEN_RESOLUTION -noreset" /usr/bin/chromedriver --port=4444 --whitelisted-ips='' --verbose