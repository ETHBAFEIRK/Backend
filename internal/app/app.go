package app

import (
	"database/sql"
	"log"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type ScraperManager struct {
	DB         *sql.DB
	LastScrape map[string]time.Time
	Mutex      sync.Mutex
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
			pool_name TEXT,
			apy REAL,
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

func (sm *ScraperManager) GetCachedRate(id string) (float64, time.Time, bool) {
	row := sm.DB.QueryRow("SELECT apy, last_scrape FROM scrape_cache WHERE id = ?", id)
	var apy float64
	var lastScrape time.Time
	err := row.Scan(&apy, &lastScrape)
	if err != nil {
		return 0, time.Time{}, false
	}
	return apy, lastScrape, true
}

func (sm *ScraperManager) SetCachedRate(id, project, inputSymbol, poolName string, apy float64, t time.Time) error {
	_, err := sm.DB.Exec(`
		INSERT INTO scrape_cache (id, project, input_symbol, pool_name, apy, last_scrape)
		VALUES (?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET apy=excluded.apy, last_scrape=excluded.last_scrape
	`, id, project, inputSymbol, poolName, apy, t)
	return err
}
