package influxdb

import (
	"fmt"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/mochahub/coinprice-scraper/app/configs"
)

func NewInfluxDBClient(secrets *configs.Secrets) (*influxdb2.Client, error) {
	influxDBURL := fmt.Sprintf("http://%s:%d", secrets.InfluxDbCredentials.Host, secrets.InfluxDbCredentials.Port)
	client := influxdb2.NewClient(influxDBURL, secrets.InfluxDbCredentials.Token)
	// Turn this into a command line script to generate the token
	//lc.Append(fx.Hook{
	//	OnStart: func(ctx context.Context) error {
	//
	//		response, err := client.Setup(
	//			ctx,
	//			secrets.InfluxDbCredentials.User,
	//			secrets.InfluxDbCredentials.Password,
	//			secrets.InfluxDbCredentials.Org,
	//			secrets.InfluxDbCredentials.Bucket,
	//			0)
	//		if err != nil {
	//			return err
	//		}
	//
	//		return nil
	//	},
	//	OnStop: func(ctx context.Context) error {
	//		client.Close()
	//		return nil
	//	},
	//})
	return &client, nil
}
