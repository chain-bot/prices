package repository

import (
	"github.com/chain-bot/prices/app/configs"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/jmoiron/sqlx"
)

type impl struct {
	db           *sqlx.DB
	influxClient *influxdb2.Client
	writeAPI     api.WriteAPI
	queryAPI     api.QueryAPI
	influxOrg    string
	ohlcvBucket  string
}

func NewRepositoryImpl(
	config *configs.Secrets,
	db *sqlx.DB,
	influxClient *influxdb2.Client,
) *impl {
	return &impl{
		db:           db,
		influxClient: influxClient,
		writeAPI:     (*influxClient).WriteAPI(config.Org, config.Bucket),
		queryAPI:     (*influxClient).QueryAPI(config.Org),
		influxOrg:    config.Org,
		ohlcvBucket:  config.Bucket,
	}
}
