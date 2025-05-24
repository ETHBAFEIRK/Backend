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
	// Only one scraper for now: Puffer
	id := scraper.PufferID
	project := scraper.PufferProject
	inputSymbol := scraper.PufferInput
	poolName := scraper.PufferPool

	// Check cache
	apy, lastScrape, found := scraperManager.GetCachedRate(id)
	shouldScrape := true
	if found && time.Since(lastScrape) < 10*time.Minute {
		shouldScrape = false
	}

	var rate model.Rate
	var err error
	if shouldScrape {
		rate, err = scraper.ScrapePuffer()
		if err == nil {
			apy = rate.APY
			_ = scraperManager.SetCachedRate(id, project, inputSymbol, poolName, apy, time.Now())
		} else if found {
			// fallback to cached
			rate = model.Rate{
				InputSymbol: inputSymbol,
				ProjectName: project,
				PoolName:    poolName,
				APY:         apy,
				ProjectLink: scraper.PufferProjectURL,
			}
		}
	} else {
		rate = model.Rate{
			InputSymbol: inputSymbol,
			ProjectName: project,
			PoolName:    poolName,
			APY:         apy,
			ProjectLink: scraper.PufferProjectURL,
		}
	}

	rates := []model.Rate{}
	if rate.ProjectName != "" {
		rates = append(rates, rate)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rates)
}
