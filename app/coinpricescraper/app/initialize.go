package app

import (
	"context"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/jmoiron/sqlx"
	"github.com/mochahub/coinprice-scraper/app/coinpricescraper/service/api"
	"github.com/mochahub/coinprice-scraper/config"
	cron "github.com/robfig/cron"
	"go.uber.org/fx"
)

func InitScrapper(
	lc fx.Lifecycle,
	secrets *config.Secrets,
	db *sqlx.DB,
	influxClient *influxdb2.Client,
	clients api.ExchangeClients) {
	var cronJob *cron.Cron
	var err error
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// First Time Run Can Potentially take Hours (back filling market data)
			if err = StartScraper(ctx, secrets, db, influxClient, clients.Clients); err != nil {
				return err
			}
			if cronJob, err = StartScrapperCron(ctx, secrets, db, influxClient, clients.Clients); err != nil {
				return err
			}
			return nil
		},
		OnStop: func(ctx context.Context) error {
			if cronJob != nil {
				cronJob.Stop()
			}
			return nil
		},
	})
}
