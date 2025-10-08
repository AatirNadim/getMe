package src

import (
	"encoding/json"
	"fmt"
	"getMeMod/server/store"
	"getMeMod/utils/logger"
	"io"
	"net/http"
)

type PutRequestBody struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func StartServer(socketPath, storePath, compactedStorePath string) error {
	l, err := createSocket(socketPath)
	if err != nil {
		return err
	}
	defer l.Close()

	storeInstance := store.NewStore(storePath, compactedStorePath)
	logger.Info("Store has been initialized at path:", storePath)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /get", func(w http.ResponseWriter, r *http.Request) {
		key := r.URL.Query().Get("key")
		if key == "" {
			http.Error(w, "missing key parameter", http.StatusBadRequest)
			return
		}
		value, found, err := storeInstance.Get(key)
		if err != nil {
			http.Error(w, fmt.Sprintf("error getting value for key '%s': %v", key, err), http.StatusInternalServerError)
			return
		}
		if !found {
			http.Error(w, fmt.Sprintf("key '%s' not found", key), http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "%s", value)
	})

	mux.HandleFunc("POST /put", func(w http.ResponseWriter, r *http.Request) {
		// if r.Method != http.MethodPost {
		// 	http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		// 	return
		// }

		logger.Debug("Handling PUT request, parsing form data")
		body, err := io.ReadAll(r.Body)

		if err != nil {
			logger.Error("Error reading request body:", err)
			http.Error(w, fmt.Sprintf("failed to read request body: %v", err), http.StatusBadRequest)
			return
		}

		var requestPayload PutRequestBody

		if err := json.Unmarshal(body, &requestPayload); err != nil {
			logger.Error("Error parsing JSON body:", err)
			http.Error(w, fmt.Sprintf("failed to parse JSON body: %v", err), http.StatusBadRequest)
			return
		}

		logger.Debug("Parsed request payload:", requestPayload)

		key := requestPayload.Key
		value := requestPayload.Value

		if key == "" || value == "" {
			logger.Error("Missing key or value in request")
			http.Error(w, "missing key or value parameter", http.StatusBadRequest)
			return
		}
		if err := storeInstance.Put(key, value); err != nil {
			logger.Error("Error putting value in store:", err)
			http.Error(w, fmt.Sprintf("error putting value for key '%s': %v", key, err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Successfully set value for key '%s'\n", key)
	})

	mux.HandleFunc("DELETE /delete", func(w http.ResponseWriter, r *http.Request) {
		// if r.Method != http.MethodDelete {
		// 	http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		// 	return
		// }
		key := r.URL.Query().Get("key")
		if key == "" {
			http.Error(w, "missing key parameter", http.StatusBadRequest)
			return
		}
		if err := storeInstance.Delete(key); err != nil {
			http.Error(w, fmt.Sprintf("error deleting key '%s': %v", key, err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Successfully deleted key '%s'\n", key)
	})

	mux.HandleFunc("DELETE /clearStore", func(w http.ResponseWriter, r *http.Request) {

		if err := storeInstance.Clear(); err != nil {
			http.Error(w, fmt.Sprintf("error clearing store: %v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Successfully cleared the store\n")

	})

	server := &http.Server{
		Handler: mux,
	}

	fmt.Println("Server is listening on socket:", socketPath)
	if err := server.Serve(l); err != nil {
		return fmt.Errorf("server error: %w", err)
	}

	return nil
}
