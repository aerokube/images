#!/bin/bash
CHROMEDRIVER_PORT=9515
BOOTSTRAP_PORT=4725
EMULATOR=emulator-5554
APPIUM_ARGS=${APPIUM_ARGS:-""}
PORT=${PORT:-"4444"}
DISPLAY_NUM=99
export DISPLAY=":$DISPLAY_NUM"
SCREEN_RESOLUTION=${SCREEN_RESOLUTION:-"1920x1080x24"}
SKIN=${SKIN:-"1080x1920"}
STOP=""

clean() {
  STOP="yes"
  if [ -n "$APPIUM_PID" ]; then
    kill -TERM "$APPIUM_PID"
  fi
  if [ -n "$EMULATOR_PID" ]; then
    kill -TERM "$EMULATOR_PID"
  fi
  if [ -n "$X11VNC_PID" ]; then
    kill -TERM "$X11VNC_PID"
  fi
  if [ -n "$DEVTOOLS_PID" ]; then
    kill -TERM "$DEVTOOLS_PID"
  fi
  if [ -n "$XVFB_PID" ]; then
    kill -TERM "$XVFB_PID"
  fi
}

trap clean SIGINT SIGTERM

/usr/bin/xvfb-run -e /dev/stdout -l -n "$DISPLAY_NUM" -s "-ac -screen 0 $SCREEN_RESOLUTION -noreset -listen tcp" /usr/bin/fluxbox -display "$DISPLAY" -log /root/fluxbox.log 2>/dev/null &
XVFB_PID=$!

retcode=1
until [ $retcode -eq 0 -o -n "$STOP" ]; do
  DISPLAY="$DISPLAY" wmctrl -m >/dev/null 2>&1
  retcode=$?
  if [ $retcode -ne 0 ]; then
    echo Waiting X server...
    sleep 0.1
  fi
done
if [ -n "$STOP" ]; then exit 0; fi

if [ "$ENABLE_VNC" != "true" -a "$ENABLE_VIDEO" != "true" ]; then
    EMULATOR_ARGS="-no-window"
fi
ANDROID_AVD_HOME=/root/.android/avd DISPLAY="$DISPLAY" /opt/android-sdk-linux/emulator/emulator ${EMULATOR_ARGS} -writable-system -no-boot-anim -no-audio -no-jni -avd @AVD_NAME@ -sdcard /sdcard.img -skin "$SKIN" -skindir /opt/android-sdk-linux/platforms/@PLATFORM@/skins/ -verbose -gpu swiftshader_indirect -ranchu -qemu -enable-kvm &
EMULATOR_PID=$!

if [ "$ENABLE_VNC" == "true" ]; then
    x11vnc -display "$DISPLAY" -passwd selenoid -shared -forever -loop500 -rfbport 5900 -rfbportv6 5900 -logfile /root/x11vnc.log &
    X11VNC_PID=$!
fi

while [ "`adb shell getprop sys.boot_completed | tr -d '\r' `" != "1" -a -z "$STOP" ] ; do sleep 1; done
if [ -n "$STOP" ]; then exit 0; fi

if [ -n "@CHROME_MOBILE@" ]; then
    while ip addr | grep inet | grep -q tentative > /dev/null; do sleep 0.1; done
	APPIUM_ARGS="$APPIUM_ARGS --chromedriver-port $CHROMEDRIVER_PORT"
    /usr/bin/devtools &
    DEVTOOLS_PID=$!
fi

if [ -x "/usr/bin/chromedriver" ]; then
    APPIUM_ARGS="$APPIUM_ARGS --chromedriver-executable /usr/bin/chromedriver"
fi

/opt/node_modules/.bin/appium -a 0.0.0.0 -p "$PORT" -bp "$BOOTSTRAP_PORT" -U "$EMULATOR" --platform-name Android --device-name android --log-timestamp --log-no-colors --command-timeout 90 --no-reset ${APPIUM_ARGS} --default-capabilities '{"androidNaturalOrientation": true}' &
APPIUM_PID=$!

wait
