package handler

import (
	"encoding/json"
	"example.com/rates/v2/internal/app"
	"fmt"
	"net/http"
)

var tokenIcons = map[string]string{
	"ETH":     "https://raw.githubusercontent.com/ethereum/ethereum-org-website/dev/src/data/networks/icons/eth.svg",
	"WETH":    "https://raw.githubusercontent.com/trustwallet/assets/master/blockchains/ethereum/assets/0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2/logo.png",
	"stETH":   "https://raw.githubusercontent.com/LidoFi/lido-dao/main/logo.svg",
	"wstETH":  "https://raw.githubusercontent.com/LidoFi/lido-dao/main/logo.svg",
	"ezETH":   "https://app.renzoprotocol.com/icons/ezeth.svg",
	"pzETH":   "https://app.renzoprotocol.com/icons/pzeth.svg",
	"STONE":   "https://stakestone.io/logo.svg",
	"xPufETH": "https://app.puffer.fi/icons/xpufeth.svg",
	"mstETH":  "https://app.magpiexyz.io/icons/msteth.svg",
	"weETH":   "https://app.magpiexyz.io/icons/weethe.svg",
	"egETH":   "https://app.magpiexyz.io/icons/egeth.svg",
	"inwstETH": "https://www.inceptionlrt.com/icons/inwsteth.svg",
}

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
	// Add icons to each rate
	for i := range rates {
		rates[i].FromIcon = tokenIcons[rates[i].InputSymbol]
		rates[i].ToIcon = tokenIcons[rates[i].OutputToken]
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rates)
}
