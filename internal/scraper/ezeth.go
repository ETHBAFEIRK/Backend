package scraper

import (
	"example.com/rates/v2/internal/model"
	"log"
)

const (
	EzETHProject    = "renzo"
	EzETHPool       = "renzo"
	EzETHProjectURL = "https://app.renzoprotocol.com"
	EzETHAPY        = 3.77
)

// ScrapeRenzo returns hardcoded Renzo rates (ETH, WETH, wstETH -> pzETH at 2.86).
func ScrapeRenzo() []model.Rate {
	log.Println("[scraper] Scraping Inception...")

	return []model.Rate{
		{
			InputSymbol: "ETH",
			OutputToken: "ezETH",
			ProjectName: EzETHProject,
			PoolName:    EzETHPool,
			APY:         EzETHAPY,
			ProjectLink: EzETHProjectURL,
			OutputKind:  "stake",
			Points:      "Zircuit Points: 2; ",
		},
		{
			InputSymbol: "stETH",
			OutputToken: "ezETH",
			ProjectName: EzETHProject,
			PoolName:    EzETHPool,
			APY:         EzETHAPY,
			ProjectLink: EzETHProjectURL,
			Points:      "Zircuit Points: 2; ",
			OutputKind:  "stake",
		},
		{
			InputSymbol: "ETH",
			OutputToken: "pzETH",
			ProjectName: "Renzo",
			PoolName:    "Renzo",
			APY:         2.86,
			ProjectLink: "https://app.renzoprotocol.com/",
			Points:      "Zircuit Points: 2; ",
			OutputKind:  "stake",
		},
		{
			InputSymbol: "WETH",
			OutputToken: "pzETH",
			ProjectName: "Renzo",
			PoolName:    "Renzo",
			APY:         2.86,
			ProjectLink: "https://app.renzoprotocol.com/",
			Points:      "Zircuit Points: 2; ",
			OutputKind:  "stake",
		},
		{
			InputSymbol: "wstETH",
			OutputToken: "pzETH",
			ProjectName: "Renzo",
			PoolName:    "Renzo",
			APY:         2.86,
			ProjectLink: "https://app.renzoprotocol.com/",
			Points:      "Zircuit Points: 2; ",
			OutputKind:  "stake",
		},
	}
}
