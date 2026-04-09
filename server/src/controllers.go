package src

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/AatirNadim/getMe/server/store"
	local "github.com/AatirNadim/getMe/server/store/utils"
	"github.com/AatirNadim/getMe/server/store/utils/constants"
	"github.com/AatirNadim/getMe/server/utils"

	"github.com/AatirNadim/getMe/server/utils/logger"
)

type Controllers struct {
	StoreInstance *store.Store
}

func (c *Controllers) GetController() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		key := r.URL.Query().Get("key")
		if key == "" {
			http.Error(w, "missing key parameter", http.StatusBadRequest)
			return
		}
		value, found, err := c.StoreInstance.Get(key)

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

func (c *Controllers) PutController() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

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
		if err := c.StoreInstance.Put(key, value); err != nil {
			logger.Error("Error putting value in store:", err)
			http.Error(w, fmt.Sprintf("error putting value for key '%s': %v", key, err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Successfully set value for key '%s'\n", key)
	}
}

func (c *Controllers) DeleteController() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		key := r.URL.Query().Get("key")
		if key == "" {
			http.Error(w, "missing key parameter", http.StatusBadRequest)
			return
		}
		if err := c.StoreInstance.Delete(key); err != nil {
			http.Error(w, fmt.Sprintf("error deleting key '%s': %v", key, err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Successfully deleted key '%s'\n", key)
	}
}

func (c *Controllers) ClearStoreController() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if err := c.StoreInstance.Clear(); err != nil {
			http.Error(w, fmt.Sprintf("error clearing store: %v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Successfully cleared the store\n")

	}
}

func (c *Controllers) BatchPutController() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, constants.MaxBodySize)

		body, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Error("Error reading batch-put request body:", err)
			http.Error(w, "failed to read request body", http.StatusBadRequest)
			return
		}

		var batch map[string]string
		if err := json.Unmarshal(body, &batch); err != nil {
			logger.Error("Error parsing batch-put JSON body:", err)
			http.Error(w, "failed to parse JSON body", http.StatusBadRequest)
			return
		}

		logger.Debug("\n\nParsed batch put request payload:", batch, "\n\n")

		res, err := c.StoreInstance.BatchPut(batch)

		if err != nil {
			logger.Error("Error in BatchPut operation:", err)
			http.Error(w, fmt.Sprintf("error during batch put: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(res); err != nil {
			logger.Error("Error encoding batch put response:", err)
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
			return
		}

	}
}

func (c *Controllers) BatchGetController() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, constants.MaxBodySize)

		body, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Error("Error reading batch-get request body:", err)
			http.Error(w, "failed to read request body", http.StatusBadRequest)
			return
		}

		var payload utils.BatchGetRequestBody
		if err := json.Unmarshal(body, &payload); err != nil {
			logger.Error("Error parsing batch-get JSON body:", err)
			http.Error(w, "failed to parse JSON body", http.StatusBadRequest)
			return
		}

		keys := local.DeleteDuplicateKeys(payload.Keys)

		if len(keys) == 0 {
			http.Error(w, "empty keys list", http.StatusBadRequest)
			return
		}
		if len(keys) > 10000 {
			http.Error(w, "too many keys", http.StatusBadRequest)
			return
		}

		result, err := c.StoreInstance.BatchGet(keys)
		if err != nil {
			logger.Error("Error in BatchGet operation:", err)
			http.Error(w, fmt.Sprintf("error during batch get: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(result); err != nil {
			logger.Error("Error encoding response:", err)
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}

func (c *Controllers) BatchDeleteController() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, constants.MaxBodySize)

		body, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Error("Error reading batch-delete request body:", err)
			http.Error(w, "failed to read request body", http.StatusBadRequest)
			return
		}

		var payload utils.BatchDeleteRequestBody
		if err := json.Unmarshal(body, &payload); err != nil {
			logger.Error("Error parsing batch-delete JSON body:", err)
			http.Error(w, "failed to parse JSON body", http.StatusBadRequest)
			return
		}

		keys := local.DeleteDuplicateKeys(payload.Keys)

		if len(keys) == 0 {
			http.Error(w, "empty keys list", http.StatusBadRequest)
			return
		}

		res, err := c.StoreInstance.BatchDelete(keys)
		if err != nil {
			logger.Error("Error in BatchDelete operation:", err)
			http.Error(w, fmt.Sprintf("error during batch delete: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(res); err != nil {
			logger.Error("Error encoding batch delete response:", err)
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}
