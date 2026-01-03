#!/bin/bash

# This script tears down the local getMe server environment created by init-server-local.sh.
# It performs the following steps:
# 1. Stops the background logging stack (Loki, Grafana, Alloy) via Docker Compose.
# 2. Removes the temporary log and socket directories under /tmp.
# 3. Removes the persistent data directory under /var/lib (typically requires sudo).
#
# IMPORTANT:
# - Run this script from the 'server/' directory (same assumption as init-server-local.sh).
# - Even though init-data-dir.sh transfers ownership of the data directory to the sudo user,
#   deleting it usually still requires sudo because the parent directory under /var/lib is root-owned.

set -euo pipefail

# --- Configuration ---
LOGGING_COMPOSE_FILE_PATH="./utils/logger/docker-compose.logging.yml"

DATA_DIR="/var/lib/getMeStore/dataDir"

TMP_BASE_DIR="/tmp/getMeStore"
LOG_DIR="$TMP_BASE_DIR/dumpDir"
SOCK_DIR="$TMP_BASE_DIR/sockDir"

echo -e "\n=== Tearing Down Local Server Environment ===\n"

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

remove_dir_with_sudo_fallback() {
    local path="$1"
    abort_if_unsafe_path "$path"

    if [[ ! -e "$path" ]]; then
        echo "[SKIP] Not found: $path"
        return 0
    fi

    echo "--> Removing: $path"
    if rm -rf -- "$path" 2>/dev/null; then
        return 0
    fi

    echo "[INFO] Permission denied; retrying with sudo..."
    sudo rm -rf -- "$path"
}

# --- 1. Stop Background Logging Services ---
echo "--> Stopping the logging stack via Docker Compose..."
if command -v docker >/dev/null 2>&1 && [[ -f "$LOGGING_COMPOSE_FILE_PATH" ]]; then
    set +e
    echo "Taking down logging stack containers..."
    docker compose -f "$LOGGING_COMPOSE_FILE_PATH" down --remove-orphans --volumes
    status=$?
    set -e

    if [[ $status -ne 0 ]]; then
        echo "[WARN] Failed to bring down logging stack (docker status=$status). Continuing teardown."
    else
        echo "Logging stack is stopped."
    fi
else
    echo "[SKIP] Docker or compose file not available; skipping container teardown."
fi

# --- 2. Remove Temporary Log and Socket Directories ---
echo -e "\n--> Removing temporary directories in /tmp..."
if ! remove_dir "$LOG_DIR"; then
    echo "[WARN] Could not remove $LOG_DIR (permissions?). You may need to run: sudo rm -rf '$LOG_DIR'"
fi

if ! remove_dir "$SOCK_DIR"; then
    echo "[WARN] Could not remove $SOCK_DIR (permissions?). You may need to run: sudo rm -rf '$SOCK_DIR'"
fi

# Try removing the shared base directory if empty.
if [[ -d "$TMP_BASE_DIR" ]]; then
    rmdir -- "$TMP_BASE_DIR" 2>/dev/null || true
fi

echo "Temporary directories teardown complete."

# --- 3. Remove Persistent Data Directory ---
echo -e "\n--> Removing persistent data directory in /var/lib..."
remove_dir_with_sudo_fallback "$DATA_DIR"
echo "Persistent data directory teardown complete."

echo -e "\n=== Teardown Complete ===\n"
