package app

import (
	"database/sql"
	"log"
	"sync"
	"time"

	"example.com/rates/v2/internal/model"
	"example.com/rates/v2/internal/scraper"
	_ "github.com/mattn/go-sqlite3"
)

// ScraperManager manages scraping and caching.
type ScraperManager struct {
	DB         *sql.DB
	LastScrape map[string]time.Time
	Mutex      sync.Mutex
	ExchangeDB *sql.DB
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
				rates := scraper.ScrapeRenzo()
				for _, rate := range rates {
					id := rate.ProjectName + ":" + rate.InputSymbol + ":" + rate.PoolName
					_ = sm.SetCachedRateFull(id, rate, time.Now())
				}
			}
			{
				rates, err := scraper.ScrapeKelp()
				if err != nil {
					log.Printf("[scraper] Kelp: error: %v", err)
					continue
				}
				for _, rate := range rates {
					id := rate.ProjectName + ":" + rate.InputSymbol + ":" + rate.PoolName
					_ = sm.SetCachedRateFull(id, rate, time.Now())
				}
			}
			// Eigenpie
			{
				rates, err := scraper.ScrapeEigenpie()
				if err != nil {
					log.Printf("[scraper] Eigenpie: error: %v", err)
				} else {
					for _, rate := range rates {
						id := rate.ProjectName + ":" + rate.InputSymbol + ":" + rate.PoolName
						_ = sm.SetCachedRateFull(id, rate, time.Now())
					}
				}
			}
			// Stakestone
			{
				id := scraper.StakestoneID
				rate, err := scraper.ScrapeStakestone()
				if err == nil {
					_ = sm.SetCachedRateFull(id, rate, time.Now())
				}
			}
			// Lido
			{
				rates, err := scraper.ScrapeLido()
				if err == nil {
					for _, rate := range rates {
						id := rate.ProjectName + ":" + rate.InputSymbol + ":" + rate.PoolName + ":" + rate.OutputToken
						_ = sm.SetCachedRateFull(id, rate, time.Now())
					}
				}
			}
			// Zuit pairs (dummy)
			{
				pairs := scraper.ScrapeZuitPairs()
				for _, pair := range pairs {
					sm.updatePair(pair)
				}
			}
			time.Sleep(10 * time.Minute)
		}
	}()
}

