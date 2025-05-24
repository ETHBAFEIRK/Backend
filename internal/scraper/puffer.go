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
	PufferID         = "puffer:ETH:Puffer Staking"
	PufferProject    = "Puffer"
	PufferInput      = "ETH"
	PufferPool       = "Puffer Staking"
	PufferProjectURL = "https://app.puffer.fi/"
	PufferAPI        = "https://api.puffer.fi/backend-for-frontend/tvl/all"
)

type PufferAPIResponse struct {
	APY string `json:"apy"`
}

func ScrapePuffer() (model.Rate, error) {
	log.Println("[scraper] Scraping Puffer...")
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", PufferAPI, nil)
	if err != nil {
		log.Printf("[scraper] Puffer: request error: %v", err)
		return model.Rate{}, err
	}
	req.Header.Set("accept", "application/json, text/plain, */*")
	req.Header.Set("origin", "https://app.puffer.fi")
	req.Header.Set("referer", "https://app.puffer.fi/")
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[scraper] Puffer: HTTP error: %v", err)
		return model.Rate{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Printf("[scraper] Puffer: unexpected status: %d", resp.StatusCode)
		return model.Rate{}, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[scraper] Puffer: read error: %v", err)
		return model.Rate{}, err
	}
	var apiResp PufferAPIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		log.Printf("[scraper] Puffer: unmarshal error: %v", err)
		return model.Rate{}, err
	}
	if apiResp.APY == "" {
		log.Printf("[scraper] Puffer: apy not found in response")
		return model.Rate{}, errors.New("apy not found in response")
	}
	var apy float64
	if _, err := fmt.Sscanf(apiResp.APY, "%f", &apy); err != nil {
		log.Printf("[scraper] Puffer: apy parse error: %v", err)
		return model.Rate{}, err
	}
	log.Printf("[scraper] Puffer: APY scraped: %.4f", apy)
	return model.Rate{
		InputSymbol: PufferInput,
		OutputToken: "xPufETH",
		ProjectName: PufferProject,
		PoolName:    PufferPool,
		APY:         apy,
		ProjectLink: PufferProjectURL,
		Points:      "Puffer Points: 1; Zircuit Points: 1; EigenLayer Points: 1",
	}, nil
}
