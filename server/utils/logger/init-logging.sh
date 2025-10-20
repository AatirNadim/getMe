#!/bin/bash

# This script prepares the host environment and launches the Grafana logging stack.

# Define the log directory that the Go application writes to and Alloy reads from.
LOG_DIR="/tmp/getMeStore/dump"

echo "--- Initializing Logging Stack ---"

# Step 1: Ensure the log directory exists on the host.
# The 'grafana-alloy' container requires this directory to be present for its volume mount.
# We use 'mkdir -p' to create it without causing an error if it already exists.
echo "Ensuring log directory exists: $LOG_DIR"
mkdir -p "$LOG_DIR"

# Step 2: Launch the Grafana, Loki, and Alloy services using Docker Compose.
# '-d' runs the containers in detached mode (in the background).
echo "Starting Grafana, Loki, and Alloy containers..."
docker-compose up -d

echo ""
echo "--- Logging Stack Initialized ---"
echo "Grafana should be available at: http://localhost:3000"
echo "Logs are being collected from: $LOG_DIR"