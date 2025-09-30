package src

import (
	"fmt"
	"getMeMod/server/store"
	"getMeMod/utils/logger"
	"net/http"
)


func StartServer(socketPath string, storePath string) error {
	l, err := createSocket(socketPath)
	if err != nil {
		return err
	}
	defer l.Close()

	storeInstance := store.NewStore(storePath)
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
		fmt.Fprintf(w, "%s", value)
	})

	mux.HandleFunc("POST /put", func(w http.ResponseWriter, r *http.Request) {
		// if r.Method != http.MethodPost {
		// 	http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		// 	return
		// }
		key := r.URL.Query().Get("key")
		value := r.URL.Query().Get("value")
		if key == "" || value == "" {
			http.Error(w, "missing key or value parameter", http.StatusBadRequest)
			return
		}
		if err := storeInstance.Put(key, value); err != nil {
			http.Error(w, fmt.Sprintf("error putting value for key '%s': %v", key, err), http.StatusInternalServerError)
			return
		}

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

		fmt.Fprintf(w, "Successfully deleted key '%s'\n", key)
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
