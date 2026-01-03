#!/bin/bash

# This script prepares the host environment and launches the complete,
# containerized getMe application stack using Docker Compose.
#
# It performs the following steps:
# 1. Defines the host paths for the socket directory.
# 2. Creates this directory on the host machine to ensure it is
#    available for bind mounting into the containers.
# 3. Runs `docker compose up` to build the images and start all services,
#    including the getMe server and the logging infrastructure.

# Define host directories that will be mounted into containers.
SOCK_DIR="/tmp/getMeStore/sockDir"
TEARDOWN_SCRIPT="./teardown-server-docker.sh"

# --- Teardown on interrupt/termination ---
TEARDOWN_RAN=0

run_teardown() {
	if [[ "$TEARDOWN_RAN" -eq 1 ]]; then
		return 0
	fi
	TEARDOWN_RAN=1

	echo -e "\n--> Detected shutdown. Running docker teardown..."
	if [[ -x "$TEARDOWN_SCRIPT" ]]; then
		"$TEARDOWN_SCRIPT" || true
	elif [[ -f "$TEARDOWN_SCRIPT" ]]; then
		bash "$TEARDOWN_SCRIPT" || true
	else
		echo "[WARN] Teardown script not found at: $TEARDOWN_SCRIPT"
	fi
}

# NOTE: We intentionally do NOT trap EXIT here, because this init script starts
# the stack in detached mode and then exits normally.
trap run_teardown INT TERM HUP

echo -e "\n=== Initializing Containerized Server Environment ===\n"

# --- 1. Prepare Host Directories ---
echo "--> Ensuring host directories for bind mounts exist..."
mkdir -p "$SOCK_DIR"
echo "Host directories are ready."

# --- 2. Export Host User/Group IDs ---
# This is the crucial step. We export the current user's UID and GID so that
# Docker Compose can use them to run the container with the same permissions.
# This ensures that files created in bind mounts (like the socket file)
# have the correct ownership on the host.
export currentUserId=$(id -u)
export currentGroupId=$(id -g)
echo "--> Running as user $currentUserId:$currentGroupId."

# --- 3. Launch Docker Compose Stack ---
echo -e "\n--> Building and starting all services via Docker Compose..."
# `docker compose up` will read the `docker-compose.yml` file in the current
# directory, build the `get_me_store` image, and start all defined services.
# The `--build` flag ensures the image is rebuilt if the source code has changed.
docker compose up -d --build

echo -e "\n--- Docker Compose stack running in detached mode! ---\n"
