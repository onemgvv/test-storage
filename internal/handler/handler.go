package handler

import (
	"encoding/json"
	"net/http"
	"test-storage/pkg/config"
	"test-storage/pkg/storage"
	"time"
)

type SetterBody struct{
	Key string `mapstructure:"key" json:"key"`
	Value any `mapstructure:"value" json:"value"`
	TTL time.Duration `mapstructure:"ttl" json:"ttl"`
}

type Handler struct{
	storage *storage.Storage
	config *config.Config
}

func NewHandler(storage *storage.Storage, cfg *config.Config) *Handler {
	return &Handler{
		storage: storage,
		config: cfg,
	}
}

func (h *Handler) Init() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/storage/set", h.storageSetter)
	mux.HandleFunc("/api/storage/get", h.storageGetter)

	return mux
}

func (h *Handler) storageSetter(w http.ResponseWriter, r *http.Request) {
	var body SetterBody

	w.Header().Set("Content-Type", "application/json")

	if r.Method == http.MethodPost {
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return 
		}

		if body.TTL == 0 {
			body.TTL = h.config.Storage.DefaultTTL
		}

		h.storage.Add(body.Key, body.Value, body.TTL * time.Minute)

		data, _ := json.Marshal(map[string]int{ "success": 1 })
		w.Write(data)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) storageGetter(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {

		w.Header().Set("Content-Type", "application/json")

		key := r.URL.Query().Get("key")
		data := h.storage.Get(key)
		if data.Value == nil {
			w.WriteHeader(http.StatusNotFound)
			data, _ := json.Marshal(map[string]any{ "statusCode": 404, "message": "Data expired" })
			w.Write([]byte(data))
			return
		}

		resBody, _ := json.Marshal(map[string]any{"success": 1, "data": data.Value})
		w.Write(resBody)
		return 
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}