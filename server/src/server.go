package src

import (
	"fmt"
	"getMeMod/server/store"
	"getMeMod/server/store/utils/constants"
	"getMeMod/server/utils/logger"
	"net/http"
	"os"
)

func StartServer(socketPath, storePath, compactedStorePath string, loggingDisabled *bool) error {
	l, err := createSocket(socketPath)
	if err != nil {
		return err
	}
	defer l.Close()

	fmt.Println("Store is being initialized")
	
	storeInstance, err := InitializeStore(storePath, compactedStorePath, loggingDisabled)

	if err != nil {
		return fmt.Errorf("failed to initialize store: %w", err)
	}

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

func InitializeStore(storePath, compactedStorePath string, loggingDisabled *bool) (*store.Store, error) {
	if *loggingDisabled {
		logger.Disable()
		// fmt.Println("Logging disabled; running without file-backed logging")
	} else {
		// if logging is enabled, set the appropriate flag and initialize the logging dir
		logger.Enable()
		fmt.Println("Creating the log file in the executable directory")
		if err := logger.Initialize(constants.LogsDirName); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
			os.Exit(1)
		}
		defer logger.Close()
	}

	storeInstance := store.NewStore(storePath, compactedStorePath)
	logger.Info("Store has been initialized at path:", storePath)
	return storeInstance, nil
}
