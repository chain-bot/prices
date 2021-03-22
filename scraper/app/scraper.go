package app

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/mochahub/coinprice-scraper/scraper/repository"
	"github.com/mochahub/coinprice-scraper/scraper/service/api"
	"log"
	"sync"
	"time"
)

// Get Prices from last_sync.last_sync to time.now()
func StartRestScraper(
	ctx context.Context,
	repo repository.Repository,
	clients []api.RestExchangeAPIClient,
) error {
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
			// Default start time is 1 day prior for ease of local development
			//startTime := time.Now().AddDate(0, 0, -1)
			startTime, _ := time.Parse(time.RFC3339, "2021-01-01T00:00:00+00:00")
			endTime := time.Now()
			err := ScrapeExchange(ctx, repo, client, startTime, endTime)
			if err != nil {
				log.Printf("ERROR SCRAPING %s\n", client.GetExchangeIdentifier())
			}
		}(&waitGroup)
	}
	waitGroup.Wait()
	return nil
}

func ScrapeExchange(
	ctx context.Context,
	repo repository.Repository,
	client api.RestExchangeAPIClient,
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
		ohlcData, err := client.GetAllOHLCMarketData(*pair, time.Minute, startTime, endTime)
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
		go repo.UpsertOHLCData(client.GetExchangeIdentifier(), pair, ohlcData...)
	}
	return nil
}

func StartSocketScraper(
	ctx context.Context,
	repo repository.Repository,
	socketClients []api.SocketExchangeAPIClient,
) error {
	log.Println("Starting Socket Scraper")
	defer log.Println("Stopping Socket Scraper")
	for index := range socketClients {
		client := socketClients[index]
		go func() {
			ScrapeSocketExchange(ctx, repo, client)
		}()
	}
	<-ctx.Done()
	return nil
}
func ScrapeSocketExchange(
	ctx context.Context,
	repo repository.Repository,
	client api.SocketExchangeAPIClient,
) error {
	pairs, err := client.GetSupportedPairs()
	if err != nil {
		return err
	}
	for index := range pairs {
		pair := pairs[index]
		go func() {
			ohlcChannel, err := client.GetOHLCMarketDataChannel(ctx, *pair, time.Minute)
			if err != nil {
				log.Println(fmt.Sprintf("Client Error %s", err))
				close(ohlcChannel)
				return
			}
			for {
				ohlcData := <-ohlcChannel
				go repo.UpsertOHLCData(client.GetExchangeIdentifier(), pair, ohlcData)
			}
		}()
	}
	<-ctx.Done()
	return nil
}
