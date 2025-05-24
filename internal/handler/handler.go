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
	var rates []model.Rate

	// Puffer
	{
		id := scraper.PufferID
		project := scraper.PufferProject
		inputSymbol := scraper.PufferInput
		poolName := scraper.PufferPool

		apy, _, found := scraperManager.GetCachedRate(id)
		if found {
			rate := model.Rate{
				InputSymbol: inputSymbol,
				OutputToken: "xpufETH",
				ProjectName: project,
				PoolName:    poolName,
				APY:         apy,
				ProjectLink: scraper.PufferProjectURL,
				Points:      "Puffer Points: 1; Zircuit Points: 1; EigenLayer Points: 1",
			}
			rates = append(rates, rate)
		}
	}

	// Inception
	{
		id := scraper.InceptionID
		project := scraper.InceptionProject
		inputSymbol := scraper.InceptionInput
		poolName := scraper.InceptionPool

		apy, _, found := scraperManager.GetCachedRate(id)
		if found {
			rate := model.Rate{
				InputSymbol: inputSymbol,
				OutputToken: "inwstETH",
				ProjectName: project,
				PoolName:    poolName,
				APY:         apy,
				ProjectLink: scraper.InceptionProjectURL,
				Points:      "Zircuit Points: 2; Mellow Points: 2; InceptionLRT Totems: 3; Symbiotic Points: 1",
			}
			rates = append(rates, rate)
		}
	}

	// ezETH (renzo)
	{
		project := "renzo"
		poolName := "renzo"
		projectLink := "https://app.renzoprotocol.com"

		for _, inputSymbol := range []string{"ETH", "stETH"} {
			id := project + ":" + inputSymbol + ":" + poolName
			cachedAPY, _, found := scraperManager.GetCachedRate(id)
			if found {
				rate := model.Rate{
					InputSymbol: inputSymbol,
					OutputToken: "ezETH",
					ProjectName: project,
					PoolName:    poolName,
					APY:         cachedAPY,
					ProjectLink: projectLink,
					Points:      "",
				}
				rates = append(rates, rate)
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rates)
}
