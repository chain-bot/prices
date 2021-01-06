package app

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/mochahub/coinprice-scraper/main/api"
	cron "github.com/robfig/cron"
	"go.uber.org/fx"
)

func InitScrapper(lc fx.Lifecycle, db *sqlx.DB, clients api.ExchangeClients) {
	var cronJob *cron.Cron
	var err error
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// First Time Run Can Potentially take Hours (back filling market data)
			if err = StartScraper(ctx, db, clients.Clients); err != nil {
				return err
			}
			if cronJob, err = StartScrapperCron(); err != nil {
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
