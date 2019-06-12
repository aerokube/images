#!/bin/bash
MAX_ATTEMPTS=5
adb root
adb devices | grep emulator | cut -f1 | while read id; do
    apks=(/usr/bin/*.apk)
    for apk in "${apks[@]}"; do
        if [ -r "$apk" ]; then
            for i in `seq 1 ${MAX_ATTEMPTS}`; do
                echo "Installing $apk (attempt #$i of $MAX_ATTEMPTS)"
                adb -s "$id" install -r "$apk" && break || sleep 15 && echo "Retrying to install $apk"
            done
        fi
    done
    adb -s "$id" emu kill -2 || true
done
rm -f /tmp/.X99-lock || true
