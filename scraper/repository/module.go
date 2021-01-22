package repository

import (
	"context"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/jmoiron/sqlx"
	"github.com/mochahub/coinprice-scraper/config"
	"github.com/mochahub/coinprice-scraper/scraper/models"
	"github.com/mochahub/coinprice-scraper/scraper/repository/impl"
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
	config *config.Secrets,
	db *sqlx.DB,
	influxClient *influxdb2.Client,
) Repository {
	return impl.NewRepositoryImpl(config, db, influxClient)
}
