package app

import (
	"context"
	"github.com/mochahub/coinprice-scraper/scraper/repository"
	"github.com/mochahub/coinprice-scraper/scraper/service/api"
	cron "github.com/robfig/cron"
	"go.uber.org/fx"
)

func InitScrapper(
	lc fx.Lifecycle,
	repo repository.Repository,
	clients api.ExchangeClients,
) {
	var cronJob *cron.Cron
	var err error
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			// First Time Run Can Potentially take Hours (back filling market data)
			if err = StartScraper(ctx, repo, clients.Clients); err != nil {
				return err
			}
			if cronJob, err = StartScrapperCron(ctx, repo, clients.Clients); err != nil {
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
