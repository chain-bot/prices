package app

import (
	"context"
	"database/sql"
	"github.com/mochahub/coinprice-scraper/scraper/repository"
	"github.com/mochahub/coinprice-scraper/scraper/service/api"
	"log"
	"sync"
	"time"
)

// Get Prices from last_sync.last_sync to time.now()
func StartScraper(
	ctx context.Context,
	repo repository.Repository,
	clients []api.ExchangeAPIClient) error {
	log.Println("Starting Scraper")
	defer log.Println("Stopping Scraper")
	var waitGroup sync.WaitGroup

	for index := range clients {
		client := clients[index]
		// TODO(Zahin): Handle Errors
		waitGroup.Add(1)
		go func(wg *sync.WaitGroup) {
			// Decrement the counter when the goroutine completes.
			defer wg.Done()
			err := ScrapeExchange(ctx, repo, client)
			if err != nil {
				log.Printf("ERROR SCRAPING %s\n", client.GetExchangeIdentifier())
			}
		}(&waitGroup)
		//go ScrapeExchange(ctx, secrets, db, influxDBClient, clients[i])
	}
	waitGroup.Wait()
	return nil
}

func ScrapeExchange(
	ctx context.Context,
	repo repository.Repository,
	client api.ExchangeAPIClient,
) error {
	pairs, err := client.GetSupportedPairs()
	if err != nil {
		return err
	}
	// Default start time is 1 day prior for ease of local development
	//startTime := time.Now().AddDate(0, 0, -1)
	startTime, _ := time.Parse(time.RFC3339, "2016-01-01T00:00:00+00:00")
	endTime := time.Now()
	for index := range pairs {
		pair := pairs[index]
		lastSync, err := repo.GetLastSync(ctx, client.GetExchangeIdentifier(), pair)
		if err != nil && err != sql.ErrNoRows {
			return err
		}
		if lastSync != nil && !lastSync.LastSyncTime.IsZero() {
			startTime = lastSync.LastSyncTime
		}
		ohlcData, err := client.GetAllOHLCMarketData(pair.RawBase, pair.RawQuote, time.Minute, startTime, endTime)
		if err != nil {
			return err
		}
		if len(ohlcData) == 0 {
			continue
		}
		// Update Last Sync, Upsert candle data
		if err := repo.UpsertLastSync(
			ctx,
			client.GetExchangeIdentifier(),
			pair,
			ohlcData[len(ohlcData)-1].StartTime); err != nil {
			return err
		}
		go repo.UpsertOHLCData(ohlcData, client.GetExchangeIdentifier(), pair)
	}
	return nil
}
