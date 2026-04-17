#!/bin/sh
set -e

# Handle shutdown signals gracefully
trap 'kill -TERM $SERVER_PID $PROXY_PID; wait $SERVER_PID $PROXY_PID' INT TERM

echo "Starting getMe-server in the background..."
/usr/local/bin/getMe-server &
SERVER_PID=$!

# Wait a brief moment to ensure the server starts and creates the socket file
sleep 2

echo "Starting getMe-proxy..."
/usr/local/bin/getMe-proxy &
PROXY_PID=$!

# Wait for background processes to finish
wait $PROXY_PID
wait $SERVER_PID
