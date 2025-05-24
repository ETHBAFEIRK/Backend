package scraper

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"example.com/rates/v2/internal/model"
)

const (
	LidoID         = "lido:ETH:Lido"
	LidoProject    = "Lido"
	LidoInput      = "ETH"
	LidoOutput     = "stETH"
	LidoPool       = "Lido"
	LidoProjectURL = "https://stake.lido.fi/"
	LidoAPI        = "https://eth-api.lido.fi/v1/protocol/steth/apr/sma"
)

type lidoAPIResponse struct {
	Data struct {
		SmaApr float64 `json:"smaApr"`
	} `json:"data"`
}

func ScrapeLido() (model.Rate, error) {
	log.Println("[scraper] Scraping Lido...")
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", LidoAPI, nil)
	if err != nil {
		log.Printf("[scraper] Lido: request error: %v", err)
		return model.Rate{}, err
	}
	req.Header.Set("accept", "application/json, text/plain, */*")
	req.Header.Set("origin", "https://stake.lido.fi")
	req.Header.Set("referer", "https://stake.lido.fi/")
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[scraper] Lido: HTTP error: %v", err)
		return model.Rate{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Printf("[scraper] Lido: unexpected status: %d", resp.StatusCode)
		return model.Rate{}, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[scraper] Lido: read error: %v", err)
		return model.Rate{}, err
	}
	var apiResp lidoAPIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		log.Printf("[scraper] Lido: unmarshal error: %v", err)
		return model.Rate{}, err
	}
	apy := apiResp.Data.SmaApr
	if apy == 0 {
		log.Printf("[scraper] Lido: smaApr not found in response")
		return model.Rate{}, errors.New("smaApr not found in response")
	}
	log.Printf("[scraper] Lido: APY scraped: %.4f", apy)
	// Return both stETH and wstETH as output tokens
	returns := []model.Rate{
		{
			InputSymbol: LidoInput,
			OutputToken: LidoOutput,
			ProjectName: LidoProject,
			PoolName:    LidoPool,
			APY:         apy,
			ProjectLink: LidoProjectURL,
			Points:      "",
		},
		{
			InputSymbol: LidoInput,
			OutputToken: "wstETH",
			ProjectName: LidoProject,
			PoolName:    LidoPool,
			APY:         apy,
			ProjectLink: LidoProjectURL,
			Points:      "",
		},
	}
	// For backward compatibility, return the stETH record as the main result, but let caller use both if needed
	return returns[0], nil
}
