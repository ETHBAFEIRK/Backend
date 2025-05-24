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

func ScrapeEzETH() []model.Rate {
	log.Println("[scraper] Scraping (hardcoded) ezETH rates...")
	points := "2x Zircuit Points"
	return []model.Rate{
		{
			InputSymbol: "ETH",
			OutputToken: "ezETH",
			ProjectName: EzETHProject,
			PoolName:    EzETHPool,
			APY:         EzETHAPY,
			ProjectLink: EzETHProjectURL,
			Points:      points,
		},
		{
			InputSymbol: "stETH",
			OutputToken: "ezETH",
			ProjectName: EzETHProject,
			PoolName:    EzETHPool,
			APY:         EzETHAPY,
			ProjectLink: EzETHProjectURL,
			Points:      points,
		},
	}
}
