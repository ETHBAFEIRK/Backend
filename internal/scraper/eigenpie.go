package scraper

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"example.com/rates/v2/internal/model"
)

const (
	EigenpieProject    = "Eigenpie"
	EigenpieProjectURL = "https://app.magpiexyz.io/eigenpie"
	EigenpieAPI        = "https://dev.api.magpiexyz.io/poolsnapshot/getEigenpiePoolSnapshot"
)

type eigenpieAPIResponse struct {
	Data struct {
		Snapshot struct {
			ChainData []struct {
				ChainId int `json:"chainId"`
				Data    struct {
					IsExpired bool `json:"isExpired"`
					Data      []struct {
						PoolId           int    `json:"poolId"`
						PoolName         string `json:"poolName"`
						StakeTokenInfo   struct {
							Symbol string `json:"symbol"`
						} `json:"stakeTokenInfo"`
						ReceiptTokenInfo struct {
							Symbol string `json:"symbol"`
						} `json:"receiptTokenInfo"`
						AprInfo struct {
							FormatValue string `json:"formatValue"`
						} `json:"aprInfo"`
					} `json:"data"`
				} `json:"data"`
			} `json:"chainData"`
		} `json:"snapshot"`
	} `json:"data"`
}

// ScrapeEigenpie fetches Eigenpie pool rates for mainnet (chainId=1).
func ScrapeEigenpie() ([]model.Rate, error) {
	log.Println("[scraper] Scraping Eigenpie...")
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", EigenpieAPI, nil)
	if err != nil {
		log.Printf("[scraper] Eigenpie: request error: %v", err)
		return nil, err
	}
	req.Header.Set("accept", "application/json, text/plain, */*")
	req.Header.Set("origin", "https://app.magpiexyz.io")
	req.Header.Set("referer", "https://app.magpiexyz.io/eigenpie")
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[scraper] Eigenpie: HTTP error: %v", err)
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Printf("[scraper] Eigenpie: unexpected status: %d", resp.StatusCode)
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[scraper] Eigenpie: read error: %v", err)
		return nil, err
	}
	var apiResp eigenpieAPIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		log.Printf("[scraper] Eigenpie: unmarshal error: %v", err)
		return nil, err
	}

	var rates []model.Rate
	allowedOutputs := map[string]struct{}{
		"mstETH": {},
		"weETH":  {},
		"wstETH": {},
		"egETH":  {},
	}
	for _, chain := range apiResp.Data.Snapshot.ChainData {
		if chain.ChainId != 1 {
			continue
		}
		for _, pool := range chain.Data.Data {
			if _, ok := allowedOutputs[pool.ReceiptTokenInfo.Symbol]; !ok {
				continue
			}
			apy, err := strconv.ParseFloat(pool.AprInfo.FormatValue, 64)
			if err != nil {
				log.Printf("[scraper] Eigenpie: failed to parse APY for %s: %v", pool.PoolName, err)
				continue
			}
			rates = append(rates, model.Rate{
				InputSymbol: pool.StakeTokenInfo.Symbol,
				OutputToken: pool.ReceiptTokenInfo.Symbol,
				ProjectName: EigenpieProject,
				PoolName:    pool.PoolName,
				APY:         apy,
				ProjectLink: EigenpieProjectURL,
				Points:      "Eigenpie Points: 1; EigenLayer Points: 1",
			})
		}
	}
	return rates, nil
}
