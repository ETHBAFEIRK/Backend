package handler

import (
	"encoding/json"
	"example.com/rates/v2/internal/app"
	"fmt"
	"net/http"
)

var tokenIcons = map[string]string{
	"ETH":      "https://static.zircuit.com/stake/app-dbbe0da3d/_next/static/media/eth-logo.9c7e160a.svg",
	"WETH":     "https://static.zircuit.com/stake/app-dbbe0da3d/_next/static/media/eth-logo.9c7e160a.svg:",
	"stETH":    "https://static.zircuit.com/stake/app-dbbe0da3d/_next/static/media/wsteth-logo.70d80504.svg",
	"wstETH":   "https://static.zircuit.com/stake/app-dbbe0da3d/_next/static/media/wsteth-logo.70d80504.svg",
	"ezETH":    "https://static.zircuit.com/stake/app-dbbe0da3d/_next/static/media/ezeth-logo.6809574f.svg",
	"pzETH":    "https://static.zircuit.com/stake/app-dbbe0da3d/_next/static/media/pzeth-logo.c24e47cd.svg",
	"STONE":    "https://static.zircuit.com/stake/app-dbbe0da3d/_next/static/media/stone-logo.ee085a0a.svg",
	"xPufETH":  "https://static.zircuit.com/stake/app-dbbe0da3d/_next/static/media/xpufeth-logo.1bfb3c5a.svg",
	"mstETH":   "https://static.zircuit.com/stake/app-dbbe0da3d/_next/static/media/msteth-logo.70d80504.svg",
	"weETH":    "https://static.zircuit.com/stake/app-dbbe0da3d/_next/static/media/weeth-logo.209d6604.svg",
	"egETH":    "https://static.zircuit.com/stake/app-dbbe0da3d/_next/static/media/egeth-logo.bd7e9357.svg",
	"inwstETH": "https://static.zircuit.com/stake/app-dbbe0da3d/_next/static/media/inwsteths-logo.2406ea8b.svg",
	"rsETH":    "https://static.zircuit.com/stake/app-dbbe0da3d/_next/static/media/rseth-logo.948cf45f.svg",
	"LsETH":    "https://static.zircuit.com/stake/app-dbbe0da3d/_next/static/media/lseth-logo.6dab9ca0.svg",
	"USDC":     "https://static.zircuit.com/stake/app-dbbe0da3d/_next/static/media/usdc-logo.ffc33eac.svg",
	"USDT":     "https://static.zircuit.com/stake/app-dbbe0da3d/_next/static/media/usdt-logo.b1f7c50b.svg",
	"USDe":     "https://static.zircuit.com/stake/app-dbbe0da3d/_next/static/media/usde-logo.c17debe1.svg",
	"FBTC":     "https://static.zircuit.com/stake/app-dbbe0da3d/_next/static/media/fbtc-logo.51f2d301.svg",
	"LBTC":     "https://static.zircuit.com/stake/app-dbbe0da3d/_next/static/media/lbtc-logo.fd05641f.svg",
	"mBTC":     "https://static.zircuit.com/stake/app-dbbe0da3d/_next/static/media/mbtc-logo.3e52220b.svg",
	"pumpBTC":  "https://static.zircuit.com/stake/app-dbbe0da3d/_next/static/media/pumpbtc-logo.af55710d.svg",
	"mswETH":   "https://static.zircuit.com/stake/app-dbbe0da3d/_next/static/media/msweth-logo.e4de1bfd.svg",
	"mwBETH":   "https://static.zircuit.com/stake/app-dbbe0da3d/_next/static/media/mwbeth-logo.857f1a84.svg",
	"mETH":     "https://static.zircuit.com/stake/app-dbbe0da3d/_next/static/media/meth-logo.2d380f4a.svg",
	"rstETH":   "https://static.zircuit.com/stake/app-dbbe0da3d/_next/static/media/rsteth-logo.9cef011b.svg",
	"steakLRT": "https://static.zircuit.com/stake/app-dbbe0da3d/_next/static/media/steaklrt-logo.9cef011b.svg",
	"Re7LRT":   "https://static.zircuit.com/stake/app-dbbe0da3d/_next/static/media/re7lrt-logo.9cef011b.svg",
	"amphrETH": "https://static.zircuit.com/stake/app-dbbe0da3d/_next/static/media/amphreth-logo.9cef011b.svg",
	"rswETH":   "https://static.zircuit.com/stake/app-dbbe0da3d/_next/static/media/rsweth-logo.996367c4.svg",
	"swETH":    "https://static.zircuit.com/stake/app-dbbe0da3d/_next/static/media/sweth-logo.037fa270.svg",
	"weETHs":   "https://static.zircuit.com/stake/app-dbbe0da3d/_next/static/media/weeths-logo.70e83562.svg",
}

// Output token kind mapping: stake or restake
var tokenKind = map[string]string{
	// APY < 3 as of 2025-05-24:
	"xPufETH":  "restake",
	"inwstETH": "restake",
	"pzETH":    "restake",
	"STONE":    "stake",
	"stETH":    "stake",
	"egETH":    "restake",
	// APY >= 3 (for reference, not used in this map):
	// "ezETH": "restake",
	// "rsETH": "restake",
	// "mstETH": "restake",
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

	// Add icons to each rate, but do NOT override OutputKind (use saved kind from DB)
	for i := range rates {
		rates[i].FromIcon = tokenIcons[rates[i].InputSymbol]
		rates[i].ToIcon = tokenIcons[rates[i].OutputToken]
		// OutputKind is already set from DB, do not override
	}

	// --- Begin: filter out isolated pairs ---
	// Build set of destination tokens (those with icons)
	destinations := make(map[string]struct{})
	for token := range tokenIcons {
		destinations[token] = struct{}{}
	}

	// Build adjacency list for the graph
	adj := make(map[string][]string)
	for _, rate := range rates {
		adj[rate.InputSymbol] = append(adj[rate.InputSymbol], rate.OutputToken)
	}

	// Find all tokens that can reach a destination via BFS
	reachable := make(map[string]struct{})
	queue := make([]string, 0)
	// Start from all destination tokens
	for dest := range destinations {
		queue = append(queue, dest)
		reachable[dest] = struct{}{}
	}

	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]
		for token, outs := range adj {
			for _, out := range outs {
				if out == curr {
					if _, seen := reachable[token]; !seen {
						reachable[token] = struct{}{}
						queue = append(queue, token)
					}
				}
			}
		}
	}

	// Filter rates: keep only those where output is reachable from some input to a destination
	filtered := make([]struct {
		InputSymbol string  `json:"input_symbol"`
		OutputToken string  `json:"output_token"`
		ProjectName string  `json:"project_name"`
		PoolName    string  `json:"pool_name"`
		APY         float64 `json:"apy"`
		ProjectLink string  `json:"project_link"`
		Points      string  `json:"points"`
		FromIcon    string  `json:"from_icon"`
		ToIcon      string  `json:"to_icon"`
		OutputKind  string  `json:"output_kind"`
	}, 0, len(rates))
	for _, rate := range rates {
		// If input or output is in reachable set, keep
		if _, ok := reachable[rate.InputSymbol]; ok {
			filtered = append(filtered, rate)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(filtered)
}
