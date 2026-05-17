package src

import (
	"fmt"
	"net/http"
	"os"

	"github.com/AatirNadim/getMe/server/store"
	"github.com/AatirNadim/getMe/server/store/utils/constants"
	"github.com/AatirNadim/getMe/server/utils/logger"
)

func StartServer(socketPath, storePath, compactedStorePath string, loggingDisabled *bool, logToStdout *bool) error {
	l, err := createSocket(socketPath)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := l.Close(); closeErr != nil {
			fmt.Fprintf(os.Stderr, "Failed to close listener: %v\n", closeErr)
		}
	}()

	fmt.Println("Store is being initialized")

	storeInstance, err := InitializeStore(storePath, compactedStorePath, loggingDisabled, logToStdout)

	if err != nil {
		return fmt.Errorf("failed to initialize store: %w", err)
	}
	defer func() {
		if err := logger.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to close logger: %v\n", err)
		}
	}()

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

func InitializeStore(storePath, compactedStorePath string, loggingDisabled *bool, logToStdout *bool) (*store.Store, error) {
	if *loggingDisabled {
		logger.Disable()
		// fmt.Println("Logging disabled; running without file-backed logging")
	} else {
		// if logging is enabled, set the appropriate flag and initialize the logging dir
		logger.Enable()
		fmt.Println("Creating the log file in the executable directory")
		if err := logger.Initialize(constants.LogsDirName, logToStdout); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
			os.Exit(1)
		}
	}

	storeInstance := store.NewStore(storePath, compactedStorePath)
	logger.Info("Store has been initialized at path:", storePath)
	return storeInstance, nil
}
