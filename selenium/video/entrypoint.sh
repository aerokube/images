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

mkdir -p ~/.config/pulse

echo -n 'gIvST5iz2S0J1+JlXC1lD3HWvg61vDTV1xbmiGxZnjB6E3psXsjWUVQS4SRrch6rygQgtpw7qmghDFTaekt8qWiCjGvB0LNzQbvhfs1SFYDMakmIXuoqYoWFqTJ+GOXYByxpgCMylMKwpOoANEDePUCj36nwGaJNTNSjL8WBv+Bf3rJXqWnJ/43a0hUhmBBt28Dhiz6Yqowa83Y4iDRNJbxih6rB1vRNDKqRr/J9XJV+dOlM0dI+K6Vf5Ag+2LGZ3rc5sPVqgHgKK0mcNcsn+yCmO+XLQHD1K+QgL8RITs7nNeF1ikYPVgEYnc0CGzHTMvFR7JLgwL2gTXulCdwPbg=='| base64 -d>~/.config/pulse/cookie

export PULSE_SERVER=${BROWSER_CONTAINER_NAME}

if pactl info >/dev/null 2>&1; then
  exec ffmpeg -f pulse -thread_queue_size 1024 -i default -y -f x11grab -video_size ${VIDEO_SIZE} -r ${FRAME_RATE} ${INPUT_OPTIONS} -i ${BROWSER_CONTAINER_NAME}:${DISPLAY} -codec:v ${CODEC} ${PRESET} -pix_fmt yuv420p -filter:v "pad=ceil(iw/2)*2:ceil(ih/2)*2" "/data/$FILE_NAME"
else
  exec ffmpeg -y -f x11grab -video_size ${VIDEO_SIZE} -r ${FRAME_RATE} ${INPUT_OPTIONS} -i ${BROWSER_CONTAINER_NAME}:${DISPLAY} -codec:v ${CODEC} ${PRESET} -pix_fmt yuv420p -filter:v "pad=ceil(iw/2)*2:ceil(ih/2)*2" "/data/$FILE_NAME"
fi
