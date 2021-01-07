package app

import (
	"context"
	"database/sql"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/jmoiron/sqlx"
	"github.com/mochahub/coinprice-scraper/config"
	models "github.com/mochahub/coinprice-scraper/main/models/generated"
	"github.com/mochahub/coinprice-scraper/main/src/service/api"
	"github.com/mochahub/coinprice-scraper/main/src/service/api/common"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"log"
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

	for i := range clients {
		// TODO(Zahin): Handle Errors
		go ScrapeExchange(ctx, secrets, db, influxDBClient, clients[i])
	}
	// TODO: Move Below Comments to a Space Document
	// General Flow:
	// - Loop Through Implemented Exchanges
	// - Each Exchange is a goroutine
	// - Each goroutine writes to influx independently (conflicts will take latest price?)
	//	- Each exchange writes to their own bucket with expiry
	//	- A separate bucket aggregates price data from exchange buckets periodically
	// - TODO: Interface File that Each Exchange Inherits
	//	- Implement Influx Config and Functions here to abstract from exchange
	return nil
}

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
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	// Default Time if pair was never synced: Thursday, December 31, 2015 7:00:00 PM GMT-05:00
	startTime := time.Unix(1451606400, 0)
	endTime := time.Now()
	for index := range pairs {
		pair := pairs[index]
		lastSync, err := models.FindLastSync(ctx, tx, pair.NormalizedBase, pair.NormalizedQuote, client.GetExchangeIdentifier())
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
		lastSyncTime := ohlcData[len(ohlcData)-1].StartTime
		if err = upsertLastSyncWithTx(ctx, tx, pair, client.GetExchangeIdentifier(), lastSyncTime); err != nil {
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
	return tx.Commit()
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
		_ = tx.Rollback()
		return err
	}
	return nil
}

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
