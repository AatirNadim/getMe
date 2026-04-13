package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AatirNadim/getMe/commons"
	gosdk "github.com/AatirNadim/getMe/sdks/goSdk"
)

type HttpProxy struct {
	Client *gosdk.GetMeClient
}

func (h *HttpProxy) GetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "missing key parameter", http.StatusBadRequest)
		return
	}

	val, err := h.Client.Get(key)
	if err != nil {
		http.Error(w, fmt.Sprintf("error: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", val)
}

func (h *HttpProxy) PutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "failed to parse JSON body", http.StatusBadRequest)
		return
	}

	if req.Key == "" || req.Value == "" {
		http.Error(w, "missing key or value parameter", http.StatusBadRequest)
		return
	}

	if err := h.Client.Put(req.Key, req.Value); err != nil {
		http.Error(w, fmt.Sprintf("error: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Successfully set value for key '%s'\n", req.Key)
}

func (h *HttpProxy) DeleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	key := r.URL.Query().Get("key")
	if key == "" {
		http.Error(w, "missing key parameter", http.StatusBadRequest)
		return
	}

	if err := h.Client.Delete(key); err != nil {
		http.Error(w, fmt.Sprintf("error: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Successfully deleted key '%s'\n", key)
}

func (h *HttpProxy) BatchGetHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload commons.BatchGetRequestBody
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "failed to parse JSON body", http.StatusBadRequest)
		return
	}

	if len(payload.Keys) == 0 {
		http.Error(w, "empty keys list", http.StatusBadRequest)
		return
	}

	res, err := h.Client.BatchGetForPayload(payload)
	if err != nil {
		http.Error(w, fmt.Sprintf("error: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *HttpProxy) BatchPutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload []commons.KeyValue
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "failed to parse JSON body", http.StatusBadRequest)
		return
	}

	if len(payload) == 0 {
		http.Error(w, "empty payload", http.StatusBadRequest)
		return
	}

	// commented out for now, keeping in mind the sheer potential volume of payload

	// for _, kv := range payload {
	// 	if kv.Key == "" || kv.Value == "" {
	// 		http.Error(w, "missing key or value in one or more payload items", http.StatusBadRequest)
	// 		return
	// 	}
	// }

	res, err := h.Client.BatchPutForPayload(payload)
	if err != nil {
		http.Error(w, fmt.Sprintf("error: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *HttpProxy) BatchDeleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var payload commons.BatchDeleteRequestBody
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "failed to parse JSON body", http.StatusBadRequest)
		return
	}

	if len(payload.Keys) == 0 {
		http.Error(w, "empty keys list", http.StatusBadRequest)
		return
	}

	res, err := h.Client.BatchDeleteForPayload(payload)
	if err != nil {
		http.Error(w, fmt.Sprintf("error: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *HttpProxy) ClearStoreHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodDelete {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := h.Client.ClearStore(); err != nil {
		http.Error(w, fmt.Sprintf("error: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Successfully cleared the store\n")
}
