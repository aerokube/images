#!/bin/sh
VIDEO_SIZE=${VIDEO_SIZE:-"1920x1080"}
BROWSER_CONTAINER_NAME=${BROWSER_CONTAINER_NAME:-"browser"}
DISPLAY=${DISPLAY:-"99"}
FILE_NAME=${FILE_NAME:-"video.mp4"}
FRAME_RATE=${FRAME_RATE:-"12"}
CODEC=${CODEC:-"libx264"}
PRESET=${PRESET:-""}
if [ "$CODEC" == "libx264" -a -n "$PRESET" ]; then
    PRESET="-preset $PRESET"
fi
INPUT_OPTIONS=${INPUT_OPTIONS:-""}
HIDE_CURSOR=${HIDE_CURSOR:-""}
if [ -n "$HIDE_CURSOR" ]; then
    INPUT_OPTIONS="$INPUT_OPTIONS -draw_mouse 0"
fi
retcode=1
max_attempts=300
attempts=0
echo 'Waiting for display to open...'
until [ $retcode -eq 0 -o $attempts -eq $max_attempts ]; do
	xset -display ${BROWSER_CONTAINER_NAME}:${DISPLAY} b off > /dev/null 2>&1
	retcode=$?
	if [ $retcode -ne 0 ]; then
		echo 'Sleeping before next attempt...'
		sleep 0.1
	fi
	attempts=$((attempts+1))
done
exec ffmpeg -y -f x11grab -video_size ${VIDEO_SIZE} -r ${FRAME_RATE} ${INPUT_OPTIONS} -i ${BROWSER_CONTAINER_NAME}:${DISPLAY} -codec:v ${CODEC} ${PRESET} -pix_fmt yuv420p "/data/$FILE_NAME"
