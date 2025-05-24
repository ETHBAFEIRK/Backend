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
	KelpProject    = "kelp"
	KelpPool       = "kelp"
	KelpProjectURL = "https://universe.kelpdao.xyz/"
	KelpAPI        = "https://universe.kelpdao.xyz/rseth/totalApy"
)

type kelpAPIResponse struct {
	TotalAPY float64 `json:"totalAPY"`
}

func ScrapeKelp() ([]model.Rate, error) {
	log.Println("[scraper] Scraping Kelp...")
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", KelpAPI, nil)
	if err != nil {
		log.Printf("[scraper] Kelp: request error: %v", err)
		return nil, err
	}
	req.Header.Set("accept", "application/json, text/plain, */*")
	req.Header.Set("origin", "https://universe.kelpdao.xyz")
	req.Header.Set("referer", "https://universe.kelpdao.xyz/")
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[scraper] Kelp: HTTP error: %v", err)
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Printf("[scraper] Kelp: unexpected status: %d", resp.StatusCode)
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[scraper] Kelp: read error: %v", err)
		return nil, err
	}
	var apiResp kelpAPIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		log.Printf("[scraper] Kelp: unmarshal error: %v", err)
		return nil, err
	}
	if apiResp.TotalAPY == 0 {
		log.Printf("[scraper] Kelp: totalAPY not found in response")
		return nil, errors.New("totalAPY not found in response")
	}
	log.Printf("[scraper] Kelp: APY scraped: %.4f", apiResp.TotalAPY)
	inputs := []string{"ETHx", "ETH", "stETH"}
	var rates []model.Rate
	for _, input := range inputs {
		rates = append(rates, model.Rate{
			InputSymbol: input,
			OutputToken: "rsETH",
			ProjectName: KelpProject,
			PoolName:    KelpPool,
			APY:         apiResp.TotalAPY,
			ProjectLink: KelpProjectURL,
			Points:      "Zircuit Points: 2; Kelp Miles: 2; EigenLayer Points: 1",
		})
	}
	return rates, nil
}
