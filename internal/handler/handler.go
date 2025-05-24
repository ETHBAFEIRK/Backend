package handler

import (
	"encoding/json"
	"example.com/rates/v2/internal/app"
	"example.com/rates/v2/internal/model"
	"example.com/rates/v2/internal/scraper"
	"fmt"
	"net/http"
)

var scraperManager *app.ScraperManager

func InitScraperManager(sm *app.ScraperManager) {
	scraperManager = sm
}

func Healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "ok")
}

func Home(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Welcome to the Rates API v2")
}

func Rates(w http.ResponseWriter, r *http.Request) {
	rates, err := scraperManager.GetAllRates()
	if err != nil {
		http.Error(w, "Failed to fetch rates", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rates)
}
