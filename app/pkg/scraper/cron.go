package scraper

import (
	"context"
	"github.com/chain-bot/prices/app/configs"
	"github.com/chain-bot/prices/app/internal/repository"
	"github.com/chain-bot/prices/app/pkg/models"
	"github.com/robfig/cron"
)

func StartScrapperCron(
	ctx context.Context,
	repo repository.Repository,
	clients []models.ExchangeAPIClient,
	secrets *configs.Secrets,
) (*cron.Cron, error) {
	c := cron.New()
	err := c.AddFunc("@every 1m", func() {
		_ = StartScraper(ctx, repo, clients, secrets)
	})
	if err != nil {
		return nil, err
	}
	c.Start()
	return c, nil
}
