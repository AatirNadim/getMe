package src

import (
	"fmt"
	"getMeMod/server/store"
	"getMeMod/utils/logger"
	"net/http"
)


func StartServer(socketPath, storePath, compactedStorePath string) error {
	l, err := createSocket(socketPath)
	if err != nil {
		return err
	}
	defer l.Close()

	storeInstance := store.NewStore(storePath, compactedStorePath)
	logger.Info("Store has been initialized at path:", storePath)

	mux := muxHandler(storeInstance)

	server := &http.Server{
		Handler: mux,
	}

	fmt.Println("Server is listening on socket:", socketPath)
	if err := server.Serve(l); err != nil {
		return fmt.Errorf("server error: %w", err)
	}

	return nil
}

