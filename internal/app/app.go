package app

import (
	"database/sql"
	"log"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"example.com/rates/v2/internal/model"
	"example.com/rates/v2/internal/scraper"
)

type ScraperManager struct {
	DB         *sql.DB
	LastScrape map[string]time.Time
	Mutex      sync.Mutex
}

// StartBackgroundScraping launches a goroutine that scrapes Puffer every 10 minutes.
func (sm *ScraperManager) StartBackgroundScraping() {
	go func() {
		for {
			// Puffer
			{
				id := scraper.PufferID
				rate, err := scraper.ScrapePuffer()
				if err == nil {
					_ = sm.SetCachedRateFull(id, rate, time.Now())
				}
			}
			// Inception
			{
				id := scraper.InceptionID
				rate, err := scraper.ScrapeInception()
				if err == nil {
					_ = sm.SetCachedRateFull(id, rate, time.Now())
				}
			}
			// ezETH (hardcoded, just update cache)
			{
				rates := scraper.ScrapeEzETH()
				for _, rate := range rates {
					id := rate.ProjectName + ":" + rate.InputSymbol + ":" + rate.PoolName
					_ = sm.SetCachedRateFull(id, rate, time.Now())
				}
			}
			time.Sleep(10 * time.Minute)
		}
	}()
}

func NewScraperManager(dbPath string) *ScraperManager {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("failed to open sqlite db: %v", err)
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS scrape_cache (
			id TEXT PRIMARY KEY,
			project TEXT,
			input_symbol TEXT,
			output_token TEXT,
			pool_name TEXT,
			apy REAL,
			project_link TEXT,
			points TEXT,
			last_scrape TIMESTAMP
		)
	`)
	if err != nil {
		log.Fatalf("failed to create scrape_cache table: %v", err)
	}
	return &ScraperManager{
		DB:         db,
		LastScrape: make(map[string]time.Time),
	}
}

// Returns (model.Rate, lastScrape, found)
func (sm *ScraperManager) GetCachedRate(id string) (model.Rate, time.Time, bool) {
	row := sm.DB.QueryRow(`SELECT project, input_symbol, output_token, pool_name, apy, project_link, points, last_scrape FROM scrape_cache WHERE id = ?`, id)
	var rate model.Rate
	var lastScrape time.Time
	err := row.Scan(&rate.ProjectName, &rate.InputSymbol, &rate.OutputToken, &rate.PoolName, &rate.APY, &rate.ProjectLink, &rate.Points, &lastScrape)
	if err != nil {
		return model.Rate{}, time.Time{}, false
	}
	return rate, lastScrape, true
}

func (sm *ScraperManager) GetAllRates() ([]model.Rate, error) {
	rows, err := sm.DB.Query(`SELECT project, input_symbol, output_token, pool_name, apy, project_link, points FROM scrape_cache`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var rates []model.Rate
	for rows.Next() {
		var rate model.Rate
		err := rows.Scan(&rate.ProjectName, &rate.InputSymbol, &rate.OutputToken, &rate.PoolName, &rate.APY, &rate.ProjectLink, &rate.Points)
		if err != nil {
			continue
		}
		rates = append(rates, rate)
	}
	return rates, nil
}

func (sm *ScraperManager) SetCachedRateFull(id string, rate model.Rate, t time.Time) error {
	_, err := sm.DB.Exec(`
		INSERT INTO scrape_cache (id, project, input_symbol, output_token, pool_name, apy, project_link, points, last_scrape)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET 
			project=excluded.project,
			input_symbol=excluded.input_symbol,
			output_token=excluded.output_token,
			pool_name=excluded.pool_name,
			apy=excluded.apy,
			project_link=excluded.project_link,
			points=excluded.points,
			last_scrape=excluded.last_scrape
	`, id, rate.ProjectName, rate.InputSymbol, rate.OutputToken, rate.PoolName, rate.APY, rate.ProjectLink, rate.Points, t)
	return err
}
