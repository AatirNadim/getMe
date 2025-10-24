#!/bin/bash

# NOTE: Be very careful with the user you run this script as. 
# This script creates a logging dir inside the /tmp dir, 
# to which the application will write logs and the logging stack will read from. 
# If you run this script as root, the log dir will be owned by root, 
# and the application (if run as a non-root user) will not be able to write logs to it.

# This script initializes the socket directory for the getMe application.
SOCK_DIR="/tmp/getMeStore/sockDir"

echo "--- Initializing Socket Directory ---"
mkdir -p "$SOCK_DIR"
echo "Socket directory created at: $SOCK_DIR"