package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"example.com/rates/v2/internal/app"
	"example.com/rates/v2/internal/model"
	"example.com/rates/v2/internal/scraper"
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
	id := scraper.PufferID
	project := scraper.PufferProject
	inputSymbol := scraper.PufferInput
	poolName := scraper.PufferPool

	apy, lastScrape, found := scraperManager.GetCachedRate(id)
	var rates []model.Rate

	if found {
		rate := model.Rate{
			InputSymbol: inputSymbol,
			ProjectName: project,
			PoolName:    poolName,
			APY:         apy,
			ProjectLink: scraper.PufferProjectURL,
		}
		rates = append(rates, rate)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rates)
}
