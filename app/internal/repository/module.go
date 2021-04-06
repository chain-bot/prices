package repository

import (
	"context"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/jmoiron/sqlx"
	"github.com/mochahub/coinprice-scraper/app/configs"
	"github.com/mochahub/coinprice-scraper/app/pkg/models"
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
	UpsertOHLCData(
		ohlcData []*models.OHLCMarketData,
		exchange string,
		pair *models.Symbol,
	)
}

func NewRepository(
	config *configs.Secrets,
	db *sqlx.DB,
	influxClient *influxdb2.Client,
) Repository {
	return NewRepositoryImpl(config, db, influxClient)
}
