#!/bin/bash
SCREEN_RESOLUTION=${SCREEN_RESOLUTION:-"1920x1080x24"}
DISPLAY_NUM=99
URLS=${URLS:-'["http://172.17.0.1:4444/"]'}
VNC_PASSWORD=${VNC_PASSWORD:-"password"}
export URLS
export VNC_PASSWORD
export DISPLAY=":$DISPLAY_NUM"

VERBOSE=${VERBOSE:-""}
if [ -n "$VERBOSE" ]; then
    sed -i 's|@@DRIVER_ARGS@@|, "--log", "debug"|g' /etc/selenoid/browsers.json
fi

clean() {
  if [ -n "$XVFB_PID" ]; then
    kill -TERM "$XVFB_PID"
  fi
  if [ -n "$X11VNC_PID" ]; then
    kill -TERM "$X11VNC_PID"
  fi
  if [ -n "$RPROXY_PID" ]; then
    kill -TERM "$RPROXY_PID"
  fi

}

trap clean SIGINT SIGTERM

URL=$(choose $URL)
HOST=$(echo $URL | sed -e  's@\(.*\)://@@' -e 's@[:/]\(.*\)$@@')

/bin/bash -c 'echo -e "$VNC_PASSWORD\n$VNC_PASSWORD\nn\n" | vncpasswd file.passwd'>/dev/null

xvfb-run -l -n $DISPLAY_NUM -s "-ac -screen 0 $SCREEN_RESOLUTION -noreset -listen tcp" \
	vncviewer passwd=file.passwd $HOST &
XVFB_PID=$!

retcode=1
until [ $retcode -eq 0 ]; do
  xdpyinfo -display $DISPLAY >/dev/null 2>&1
  retcode=$?
  if [ $retcode -ne 0 ]; then
    echo Waiting xvfb...
    sleep 1
  fi
done

x11vnc -display $DISPLAY -passwd selenoid -shared -forever -loop500 -rfbport 5900 -rfbportv6 5900 -logfile /dev/null &
X11VNC_PID=$!

rproxy $URL &
RPROXI_PID=$!

wait

