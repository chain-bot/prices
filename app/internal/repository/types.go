package repository

import (
	"github.com/chain-bot/prices/app/configs"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/jmoiron/sqlx"
)

type RepositoryImpl struct {
	db           *sqlx.DB
	influxClient *influxdb2.Client
	influxOrg    string
	ohlcvBucket  string
}

func NewRepositoryImpl(
	config *configs.Secrets,
	db *sqlx.DB,
	influxClient *influxdb2.Client,
) *RepositoryImpl {
	return &RepositoryImpl{
		db:           db,
		influxClient: influxClient,
		influxOrg:    config.Org,
		ohlcvBucket:  config.Bucket,
	}
}
