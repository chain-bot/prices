package impl

import (
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/jmoiron/sqlx"
	"github.com/mochahub/coinprice-scraper/config"
)

type RepositoryImpl struct {
	db           *sqlx.DB
	influxClient *influxdb2.Client
	influxOrg    string
	ohlcBucket   string
}

func NewRepositoryImpl(
	config *config.Secrets,
	db *sqlx.DB,
	influxClient *influxdb2.Client,
) *RepositoryImpl {
	return &RepositoryImpl{
		db:           db,
		influxClient: influxClient,
		influxOrg:    config.Org,
		ohlcBucket:   config.Bucket,
	}
}
