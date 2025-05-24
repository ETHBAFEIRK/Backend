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
	InceptionID         = "inception:wstETH:inception"
	InceptionProject    = "inception"
	InceptionInput      = "wstETH"
	InceptionOutput     = "inwstETH"
	InceptionPool       = "inception"
	InceptionProjectURL = "https://www.inceptionlrt.com/"
	InceptionAPI        = "https://bff.prod.inceptionlrt.com/stakingwatch/apr_7d"
)

type inceptionAPRItem struct {
	PoolID  string `json:"poolId"`
	TokenID string `json:"tokenId"`
	Value   string `json:"value"`
}

func ScrapeInception() (model.Rate, error) {
	log.Println("[scraper] Scraping Inception...")
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", InceptionAPI, nil)
	if err != nil {
		log.Printf("[scraper] Inception: request error: %v", err)
		return model.Rate{}, err
	}
	req.Header.Set("accept", "*/*")
	req.Header.Set("origin", "https://www.inceptionlrt.com")
	req.Header.Set("referer", "https://www.inceptionlrt.com/")
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[scraper] Inception: HTTP error: %v", err)
		return model.Rate{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Printf("[scraper] Inception: unexpected status: %d", resp.StatusCode)
		return model.Rate{}, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[scraper] Inception: read error: %v", err)
		return model.Rate{}, err
	}
	var items []inceptionAPRItem
	if err := json.Unmarshal(body, &items); err != nil {
		log.Printf("[scraper] Inception: unmarshal error: %v", err)
		return model.Rate{}, err
	}
	var valueStr string
	for _, item := range items {
		if item.PoolID == "lido" && item.TokenID == "steth" {
			valueStr = item.Value
			break
		}
	}
	if valueStr == "" {
		log.Printf("[scraper] Inception: lido/steth value not found")
		return model.Rate{}, errors.New("lido/steth value not found")
	}
	var apy float64
	if _, err := fmt.Sscanf(valueStr, "%f", &apy); err != nil {
		log.Printf("[scraper] Inception: apy parse error: %v", err)
		return model.Rate{}, err
	}
	log.Printf("[scraper] Inception: APY scraped: %.4f", apy)
	return model.Rate{
		InputSymbol: InceptionInput,
		OutputToken: InceptionOutput,
		ProjectName: InceptionProject,
		PoolName:    InceptionPool,
		APY:         apy,
		ProjectLink: InceptionProjectURL,
		Points:      "Zircuit Points: 2; Mellow Points: 2; InceptionLRT Totems: 3; Symbiotic Points: 1",
		OutputKind:  "restake",
	}, nil
}
