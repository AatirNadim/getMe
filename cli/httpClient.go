package cli

import (
	"context"
	"getMeMod/store/logger"
	"net"
	"net/http"
)

func createHttpClient(socketPath string) (*http.Client, error) {

	logger.Info("Creating HTTP client with socket path:", socketPath)
	httpClient := &http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", socketPath)
			},
		},
	}
	return httpClient, nil
}
