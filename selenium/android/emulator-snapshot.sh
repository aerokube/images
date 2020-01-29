#!/bin/bash
MAX_ATTEMPTS=5
adb root
adb devices | grep emulator | cut -f1 | while read id; do
    adb -s "$id" emu kill -2 || true
done
rm -f /tmp/.X99-lock || true
