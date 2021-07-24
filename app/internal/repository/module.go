package repository

import (
	"context"
	"github.com/chain-bot/prices/app/configs"
	"github.com/chain-bot/prices/app/pkg/models"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/jmoiron/sqlx"
	"github.com/volatiletech/null/v8"
	"time"
)

type Repository interface {
	GetLastSync(
		ctx context.Context,
		exchange string,
		pair *models.Symbol,
	) (*models.LastSync, error)
	UpsertLastSync(
		ctx context.Context,
		exchange string,
		pair *models.Symbol,
		lastSyncTime time.Time,
	) error
	UpsertOHLCVData(
		ohlcvData []*models.OHLCVMarketData,
		exchange string,
		pair *models.Symbol,
	)
	GetOHLCVData(
		base string,
		quote null.String,
		exchange null.String,
		start time.Time,
		end time.Time,
	) ([]models.OHLCVMarketData, error)
}

func NewRepository(
	config *configs.Secrets,
	db *sqlx.DB,
	influxClient *influxdb2.Client,
) Repository {
	return NewRepositoryImpl(config, db, influxClient)
}
