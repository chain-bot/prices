package scraper

import (
	"context"
	"github.com/mochahub/coinprice-scraper/app/internal/repository"
	"github.com/mochahub/coinprice-scraper/app/pkg/models"
	"github.com/robfig/cron"
)

func StartScrapperCron(
	ctx context.Context,
	repo repository.Repository,
	clients []models.ExchangeAPIClient,
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