func (sm *ScraperManager) updatePair(pair scraper.ExchangePair) {
	_, err := sm.ExchangeDB.Exec(`
						INSERT INTO exchange_pairs (id, pool, token1, token2, type, liquidity, liquidity_formatted, apr, apr_formatted, exchange_id)
						VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
						ON CONFLICT(id) DO UPDATE SET
							pool=excluded.pool,
							token1=excluded.token1,
							token2=excluded.token2,
							type=excluded.type,
							liquidity=excluded.liquidity,
							liquidity_formatted=excluded.liquidity_formatted,
							apr=excluded.apr,
							apr_formatted=excluded.apr_formatted,
							exchange_id=excluded.exchange_id
					`, pair.ID, pair.Pool, pair.Token1, pair.Token2, pair.Type, pair.Liquidity, pair.LiquidityFormatted, pair.APR, pair.APRFormatted, pair.ExchangeID)
	if err != nil {
		log.Printf("[exchange_pairs] failed to upsert: %v", err)
	}
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
			output_kind TEXT,
			last_scrape TIMESTAMP
		)
	`)
	if err != nil {
		log.Fatalf("failed to create scrape_cache table: %v", err)
	}

	// Open or create exchange_pairs.db
	exdb, err := sql.Open("sqlite3", "exchange_pairs.db")
	if err != nil {
		log.Fatalf("failed to open exchange_pairs.db: %v", err)
	}
	_, err = exdb.Exec(`
		CREATE TABLE IF NOT EXISTS exchange_pairs (
			id INTEGER PRIMARY KEY,
			pool TEXT,
			token1 TEXT,
			token2 TEXT,
			type TEXT,
			liquidity REAL,
			liquidity_formatted TEXT,
			apr REAL,
			apr_formatted TEXT,
			exchange_id TEXT
		)
	`)
	if err != nil {
		log.Fatalf("failed to create exchange_pairs table: %v", err)
	}

	return &ScraperManager{
		DB:         db,
		LastScrape: make(map[string]time.Time),
		ExchangeDB: exdb,
	}
}

// Returns (model.Rate, lastScrape, found)
func (sm *ScraperManager) GetCachedRate(id string) (model.Rate, time.Time, bool) {
	row := sm.DB.QueryRow(`SELECT project, input_symbol, output_token, pool_name, apy, project_link, points, output_kind, last_scrape FROM scrape_cache WHERE id = ?`, id)
	var rate model.Rate
	var lastScrape time.Time
	err := row.Scan(&rate.ProjectName, &rate.InputSymbol, &rate.OutputToken, &rate.PoolName, &rate.APY, &rate.ProjectLink, &rate.Points, &rate.OutputKind, &lastScrape)
	if err != nil {
		return model.Rate{}, time.Time{}, false
	}
	return rate, lastScrape, true
}

func (sm *ScraperManager) GetAllRates() ([]model.Rate, error) {
	rows, err := sm.DB.Query(`SELECT project, input_symbol, output_token, pool_name, apy, project_link, points, output_kind FROM scrape_cache`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var rates []model.Rate
	for rows.Next() {
		var rate model.Rate
		err := rows.Scan(&rate.ProjectName, &rate.InputSymbol, &rate.OutputToken, &rate.PoolName, &rate.APY, &rate.ProjectLink, &rate.Points, &rate.OutputKind)
		if err != nil {
			continue
		}
		rates = append(rates, rate)
	}

	// Also add all pairs from exchange_pairs as rates
	pairRows, err := sm.ExchangeDB.Query(`SELECT token1, token2, exchange_id FROM exchange_pairs`)
	if err == nil {
		defer pairRows.Close()
		for pairRows.Next() {
			var token1, token2, exchangeID string
			err := pairRows.Scan(&token1, &token2, &exchangeID)
			if err != nil {
				continue
			}
			rates = append(rates, model.Rate{
				InputSymbol: token1,
				OutputToken: token2,
				ProjectName: exchangeID,
				PoolName:    "",
				APY:         0,
				ProjectLink: "",
				Points:      "",
				OutputKind:  "swap",
			})
		}
	}

	return rates, nil
}

// GetPointsForRate returns the points string for a given rate from the DB, or "" if not found.
func (sm *ScraperManager) GetPointsForRate(rate model.Rate) string {
	row := sm.DB.QueryRow(`SELECT points FROM scrape_cache WHERE project = ? AND input_symbol = ? AND output_token = ? AND pool_name = ?`,
		rate.ProjectName, rate.InputSymbol, rate.OutputToken, rate.PoolName)
	var points string
	err := row.Scan(&points)
	if err != nil {
		return ""
	}
	return points
}

func (sm *ScraperManager) SetCachedRateFull(id string, rate model.Rate, t time.Time) error {
	log.Printf("[db] Update: input=%s, output=%s, apy=%.4f", rate.InputSymbol, rate.OutputToken, rate.APY)
	_, err := sm.DB.Exec(`
		INSERT INTO scrape_cache (id, project, input_symbol, output_token, pool_name, apy, project_link, points, output_kind, last_scrape)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET 
			project=excluded.project,
			input_symbol=excluded.input_symbol,
			output_token=excluded.output_token,
			pool_name=excluded.pool_name,
			apy=excluded.apy,
			project_link=excluded.project_link,
			points=excluded.points,
			output_kind=excluded.output_kind,
			last_scrape=excluded.last_scrape
	`, id, rate.ProjectName, rate.InputSymbol, rate.OutputToken, rate.PoolName, rate.APY, rate.ProjectLink, rate.Points, rate.OutputKind, t)
	return err
}
