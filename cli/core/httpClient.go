package core

import (
	"context"
	"net"
	"net/http"

	"github.com/AatirNadim/getMe/utils/logger"
)

func CreateHttpClient(socketPath string) (*http.Client, error) {

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
