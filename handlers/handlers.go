package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"url-shortener/storage"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	tmpl := `
	<!DOCTYPE html>
	<html>
	<head>
		<title>URL Shortener</title>
	</head>
	<body>
		<h1>URL shortener</h1>
		<form action="/shorten" method="post">
			<input type="url" name="url" placeholder="Введите URL" required style="width: 300px; padding: 8px;">
			<button type="submit" style="padding: 8px 16px;">Shorten</button>
		</form>
	</body>
	</html>
	`

	t, _ := template.New("home").Parse(tmpl)
	t.Execute(w, nil)
}

func ShortenHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	longURL := r.FormValue("url")
	if longURL == "" {
		http.Error(w, "URL is required", http.StatusBadRequest)
		return
	}

	shortCode := storage.Storage.Save(longURL)
	shortURL := fmt.Sprintf("http://localhost:5000/r/%s", shortCode)

	response := map[string]string{
		"original_url": longURL,
		"short_url":    shortURL,
		"short_code":   shortCode,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func RedirectHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/r/")
	parts := strings.Split(path, "/")
	if len(parts) == 0 {
		http.NotFound(w, r)
		return
	}

	shortCode := parts[0]

	longURL, exists := storage.Storage.Get(shortCode)
	if !exists {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, longURL, http.StatusFound)
}
