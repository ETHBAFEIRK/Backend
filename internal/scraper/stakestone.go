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
	StakestoneID         = "stakestone:ETH:Stakestone"
	StakestoneProject    = "Stakestone"
	StakestoneInput      = "ETH"
	StakestoneOutput     = "STONE"
	StakestonePool       = "Stakestone"
	StakestoneProjectURL = "https://stakestone.io/"
	StakestoneAPI        = "https://eth-api.lido.fi/v1/protocol/steth/apr/sma"
)

type stakestoneAPIResponse struct {
	Data struct {
		SmaApr float64 `json:"smaApr"`
	} `json:"data"`
}

func ScrapeStakestone() (model.Rate, error) {
	log.Println("[scraper] Scraping Stakestone...")
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", StakestoneAPI, nil)
	if err != nil {
		log.Printf("[scraper] Stakestone: request error: %v", err)
		return model.Rate{}, err
	}
	req.Header.Set("accept", "application/json, text/plain, */*")
	req.Header.Set("origin", "https://stakestone.io")
	req.Header.Set("referer", "https://stakestone.io/")
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[scraper] Stakestone: HTTP error: %v", err)
		return model.Rate{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Printf("[scraper] Stakestone: unexpected status: %d", resp.StatusCode)
		return model.Rate{}, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[scraper] Stakestone: read error: %v", err)
		return model.Rate{}, err
	}
	var apiResp stakestoneAPIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		log.Printf("[scraper] Stakestone: unmarshal error: %v", err)
		return model.Rate{}, err
	}
	apy := apiResp.Data.SmaApr
	if apy == 0 {
		log.Printf("[scraper] Stakestone: smaApr not found in response")
		return model.Rate{}, errors.New("smaApr not found in response")
	}
	log.Printf("[scraper] Stakestone: APY scraped: %.4f", apy)
	return model.Rate{
		InputSymbol: StakestoneInput,
		OutputToken: StakestoneOutput,
		ProjectName: StakestoneProject,
		PoolName:    StakestonePool,
		APY:         apy,
		ProjectLink: StakestoneProjectURL,
		Points:      "Zircuit Points: 2; Stone Points: 2",
	}, nil
}
