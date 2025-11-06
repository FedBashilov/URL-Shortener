package main

import (
	"fmt"
	"log"
	"net/http"
	"url-shortener/handlers"
)

func main() {
	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/shorten", handlers.ShortenHandler)
	http.HandleFunc("/r/", handlers.RedirectHandler)

	fmt.Println("Server starting on :5000...")
	log.Fatal(http.ListenAndServe(":5000", nil))
}
