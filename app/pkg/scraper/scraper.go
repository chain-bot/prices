package scraper

import (
	"context"
	"database/sql"
	"sync"
	"time"

	"github.com/chain-bot/prices/app/configs"
	"github.com/chain-bot/prices/app/internal/repository"
	"github.com/chain-bot/prices/app/pkg/models"
	log "github.com/sirupsen/logrus"
)

// StartScraper Get Prices from last_sync.last_sync to time.now()
func StartScraper(
	ctx context.Context,
	repo repository.Repository,
	clients []models.ExchangeAPIClient,
	secrets *configs.Secrets,
) error {
	log.Info("starting scraper")
	defer log.Printf("stopping Scraper")
	var waitGroup sync.WaitGroup

	for index := range clients {
		client := clients[index]
		// TODO(Zahin): Handle Errors
		waitGroup.Add(1)
		go func(wg *sync.WaitGroup) {
			// Decrement the counter when the goroutine completes.
			defer wg.Done()
			// Default start time is 1 day prior for ease of local development
			startTime := time.Now().AddDate(0, 0, -1)
			if !secrets.IsLocal() {
				startTime = time.Time{}
			}
			//startTime, _ := time.Parse(time.RFC3339, "2021-03-01T00:00:00+00:00")
			endTime := time.Now()
			err := ScrapeExchange(ctx, repo, client, startTime, endTime)
			if err != nil {
				log.WithFields(log.Fields{
					"err":    err.Error(),
					"client": client.GetExchangeIdentifier(),
				}).Errorf("scrape excahnge")
			}
		}(&waitGroup)
	}
	waitGroup.Wait()
	return nil
}

func ScrapeExchange(
	ctx context.Context,
	repo repository.Repository,
	client models.ExchangeAPIClient,
	startTime time.Time,
	endTime time.Time,
) error {
	pairs, err := client.GetSupportedPairs()
	if err != nil {
		return err
	}
	for index := range pairs {
		pair := pairs[index]
		lastSync, err := repo.GetLastSync(ctx, client.GetExchangeIdentifier(), pair)
		if err != nil && err != sql.ErrNoRows {
			return err
		}
		if lastSync != nil && !lastSync.LastSyncTime.IsZero() {
			startTime = lastSync.LastSyncTime
		}
		ohlcvMarketData, err := client.GetAllOHLCVMarketData(*pair, time.Minute, startTime, endTime)
		if err != nil {
			return err
		}
		if len(ohlcvMarketData) == 0 {
			continue
		}
		// Update Last Sync, Upsert candle data
		if err := repo.UpsertLastSync(
			ctx,
			client.GetExchangeIdentifier(),
			pair,
			ohlcvMarketData[len(ohlcvMarketData)-1].StartTime); err != nil {
			return err
		}
		go repo.UpsertOHLCVData(ohlcvMarketData, client.GetExchangeIdentifier(), pair)
	}
	return nil
}
