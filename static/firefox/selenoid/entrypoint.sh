#!/bin/bash
SCREEN_RESOLUTION=${SCREEN_RESOLUTION:-"1920x1080x24"}
DISPLAY_NUM=99
export DISPLAY=":$DISPLAY_NUM"

QUIET=${QUIET:-""}
if [ -z "$QUIET" ]; then
    sed -i 's|@@DRIVER_ARGS@@|, "--log", "debug"|g' /home/selenium/browsers.json
else
    sed -i 's|@@DRIVER_ARGS@@||g' /home/selenium/browsers.json
fi

clean() {
  if [ -n "$FILESERVER_PID" ]; then
    kill -TERM "$FILESERVER_PID"
  fi
  if [ -n "$XSELD_PID" ]; then
    kill -TERM "$XSELD_PID"
  fi
  if [ -n "$XVFB_PID" ]; then
    kill -TERM "$XVFB_PID"
  fi
  if [ -n "$SELENOID_PID" ]; then
    kill -TERM "$SELENOID_PID"
  fi
  if [ -n "$X11VNC_PID" ]; then
    kill -TERM "$X11VNC_PID"
  fi
}

trap clean SIGINT SIGTERM

/usr/bin/fileserver &
FILESERVER_PID=$!

DISPLAY="$DISPLAY" /usr/bin/xseld &
XSELD_PID=$!

/usr/bin/xvfb-run -l -n "$DISPLAY_NUM" -s "-ac -screen 0 $SCREEN_RESOLUTION -noreset -listen tcp" /usr/bin/fluxbox -display "$DISPLAY" -log /dev/null 2>/dev/null &
XVFB_PID=$!

retcode=1
until [ $retcode -eq 0 ]; do
  DISPLAY="$DISPLAY" wmctrl -m >/dev/null 2>&1
  retcode=$?
  if [ $retcode -ne 0 ]; then
    echo Waiting X server...
    sleep 0.1
  fi
done

if [ "$ENABLE_VNC" == "true" ]; then
    x11vnc -display "$DISPLAY" -passwd selenoid -shared -forever -loop500 -rfbport 5900 -rfbportv6 5900 -logfile /dev/null &
    X11VNC_PID=$!
fi

DISPLAY="$DISPLAY" /usr/bin/selenoid -conf /home/selenium/browsers.json -disable-docker -timeout 1h -max-timeout 24h -enable-file-upload -capture-driver-logs &
SELENOID_PID=$!

wait
