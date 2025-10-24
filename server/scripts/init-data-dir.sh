#!/bin/bash

# This script performs the one-time setup required to run the getMe server
# directly on a host machine. It creates necessary directories in /var/lib
# and sets the correct permissions.
#
# This script MUST be run with sudo.

# --- Check for root privileges ---
if [ "$EUID" -ne 0 ]; then
  echo "Please run this script with sudo."
  exit 1
fi

# --- Configuration ---
# Define the base directory for all application data.
BASE_DIR="/var/lib/getMeStore/dataDir"
STORE_DIR="$BASE_DIR/segments"
COMPACTED_DIR="$BASE_DIR/compactedSegments"

# Get the user who invoked sudo, not 'root'.
# This ensures we give ownership to the actual user.
OWNER_USER=${SUDO_USER:-$(whoami)}

echo "--- Setting up getMe Server Environment ---"

# --- Step 1: Create Directories ---
echo "Creating directories in $BASE_DIR..."
mkdir -p "$STORE_DIR"
mkdir -p "$COMPACTED_DIR"
echo "Directories created."

# --- Step 2: Set Ownership ---
echo "Setting ownership of $BASE_DIR to user '$OWNER_USER'..."
chown -R "$OWNER_USER:$OWNER_USER" "$BASE_DIR"
echo "Ownership set."

echo ""
echo "--- Setup Complete ---"
echo "You can now run the 'getMe' server application as the '$OWNER_USER' user."