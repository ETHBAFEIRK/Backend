package main

import (
	"log"
	"net/http"

	"example.com/rates/v2/internal/app"
	"example.com/rates/v2/internal/handler"
)

func main() {
	scraperManager := app.NewScraperManager("scrape_cache.db")
	handler.InitScraperManager(scraperManager)
	scraperManager.StartBackgroundScraping()

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", handler.Healthz)
	mux.HandleFunc("/rates", handler.Rates)
	mux.HandleFunc("/api/rates", handler.Rates)
	mux.HandleFunc("/", handler.Home)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
