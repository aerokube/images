#!/bin/bash
SCREEN_RESOLUTION=${SCREEN_RESOLUTION:-"1920x1080x24"}
DISPLAY=99
/usr/bin/xvfb-run -l -n "$DISPLAY" -s "-ac -screen 0 $SCREEN_RESOLUTION -noreset -listen tcp" /usr/bin/operadriver --port=4444 --whitelisted-ips='' --verbose &
XVFB_PID=$!

retcode=1
until [ $retcode -eq 0 ]; do
  xdpyinfo -display :$DISPLAY >/dev/null 2>&1
  retcode=$?
  if [ $retcode -ne 0 ]; then
    echo Waiting xvfb...
    sleep 0.1
  fi
done

fluxbox -display :$DISPLAY &
wait $XVFB_PID
