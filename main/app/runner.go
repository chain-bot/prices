package app

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/mochahub/coinprice-scraper/main/api"
	"github.com/mochahub/coinprice-scraper/main/api/common"
	"github.com/mochahub/coinprice-scraper/main/utils"
	models "github.com/mochahub/coinprice-scraper/models/generated"
	"log"
	"time"
)

// Get Prices from last_sync.last_sync to time.now()
func StartScraper(ctx context.Context, db *sqlx.DB, clients []api.ExchangeAPIClient) error {
	log.Println("Starting Scraper")
	defer log.Println("Stopping Scraper")

	for i := range clients {
		client := clients[i]
		//tx, err := db.BeginTx(ctx, nil)
		//if err != nil {
		//	return err
		//}
		// TODO(Zahin): Handle Errors
		go ScrapeClient(ctx, db, client)
		//if err := tx.Commit(); err != nil {
		//	return err
		//}
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

func ScrapeClient(ctx context.Context, db *sqlx.DB, client api.ExchangeAPIClient) error {
	pairs, err := client.GetSupportedPairs()
	if err != nil {
		return err
	}
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	for index := range pairs {
		pair := pairs[index]
		lastSync, err := models.FindLastSync(ctx, tx, pair.NormalizedBase, pair.NormalizedQuote, client.GetExchangeIdentifier())
		if err != nil && err != sql.ErrNoRows {
			return err
		}
		startTime := time.Now().AddDate(0, 0, -1)
		if lastSync != nil && lastSync.LastSync.Valid {
			startTime = lastSync.LastSync.Time
		}
		ohlcData, err := client.GetAllOHLCMarketData(pair.RawBase, pair.RawQuote, common.Minute, startTime, time.Time{})
		if err != nil {
			return err
		}
		fmt.Println(utils.PrettyJSON(ohlcData))
		// when saving data use the Tx
	}
	return tx.Commit()
}
