package influxdb

import (
	"fmt"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/mochahub/coinprice-scraper/app/configs"
)

func NewInfluxDBClient(secrets *configs.Secrets) (*influxdb2.Client, error) {
	influxDBURL := fmt.Sprintf("http://%s:%d", secrets.InfluxDbCredentials.Host, secrets.InfluxDbCredentials.Port)
	client := influxdb2.NewClient(influxDBURL, secrets.InfluxDbCredentials.Token)
	return &client, nil
}
