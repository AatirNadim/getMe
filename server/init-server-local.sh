#!/bin/bash

# This script sets up and runs the getMe server in a local development environment.
# It performs the following steps:
# 1. Builds the Go project. If the build fails, it exits.
# 2. Runs the data directory setup script with sudo to create persistent storage in /var/lib.
# 3. Runs scripts to create temporary log and socket directories in /tmp.
# 4. Starts the background logging services (Loki, Grafana, Alloy) using Docker Compose.
# 5. Runs the main getMe server application in the foreground.

# --- Configuration ---
OUT_DIR="dist"
SCRIPTS_BASE_PATH="./scripts"
LOGGING_COMPOSE_FILE_PATH="./utils/logger/docker-compose.logging.yml"
DATA_DIR_SCRIPT="$SCRIPTS_BASE_PATH/init-data-dir.sh"
LOGS_DIR_SCRIPT="$SCRIPTS_BASE_PATH/init-logs-dir.sh"
SOCK_DIR_SCRIPT="$SCRIPTS_BASE_PATH/init-sock-dir.sh"

echo -e "\n=== Initializing Local Server Environment ===\n"

# --- 1. Build the Go Project ---
echo "--> Building the Go project..."
mkdir -p "$OUT_DIR"
if ! go build -o "$OUT_DIR/getMeMod" .; then
    echo -e "\n[ERROR] Go build failed. Please check the compilation errors above."
    exit 1
fi
echo "Build complete."

# --- 2. Initialize Persistent Data Directory ---
echo -e "\n--> Setting up persistent data directory in /var/lib..."
# This requires sudo to create the directory and set permissions.
sudo "$DATA_DIR_SCRIPT"
echo "Data directory is ready."

# --- 3. Initialize Temporary Log and Socket Directories ---
echo -e "\n--> Setting up temporary directories in /tmp..."
"$LOGS_DIR_SCRIPT"
"$SOCK_DIR_SCRIPT"
echo "Log and socket directories are ready."

# --- 4. Set up Background Logging Services ---
echo -e "\n--> Starting the logging stack via Docker Compose..."
# The docker-compose command is run with '-d' to start containers in the background.
docker compose -f "$LOGGING_COMPOSE_FILE_PATH" up -d
echo "Logging stack is running in the background."

# --- 5. Run the Main Application ---
echo -e "\n--> Starting the getMe server..."
# This is the last command. It runs the server in the foreground.
# The script will end when you stop the server (e.g., with Ctrl+C).
./"$OUT_DIR/getMeMod"

echo -e "\n--- getMe server has stopped. ---"