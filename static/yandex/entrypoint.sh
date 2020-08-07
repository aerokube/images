#!/bin/bash
SCREEN_RESOLUTION=${SCREEN_RESOLUTION:-"1920x1080x24"}
ENABLE_WINDOW_MANAGER=${ENABLE_WINDOW_MANAGER:-""}
DISPLAY_NUM=99
export DISPLAY=":$DISPLAY_NUM"

QUIET=${QUIET:-""}
DRIVER_ARGS=""
if [ -z "$QUIET" ]; then
    DRIVER_ARGS="--verbose"
fi

ROOT_CA_PATH=${ROOT_CA_PATH:-""}
if [ -r "$ROOT_CA_PATH" ]; then
    ROOT_CA=$(<"$ROOT_CA_PATH")
    ROOT_CA_NAME=$(basename "$ROOT_CA_PATH" | sed -e 's|.crt||g')
fi

ROOT_CA=${ROOT_CA:-""}
if [ -n "$ROOT_CA" ]; then
    ROOT_CA_NAME=${ROOT_CA_NAME:-"UserRootCA"}
    mkdir -p ~/.pki/nssdb
    certutil -d "sql:$HOME/.pki/nssdb" -N --empty-password
    echo "$ROOT_CA" | certutil -d "sql:$HOME/.pki/nssdb" -A -t TC -n "$ROOT_CA_NAME"
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
  if [ -n "$DRIVER_PID" ]; then
    kill -TERM "$DRIVER_PID"
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

/usr/bin/xvfb-run -l -n "$DISPLAY_NUM" -s "-ac -screen 0 $SCREEN_RESOLUTION -noreset -listen tcp" /usr/bin/fluxbox -display "$DISPLAY" -log /tmp/fluxbox.log  2>/dev/null &
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
    x11vnc -display "$DISPLAY" -passwd selenoid -shared -forever -loop500 -rfbport 5900 -rfbportv6 5900 -logfile /tmp/x11vnc.log &
    X11VNC_PID=$!
fi

DISPLAY="$DISPLAY" /usr/bin/yandexdriver --port=4444 --whitelisted-ips='' ${DRIVER_ARGS} &
DRIVER_PID=$!

wait
