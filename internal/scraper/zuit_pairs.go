package scraper

type ExchangePair struct {
	ID                int
	Pool              string
	Token1            string
	Token2            string
	Type              string
	Liquidity         float64
	LiquidityFormatted string
	APR               float64
	APRFormatted      string
	ExchangeID        string
}

func ScrapeZuitPairs() []ExchangePair {
	return []ExchangePair{
		{
			ID:                 1,
			Pool:               "ETH/ZRC",
			Token1:             "ETH",
			Token2:             "ZRC",
			Type:               "Base",
			Liquidity:          470016,
			LiquidityFormatted: "$470,016",
			APR:                2.8,
			APRFormatted:       "2.8%",
			ExchangeID:         "zuit",
		},
		{
			ID:                 2,
			Pool:               "USDC.e/USDT",
			Token1:             "USDC.e",
			Token2:             "USDT",
			Type:               "Stable",
			Liquidity:          194325,
			LiquidityFormatted: "$194,325",
			APR:                0.1,
			APRFormatted:       "0.1%",
			ExchangeID:         "zuit",
		},
		{
			ID:                 3,
			Pool:               "WBTC.e/ETH",
			Token1:             "WBTC.e",
			Token2:             "ETH",
			Type:               "Oasis",
			Liquidity:          100045,
			LiquidityFormatted: "$100,045",
			APR:                0.6,
			APRFormatted:       "0.6%",
			ExchangeID:         "zuit",
		},
	}
}
