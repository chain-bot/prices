package app

import (
	"context"
	"github.com/mochahub/coinprice-scraper/scraper/repository"
	"github.com/mochahub/coinprice-scraper/scraper/service/api"
	"github.com/robfig/cron"
)

func StartScrapperCron(
	ctx context.Context,
	repo repository.Repository,
	clients []api.RestExchangeAPIClient,
) (*cron.Cron, error) {
	c := cron.New()
	err := c.AddFunc("@every 1m", func() {
		_ = StartScraper(ctx, repo, clients)
	})
	if err != nil {
		return nil, err
	}
	c.Start()
	return c, nil
}
