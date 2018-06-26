#!/bin/bash
SCREEN_RESOLUTION=${SCREEN_RESOLUTION:-"1920x1080x24"}
ENABLE_WINDOW_MANAGER=${ENABLE_WINDOW_MANAGER:-""}
DISPLAY=99

clean() {
  if [ -n "$FILESERVER_PID" ]; then
    kill -TERM "$FILESERVER_PID"
  fi
  if [ -n "$XVFB_PID" ]; then
    kill -TERM "$XVFB_PID"
  fi
  if [ -n "$FLUXBOX_PID" ]; then
    kill -TERM "$FLUXBOX_PID"
  fi
  if [ -n "$X11VNC_PID" ]; then
    kill -TERM "$X11VNC_PID"
  fi
}

trap clean SIGINT SIGTERM

/usr/bin/fileserver &
FILESERVER_PID=$!

/usr/bin/xvfb-run -l -n "$DISPLAY" -s "-ac -screen 0 $SCREEN_RESOLUTION -noreset -listen tcp" /usr/bin/chromedriver --port=4444 --whitelisted-ips='' --verbose &
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

if [ -n "$ENABLE_WINDOW_MANAGER" ]; then
    fluxbox -display :$DISPLAY &
    FLUXBOX_PID=$!
fi

if [ "$ENABLE_VNC" == "true" ]; then
    x11vnc -display ":$DISPLAY" -passwd selenoid -shared -forever -loop500 -rfbport 5900 -rfbportv6 5900 -logfile /home/selenium/x11vnc.log &
    X11VNC_PID=$!
fi

wait

