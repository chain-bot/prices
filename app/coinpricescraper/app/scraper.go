package app

import (
	"context"
	"database/sql"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/jmoiron/sqlx"
	"github.com/mochahub/coinprice-scraper/app/coinpricescraper/service/api"
	"github.com/mochahub/coinprice-scraper/app/coinpricescraper/service/api/common"
	models "github.com/mochahub/coinprice-scraper/app/models/generated"
	"github.com/mochahub/coinprice-scraper/config"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"log"
	"sync"
	"time"
)

// Get Prices from last_sync.last_sync to time.now()
func StartScraper(
	ctx context.Context,
	secrets *config.Secrets,
	db *sqlx.DB,
	influxDBClient *influxdb2.Client,
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
			err := ScrapeExchange(ctx, secrets, db, influxDBClient, client)
			if err != nil {
				log.Printf("ERROR SCRAPING %s\n", client.GetExchangeIdentifier())
			}
		}(&waitGroup)
		//go ScrapeExchange(ctx, secrets, db, influxDBClient, clients[i])
	}
	waitGroup.Wait()
	return nil
}

// TODO(Zahin): Abstract the use of sqlboiler & psql away
// We should be able to switch from sqlboiler/psql without changing this file
func ScrapeExchange(
	ctx context.Context,
	secrets *config.Secrets,
	db *sqlx.DB,
	influxDBClient *influxdb2.Client,
	client api.ExchangeAPIClient,
) error {
	pairs, err := client.GetSupportedPairs()
	if err != nil {
		return err
	}
	// Default Time if pair was never synced: Friday, July 1, 2016 12:00:00 AM
	startTime := time.Unix(1467331200, 0)
	endTime := time.Now()
	for index := range pairs {
		pair := pairs[index]
		lastSync, err := models.FindLastSync(ctx, db, pair.NormalizedBase, pair.NormalizedQuote, client.GetExchangeIdentifier())
		if err != nil && err != sql.ErrNoRows {
			return err
		}
		if lastSync != nil && lastSync.LastSync.Valid {
			startTime = lastSync.LastSync.Time
		}
		ohlcData, err := client.GetAllOHLCMarketData(pair.RawBase, pair.RawQuote, common.Minute, startTime, endTime)
		if err != nil {
			return err
		}
		// Save last sync to psql
		if len(ohlcData) == 0 {
			continue
		}
		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}
		lastSyncTime := ohlcData[len(ohlcData)-1].StartTime
		if err = upsertLastSyncWithTx(ctx, tx, pair, client.GetExchangeIdentifier(), lastSyncTime); err != nil {
			_ = tx.Rollback()
			return err
		}
		if err = tx.Commit(); err != nil {
			return err
		}
		go upsertOHLCData(
			ohlcData,
			*influxDBClient,
			secrets.InfluxDbCredentials.Org,
			secrets.InfluxDbCredentials.Bucket,
			pair.NormalizedBase,
			pair.NormalizedQuote,
			client.GetExchangeIdentifier())
	}
	return nil
}

func upsertLastSyncWithTx(
	ctx context.Context,
	tx *sql.Tx,
	pair *common.Symbol,
	exchange string,
	lastSyncTime time.Time,
) error {
	lastSync := &models.LastSync{
		BaseAsset:  pair.NormalizedBase,
		QuoteAsset: pair.NormalizedQuote,
		Exchange:   exchange,
		LastSync:   null.TimeFrom(lastSyncTime),
	}
	if err := lastSync.Upsert(ctx, tx, true,
		[]string{models.LastSyncColumns.BaseAsset, models.LastSyncColumns.QuoteAsset, models.LastSyncColumns.Exchange},
		boil.Whitelist(models.LastSyncColumns.LastSync),
		boil.Infer()); err != nil {
		return err
	}
	return nil
}

// TODO(Zahin): Abstract methods directly calling influxDB
// We should be able to switch from influxDB without changing this file
func upsertOHLCData(
	ohlcData []*common.OHLCMarketData,
	influxDBClient influxdb2.Client,
	org, bucket, base, quote, exchange string,
) {
	writeAPI := influxDBClient.WriteAPI(org, bucket)
	tags := map[string]string{
		"quote":    quote,
		"exchange": exchange,
	}
	for index := range ohlcData {
		ohlc := ohlcData[index]
		fields := map[string]interface{}{
			"open":   ohlc.OpenPrice,
			"high":   ohlc.HighPrice,
			"low":    ohlc.LowPrice,
			"close":  ohlc.ClosePrice,
			"volume": ohlc.Volume,
		}
		p := influxdb2.NewPoint(
			base,
			tags,
			fields,
			ohlc.StartTime)
		writeAPI.WritePoint(p)
	}
}
