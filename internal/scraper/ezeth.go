package scraper

import "example.com/rates/v2/internal/model"

// ScrapeRenzo returns hardcoded Renzo rates (ETH, WETH, wstETH -> pzETH at 2.86).
func ScrapeRenzo() []model.Rate {
	return []model.Rate{
		{
			InputSymbol: "ETH",
			OutputToken: "pzETH",
			ProjectName: "Renzo",
			PoolName:    "Renzo",
			APY:         2.86,
			ProjectLink: "https://app.renzoprotocol.com/",
			Points:      "EigenLayer Points: 1; Renzo Points: 1",
		},
		{
			InputSymbol: "WETH",
			OutputToken: "pzETH",
			ProjectName: "Renzo",
			PoolName:    "Renzo",
			APY:         2.86,
			ProjectLink: "https://app.renzoprotocol.com/",
			Points:      "EigenLayer Points: 1; Renzo Points: 1",
		},
		{
			InputSymbol: "wstETH",
			OutputToken: "pzETH",
			ProjectName: "Renzo",
			PoolName:    "Renzo",
			APY:         2.86,
			ProjectLink: "https://app.renzoprotocol.com/",
			Points:      "EigenLayer Points: 1; Renzo Points: 1",
		},
	}
}
