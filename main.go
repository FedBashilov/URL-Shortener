package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"url-shortener/handlers"
	"url-shortener/storage"
)

func main() {
	connString := "postgres://user:password@localhost:5432/url_shortener?sslmode=disable"
	postgresStorage, err := storage.NewPostgresStorage(connString)
	if err != nil {
		log.Fatalf("❌ Failed to connect to PostgreSQL: %v", err)
	}

	storage.Postgres = postgresStorage
	storage.Memory = storage.Storage

	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/shorten", handlers.ShortenHandler)
	http.HandleFunc("/r/", handlers.RedirectHandler)
	http.HandleFunc("/stats/", handlers.StatsHandler)

	fmt.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

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
