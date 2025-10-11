package src

import (
	"encoding/json"
	"fmt"
	"getMeMod/server/store"
	"getMeMod/server/utils"
	"getMeMod/utils/logger"
	"io"
	"net/http"
)


type Controllers struct {}


func GetController(storeInstance *store.Store) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
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
	}
}

func PutController(storeInstance *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		var requestPayload utils.PutRequestBody	

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
	}
}

func DeleteController (storeInstance *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
	}
}


func ClearStoreController(storeInstance *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if err := storeInstance.Clear(); err != nil {
			http.Error(w, fmt.Sprintf("error clearing store: %v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Successfully cleared the store\n")

	}
}

func BatchSetCController(storeInstance *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Error("Error reading batch-set request body:", err)
			http.Error(w, "failed to read request body", http.StatusBadRequest)
			return
		}

		var batch map[string]string
		if err := json.Unmarshal(body, &batch); err != nil {
			logger.Error("Error parsing batch-set JSON body:", err)
			http.Error(w, "failed to parse JSON body", http.StatusBadRequest)
			return
		}

		if err := storeInstance.BatchSet(batch); err != nil {
			logger.Error("Error in BatchSet operation:", err)
			http.Error(w, fmt.Sprintf("error during batch set: %v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Batch set operation successful")
	}
}
