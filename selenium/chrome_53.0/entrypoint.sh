#!/bin/bash

NODE_PORT=${NODE_PORT:-4444}
CHROMEDRIVER_PATH="/usr/bin/chromedriver "
CHROMEDRIVER_ARGS="--port=$NODE_PORT --whitelisted-ips="" --verbose"

echo "Starting Chromedriver on port $NODE_PORT..."
xvfb-run -a -s '-screen 0 1280x1600x24 -noreset' $CHROMEDRIVER_PATH $CHROMEDRIVER_ARGS
