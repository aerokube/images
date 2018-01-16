#!/bin/bash
CHROMEDRIVER_PORT=9515
BOOTSTRAP_PORT=4725
EMULATOR=emulator-5554
CHROME_ARGS=""
TIMEOUT=${TIMEOUT:-"60"}
SKIN=${SKIN:-"WXGA800"}

if [ -n "$DENSITY" ]; then
    echo "hw.lcd.density=$DENSITY" >> /root/.android/avd/android6.0-1.avd/config.ini
fi

Xorg :0 &
sleep 1
x11vnc -display ":0" -passwd selenoid -shared -forever -loop500 -rfbport 5900 -rfbportv6 5900 -logfile /var/log/x11vnc.log &
ANDROID_AVD_HOME=/root/.android/avd DISPLAY=:0 /opt/android-sdk-linux/emulator/emulator -no-audio -no-jni -avd android6.0-1 -sdcard /sdcard.img -skin "$SKIN" -skindir /opt/android-sdk-linux/platforms/android-23/skins/ -verbose -gpu swiftshader -qemu -enable-kvm &
sleep 30

#if [ -n "$INSTALL_APP" ]; then
#	for i in 1 2 3 4 5; do adb install /chrome.apk && break || sleep 15 && echo "retrying install app"; done
#fi

sleep 5
#if [ -z "$ADD_APP" ]; then
#	CHROME_ARGS="--chromedriver-port $CHROMEDRIVER_PORT --app-pkg \"com.android.chrome\" --app-activity \"com.google.android.apps.chrome.Main\" --no-reset"
#fi
/opt/node_modules/.bin/appium -a 0.0.0.0 -p 4444 -bp $BOOTSTRAP_PORT -U $EMULATOR --platform-name Android --device-name android --log-timestamp --log-no-colors --default-capabilities "{\"newCommandTimeout\": \"$TIMEOUT\"}" $CHROME_ARGS
