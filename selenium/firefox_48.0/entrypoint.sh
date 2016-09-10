#!/bin/bash

NODE_PORT=${NODE_PORT:-4444}
GECKODRIVER_PATH="/usr/bin/geckodriver"
GECKODRIVER_ARGS="--host :: --port $NODE_PORT --log debug"

echo "Starting Geckodriver on port $NODE_PORT..."
xvfb-run -a -s '-screen 0 1280x1600x24 -noreset' $GECKODRIVER_PATH $GECKODRIVER_ARGS
