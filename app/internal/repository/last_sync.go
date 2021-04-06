package repository

import (
	"context"
	generated "github.com/mochahub/coinprice-scraper/app/internal/data/psql/generated"
	"github.com/mochahub/coinprice-scraper/app/pkg/models"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"time"
)

func (repo *RepositoryImpl) GetLastSync(
	ctx context.Context,
	exchange string,
	pair *models.Symbol,
) (*models.LastSync, error) {
	lastSync, err := generated.FindLastSync(ctx, repo.db, pair.NormalizedBase, pair.NormalizedQuote, exchange)
	if err != nil {
		return nil, err
	}
	return &models.LastSync{
		BaseAsset:    lastSync.BaseAsset,
		QuoteAsset:   lastSync.QuoteAsset,
		Exchange:     lastSync.Exchange,
		LastSyncTime: lastSync.LastSync.Time,
	}, nil
}

func (repo *RepositoryImpl) UpsertLastSync(
	ctx context.Context,
	exchange string,
	pair *models.Symbol,
	lastSyncTime time.Time,
) error {
	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	lastSync := &generated.LastSync{
		BaseAsset:  pair.NormalizedBase,
		QuoteAsset: pair.NormalizedQuote,
		Exchange:   exchange,
		LastSync:   null.TimeFrom(lastSyncTime),
	}
	if err := lastSync.Upsert(ctx, tx, true,
		[]string{generated.LastSyncColumns.BaseAsset, generated.LastSyncColumns.QuoteAsset, generated.LastSyncColumns.Exchange},
		boil.Whitelist(generated.LastSyncColumns.LastSync),
		boil.Infer()); err != nil {
		_ = tx.Rollback()
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}
