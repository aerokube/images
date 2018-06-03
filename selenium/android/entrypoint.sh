#!/bin/bash
mkdir -p /etc/appium
CONFIG=/etc/appium/appium.json
CHROMEDRIVER_PORT=9515
BOOTSTRAP_PORT=4725
EMULATOR=emulator-5554
APPIUM_ARGS=""
PORT=${PORT:-"4444"}
DISPLAY=99
SCREEN_RESOLUTION=${SCREEN_RESOLUTION:-"1920x1080x24"}
SKIN=${SKIN:-"1080x1920"}

clean() {
  if [ -n "$XVFB_PID" ]; then
    kill -TERM "$XVFB_PID"
  fi
  if [ -n "$X11VNC_PID" ]; then
    kill -TERM "$X11VNC_PID"
  fi
  if [ -n "$APPIUM_PID" ]; then
    kill -TERM "$APPIUM_PID"
  fi
}

trap clean SIGINT SIGTERM

if [ "$ENABLE_VNC" != "true" ]; then
    EMULATOR_ARGS="-no-window"
fi
/usr/bin/xvfb-run -l -n "$DISPLAY" -s "-ac -screen 0 $SCREEN_RESOLUTION -noreset -listen tcp" /bin/sh -c "ANDROID_AVD_HOME=/root/.android/avd DISPLAY=:$DISPLAY /opt/android-sdk-linux/emulator/emulator $EMULATOR_ARGS -no-boot-anim -no-audio -no-jni -avd @AVD_NAME@ -sdcard /sdcard.img -skin $SKIN -skindir /opt/android-sdk-linux/platforms/@PLATFORM@/skins/ -verbose -gpu swiftshader_indirect -qemu -enable-kvm" &
XVFB_PID=$!

while [ "`adb shell getprop sys.boot_completed | tr -d '\r' `" != "1" ] ; do sleep 1; done

if [ -n "@CHROME_MOBILE@" ]; then
	APPIUM_ARGS="$APPIUM_ARGS --chromedriver-port $CHROMEDRIVER_PORT --app-pkg \"com.android.chrome\" --app-activity \"com.google.android.apps.chrome.Main\""
fi

if [ -x "/usr/bin/chromedriver" ]; then
    APPIUM_ARGS="$APPIUM_ARGS --chromedriver-executable /usr/bin/chromedriver"
fi
if [ "$ENABLE_VNC" == "true" ]; then
    x11vnc -display ":$DISPLAY" -passwd selenoid -shared -forever -loop500 -rfbport 5900 -rfbportv6 5900 -logfile /root/x11vnc.log &
    X11VNC_PID=$!
fi

/opt/node_modules/.bin/appium -a 0.0.0.0 -p "$PORT" -bp "$BOOTSTRAP_PORT" -U "$EMULATOR" --platform-name Android --device-name android --log-timestamp --log-no-colors --command-timeout 90 --no-reset ${APPIUM_ARGS} &
APPIUM_PID=$!

wait