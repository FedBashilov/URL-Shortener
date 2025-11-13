package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"url-shortener/storage"
)

func StatsHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/stats/")
	parts := strings.Split(path, "/")
	if len(parts) == 0 {
		http.Error(w, "Short code required", http.StatusBadRequest)
		return
	}

	shortCode := parts[0]

	stats, err := storage.Postgres.GetStats(shortCode)
	if err != nil {
		http.Error(w, "Failed to get stats", http.StatusInternalServerError)
		return
	}

	if stats == nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
