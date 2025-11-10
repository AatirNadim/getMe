package src

import (
	"fmt"
	// "getMeMod/server/utils/logger"
	"net"
	"os"
)

func createSocket(socketPath string) (net.Listener, error) {

	if err := os.RemoveAll(socketPath); err != nil {
		return nil, fmt.Errorf("failed to remove existing socket file at %s: %w", socketPath, err)
	}

	fmt.Println("Removed existing socket file at:", socketPath)

	if err := os.MkdirAll("/tmp/getMeStore/sockDir", 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory for socket at /tmp/getMeStore/sockDir: %w", err)
	}
	fmt.Println("directory created or already present for socket at:", socketPath)

	l, err := net.Listen("unix", socketPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create unix socket at %s: %w", socketPath, err)
	}
	fmt.Println("Socket created at:", socketPath)

	return l, nil
}
