#!/bin/sh
VIDEO_SIZE=${VIDEO_SIZE:-"1920x1080"}
BROWSER_CONTAINER_NAME=${BROWSER_CONTAINER_NAME:-"browser"}
DISPLAY=${DISPLAY:-"99"}
FILE_NAME=${FILE_NAME:-"video.mp4"}
RATE=${FRAME_RATE:-"12"}
for i in {0..50}
do
	nc -z ${BROWSER_CONTAINER_NAME} 60${DISPLAY}
	if [ $? -ne 0 ]
	then
		echo 'wait...'
		sleep 0.1
	else	
		break
	fi
done
exec ffmpeg -y -f x11grab -video_size $VIDEO_SIZE -r ${RATE} -i $BROWSER_CONTAINER_NAME:$DISPLAY -codec:v libx264 "/data/$FILE_NAME"
