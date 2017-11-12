#!/bin/bash
VIDEO_SIZE=${VIDEO_SIZE:-"1920x1080"}
DISPLAY=${DISPLAY:-"99"}
FILENAME=${FILENAME:-"video.mp4"}
exec ffmpeg -y -f x11grab -video_size "$VIDEO_SIZE" -i "browser:$DISPLAY" -codec:v libx264 -r 12 "/data/$FILENAME"
