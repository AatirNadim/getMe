#!/bin/bash

# Define the vars
OUT_DIR="dist"
LOGGING_COMPOSE_FILE_PATH="./utils/logger/docker-compose.yml"
# Define the log directory that the Go application writes to and Alloy reads from.
LOG_DIR="/tmp/getMeStore/dump"

# This script initializes the server environment.
echo -e "\n=== Initializing Server Environment ===\n"

# --- 1. Build the Go Project ---
echo -e "\n--- Building the getMe server project ---\n"

echo "Creating output directory: $OUT_DIR"
mkdir -p $OUT_DIR

echo "Building the Go project..."
go build -o $OUT_DIR/getMeMod .
echo "Build complete."

# --- 2. Set up Background Services (Logging) ---
echo -e "\n--- Setting up the logging stack ---\n"

echo "Ensuring log directory exists: $LOG_DIR"
mkdir -p "$LOG_DIR"

# The docker-compose command is run first with '-d' to start the containers
# in the background, allowing the script to proceed immediately.
docker compose -f $LOGGING_COMPOSE_FILE_PATH up -d

echo -e "\n--- Logging stack initialized! ---\n"

# --- 3. Run the Main Application ---
echo -e "\n--- Spinning up the getMe server ---\n"

# This is the last command. It runs the server in the foreground,
# so you will see its log output directly in this terminal.
# The script will end when you stop the server (e.g., with Ctrl+C).
./$OUT_DIR/getMeMod

echo -e "\n--- getMe server has stopped. ---"