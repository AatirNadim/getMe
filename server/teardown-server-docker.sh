#!/bin/bash

# This script tears down the containerized getMe server environment created by init-server-docker.sh.
# It performs the following steps:
# 1. Stops the Docker Compose stack (getMe server + logging stack).
# 2. Removes named volumes created by docker-compose.yml.
# 3. Removes the host socket directory used for bind mounting.
#
# IMPORTANT:
# - Run this script from the 'server/' directory (same assumption as init-server-docker.sh).

set -euo pipefail

SOCK_DIR="/tmp/getMeStore/sockDir"
TMP_BASE_DIR="/tmp/getMeStore"

echo -e "\n=== Tearing Down Containerized Server Environment ===\n"

abort_if_unsafe_path() {
    local path="$1"
    if [[ -z "$path" || "$path" == "/" ]]; then
        echo "[ERROR] Refusing to operate on an unsafe path ('$path')."
        exit 1
    fi
}

remove_dir() {
    local path="$1"
    abort_if_unsafe_path "$path"

    if [[ ! -e "$path" ]]; then
        echo "[SKIP] Not found: $path"
        return 0
    fi

    echo "--> Removing: $path"
    rm -rf -- "$path" 2>/dev/null || return 1
}

# --- 1. Stop Docker Compose Stack ---
echo "--> Stopping Docker Compose stack..."
if command -v docker >/dev/null 2>&1 && [[ -f "./docker-compose.yml" ]]; then
    set +e
    docker compose down --remove-orphans --volumes
    status=$?
    set -e

    if [[ $status -ne 0 ]]; then
        echo "[WARN] Failed to bring down stack (docker status=$status). Continuing teardown."
    else
        echo "Docker Compose stack is stopped."
    fi
else
    echo "[SKIP] Docker or docker-compose.yml not available; skipping container teardown."
fi

# --- 2. Remove Host Socket Directory ---
echo -e "\n--> Removing host socket directory..."
if ! remove_dir "$SOCK_DIR"; then
    echo "[WARN] Could not remove $SOCK_DIR (permissions?). You may need to run: sudo rm -rf '$SOCK_DIR'"
fi

# Try removing the shared base directory if empty.
if [[ -d "$TMP_BASE_DIR" ]]; then
    rmdir -- "$TMP_BASE_DIR" 2>/dev/null || true
fi

echo -e "\n=== Teardown Complete ===\n"
