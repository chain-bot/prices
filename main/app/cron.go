package app

import (
	"context"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/jmoiron/sqlx"
	"github.com/mochahub/coinprice-scraper/config"
	"github.com/mochahub/coinprice-scraper/main/src/service/api"
	"github.com/robfig/cron"
)

func StartScrapperCron(
	ctx context.Context,
	secrets *config.Secrets,
	db *sqlx.DB,
	influxDBClient *influxdb2.Client,
	clients []api.ExchangeAPIClient,
) (*cron.Cron, error) {
	c := cron.New()
	err := c.AddFunc("@every 1m", func() {
		_ = StartScraper(ctx, secrets, db, influxDBClient, clients)
	})
	if err != nil {
		return nil, err
	}
	c.Start()
	return c, nil
}

func StartScraperWrapper() {

}
