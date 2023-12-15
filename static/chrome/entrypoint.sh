#!/bin/bash
SCREEN_RESOLUTION=${SCREEN_RESOLUTION:-"1920x1080x24"}
DISPLAY_NUM=99
export DISPLAY=":$DISPLAY_NUM"

VERBOSE=${VERBOSE:-""}
DRIVER_ARGS=${DRIVER_ARGS:-""}
if [ -n "$VERBOSE" ]; then
    DRIVER_ARGS="$DRIVER_ARGS --verbose"
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
  if [ -n "$DEVTOOLS_PID" ]; then
    kill -TERM "$DEVTOOLS_PID"
  fi
  if [ -n "$PULSE_PID" ]; then
    kill -TERM "$PULSE_PID"
  fi
}

trap clean SIGINT SIGTERM

if env | grep -q ROOT_CA_; then
  mkdir -p $HOME/.pki/nssdb
  certutil -N --empty-password -d sql:$HOME/.pki/nssdb
  for e in $(env | grep ROOT_CA_ | sed -e 's/=.*$//'); do
    certname=$(echo -n $e | sed -e 's/ROOT_CA_//')
    echo ${!e} | base64 -d >/tmp/cert.pem
    certutil -A -n ${certname} -t "TC,C,T" -i /tmp/cert.pem -d sql:$HOME/.pki/nssdb
    if cat tmp/cert.pem | grep -q "PRIVATE KEY"; then
      PRIVATE_KEY_PASS=${PRIVATE_KEY_PASS:-\'\'}
      openssl pkcs12 -export -in /tmp/cert.pem -clcerts -nodes -out /tmp/key.p12 -passout pass:${PRIVATE_KEY_PASS} -passin pass:${PRIVATE_KEY_PASS}
      pk12util -d sql:$HOME/.pki/nssdb -i /tmp/key.p12 -W ${PRIVATE_KEY_PASS}
      rm /tmp/key.p12
    fi
    rm /tmp/cert.pem
  done
fi

if env | grep -q CH_POLICY_; then
  for p in $(env | grep CH_POLICY_ | sed 's/CH_POLICY_//'); do
    jsonkey=$(echo $p | sed 's/=.*//')
    jsonvalue=$(echo $p | sed 's/^.*=//')
    cat <<< $(jq --arg key $jsonkey --argjson value $jsonvalue '.[$key] = $value' /etc/opt/chrome/policies/managed/policies.json) > /etc/opt/chrome/policies/managed/policies.json
  done
fi

/usr/bin/fileserver &
FILESERVER_PID=$!

/usr/bin/devtools &
DEVTOOLS_PID=$!

DISPLAY="$DISPLAY" /usr/bin/xseld &
XSELD_PID=$!

while ip addr | grep inet | grep -q tentative > /dev/null; do sleep 0.1; done

mkdir -p ~/.config/pulse
echo -n 'gIvST5iz2S0J1+JlXC1lD3HWvg61vDTV1xbmiGxZnjB6E3psXsjWUVQS4SRrch6rygQgtpw7qmghDFTaekt8qWiCjGvB0LNzQbvhfs1SFYDMakmIXuoqYoWFqTJ+GOXYByxpgCMylMKwpOoANEDePUCj36nwGaJNTNSjL8WBv+Bf3rJXqWnJ/43a0hUhmBBt28Dhiz6Yqowa83Y4iDRNJbxih6rB1vRNDKqRr/J9XJV+dOlM0dI+K6Vf5Ag+2LGZ3rc5sPVqgHgKK0mcNcsn+yCmO+XLQHD1K+QgL8RITs7nNeF1ikYPVgEYnc0CGzHTMvFR7JLgwL2gTXulCdwPbg=='| base64 -d>~/.config/pulse/cookie
pulseaudio --start --exit-idle-time=-1
pactl load-module module-native-protocol-tcp
PULSE_PID=$(ps --no-headers -C pulseaudio -o pid | sed -r 's/( )+//g')

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

DISPLAY="$DISPLAY" /usr/bin/chromedriver --port=4444 --allowed-ips='' --allowed-origins='*' ${DRIVER_ARGS} &
DRIVER_PID=$!

echo Moving mouse to position 0 0
DISPLAY="$DISPLAY" /usr/bin/xdotool mousemove --sync 0 0

wait
