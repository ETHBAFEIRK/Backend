package scraper

import "example.com/rates/v2/internal/model"

const (
	EzETHProject    = "renzo"
	EzETHPool       = "renzo"
	EzETHProjectURL = "https://app.renzoprotocol.com"
	EzETHAPY        = 3.77
)

// ScrapeRenzo returns hardcoded Renzo rates (ETH, WETH, wstETH -> pzETH at 2.86).
func ScrapeRenzo() []model.Rate {

	return []model.Rate{
		{
			InputSymbol: "ETH",
			OutputToken: "ezETH",
			ProjectName: EzETHProject,
			PoolName:    EzETHPool,
			APY:         EzETHAPY,
			ProjectLink: EzETHProjectURL,
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
		},
		{
			InputSymbol: "ETH",
			OutputToken: "pzETH",
			ProjectName: "Renzo",
			PoolName:    "Renzo",
			APY:         2.86,
			ProjectLink: "https://app.renzoprotocol.com/",
			Points:      "Zircuit Points: 2; ",
		},
		{
			InputSymbol: "WETH",
			OutputToken: "pzETH",
			ProjectName: "Renzo",
			PoolName:    "Renzo",
			APY:         2.86,
			ProjectLink: "https://app.renzoprotocol.com/",
			Points:      "Zircuit Points: 2; ",
		},
		{
			InputSymbol: "wstETH",
			OutputToken: "pzETH",
			ProjectName: "Renzo",
			PoolName:    "Renzo",
			APY:         2.86,
			ProjectLink: "https://app.renzoprotocol.com/",
			Points:      "Zircuit Points: 2; ",
		},
	}
}
