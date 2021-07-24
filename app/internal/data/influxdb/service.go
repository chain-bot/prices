package influxdb

import (
	"fmt"

	"github.com/chain-bot/prices/app/configs"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

func NewInfluxDBClient(secrets *configs.Secrets) (*influxdb2.Client, error) {
	influxDBURL := fmt.Sprintf("http://%s:%d", secrets.InfluxDbCredentials.Host, secrets.InfluxDbCredentials.Port)
	client := influxdb2.NewClient(influxDBURL, secrets.InfluxDbCredentials.Token)
	return &client, nil
}
